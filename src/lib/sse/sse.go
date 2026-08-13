// Sango Syllabic Encoding (SSE)
//
// This file contains utilities to convert between SSE and various input/output formats.
// UTF8, Unicode diacritics (two flavors: precomposed and combining), and the lossy
// standard Sango orthography (no vowel height) make it too cumbersome to use these latter
// formats internally, which are all done in SSE.

package sse

import (
	"encoding/binary"
	"fmt"
	"strings"
)

type SSE uint64

func IsUnicodeSSE(sse SSE) bool              { return isUnicodeSSE(sse) }
func IsSangoSSE(sse SSE) bool                { return isSangoSSE(sse) }
func UTF8FromSSEs(sses []SSE) string         { return utf8FromSSEs(sses) }
func UTF8FromSSE(sse SSE) string             { return utf8FromSSE(sse) }
func UTF8FromUnicodeSSE(sse SSE) string      { return utf8FromUnicodeSSE(sse) }
func UTF8FromSangoSSE(sse SSE) string        { return utf8FromSangoSSE(sse) }
func CanonicalFromSSEs(sses []SSE) string    { return canonicalFromSSEs(sses) }
func CanonicalFromSSE(sse SSE) string        { return canonicalFromSSE(sse) }
func CanonicalFromUnicodeSSE(sse SSE) string { return canonicalFromUnicodeSSE(sse) }
func CanonicalFromSangoSSE(sse SSE) string   { return canonicalFromSangoSSE(sse) }

/////////////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

func isUnicodeSSE(sse SSE) bool {
	return !IsSangoSSE(sse)
}

func isSangoSSE(sse SSE) bool {
	return (sse>>63)%2 != 0
}

func extractBits(src uint64, lsb int, msb int) uint64 {
	if msb < lsb || lsb > 63 || msb < 0 {
		return uint64(0)
	}
	if lsb < 0 {
		lsb = 0
	}
	if msb > 63 {
		msb = 63
	}
	src >>= lsb
	shift := msb + 1 - lsb
	if shift < 63 {
		src %= (1 << shift)
	}
	return src
}

/////////////////////////////////////////////////////////////////////////////////////////

func utf8FromSSEs(sses []SSE) string {
	var s strings.Builder
	for _, sse := range sses {
		s.WriteString(UTF8FromSSE(sse))
	}
	return s.String()
}

func utf8FromSSE(sse SSE) string {
	if IsSangoSSE(sse) {
		return UTF8FromSangoSSE(sse)
	}
	return UTF8FromUnicodeSSE(sse)
}

func utf8FromUnicodeSSE(sse SSE) string {
	if !IsUnicodeSSE(sse) {
		panic("sse is not Unicode")
	}
	var bWithNulls [8]byte
	binary.BigEndian.PutUint64(bWithNulls[:], uint64(sse))
	var bWithoutNulls []byte
	for _, b := range bWithNulls {
		if b != 0 {
			bWithoutNulls = append(bWithoutNulls, b)
		}
	}
	return string(bWithoutNulls)
}

func utf8FromSangoSSE(sse SSE) string {
	if !IsSangoSSE(sse) {
		panic("sse is not Unicode")
	}
	src := uint64(sse)
	var ss strings.Builder
	if extractBits(src, 62, 62) != 0 {
		ss.WriteString(" ")
	}
	nonemptyUTF8 := false
	for k := range 5 {
		s := ""
		lsb := (4 - k) * 12
		msb := lsb + 11
		syllable := extractBits(src, lsb, msb)
		// Hyphen prefix
		if extractBits(syllable, 11, 11) != 0 {
			s += "-"
		}
		// Consonant
		switch extractBits(syllable, 6, 10) {
		case 0b00000:
		case 0b00001:
			s += "b"
		case 0b00010:
			s += "v"
		case 0b00011:
			s += "y"
		case 0b00100:
			s += "d"
		case 0b00101:
			s += "z"
		case 0b00110:
			s += "q"
		case 0b00111:
			s += "g"
		case 0b01000:
			s += "h"
		case 0b01001:
			s += "p"
		case 0b01010:
			s += "f"
		case 0b01011:
			s += "l"
		case 0b01100:
			s += "t"
		case 0b01101:
			s += "s"
		case 0b01110:
			s += "kp"
		case 0b01111:
			s += "k"
		case 0b10000:
			s += "w"
		case 0b10001:
			s += "mb"
		case 0b10010:
			s += "mv"
		case 0b10011:
			s += "ny"
		case 0b10100:
			s += "nd"
		case 0b10101:
			s += "nz"
		case 0b10110:
			s += "ngb"
		case 0b10111:
			s += "ng"
		case 0b11000:
			s += "n"
		case 0b11001:
			s += "mp"
		case 0b11010:
			s += "m"
		case 0b11011:
			s += "r"
		case 0b11100:
			continue
		case 0b11101:
			continue
		case 0b11110:
			continue
		case 0b11111:
			continue
		default:
			panic("Bad consonant bits extraction")
		}
		// Pitch
		switch extractBits(syllable, 0, 1) {
		case 0b00:
			// Low Vowel
			switch extractBits(syllable, 2, 5) {
			case 0b0000:
				s += "a"
			case 0b0001:
				s += "aN"
			case 0b0010:
				s += "e"
			case 0b0011:
				s += "eN"
			case 0b0100:
				s += "i"
			case 0b0101:
				s += "iN"
			case 0b0110:
				s += "o"
			case 0b0111:
				s += "oN"
			case 0b1000:
				s += "ɛ"
			case 0b1001:
				s += "ɔ"
			case 0b1010:
				s += "u"
			case 0b1011:
				s += "uN"
			case 0b1100:
				s += "x"
			case 0b1101:
				s += "c"
			case 0b1110:
				continue
			case 0b1111:
				continue
			default:
				panic("Bad vowel bits extraction")
			}
		case 0b01:
			// Mid Vowel
			switch extractBits(syllable, 2, 5) {
			case 0b0000:
				s += "ä"
			case 0b0001:
				s += "äN"
			case 0b0010:
				s += "ë"
			case 0b0011:
				s += "ëN"
			case 0b0100:
				s += "ï"
			case 0b0101:
				s += "ïN"
			case 0b0110:
				s += "ö"
			case 0b0111:
				s += "öN"
			case 0b1000:
				s += "ɛ̈"
			case 0b1001:
				s += "ɔ̈"
			case 0b1010:
				s += "ü"
			case 0b1011:
				s += "üN"
			case 0b1100:
				s += "ẍ"
			case 0b1101:
				s += "c̈"
			case 0b1110:
				continue
			case 0b1111:
				continue
			default:
				panic("Bad vowel bits extraction")
			}
		case 0b10:
			// High Vowel
			switch extractBits(syllable, 2, 5) {
			case 0b0000:
				s += "â"
			case 0b0001:
				s += "âN"
			case 0b0010:
				s += "ê"
			case 0b0011:
				s += "êN"
			case 0b0100:
				s += "î"
			case 0b0101:
				s += "îN"
			case 0b0110:
				s += "ô"
			case 0b0111:
				s += "ôN"
			case 0b1000:
				s += "ɛ̂"
			case 0b1001:
				s += "ɔ̂"
			case 0b1010:
				s += "û"
			case 0b1011:
				s += "ûN"
			case 0b1100:
				s += "x̂"
			case 0b1101:
				s += "ĉ"
			case 0b1110:
				continue
			case 0b1111:
				continue
			default:
				panic("Bad vowel bits extraction")
			}
		case 0b11:
			// Unknown Pitch Vowel
			switch extractBits(syllable, 2, 5) {
			case 0b0000:
				s += "ạ"
			case 0b0001:
				s += "ạN"
			case 0b0010:
				s += "ẹ"
			case 0b0011:
				s += "ẹN"
			case 0b0100:
				s += "ị"
			case 0b0101:
				s += "ịN"
			case 0b0110:
				s += "ọ"
			case 0b0111:
				s += "ọN"
			case 0b1000:
				s += "ɛ̣"
			case 0b1001:
				s += "ɔ̣"
			case 0b1010:
				s += "ụ"
			case 0b1011:
				s += "ụN"
			case 0b1100:
				s += "x̣"
			case 0b1101:
				s += "c̣"
			case 0b1110:
				continue
			case 0b1111:
				continue
			default:
				panic("Bad vowel bits extraction")
			}
		default:
			panic("Bad pitch bits extraction")
		}
		if s != "" {
			nonemptyUTF8 = true
			ss.WriteString(s)
		}
	}
	if nonemptyUTF8 {
		return ss.String()
	}
	return ""
}

//////////////////////////////////////////////////////////////////////////////

func canonicalFromSSEs(sses []SSE) string {
	var s strings.Builder
	for _, sse := range sses {
		s.WriteString(CanonicalFromSSE(sse))
	}
	return s.String()
}

func canonicalFromSSE(sse SSE) string {
	if IsSangoSSE(sse) {
		return CanonicalFromSangoSSE(sse)
	}
	return CanonicalFromUnicodeSSE(sse)
}

func canonicalFromUnicodeSSE(sse SSE) string {
	if !IsUnicodeSSE(sse) {
		panic("sse is not Unicode")
	}
	var s strings.Builder
	var bWithNulls [8]byte
	binary.BigEndian.PutUint64(bWithNulls[:], uint64(sse))
	for _, b := range bWithNulls {
		if b != 0 {
			s.WriteString(fmt.Sprintf("%02X", b))
		}
	}
	return s.String()
}

func canonicalFromSangoSSE(sse SSE) string {
	if !IsSangoSSE(sse) {
		panic("sse is not Unicode")
	}
	src := uint64(sse)
	var ss strings.Builder
	if extractBits(src, 62, 62) != 0 {
		ss.WriteString(" ")
	}
	nonemptyCanonical := false
	for k := range 5 {
		s := ""
		lsb := (4 - k) * 12
		msb := lsb + 11
		syllable := extractBits(src, lsb, msb)
		// Hyphen prefix
		if extractBits(syllable, 11, 11) != 0 {
			s += "-"
		}
		// Consonant
		switch extractBits(syllable, 6, 10) {
		case 0b00000:
			s += "h"
		case 0b00001:
			s += "b"
		case 0b00010:
			s += "v"
		case 0b00011:
			s += "y"
		case 0b00100:
			s += "d"
		case 0b00101:
			s += "z"
		case 0b00110:
			s += "q"
		case 0b00111:
			s += "g"
		case 0b01000:
			s += "H"
		case 0b01001:
			s += "p"
		case 0b01010:
			s += "f"
		case 0b01011:
			s += "l"
		case 0b01100:
			s += "t"
		case 0b01101:
			s += "s"
		case 0b01110:
			s += "K"
		case 0b01111:
			s += "k"
		case 0b10000:
			s += "w"
		case 0b10001:
			s += "B"
		case 0b10010:
			s += "V"
		case 0b10011:
			s += "Y"
		case 0b10100:
			s += "D"
		case 0b10101:
			s += "Z"
		case 0b10110:
			s += "Q"
		case 0b10111:
			s += "G"
		case 0b11000:
			s += "n"
		case 0b11001:
			s += "P"
		case 0b11010:
			s += "m"
		case 0b11011:
			s += "r"
		case 0b11100:
			continue
		case 0b11101:
			continue
		case 0b11110:
			continue
		case 0b11111:
			continue
		default:
			panic("Bad consonant bits extraction")
		}
		// Vowel
		switch extractBits(syllable, 2, 5) {
		case 0b0000:
			s += "a"
		case 0b0001:
			s += "A"
		case 0b0010:
			s += "e"
		case 0b0011:
			s += "E"
		case 0b0100:
			s += "i"
		case 0b0101:
			s += "I"
		case 0b0110:
			s += "o"
		case 0b0111:
			s += "O"
		case 0b1000:
			s += "x"
		case 0b1001:
			s += "c"
		case 0b1010:
			s += "u"
		case 0b1011:
			s += "U"
		case 0b1100:
			s += "X"
		case 0b1101:
			s += "C"
		case 0b1110:
			continue
		case 0b1111:
			continue
		default:
			panic("Bad vowel bits extraction")
		}
		// Pitch
		switch extractBits(syllable, 0, 1) {
		case 0b00:
			s += "_"
		case 0b01:
			s += ":"
		case 0b10:
			s += "^"
		case 0b11:
			s += "="
		default:
			panic("Bad pitch bits extraction")
		}
		if s != "" {
			nonemptyCanonical = true
			ss.WriteString(s)
		}
	}
	if nonemptyCanonical {
		return ss.String()
	}
	return ""
}
