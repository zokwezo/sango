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

type WriteUtf8Options struct {
	ForSpaceUse, ForHyphenUse                    string
	WithShift, WithHeight, WithNTilde, WithPitch bool
}

var (
	AsToneless   = WriteUtf8Options{ForSpaceUse: "", ForHyphenUse: "", WithShift: false, WithHeight: false, WithNTilde: false, WithPitch: false}
	AsHeightless = WriteUtf8Options{ForSpaceUse: " ", ForHyphenUse: "-", WithShift: true, WithHeight: false, WithNTilde: false, WithPitch: true}
	AsLemma      = WriteUtf8Options{ForSpaceUse: " ", ForHyphenUse: "-", WithShift: true, WithHeight: true, WithNTilde: false, WithPitch: true}
	AsUtf8       = WriteUtf8Options{ForSpaceUse: " ", ForHyphenUse: "-", WithShift: true, WithHeight: true, WithNTilde: true, WithPitch: true}
)

type FromUtf8Options struct {
	TreatClosedVowelAsUnknownHeight bool
	TreatLowPitchAsUnknownPitch     bool
}

var (
	FromToneless   = FromUtf8Options{TreatClosedVowelAsUnknownHeight: true, TreatLowPitchAsUnknownPitch: true}
	FromHeightless = FromUtf8Options{TreatClosedVowelAsUnknownHeight: true, TreatLowPitchAsUnknownPitch: false}
	FromLemma      = FromUtf8Options{TreatClosedVowelAsUnknownHeight: false, TreatLowPitchAsUnknownPitch: false}
)

func (sse SSE) WriteUtf8To(s *strings.Builder, options WriteUtf8Options) {
	writeUtf8To(s, uint64(sse), options)
}

func (sse SSE) WriteAsUtf8To(s *strings.Builder) {
	writeUtf8To(s, uint64(sse), AsUtf8)
}

func (sse SSE) WriteAsTonelessTo(s *strings.Builder) {
	writeUtf8To(s, uint64(sse), AsToneless)
}

func (sse SSE) WriteAsHeightlessTo(s *strings.Builder) {
	writeUtf8To(s, uint64(sse), AsHeightless)
}

func (sse SSE) WriteAsLemmaTo(s *strings.Builder) {
	writeUtf8To(s, uint64(sse), AsLemma)
}

func (sse SSE) WriteAsCanonicalTo(s *strings.Builder) {
	writeAsCanonicalTo(s, uint64(sse))
}

func CanonicalToSSEs(s string) ([]SSE, error) {
	return canonicalToSSEs(s)
}

func Utf8ToSSEs(s string, options FromUtf8Options) ([]SSE, error) {
	return utf8ToSSEs(s, options)
}

func UnpadRight(word uint64) uint64 {
	w := word
	if w&0x_8_000_000_000_000_000 != 0 { // Sango word
		for w&0x000_000_000_000_FFF == 0 {
			w >>= 12
			if w == 0 {
				// Not a valid Sango SSE
				return word
			}
		}
	} else { // Unicode word
		for w&0x0000_0000_0000_FFFF == 0 {
			w >>= 16
			if w == 0 {
				// Not a valid Unicode SSE
				return word
			}
		}
	}
	return w
}

func PadRight(word uint64) uint64 {
	w := word
	for w&0x_F_FF0_000_000_000_000 == 0 {
		w <<= 12
		if w == 0 {
			// Not a valid Sango SSE
			break
		}
	}
	if w&0x_8000_0000_0000_0000 != 0 {
		return w // good Sango word
	}
	w = word
	// NOTE: the first byte must be in [U+0000, U+7FFF]
	// so as not to be interpreted as a Sango word.
	for w&0xFFFF_8000_0000_0000 == 0 {
		w <<= 16
		if w == 0 {
			// Not a valid Unicode SSE either
			return word
		}
	}
	return w // good Unicode
}

//////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

func writeUtf8To(s *strings.Builder, b uint64, options WriteUtf8Options) {
	if (b >> 63) == 0 { // up to 4 unicode runes
		rr := [4]rune{}
		for k := range 4 {
			rr[3-k] = rune(b & 0xFFFF)
			b >>= 16
		}
		for _, r := range rr {
			if r != 0x0 && r != 0x10 {
				s.WriteRune(r)
			}
		}
	} else { // up to 5 Sango syllables
		p0 := uint16(b >> 60 << 12)
		p := p0               // all but the first syllable
		p &= ^PrefixCode_MASK // force no space
		if getShiftCode(p) == ShiftCode_Title {
			p &= ^ShiftCode_MASK         // force no shift
			p |= uint16(ShiftCode_lower) // set lowercase
		}
		cc := [5]uint16{}
		for k := range 4 {
			cc[4-k] = uint16(b&0xFFF) | p
			b >>= 12
		}
		cc[0] = uint16(b&0xFFF) | p0
		for _, c := range cc {
			s.WriteString(utf8FromSangoCodeValue(c, options))
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
			if r != 0x0 && r != 0x10 {
				s.WriteString(fmt.Sprintf("%U", r))
			}
		}
	} else { // up to 5 Sango syllables
		p0 := uint16(b >> 60)
		p := p0             // all but the first syllable
		if p&0b11 == 0b10 { // if Titlecase
			p &= 0b1000 // force no space
			p |= 0b0001 // force lowercase
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
			s.WriteString(canonicalFromSangoCodeValue(c))
		}
	}
}

// We are okay with truncating a Unicode rune from 32 bits to 16 bits which
// is enough to express the entire Basic Multilingual Plane (BMP), since
// higher planes have only obscure glyphs. However, forgoing even one more bit
// would rule out interesting runes such as CJK glyphs and yet the MSB is needed
// to determine whether a code stores a Unicode rune or a Sango syllable.
// Consequently, the uint16 value must be supplemented with a separate bool
// to indicate the code variant type.
type sseCode struct {
	value   uint16 // if isSango is true, the MSB must be set to 1
	isSango bool
}

// Returns the accumulated codes and the index where processing stopped.
// The latter will equal len(s) on success, else the index of the first error.
func canonicalToCodes(s string) ([]sseCode, int) {
	// Initialize return values
	var codes []sseCode

	// Partition string into codes.
	indexes := canonicalRE.FindAllStringSubmatchIndex(s, -1)

	// Bail out if the start of string doesn't match.
	n := len(indexes)
	if n == 0 || indexes[0][0] != 0 {
		return codes, 0
	}

	for k := range n {
		ii := indexes[k]
		if len(ii) != 14 {
			panic("Bad index length")
		}

		// Bail out if there is a gap between codes.
		if k > 0 && indexes[k][0] != indexes[k-1][1] {
			return codes, indexes[k-1][1]
		}

		if ii[2] != -1 { // Unicode code
			value, err := strconv.ParseUint(s[ii[2]:ii[3]], 16, 16)
			if err != nil {
				return codes, ii[0]
			}
			if value != 0x0 && value != 0x10 {
				codes = append(codes, sseCode{value: uint16(value), isSango: false})
			}
		} else { // Sango syllable
			affix := s[ii[4]:ii[5]]
			shift := s[ii[6]:ii[7]]
			consonant := s[ii[8]:ii[9]]
			vowel := s[ii[10]:ii[11]]
			pitch := s[ii[12]:ii[13]]
			value, err := canonicalToSangoCodeValue(affix, shift, consonant, vowel, pitch)
			if err != nil {
				return codes, ii[0]
			}
			if !IsValid(value) {
				panic("bad value returned from canonicalToSangoCodeValue")
			}
			codes = append(codes, sseCode{value: value, isSango: true})
		}
	}
	return codes, len(s)
}

func utf8ToCodes(s string, options FromUtf8Options) ([]sseCode, int) {
	// Initialize return values
	var codes []sseCode

	// Partition string into codes.
	indexes := utf8RE.FindAllStringSubmatchIndex(s, -1)

	// Bail out if the start of string doesn't match.
	n := len(indexes)
	if n == 0 || indexes[0][0] != 0 {
		return codes, 0
	}

	for k := range n {
		ii := indexes[k]
		if len(ii) != 10 {
			panic("Bad index length")
		}

		// Bail out if there is a gap between codes.
		if k > 0 && indexes[k][0] != indexes[k-1][1] {
			return codes, indexes[k-1][1]
		}

		if ii[8] != -1 { // Unicode code
			value := s[ii[8]:ii[9]]
			if value != "\x00" && value != "\x10" {
				v := []rune(value)
				if len(v) != 1 || v[0] > 0xFFFF {
					return codes, ii[8]
				}
				codes = append(codes, sseCode{value: uint16(v[0]), isSango: false})
			}
		} else { // Sango syllable
			affix := s[ii[2]:ii[3]]
			consonant := s[ii[4]:ii[5]]
			vowel := s[ii[6]:ii[7]]
			value, err := utf8ToSangoCodeValue(affix, consonant, vowel, options)
			if err != nil {
				return codes, ii[0]
			}
			if !IsValid(value) {
				panic("bad value returned from utf8ToSangoCodeValue")
			}
			codes = append(codes, sseCode{value: value, isSango: true})
		}
	}
	return codes, len(s)
}

func codesToSSEs(codes []sseCode) []SSE {
	var sses []SSE
	var sse uint64
	var prevIsSango bool
	numCodesSaved := 0
	var msb4 uint64
	flush := func() {
		if prevIsSango {
			sse <<= (60 - 12*numCodesSaved)
		} else {
			sse <<= (64 - 16*numCodesSaved)
		}
		sse |= msb4
		sses = append(sses, SSE(sse))
		sse = 0
		numCodesSaved = 0
	}
	for _, code := range codes {
		if code.isSango {
			if !IsValid(code.value) {
				continue
			}
		}
		for count := range 10 {
			if count > 5 {
				panic("infinite loop")
			}
			if numCodesSaved == 0 {
				prevIsSango = code.isSango
				msb4 = 0
				if code.isSango {
					msb4 = uint64(code.value) >> 12 << 60
				}
			}
			if code.isSango && numCodesSaved == 5 ||
				!code.isSango && numCodesSaved == 4 {
				// current SSE is full, flush buffer and restart
				sse |= msb4
				sses = append(sses, SSE(sse))
				sse = 0
				numCodesSaved = 0
				continue // restart loop
			}
			if numCodesSaved != 0 && (prevIsSango != code.isSango || code.isSango && getPrefixCode(code.value) == PrefixCode_Space) {
				// can't mix unicode and sango nor two sango words in one SSE
				flush()
				continue // restart loop
			}
			switch code.isSango {
			case false:
				if code.value > 0x7FFF && numCodesSaved == 0 {
					// Large unicode will not fit into the most significant byte since
					// the first bit is reserved for the kind code.
					// Push instead U+0010 as a placeholder and retry.
					// NOTE: We push U+0010 instead of U+0000 so that
					// PadRight(UnpadRight(u)) is successful without confusing the
					// result of UnpadRight as a Sango code.
					sse <<= 16
					sse |= 0x0010
					numCodesSaved = 1
					continue // restart loop
				}
				sse <<= 16
				sse |= uint64(code.value & 0xFFFF)
			case true:
				sse <<= 12
				sse |= uint64(code.value & 0xFFF)
			}
			numCodesSaved += 1
			break // move on to the next code
		}
	}
	// flush buffer
	if numCodesSaved != 0 {
		flush()
	}
	return sses
}

func canonicalToSSEs(s string) ([]SSE, error) {
	var err error
	codes, b := canonicalToCodes(s)
	n := len(s)
	if b != n {
		e := b + 10
		etc := "..."
		if e >= n {
			e = n
			etc = ""
		}
		err = fmt.Errorf("cannot parse Canonical string starting at s[%v:] = %q", b, s[b:e]+etc)
	}
	return codesToSSEs(codes), err
}

func utf8ToSSEs(s string, options FromUtf8Options) ([]SSE, error) {
	var err error
	codes, b := utf8ToCodes(s, options)
	n := len(s)
	if b != n {
		e := b + 10
		etc := "..."
		if e >= n {
			e = n
			etc = ""
		}
		err = fmt.Errorf("cannot parse Utf8 string starting at s[%v:] = %q", b, s[b:e]+etc)
	}
	return codesToSSEs(codes), err
}
