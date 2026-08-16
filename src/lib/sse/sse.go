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
	"strings"
)

type SSE uint64

func (sse SSE) WriteAsUTF8To(s *strings.Builder) {
	writeAsUTF8To(s, uint64(sse))
}

func (sse SSE) WriteAsCanonicalTo(s *strings.Builder) {
	writeAsCanonicalTo(s, uint64(sse))
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
