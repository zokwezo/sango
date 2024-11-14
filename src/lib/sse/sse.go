package sse

import (
	"bytes"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/rivo/uniseg"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/unicode/norm"
)

//////////////////////////////////////////////////////////////////////////////
// TYPES

type SSEtoken uint16

type SSE struct {
	t SSEtoken // see README.md
	b []byte   // UTF8
	g [][]byte // UTF8/glyph
	r []rune   // Unicode
	s string   // UTF8 as a string
	l string   // language code: one of {"", "en", "fr", "sg"}
	n int      // num SSEtokens left
}

type UnicodeSSE struct {
	SSE // invariant: t >> 14 & 3 == 0
}

type AsciiSSE struct {
	SSE // invariant: t >> 14 & 3 == 1
}

type SangoSSE struct {
	SSE // invariant: t  >> 15 & 1 == 1
}

const (
	// CaseEnum = SSEtoken >> 11 & 3
	// SSEtoken = CaseEnum & 3 << 11
	SangoCaseLower  = 0
	SangoCaseTitle  = 1
	SangoCaseHyphen = 2
	SangoCaseUpper  = 3

	// PitchEnum = SSEtoken & 3
	// SSEtoken = PitchEnum & 3
	SangoPitchUnknown = 0
	SangoPitchLow     = 1
	SangoPitchMid     = 2
	SangoPitchHigh    = 3
)

type SSEInterface interface {
	SSEToken() SSEtoken
	Bytes() []byte
	Glyphs() [][]byte
	Runes() []rune
	FullString() string
	String() string
	LanguageCode() string
	NumTokensLeftInWord() int
}

type SangoSSEInterface interface {
	Case() int
	Pitch() int
	Consonants() []rune
	Vowel() []rune
}

//////////////////////////////////////////////////////////////////////////////
// CONSTRUCTORS

func (t SSEtoken) AsUnicodeSSE() UnicodeSSE {
	// t = 0b_00_UUUUUUUUUUUUUU where
	// U = Unicode (U+0000 - U+3FFF)
	if t&0b_11_00000000000000 != 0b_00_00000000000000 {
		panic("SSEtoken is not a UnicodeSSE")
	}
	var sse UnicodeSSE
	sse.t = t
	sse.r = append(sse.r, rune(t&0b_00_11111111111111))
	for _, r := range sse.r {
		sse.b = utf8.AppendRune(sse.b, r)
	}
	sse.g = append(sse.g, sse.b)
	sse.s = string(sse.b)
	sse.n = 0
	return sse
}

func (t SSEtoken) AsAsciiSSE() AsciiSSE {
	// t = 0b_01_L_NNNNN_AAAAAAAA where
	// L = Language: 0 = English (en), 1 = French (fr)
	// N = min(31,m), m = # letters left excluding this one
	// A = ASCII code (U+0000 - U+00FF)
	if t&0b_11_0_00000_00000000 != 0b_01_0_00000_00000000 {
		panic("SSEtoken is not a AsciiSSE")
	}
	var sse AsciiSSE
	sse.t = t
	sse.r = append(sse.r, rune(t&0b_000_00000_11111111))
	for _, r := range sse.r {
		sse.b = utf8.AppendRune(sse.b, r)
	}
	sse.g = append(sse.g, sse.b)
	sse.s = string(sse.b)
	switch t & 0b_00_1_00000_00000000 {
	case 0b_00_0_00000_00000000:
		sse.l = "en"
	case 0b_00_1_00000_00000000:
		sse.l = "fr"
	}
	sse.n = int(t) & 0b_00_0_11111_00000000 >> 8
	return sse
}

func (t SSEtoken) AsSangoSSE() SangoSSE {
	// t = 0b_1_SS_XX_CCCCC_VVVV_PP where
	// S = min(3,m), m = # syllables left excluding this one
	// X = Case : 00=lowercase, 01=Titlecase, 10=-prefixed, 11=UPPERCASE
	// C = Consonant (first 3 on left below, last 2 on top)
	// V = Vowel     (first 2 on left below, last 2 on top)
	// P = Pitch: 00=Unknown, 01=LowPitch  , 10=MidPitch  , 11=HighPitch
	// where CCCCC and VVVV are set as follows (MSB on left, LSB on top):
	// +-----+----+----+----+----+
	// | Bit | 00 | 01 | 10 | 11 |
	// +-----+----+----+----+----+
	// | 000 |    | b  | d  | f  |
	// | 001 | g  | gb | h  | k  |
	// | 010 | kp | l  | m  | mb |
	// | 011 | mp | mv | n  | nd |
	// | 100 | ng | ngb| ny | nz |
	// | 101 | p  | r  | s  | t  |
	// | 110 | v  | w  | y  | z  |
	// +-----+----+----+----+----+
	// |  00 |    | a  | an | ə  |
	// |  01 | ɛ  | e  | en | i  |
	// |  10 | in | ø  | ɔ  | o  |
	// |  11 | on | u  | un | —— |
	// +-----+----+----+----+----+
	if t&0b_1_00_00_00_00000_0000 != 0b_1_00_00_00000_0000_00 {
		panic("SSEtoken is not a SangoSSE")
	}
	var sse SangoSSE
	sse.t = t
	sse.n = int(t) & 0b_0_11_00_00000_0000_00 >> 13
	sse.l = "sg"
	sse.r = []rune{}
	if cons, found := sseTokenToCons[t&0b_0_00_00_11111_0000_00]; found {
		sse.r = append(sse.r, cons...)
	}
	if vowel, found := sseTokenToVowel[t&0b_0_00_00_00000_1111_11]; found {
		sse.r = append(sse.r, vowel...)
	}
	if len(sse.r) > 0 {
		switch t & 0b_0_00_11_00000_0000_00 >> 11 {
		case 0: // lowercase
		case 1: // Titlecase
			sse.r[0] = unicode.ToUpper(sse.r[0])
		case 2: // hyphen-prefixed
			sse.r = append([]rune{'-'}, sse.r...)
		case 3: // UPPERCASE
			for k, _ := range sse.r {
				sse.r[k] = unicode.ToUpper(sse.r[k])
			}
		}
	}
	for _, r := range sse.r {
		sse.b = utf8.AppendRune(sse.b, r)
	}
	sse.s = string(sse.b)
	grapheme := uniseg.NewGraphemes(sse.s)
	for grapheme.Next() {
		sse.g = append(sse.g, grapheme.Bytes())
	}
	// Verify that the glyphs comprise the byte string.
	// TODO: delete once this works
	bb := bytes.Join(sse.g, []byte{})
	if len(bb) != len(sse.b) {
		panic("len(flatten(sse.g)) != len(sse.b)")
	}
	for k, v := range bb {
		if v != sse.b[k] {
			panic("flatten(sse.g) != sse.b for some element")
		}
	}
	return sse
}

func (t SSEtoken) AsSSE() SSEInterface {
	switch {
	case t < 16364:
		return t.AsUnicodeSSE()
	case t < 32768:
		return t.AsAsciiSSE()
	}
	return t.AsSangoSSE()
}

func MakeUnicodeSSE(r rune) UnicodeSSE {
	return SSEtoken(max(0, r) & 0b_00_11111111111111).AsUnicodeSSE()
}

func MakeAsciiSSE(r rune, isFrench bool, numLettersLeft int) AsciiSSE {
	s := 0b_00_0_00000_11111111 & uint16(max(0, r))
	s |= 0b_01_0_00000_00000000
	s |= uint16(min(31, max(0, numLettersLeft)) << 8)
	if isFrench {
		s |= 0b_00_1_00000_00000000
	}
	return SSEtoken(s).AsAsciiSSE()
}

func MakeSangoSSE(syllable string, numLettersLeft int) (sangoSSE SangoSSE, wordLeft string) {
	var t SSEtoken
	syllable = norm.NFD.String(syllable)
	wordLeft, t = encodeLastSyllable(syllable, numLettersLeft)
	return t.AsSangoSSE(), norm.NFC.String(wordLeft)
}

func EncodeSangoWord(word string) []SangoSSE {
	sangoSSEs := []SangoSSE{}
	if len(word) != 0 {
		word = norm.NFD.String(word)
		for numSyllablesLeft := 0; len(word) > 0; numSyllablesLeft++ {
			var sse SSEtoken
			word, sse = encodeLastSyllable(word, numSyllablesLeft)
			sangoSSEs = append(sangoSSEs, sse.AsSangoSSE())
		}
		slices.Reverse(sangoSSEs)
	}
	return sangoSSEs
}

//////////////////////////////////////////////////////////////////////////////
// GETTERS

func (sse SSE) SSEToken() SSEtoken {
	return sse.t
}

func (sse SSE) Bytes() []byte {
	return sse.b
}

func (sse SangoSSE) Case() int {
	if sse.t>>15&1 != 0 {
		return -1
	}
	return int(sse.t >> 11 & 3)
}

func (sse SangoSSE) Pitch() int {
	if sse.t>>15&1 != 0 {
		return -1
	}
	return int(sse.t & 3)
}

func (sse SangoSSE) Consonants() []rune {
	if sse.t>>15&1 != 0 {
		return []rune{}
	}
	return []rune{rune(sse.t >> 6 & 31)}
}

func (sse SangoSSE) Vowel() []rune {
	if sse.t>>15&1 != 0 {
		return []rune{}
	}
	return sseTokenToVowel[sse.t&0b_0_00_00_00000_1111_11]
}

func (sse SSE) Glyphs() [][]byte {
	return sse.g
}

func (sse SSE) Runes() []rune {
	return sse.r
}

func (sse SSE) FullString() string {
	s := "{token="
	switch sse.t >> 14 & 3 {
	case 0:
		s += fmt.Sprintf("0b_%02b_%014b",
			(sse.t>>14)%(1<<2),
			(sse.t>>0)%(1<<14))
	case 1:
		s += fmt.Sprintf("0b_%02b_%01b_%05b_%08b",
			(sse.t>>14)%(1<<2),
			(sse.t>>13)%(1<<1),
			(sse.t>>8)%(1<<5),
			(sse.t>>0)%(1<<8))
	default:
		s += fmt.Sprintf("0b_%01b_%02b_%02b_%05b_%04b_%02b",
			(sse.t>>15)%(1<<1),
			(sse.t>>13)%(1<<2),
			(sse.t>>11)%(1<<2),
			(sse.t>>6)%(1<<5),
			(sse.t>>2)%(1<<4),
			(sse.t>>0)%(1<<2))
	}
	s += fmt.Sprintf(" UTF8=%v UTF8/glyph=%v runes=%v string=%q lang=%s numTokensLeft=%v}",
		sse.b, sse.g, sse.r, sse.s, sse.l, sse.n)
	return s
}

func (sse SSE) String() string {
	return sse.s
}

func (sse SSE) LanguageCode() string {
	return sse.l
}

func (sse SSE) NumTokensLeftInWord() int {
	return sse.n
}

//////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

var (
	consToSSEtoken = map[string]SSEtoken{
		"":    0b_0_00_00_00000_0000_00,
		"b":   0b_0_00_00_00001_0000_00,
		"d":   0b_0_00_00_00010_0000_00,
		"f":   0b_0_00_00_00011_0000_00,
		"g":   0b_0_00_00_00100_0000_00,
		"gb":  0b_0_00_00_00101_0000_00,
		"h":   0b_0_00_00_00110_0000_00,
		"k":   0b_0_00_00_00111_0000_00,
		"kp":  0b_0_00_00_01000_0000_00,
		"l":   0b_0_00_00_01001_0000_00,
		"m":   0b_0_00_00_01010_0000_00,
		"mb":  0b_0_00_00_01011_0000_00,
		"mp":  0b_0_00_00_01100_0000_00,
		"mv":  0b_0_00_00_01101_0000_00,
		"n":   0b_0_00_00_01110_0000_00,
		"nd":  0b_0_00_00_01111_0000_00,
		"ng":  0b_0_00_00_10000_0000_00,
		"ngb": 0b_0_00_00_10001_0000_00,
		"ny":  0b_0_00_00_10010_0000_00,
		"nz":  0b_0_00_00_10011_0000_00,
		"p":   0b_0_00_00_10100_0000_00,
		"r":   0b_0_00_00_10101_0000_00,
		"s":   0b_0_00_00_10110_0000_00,
		"t":   0b_0_00_00_10111_0000_00,
		"v":   0b_0_00_00_11000_0000_00,
		"w":   0b_0_00_00_11001_0000_00,
		"y":   0b_0_00_00_11010_0000_00,
		"z":   0b_0_00_00_11011_0000_00,
	}

	vowelToSSEtoken = map[string]SSEtoken{
		" ̣": 0b_0_00_00_00000_0000_00,
		"ạ":  0b_0_00_00_00000_0001_00,
		"ạn": 0b_0_00_00_00000_0010_00,
		"ə̣": 0b_0_00_00_00000_0011_00,
		"ɛ̣": 0b_0_00_00_00000_0100_00,
		"ẹ":  0b_0_00_00_00000_0101_00,
		"ẹn": 0b_0_00_00_00000_0110_00,
		"ị":  0b_0_00_00_00000_0111_00,
		"ịn": 0b_0_00_00_00000_1000_00,
		"ø̣": 0b_0_00_00_00000_1001_00,
		"ɔ̣": 0b_0_00_00_00000_1010_00,
		"ọ":  0b_0_00_00_00000_1011_00,
		"ọn": 0b_0_00_00_00000_1100_00,
		"ụ":  0b_0_00_00_00000_1101_00,
		"ụn": 0b_0_00_00_00000_1110_00,

		" ":  0b_0_00_00_00000_0000_01,
		"a":  0b_0_00_00_00000_0001_01,
		"an": 0b_0_00_00_00000_0010_01,
		"ə":  0b_0_00_00_00000_0011_01,
		"ɛ":  0b_0_00_00_00000_0100_01,
		"e":  0b_0_00_00_00000_0101_01,
		"en": 0b_0_00_00_00000_0110_01,
		"i":  0b_0_00_00_00000_0111_01,
		"in": 0b_0_00_00_00000_1000_01,
		"ø":  0b_0_00_00_00000_1001_01,
		"ɔ":  0b_0_00_00_00000_1010_01,
		"o":  0b_0_00_00_00000_1011_01,
		"on": 0b_0_00_00_00000_1100_01,
		"u":  0b_0_00_00_00000_1101_01,
		"un": 0b_0_00_00_00000_1110_01,

		"¨":  0b_0_00_00_00000_0000_10,
		"ä":  0b_0_00_00_00000_0001_10,
		"än": 0b_0_00_00_00000_0010_10,
		"ə̈": 0b_0_00_00_00000_0011_10,
		"ɛ̈": 0b_0_00_00_00000_0100_10,
		"ë":  0b_0_00_00_00000_0101_10,
		"ën": 0b_0_00_00_00000_0110_10,
		"ï":  0b_0_00_00_00000_0111_10,
		"ïn": 0b_0_00_00_00000_1000_10,
		"ø̈": 0b_0_00_00_00000_1001_10,
		"ɔ̈": 0b_0_00_00_00000_1010_10,
		"ö":  0b_0_00_00_00000_1011_10,
		"ön": 0b_0_00_00_00000_1100_10,
		"ü":  0b_0_00_00_00000_1101_10,
		"ün": 0b_0_00_00_00000_1110_10,

		"^":  0b_0_00_00_00000_0000_11,
		"â":  0b_0_00_00_00000_0001_11,
		"ân": 0b_0_00_00_00000_0010_11,
		"ə̂": 0b_0_00_00_00000_0011_11,
		"ɛ̂": 0b_0_00_00_00000_0100_11,
		"ê":  0b_0_00_00_00000_0101_11,
		"ên": 0b_0_00_00_00000_0110_11,
		"î":  0b_0_00_00_00000_0111_11,
		"în": 0b_0_00_00_00000_1000_11,
		"ø̂": 0b_0_00_00_00000_1001_11,
		"ɔ̂": 0b_0_00_00_00000_1010_11,
		"ô":  0b_0_00_00_00000_1011_11,
		"ôn": 0b_0_00_00_00000_1100_11,
		"û":  0b_0_00_00_00000_1101_11,
		"ûn": 0b_0_00_00_00000_1110_11,
	}

	sseTokenToCons = func() map[SSEtoken][]rune {
		m := make(map[SSEtoken][]rune)
		for cons, sse := range consToSSEtoken {
			m[sse] = []rune(cons)
		}
		return m
	}()

	sseTokenToVowel = func() map[SSEtoken][]rune {
		m := make(map[SSEtoken][]rune)
		for vowel, sse := range vowelToSSEtoken {
			m[sse] = []rune(vowel)
		}
		return m
	}()
)

// Matches final syllable of a Sango word in NFD form.
var (
	titleCaser = cases.Title(language.Und)
	upperCaser = cases.Upper(language.Und)

	sangoSyllableRE = regexp.MustCompile(`(?i)([-]?)` +
		`(ngb|gb|kp|mb|mp|mv|nd|ng|ny|nz|h|w|r|l|y|m|b|p|k|g|n|d|t|s|z|v|f|)` +
		`([aeiou\xF8\x{0254}\x{0259}\x{025B}])` +
		`([\x{0302}\x{0308}\x{0323}]?)(n?)$`)
)

// PRE: word must be in norm.NFD format.
// POST: newWord will be in norm.NFD format.
func encodeLastSyllable(word string, numSyllablesLeft int) (newWord string, sse SSEtoken) {
	span := sangoSyllableRE.FindStringSubmatchIndex(word)
	if span == nil || len(span) != 12 || span[0] < 0 || span[1] <= span[0] {
		return
	}
	sse = 0b_1_00_00_00000_0000_00
	sse |= 0b_0_11_00_00000_0000_00 & SSEtoken(min(3, numSyllablesLeft)<<13)
	newWord = word[:span[0]]
	syllable := string(word[span[0]:span[1]])
	hyphen := string(word[span[2]:span[3]])
	cons := string(strings.ToLower(word[span[4]:span[5]]))
	vowel := string(strings.ToLower(word[span[6]:span[7]]))
	pitch := string(strings.ToLower(word[span[8]:span[9]]))
	nasal := string(strings.ToLower(word[span[10]:span[11]]))

	if syllable == "" {
		return
	}
	if hyphen == "-" {
		sse |= 0b_0_00_10_00000_0000_00
	} else if syllable == titleCaser.String(syllable) {
		sse |= 0b_0_00_01_00000_0000_00
	} else if syllable == upperCaser.String(syllable) {
		sse |= 0b_0_00_11_00000_0000_00
	}

	switch pitch {
	case "\u0323": // dot below, e.g. ọ
		sse |= 0b_0_00_00_00000_0000_00
	case "\u0308": // diaeresis above, e.g. ö
		sse |= 0b_0_00_00_00000_0000_10
	case "\u0302": // circumflex above, e.g. ô
		sse |= 0b_0_00_00_00000_0000_11
	default: // no accent, e.g. o
		sse |= 0b_0_00_00_00000_0000_01
	}
	if sseCons, isFound := consToSSEtoken[cons]; isFound {
		sse |= sseCons
	}
	if sseVowel, isFound := vowelToSSEtoken[norm.NFC.String(vowel+pitch+nasal)]; isFound {
		sse |= sseVowel
	}

	return
}
