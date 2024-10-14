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
		{"sg", "Bɛ̂-bïn", []SSE{0xbed3, 0x94dd}},
		{"sg", "bəbị", []SSE{0xa8db, 0x8ad5}},
		{"en", "Hello", []SSE{0x4448, 0x4365, 0x426c, 0x416c, 0x406f}},
		{"fr", "c'est", []SSE{0x6463, 0x6327, 0x6265, 0x6173, 0x6074}},
		{"", "?!...$", []SSE{0x003f, 0x0021, 0x2026, 0x0024}},
	}
)

func TestEncodeWord(t *testing.T) {
	for _, v := range encodeWordTestCases {
		actually := encodeWord(v.lang, []byte(v.word))
		nActual := len(actually)
		nExpect := len(v.expect)
		if nActual != nExpect {
			t.Errorf("word        = %q\n", v.word)
			t.Errorf("len(actual) = %v\n", nActual)
			t.Errorf("len(expect) = %v\n", nExpect)
		}
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
	expect := `There are 62 tokens:
#0: 0b_1_00_11_10_01100_0101
#1: 0b_00_00000000100000
#2: 0b_1_01_01_01_11001_1011
#3: 0b_1_00_01_00_11000_0011
#4: 0b_00_00000000111010
#5: 0b_00_00000000100000
#6: 0b_00_00000001010100
#7: 0b_00_00000001101000
#8: 0b_00_00000001100101
#9: 0b_00_00000000100000
#10: 0b_00_00000001110000
#11: 0b_00_00000001101000
#12: 0b_00_00000001110010
#13: 0b_00_00000001100001
#14: 0b_00_00000001110011
#15: 0b_00_00000001100101
#16: 0b_00_00000000100000
#17: 0b_00_00000010101011
#18: 0b_1_01_01_00_00000_0100
#19: 0b_1_00_01_10_10011_1110
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
#57: 0b_1_01_11_00_11001_0100
#58: 0b_1_00_01_11_00000_0100
#59: 0b_00_00000000100000
#60: 0b_1_01_01_10_11001_0011
#61: 0b_1_00_01_00_11000_0011
`
	if actual != expect {
		t.Errorf("actual: %s\n", actual)
		t.Errorf("expect: %s\n", expect)
	}
}
