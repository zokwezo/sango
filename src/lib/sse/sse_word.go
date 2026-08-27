// Implements operations on 64-bit word(s)

package sse

import (
	"fmt"
	"strconv"
	"strings"
)

func unpadRight(word uint64) uint64 {
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

func padRight(word uint64) uint64 {
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
	if n == 0 || indexes[0][canonicalRE_WholeBegin] != 0 {
		return codes, 0
	}

	for k := range n {
		ii := indexes[k]
		if len(ii) != 14 {
			panic("Bad index length")
		}

		// Bail out if there is a gap between codes.
		if k > 0 && indexes[k][canonicalRE_WholeBegin] != indexes[k-1][canonicalRE_WholeEnd] {
			return codes, indexes[k-1][canonicalRE_WholeEnd]
		}

		if ii[canonicalRE_UnicodeHexBegin] != -1 { // Unicode code
			value, err := strconv.ParseUint(s[ii[canonicalRE_UnicodeHexBegin]:ii[canonicalRE_UnicodeHexEnd]], 16, 16)
			if err != nil {
				return codes, ii[canonicalRE_WholeBegin]
			}
			if value != 0x0 && value != 0x10 {
				codes = append(codes, sseCode{value: uint16(value), isSango: false})
			}
		} else { // Sango syllable
			affix := s[ii[canonicalRE_AffixBegin]:ii[canonicalRE_AffixEnd]]
			shift := s[ii[canonicalRE_ShiftBegin]:ii[canonicalRE_ShiftEnd]]
			consonant := s[ii[canonicalRE_ConsonantBegin]:ii[canonicalRE_ConsonantEnd]]
			vowel := s[ii[canonicalRE_VowelBegin]:ii[11]]
			pitch := s[ii[canonicalRE_PitchBegin]:ii[canonicalRE_PitchEnd]]
			value, err := canonicalToSangoCodeValue(affix, shift, consonant, vowel, pitch)
			if err != nil {
				return codes, ii[canonicalRE_WholeBegin]
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
	if n == 0 || indexes[0][utf8RE_WholeBegin] != 0 {
		return codes, 0
	}

	// Loop through each span (Unicode glyph or Sango syllable) and convert to a code.
	for k := range n {
		// Use pointers because we may need to move a nasal N to the following consonant.
		iiCurr := &indexes[k]
		if len(*iiCurr) != 18 {
			panic("Bad index length")
		}
		w0Curr := &(*iiCurr)[utf8RE_WholeBegin]
		w1Curr := &(*iiCurr)[utf8RE_WholeEnd]
		u0Curr := &(*iiCurr)[utf8RE_UnicodeBegin]
		u1Curr := &(*iiCurr)[utf8RE_UnicodeEnd]
		a0Curr := &(*iiCurr)[utf8RE_AffixBegin]
		a1Curr := &(*iiCurr)[utf8RE_AffixEnd]
		c0Curr := &(*iiCurr)[utf8RE_ConsonantBegin]
		c1Curr := &(*iiCurr)[utf8RE_ConsonantEnd]
		v0Curr := &(*iiCurr)[utf8RE_OpenVowelBegin]
		v1Curr := &(*iiCurr)[utf8RE_OpenVowelEnd]
		if *v0Curr == -1 || *v1Curr == -1 {
			v0Curr = &(*iiCurr)[utf8RE_CloseVowelBegin]
			v1Curr = &(*iiCurr)[utf8RE_CloseVowelEnd]
		}
		x1Curr := &(*iiCurr)[utf8RE_VowelPlusNasalEnd]
		e1Curr := &(*iiCurr)[utf8RE_CloseVowelPlusNasalEnd]
		n0Curr := &(*iiCurr)[utf8RE_NasalBegin]
		n1Curr := &(*iiCurr)[utf8RE_NasalEnd]

		// If an n was falsely stolen from the following consonant to make a nasal, put it back.
		if k+1 < n { // if there is a next interval...
			iiNext := &indexes[k+1]
			if len(*iiNext) != 18 {
				panic("Bad index length")
			}
			w0Next := &(*iiNext)[utf8RE_WholeBegin]
			if *w1Curr != *w0Next {
				// Bail out, there is a gap between codes.
				return codes, *w1Curr
			}
			u0Next := &(*iiNext)[utf8RE_UnicodeBegin]
			u1Next := &(*iiNext)[utf8RE_UnicodeEnd]
			a0Next := &(*iiNext)[utf8RE_AffixBegin]
			a1Next := &(*iiNext)[utf8RE_AffixEnd]
			c0Next := &(*iiNext)[utf8RE_ConsonantBegin]
			c1Next := &(*iiNext)[utf8RE_ConsonantEnd]
			v0Next := &(*iiNext)[utf8RE_OpenVowelBegin]
			v1Next := &(*iiNext)[utf8RE_OpenVowelEnd]
			if *v0Next == -1 || *v1Next == -1 {
				v0Next = &(*iiNext)[utf8RE_CloseVowelBegin]
				v1Next = &(*iiNext)[utf8RE_CloseVowelEnd]
			}
			// If the consonant ends in "n" and the next syllable has no affix
			// and a consonant that is one of "", "d", "g", "gb", "y", "z", then move the
			// falsely attributed nasal N to the start of the next syllable's consonant.
			if *u0Curr == -1 && *u1Curr == -1 && *u0Next == -1 && *u1Next == -1 && *n0Curr != -1 && *n1Curr > *n0Curr {
				nasalCurr := s[*n0Curr:*n1Curr]
				switch nasalCurr {
				case "N":
					fallthrough
				case "n":
					affixNext := s[*a0Next:*a1Next]
					consonantNext := s[*c0Next:*c1Next]
					if affixNext == "" {
						if *n1Curr != *c0Next {
							panic("n1Curr != c0Next")
						}
						switch consonantNext {
						case "":
							fallthrough
						case "D":
							fallthrough
						case "G":
							fallthrough
						case "GB":
							fallthrough
						case "Gb":
							fallthrough
						case "Y":
							fallthrough
						case "Z":
							fallthrough
						case "d":
							fallthrough
						case "g":
							fallthrough
						case "gB":
							fallthrough
						case "gb":
							fallthrough
						case "y":
							fallthrough
						case "z":
							// Move nasal to start of next syllable
							*w1Curr -= 1
							*x1Curr -= 1
							*e1Curr -= 1
							*n1Curr -= 1
							*w0Next -= 1
							*c0Next -= 1
						}
					}
				}
			}
		}

		// Convert to a code.
		switch {
		case *u0Curr != -1 && *u1Curr != -1:
			v := []rune(s[*u0Curr:*u1Curr])
			switch len(v) {
			case 1:
				if unicode := v[0]; unicode != 0 && unicode != 16 && unicode <= 0xFFFF {
					codes = append(codes, sseCode{value: uint16(unicode), isSango: false})
				} else {
					return codes, *u0Curr
				}
			case 0:
				panic("v is empty") // should not have passed the regexp!
			default:
				panic("Multi-rune unicode") // should not have passed the regexp!
			}
		case *u0Curr == -1 && *u1Curr == -1:
			affix := s[*a0Curr:*a1Curr]
			consonant := s[*c0Curr:*c1Curr]
			vowel := s[*v0Curr:*v1Curr]
			nasalCurr := ""
			if *n0Curr != -1 && *n1Curr != -1 {
				nasalCurr = s[*n0Curr:*n1Curr]
			}
			value, err := utf8ToSangoCodeValue(affix, consonant, vowel, nasalCurr, options)
			if err != nil {
				return codes, *w0Curr
			}
			if !IsValid(value) {
				panic("bad value returned from utf8ToSangoCodeValue")
			}
			codes = append(codes, sseCode{value: value, isSango: true})
		default:
			panic("Bad kind")
		}
	}
	return codes, len(s)
}

func codesToSSEs(codes []sseCode) []SSE {
	const (
		msb4Mask  uint64 = 0xF000_0000_0000_0000
		sangoMask uint64 = 0x8000_0000_0000_0000
		spaceMask uint64 = 0x4000_0000_0000_0000
		shiftMask uint64 = 0x3000_0000_0000_0000
	)
	var sses []SSE
	var sse uint64
	numCodesSaved := 0
	var msb4 uint64
	sseIsSango := func() bool { return msb4&sangoMask != 0 }
	flush := func() {
		if prevIsSango := sseIsSango(); prevIsSango {
			sse <<= (60 - 12*numCodesSaved)
			sse |= msb4
		} else {
			sse <<= (64 - 16*numCodesSaved)
		}
		sses = append(sses, SSE(sse))
		sse = 0
		msb4 = 0
		numCodesSaved = 0
	}
	for _, code := range codes {
		prevIsSango := sseIsSango()
		if !IsValid(code.value) {
			continue
		}
		// This loop packs 4 Unicode or 5 Sango syllables per SSE.
		// If the SSE is already full, flush it and start a new one.
		// Also we can't mix Unicode and sango nor two Sango words in one SSE.
		if code.isSango && numCodesSaved == 5 ||
			!code.isSango && numCodesSaved == 4 ||
			numCodesSaved != 0 && (code.isSango != prevIsSango ||
				code.isSango && getPrefixCode(code.value) == PrefixCode_Space) {
			flush()
		}

		prevIsSango = code.isSango
		switch code.isSango {
		case false:
			msb4 &= ^sangoMask
			if code.value > 0x7FFF && numCodesSaved == 0 {
				// Large unicode will not fit into the most significant byte since
				// the first bit is reserved for the kind code.
				// Push instead U+0010 as a placeholder and retry.
				// NOTE: We push U+0010 instead of U+0000 so that
				// padRight(unpadRight(u)) is successful without confusing the
				// result of unpadRight as a Sango code.
				sse <<= 16
				sse |= 0x0010
				numCodesSaved += 1
			}
			sse <<= 16
			sse |= uint64(code.value & 0xFFFF)
		case true:
			msb4 |= uint64(code.value>>14) << 62
			prevShift := shiftMask & msb4
			currShift := shiftMask & (uint64(code.value) << 48)
			if currShift > prevShift {
				msb4 &= ^shiftMask
				msb4 |= currShift
			}
			sse <<= 12
			sse |= uint64(code.value & 0xFFF)
		}
		numCodesSaved += 1
	}
	if numCodesSaved != 0 {
		flush()
	}
	return sses
}

func toSSEs(s string, toCodes func(string) ([]sseCode, int)) ([]SSE, error) {
	var err error
	codes, b := toCodes(s)
	n := len(s)
	if b != n {
		e := b + 10
		etc := "..."
		if e >= n {
			e = n
			etc = ""
		}
		err = fmt.Errorf("cannot parse string starting at s[%v:] = %q", b, s[b:e]+etc)
	}
	return codesToSSEs(codes), err
}
