// Sango Syllabic Encoding (SSE)
//
// Encoding for up to 4 unicode runes or a phonemically-valid Sango syllable.
//
// This is used internally because UTF8 (with both precomposed and combining
// diacritics) and the lossy standard Sango orthography (no vowel height)
// make it too cumbersome to use these latter formats.

package sse

import (
	"fmt"
	"strconv"
	"strings"
)

type SSE uint64

func (sse SSE) WriteAsUTF8To(s *strings.Builder) {
	writeAsUTF8To(s, uint64(sse))
}

func (sse SSE) WriteAsCanonicalTo(s *strings.Builder) {
	writeAsCanonicalTo(s, uint64(sse))
}

func CanonicalToSyllables(s string) ([]uint16, int) {
	return canonicalToSyllables(s)
}

func SyllablesToSSEs(syllables []uint16) []SSE {
	return syllablesToSSEs(syllables)
}

//////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

func writeAsUTF8To(s *strings.Builder, b uint64) {
	if (b >> 63) == 0 { // up to 4 unicode runes
		rr := [4]rune{}
		for k := range 4 {
			rr[3-k] = rune(b & 0xFFFF)
			b >>= 16
		}
		for _, r := range rr {
			if r != 0 {
				s.WriteRune(r)
			}
		}
	} else { // up to 5 Sango syllables
		p0 := uint16(b >> 60)
		p := p0             // all but the first syllable
		if p&0b11 == 0b01 { // if Titlecase
			p &= 0b1000 // force no space and lowercase
		} else { // else not Titlecase
			p &= 0b1011 // force no space, preserve case
		}
		p0 <<= 12
		p <<= 12
		cc := [5]uint16{}
		for k := range 4 {
			cc[4-k] = uint16(b&0xFFF) | p
			b >>= 12
		}
		cc[0] = uint16(b&0xFFF) | p0
		for _, c := range cc {
			s.WriteString(utf8FromSangoSyllableCode(c))
		}
	}
}

func writeAsCanonicalTo(s *strings.Builder, b uint64) {
	if (b >> 63) == 0 { // up to 4 unicode runes
		rr := [4]rune{}
		for k := range 4 {
			rr[3-k] = rune(b & 0xFFFF)
			b >>= 16
		}
		for _, r := range rr {
			if r != 0 {
				s.WriteString(fmt.Sprintf("%U", r))
			}
		}
	} else { // up to 5 Sango syllables
		p0 := uint16(b >> 60)
		p := p0             // all but the first syllable
		if p&0b11 == 0b01 { // if Titlecase
			p &= 0b1000 // force no space and lowercase
		} else { // else not Titlecase
			p &= 0b1011 // force no space, preserve case
		}
		p0 <<= 12
		p <<= 12
		cc := [5]uint16{}
		for k := range 4 {
			cc[4-k] = uint16(b&0xFFF) | p
			b >>= 12
		}
		cc[0] = uint16(b&0xFFF) | p0
		for _, c := range cc {
			s.WriteString(canonicalFromSangoSyllableCode(c))
		}
	}
}

const (
	syllableType_Unicode uint64 = 0x0000_0000_0000_0000
	syllableType_Sango   uint64 = 0x8000_0000_0000_0000
	syllableType_Mask    uint64 = 0x8000_0000_0000_0000

	undefined_Unicode uint64 = 0x0000_0000_0000_0000
	undefined_Sango   uint64 = 0xFFFF_FFFF_FFFF_FFFF

	msbs_Mask uint64 = 0xF000_0000_0000_0000
)

// Returns the accumulated syllables and the index where processing stopped.
// The latter will equal len(s) on success, else the index of the first error.
func canonicalToSyllables(s string) ([]uint16, int) {
	// Initialize return values
	var syllables []uint16

	// Partition string into syllables.
	indexes := canonicalRE.FindAllStringSubmatchIndex(s, -1)

	// Bail out if the start of string doesn't match.
	n := len(indexes)
	if n == 0 || indexes[0][0] != 0 {
		return syllables, 0
	}

	// word := undefined_Sango
	// msbs := word & msbs_Mask
	// msb0 := msbs & syllableType_Mask
	// var currentSyllable int

	for k := range n {
		ii := indexes[k]
		if len(ii) != 14 {
			panic("Bad index length")
		}

		// Bail out if there is a gap between syllables.
		if k > 0 && indexes[k][0] != indexes[k-1][1] {
			return syllables, indexes[k-1][1]
		}

		if ii[2] != -1 { // Unicode syllable
			codePoint, err := strconv.ParseUint(s[ii[2]:ii[3]], 16, 16)
			if err != nil {
				return syllables, ii[0]
			}
			syllables = append(syllables, uint16(codePoint))
		} else { // Sango syllable
			affix := s[ii[4]:ii[5]]
			shift := s[ii[6]:ii[7]]
			consonant := s[ii[8]:ii[9]]
			vowel := s[ii[10]:ii[11]]
			pitch := s[ii[12]:ii[13]]
			syllable, err := canonicalToSangoSyllableCode(affix, shift, consonant, vowel, pitch)
			if err != nil {
				return syllables, ii[0]
			}
			syllables = append(syllables, syllable)
		}
	}
	return syllables, len(s)
}

func syllablesToSSEs(syllables []uint16) []SSE {
	var sses []SSE
	// TODO: Implement and update unit test
	return sses
}
