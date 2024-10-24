package sse

import (
	"bytes"
	"fmt"
	"unicode"
	"unicode/utf8"

	"github.com/rivo/uniseg"
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
	// | 000 |    | f  | r  | k  |
	// | 001 | mv | v  | ng | g  |
	// | 010 | m  | p  | l  | kp |
	// | 011 | mb | b  | ngb| gb |
	// | 100 |    | s  | y  | h  |
	// | 101 | nz | z  | ny | w  |
	// | 110 | n  | t  | nd | d  |
	// +-----+----+----+----+----+
	// |  00 |    | u  | ɔ  | ɛ  |
	// |  01 | a  | i  | o  | e  |
	// |  10 |    | uñ | ø  | ə  |
	// |  11 | añ | iñ | oñ | eñ |
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

// TODO: Remove caseEnum and pitchEnum and merge consonant and
// vowelWithNasal into a syllable, then infer everything from that.
func MakeSangoSSE(
	caseEnum, pitchEnum uint16,
	consonant, vowelWithNasal string,
	numSyllablesLeft int) SangoSSE {
	s := uint16(1 << 15)
	s |= uint16(min(3, max(0, numSyllablesLeft)) << 13)
	if caseEnum > 0 && caseEnum < 4 {
		s |= caseEnum & 3 << 11
	}
	if c, found := consToSSEToken[consonant]; found {
		s |= uint16(c)
	}
	if v, found := vowelToSSEToken[vowelWithNasal]; found {
		s |= uint16(v & 0b1111_00)
	}
	if pitchEnum > 0 && pitchEnum < 4 {
		s |= uint16(pitchEnum & 3)
	}
	return SSEtoken(s).AsSangoSSE()
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
	return int(sse.t >> 9 & 3)
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
	consToSSEToken = map[string]SSEtoken{
		"":    0b_0_00_00_00000_0000_00,
		"f":   0b_0_00_00_00001_0000_00,
		"r":   0b_0_00_00_00010_0000_00,
		"k":   0b_0_00_00_00011_0000_00,
		"mv":  0b_0_00_00_00100_0000_00,
		"v":   0b_0_00_00_00101_0000_00,
		"ng":  0b_0_00_00_00110_0000_00,
		"g":   0b_0_00_00_00111_0000_00,
		"m":   0b_0_00_00_01000_0000_00,
		"p":   0b_0_00_00_01001_0000_00,
		"l":   0b_0_00_00_01010_0000_00,
		"kp":  0b_0_00_00_01011_0000_00,
		"mb":  0b_0_00_00_01100_0000_00,
		"b":   0b_0_00_00_01101_0000_00,
		"ngb": 0b_0_00_00_01110_0000_00,
		"gb":  0b_0_00_00_01111_0000_00,
		"s":   0b_0_00_00_10001_0000_00,
		"y":   0b_0_00_00_10010_0000_00,
		"h":   0b_0_00_00_10011_0000_00,
		"nz":  0b_0_00_00_10100_0000_00,
		"z":   0b_0_00_00_10101_0000_00,
		"ny":  0b_0_00_00_10110_0000_00,
		"w":   0b_0_00_00_10111_0000_00,
		"n":   0b_0_00_00_11000_0000_00,
		"t":   0b_0_00_00_11001_0000_00,
		"nd":  0b_0_00_00_11010_0000_00,
		"d":   0b_0_00_00_11011_0000_00,
	}

	vowelToSSEToken = map[string]SSEtoken{
		".":   0b_0_00_00_00000_0000_00,
		"ụ":   0b_0_00_00_00000_0001_00,
		"ɔ̣":  0b_0_00_00_00000_0010_00,
		"ɛ̣":  0b_0_00_00_00000_0011_00,
		"ạ":   0b_0_00_00_00000_0100_00,
		"ị":   0b_0_00_00_00000_0101_00,
		"ọ":   0b_0_00_00_00000_0110_00,
		"ẹ":   0b_0_00_00_00000_0111_00,
		"ụn":  0b_0_00_00_00000_1001_00,
		"ø̣":  0b_0_00_00_00000_1010_00,
		"ə̣":  0b_0_00_00_00000_1011_00,
		"ạn":  0b_0_00_00_00000_1100_00,
		"ịn":  0b_0_00_00_00000_1101_00,
		"ọn":  0b_0_00_00_00000_1110_00,
		"ø̣n": 0b_0_00_00_00000_1110_00,
		"ẹn":  0b_0_00_00_00000_1111_00,
		"ə̣n": 0b_0_00_00_00000_1111_00,

		"":   0b_0_00_00_00000_0000_01,
		"u":  0b_0_00_00_00000_0001_01,
		"ɔ":  0b_0_00_00_00000_0010_01,
		"ɛ":  0b_0_00_00_00000_0011_01,
		"a":  0b_0_00_00_00000_0100_01,
		"i":  0b_0_00_00_00000_0101_01,
		"o":  0b_0_00_00_00000_0110_01,
		"e":  0b_0_00_00_00000_0111_01,
		"un": 0b_0_00_00_00000_1001_01,
		"ø":  0b_0_00_00_00000_1010_01,
		"ə":  0b_0_00_00_00000_1011_01,
		"an": 0b_0_00_00_00000_1100_01,
		"in": 0b_0_00_00_00000_1101_01,
		"on": 0b_0_00_00_00000_1110_01,
		"øn": 0b_0_00_00_00000_1110_01,
		"en": 0b_0_00_00_00000_1111_01,
		"ən": 0b_0_00_00_00000_1111_01,

		"¨":   0b_0_00_00_00000_0000_10,
		"ü":   0b_0_00_00_00000_0001_10,
		"ɔ̈":  0b_0_00_00_00000_0010_10,
		"ɛ̈":  0b_0_00_00_00000_0011_10,
		"ä":   0b_0_00_00_00000_0100_10,
		"ï":   0b_0_00_00_00000_0101_10,
		"ö":   0b_0_00_00_00000_0110_10,
		"ë":   0b_0_00_00_00000_0111_10,
		"ün":  0b_0_00_00_00000_1001_10,
		"ø̈":  0b_0_00_00_00000_1010_10,
		"ə̈":  0b_0_00_00_00000_1011_10,
		"än":  0b_0_00_00_00000_1100_10,
		"ïn":  0b_0_00_00_00000_1101_10,
		"ön":  0b_0_00_00_00000_1110_10,
		"ø̈n": 0b_0_00_00_00000_1110_10,
		"ën":  0b_0_00_00_00000_1111_10,
		"ə̈n": 0b_0_00_00_00000_1111_10,

		"^":   0b_0_00_00_00000_0000_11,
		"û":   0b_0_00_00_00000_0001_11,
		"ɔ̂":  0b_0_00_00_00000_0010_11,
		"ɛ̂":  0b_0_00_00_00000_0011_11,
		"â":   0b_0_00_00_00000_0100_11,
		"î":   0b_0_00_00_00000_0101_11,
		"ô":   0b_0_00_00_00000_0110_11,
		"ê":   0b_0_00_00_00000_0111_11,
		"ûn":  0b_0_00_00_00000_1001_11,
		"ø̂":  0b_0_00_00_00000_1010_11,
		"ə̂":  0b_0_00_00_00000_1011_11,
		"ân":  0b_0_00_00_00000_1100_11,
		"în":  0b_0_00_00_00000_1101_11,
		"ôn":  0b_0_00_00_00000_1110_11,
		"ø̂n": 0b_0_00_00_00000_1110_11,
		"ên":  0b_0_00_00_00000_1111_11,
		"ə̂n": 0b_0_00_00_00000_1111_11,
	}

	sseTokenToCons = func() map[SSEtoken][]rune {
		m := make(map[SSEtoken][]rune)
		for cons, sse := range consToSSEToken {
			m[sse] = []rune(cons)
		}
		return m
	}()

	sseTokenToVowel = func() map[SSEtoken][]rune {
		m := make(map[SSEtoken][]rune)
		for vowel, sse := range vowelToSSEToken {
			m[sse] = []rune(vowel)
		}
		return m
	}()
)
