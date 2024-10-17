package transcode

import (
	"bufio"
	"bytes"
	"log"
	"strings"
	"testing"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

var (
	encodeWordTestCases = []struct {
		lang   string
		word   string
		expect []SSE
	}{
		{"sg", "Bɛ̂-bïn", []SSE{
			0b_1_01_01_11_01101_0011,
			0b_1_00_10_10_01101_1101,
		}},
		{"sg", "bəbị", []SSE{
			0b_1_01_00_00_01101_1011,
			0b_1_00_00_01_01101_0101,
		}},
		{"en", "Hello", []SSE{
			0b_01_0_00100_0100_1000,
			0b_01_0_00011_0110_0101,
			0b_01_0_00010_0110_1100,
			0b_01_0_00001_0110_1100,
			0b_01_0_00000_0110_1111,
		}},
		{"fr", "c'est", []SSE{
			0b_01_1_00100_0110_0011,
			0b_01_1_00011_0010_0111,
			0b_01_1_00010_0110_0101,
			0b_01_1_00001_0111_0011,
			0b_01_1_00000_0111_0100,
		}},
		{"", "?!...$", []SSE{
			0b_00_00000000111111,
			0b_00_00000000100001,
			0b_00_10000000100110,
			0b_00_00000000100100,
		}},
	}
)

func TestEncodeWord(t *testing.T) {
	for _, v := range encodeWordTestCases {
		actually := encodeWord([]byte(v.word), v.lang)
		nActual := len(actually)
		nExpect := len(v.expect)
		if nActual != nExpect {
			t.Errorf("word        = %q\n", v.word)
			t.Errorf("len(actual) = %v\n", nActual)
			t.Errorf("len(expect) = %v\n", nExpect)
		}
		reactualWord := ""
		reexpectWord := ""
		prevLanguageCode := ""
		for k := range max(nActual, nExpect) {
			if k < nActual {
				if k < nExpect {
					if actually[k] != v.expect[k] {
						t.Errorf("word       = %q\n", v.word)
						t.Errorf("actual[%v] = %04x = %016b\n", k, actually[k], actually[k])
						t.Errorf("expect[%v] = %04x = %016b\n", k, v.expect[k], v.expect[k])
					}
				} else {
					t.Errorf("word      = %q\n", v.word)
					t.Errorf("actual[%v] = %04x = %016b\n", k, actually[k], actually[k])
					t.Errorf("expect[%v] not defined\n", k)
				}
			} else if k < nExpect {
				t.Errorf("word       = %q\n", v.word)
				t.Errorf("actual[%v] not defined\n", k)
				t.Errorf("expect[%v] = %04x = %016b\n", k, v.expect[k], v.expect[k])
			}
			sse := v.expect[k]
			serialized, languageCode, numSyllablesLeft := decodeSSE(sse)
			if languageCode != v.lang {
				t.Errorf("reactual[%v].lang = %q\n", k, languageCode)
				t.Errorf("reexpect[%v].lang = %q\n", k, v.lang)
			}
			reactualWord += string(serialized)
			reexpectWord = v.word
			if languageCode != "" && numSyllablesLeft > 0 ||
				languageCode == "" && prevLanguageCode == "" {
				prevLanguageCode = languageCode
				continue
			}
			prevLanguageCode = languageCode
			if reactualWord != v.word {
				t.Errorf("reactual[%v].word = %q\n", k, reactualWord)
				t.Errorf("reexpect[%v].word = %q\n", k, v.word)
			}
			reactualWord = ""
		}
		if reactualWord != "" {
			reexpectWord = strings.ReplaceAll(reexpectWord, "...", "…")
			reexpectWord = strings.ReplaceAll(reexpectWord, "<<", "«")
			reexpectWord = strings.ReplaceAll(reexpectWord, ">>", "»")
			reexpectWord = strings.ReplaceAll(reexpectWord, "``", "“")
			reexpectWord = strings.ReplaceAll(reexpectWord, "''", "”")
			reexpectWord = strings.ReplaceAll(reexpectWord, "---", "—")
			reexpectWord = strings.ReplaceAll(reexpectWord, "--", "–")
			if reactualWord != reexpectWord {
				t.Errorf("reactual[final].word = %q\n", reactualWord)
				t.Errorf("reexpect[final].word = %q\n", reexpectWord)
			}
		}
	}
}

func TestEncode(t *testing.T) {
	var b bytes.Buffer
	phrase := "Mbï tə̣nɛ: The phrase <<ahön ndö nî>> means...``exceeding all else''. Taâ tɛ̈nɛ"
	in := bufio.NewReader(strings.NewReader(phrase))
	out := bufio.NewWriter(&b)
	err := Encode(out, in)
	if err != nil {
		t.Errorf("error = %v", err)
	}
	actual := b.String()
	expect := `There are 62 tokens with one of the following binary formats:
0b_00_UUUUUUUUUUUUUU     = Unicode rune (each rune is its own word)
-------------------------
0b    UUUUUUUUUUUUUU     = Unicode rune value (U+0000 - U+3FFF)

0b_01_L_NNNNN_AAAAAAAA   = ASCII character (English or French only)
-------------------------
0b    L                  = Language: 0=English, 1=French
0b      NNNNN            = min(31,n), n = # runes left excluding this one
0b            AAAAAAAA   = ASCII letter value (U+00 - U+FF)

0b_1_SS_XX_PP_CCCCC_VVVV = Syllable (Sango only)
-------------------------
0b   SS                  = min(3,m), m = # syllables left excluding this one
0b      XX               = Case : 00=lowercase, 01=Titlecase, 10=-prefixed, 11=UPPERCASE
0b         PP            = Pitch: 00=Unknown, 01=LowTone  , 10=MidTone  , 11=HighTone
0b            CCCCC      = Consonant (first 3 on left below, last 2 on top)
0b                  VVVV = Vowel     (first 2 on left below, last 2 on top)

#00: 0b_1_00_01_10_01100_0101
#01: 0b_00_00000000100000
#02: 0b_1_01_00_01_11001_1011
#03: 0b_1_00_00_00_11000_0011
#04: 0b_00_00000000111010
#05: 0b_00_00000000100000
#06: 0b_00_00000001010100
#07: 0b_00_00000001101000
#08: 0b_00_00000001100101
#09: 0b_00_00000000100000
#10: 0b_00_00000001110000
#11: 0b_00_00000001101000
#12: 0b_00_00000001110010
#13: 0b_00_00000001100001
#14: 0b_00_00000001110011
#15: 0b_00_00000001100101
#16: 0b_00_00000000100000
#17: 0b_00_00000010101011
#18: 0b_1_01_00_00_00000_0100
#19: 0b_1_00_00_10_10011_1110
#20: 0b_00_00000000100000
#21: 0b_00_00000001101110
#22: 0b_00_00000001100100
#23: 0b_00_00000011110110
#24: 0b_00_00000000100000
#25: 0b_00_00000001101110
#26: 0b_00_00000011101110
#27: 0b_00_00000010111011
#28: 0b_00_00000000100000
#29: 0b_00_00000001101101
#30: 0b_00_00000001100101
#31: 0b_00_00000001100001
#32: 0b_00_00000001101110
#33: 0b_00_00000001110011
#34: 0b_00_10000000100110
#35: 0b_00_10000000011100
#36: 0b_00_00000001100101
#37: 0b_00_00000001111000
#38: 0b_00_00000001100011
#39: 0b_00_00000001100101
#40: 0b_00_00000001100101
#41: 0b_00_00000001100100
#42: 0b_00_00000001101001
#43: 0b_00_00000001101110
#44: 0b_00_00000001100111
#45: 0b_00_00000000100000
#46: 0b_00_00000001100001
#47: 0b_00_00000001101100
#48: 0b_00_00000001101100
#49: 0b_00_00000000100000
#50: 0b_00_00000001100101
#51: 0b_00_00000001101100
#52: 0b_00_00000001110011
#53: 0b_00_00000001100101
#54: 0b_00_10000000011101
#55: 0b_00_00000000101110
#56: 0b_00_00000000100000
#57: 0b_1_01_01_00_11001_0100
#58: 0b_1_00_00_11_00000_0100
#59: 0b_00_00000000100000
#60: 0b_1_01_00_10_11001_0011
#61: 0b_1_00_00_00_11000_0011
`
	if actual != expect {
		t.Errorf("actual: %s\n", actual)
		t.Errorf("expect: %s\n", expect)
	}
}
