// SSE Sango Syllable
//
// Encoding for a phonemically-valid Sango syllable.

package sse

import (
	"regexp"
	"strings"
)

// A Sango syllable comprises 16 bits, listed below from MSB to LSB.
type PrefixCode int
type CaseCode int
type InfixCode int
type ConsonantCode int
type VowelCode int
type PitchCode int

// 1 bit = 2 options
const (
	Code_NoPrefix PrefixCode = iota
	Code_SpacePrefix
)

// 2 bits = 4 options
const (
	Code_Lowercase CaseCode = iota
	Code_Titlecase
	Code_Uppercase
	Code_Invisible
)

// 1 bit = 2 options
const (
	Code_NoInfix InfixCode = iota
	Code_HyphenInfix
)

// 5 bits = 32 options
const (
	Code_h ConsonantCode = iota // unaspirated H
	Code_B
	Code_V
	Code_Y
	Code_D
	Code_Z
	Code_Q
	Code_G
	Code_H
	Code_P
	Code_F
	Code_L
	Code_T
	Code_S
	Code_KP
	Code_K
	Code_W
	Code_MB
	Code_MV
	Code_NY
	Code_ND
	Code_NZ
	Code_NGB
	Code_NG
	Code_N
	Code_MP
	Code_M
	Code_R
)

// 4 bits = 16 options
const (
	Code_ClosedA VowelCode = iota
	Code_NasalA
	Code_ClosedE
	Code_NasalE
	Code_ClosedI
	Code_NasalI
	Code_ClosedO
	Code_NasalO
	Code_OpenE
	Code_OpenO
	Code_ClosedU
	Code_NasalU
	Code_UnknownE
	Code_UnknownO
)

// 2 bits = 4 options
const (
	Code_LowPitch PitchCode = iota
	Code_MidPitch
	Code_HighPitch
	Code_UnknownPitch
)

func getPrefixCode(code uint16) PrefixCode       { return PrefixCode(extractBits(code, 14, 14)) }
func getCaseCode(code uint16) CaseCode           { return CaseCode(extractBits(code, 12, 13)) }
func getInfixCode(code uint16) InfixCode         { return InfixCode(extractBits(code, 11, 11)) }
func getConsonantCode(code uint16) ConsonantCode { return ConsonantCode(extractBits(code, 6, 10)) }
func getVowelCode(code uint16) VowelCode         { return VowelCode(extractBits(code, 2, 5)) }
func getPitchCode(code uint16) PitchCode         { return PitchCode(extractBits(code, 0, 1)) }

func utf8FromSangoSyllableCode(code uint16) string {
	if (uint16(code) & 0x8000) == 0 {
		panic("code does not represent Sango")
	}
	s := ""
	caseCode := getCaseCode(code)
	if caseCode == Code_Invisible {
		return s
	}
	switch getConsonantCode(code) {
	case Code_h:
		// omit consonant for unaspirated H
	case Code_B:
		s += "b"
	case Code_V:
		s += "v"
	case Code_Y:
		s += "y"
	case Code_D:
		s += "d"
	case Code_Z:
		s += "z"
	case Code_Q:
		s += "q"
	case Code_G:
		s += "g"
	case Code_H:
		s += "h"
	case Code_P:
		s += "p"
	case Code_F:
		s += "f"
	case Code_L:
		s += "l"
	case Code_T:
		s += "t"
	case Code_S:
		s += "s"
	case Code_KP:
		s += "kp"
	case Code_K:
		s += "k"
	case Code_W:
		s += "w"
	case Code_MB:
		s += "mb"
	case Code_MV:
		s += "mv"
	case Code_NY:
		s += "ny"
	case Code_ND:
		s += "nd"
	case Code_NZ:
		s += "nz"
	case Code_NGB:
		s += "ngb"
	case Code_NG:
		s += "ng"
	case Code_N:
		s += "n"
	case Code_MP:
		s += "mp"
	case Code_M:
		s += "m"
	case Code_R:
		s += "r"
	default:
		return ""
	}
	switch getPitchCode(code) {
	case Code_LowPitch:
		switch getVowelCode(code) {
		case Code_ClosedA:
			s += "a"
		case Code_NasalA:
			s += "aN"
		case Code_ClosedE:
			s += "e"
		case Code_NasalE:
			s += "eN"
		case Code_ClosedI:
			s += "i"
		case Code_NasalI:
			s += "iN"
		case Code_ClosedO:
			s += "o"
		case Code_NasalO:
			s += "oN"
		case Code_OpenE:
			s += "ɛ"
		case Code_OpenO:
			s += "ɔ"
		case Code_ClosedU:
			s += "u"
		case Code_NasalU:
			s += "uN"
		case Code_UnknownE:
			s += "x"
		case Code_UnknownO:
			s += "c"
		default:
			return ""
		}
	case Code_MidPitch:
		switch getVowelCode(code) {
		case Code_ClosedA:
			s += "ä"
		case Code_NasalA:
			s += "äN"
		case Code_ClosedE:
			s += "ë"
		case Code_NasalE:
			s += "ëN"
		case Code_ClosedI:
			s += "ï"
		case Code_NasalI:
			s += "ïN"
		case Code_ClosedO:
			s += "ö"
		case Code_NasalO:
			s += "öN"
		case Code_OpenE:
			s += "ɛ̈"
		case Code_OpenO:
			s += "ɔ̈"
		case Code_ClosedU:
			s += "ü"
		case Code_NasalU:
			s += "üN"
		case Code_UnknownE:
			s += "ẍ"
		case Code_UnknownO:
			s += "c̈"
		default:
			return ""
		}
	case Code_HighPitch:
		switch getVowelCode(code) {
		case Code_ClosedA:
			s += "â"
		case Code_NasalA:
			s += "âN"
		case Code_ClosedE:
			s += "ê"
		case Code_NasalE:
			s += "êN"
		case Code_ClosedI:
			s += "î"
		case Code_NasalI:
			s += "îN"
		case Code_ClosedO:
			s += "ô"
		case Code_NasalO:
			s += "ôN"
		case Code_OpenE:
			s += "Ɛ̂"
		case Code_OpenO:
			s += "ɔ̂"
		case Code_ClosedU:
			s += "û"
		case Code_NasalU:
			s += "ûN"
		case Code_UnknownE:
			s += "x̂"
		case Code_UnknownO:
			s += "ĉ"
		default:
			return ""
		}
	case Code_UnknownPitch:
		switch getVowelCode(code) {
		case Code_ClosedA:
			s += "ạ"
		case Code_NasalA:
			s += "ạN"
		case Code_ClosedE:
			s += "ẹ"
		case Code_NasalE:
			s += "ẹN"
		case Code_ClosedI:
			s += "ị"
		case Code_NasalI:
			s += "ịN"
		case Code_ClosedO:
			s += "ọ"
		case Code_NasalO:
			s += "ọN"
		case Code_OpenE:
			s += "ɛ̣"
		case Code_OpenO:
			s += "ɔ̣"
		case Code_ClosedU:
			s += "ụ"
		case Code_NasalU:
			s += "ụN"
		case Code_UnknownE:
			s += "x̣"
		case Code_UnknownO:
			s += "c̣"
		default:
			return ""
		}
	}
	switch caseCode {
	case Code_Lowercase:
		s = strings.ToLower(s)
	case Code_Titlecase:
		s = strings.ToTitle(s)
	case Code_Uppercase:
		s = strings.ToUpper(s)
	}
	if getPrefixCode(code) == Code_SpacePrefix {
		return " " + s
	}
	if getInfixCode(code) == Code_HyphenInfix {
		return "-" + s
	}
	return s
}

func canonicalFromSangoSyllableCode(code uint16) string {
	if (uint16(code) & 0x8000) == 0 {
		panic("code does not represent Sango")
	}
	s := ""
	if getPrefixCode(code) == Code_SpacePrefix {
		s += " "
	} else if getInfixCode(code) == Code_HyphenInfix {
		s += "-"
	}
	switch getCaseCode(code) {
	case Code_Lowercase:
		// No case prefix
	case Code_Titlecase:
		s += "~"
	case Code_Uppercase:
		s += "="
	case Code_Invisible:
		return ""
	}
	switch getConsonantCode(code) {
	case Code_h:
		s += "h"
	case Code_B:
		s += "b"
	case Code_V:
		s += "v"
	case Code_Y:
		s += "y"
	case Code_D:
		s += "d"
	case Code_Z:
		s += "z"
	case Code_Q:
		s += "q"
	case Code_G:
		s += "g"
	case Code_H:
		s += "H"
	case Code_P:
		s += "p"
	case Code_F:
		s += "f"
	case Code_L:
		s += "l"
	case Code_T:
		s += "t"
	case Code_S:
		s += "s"
	case Code_KP:
		s += "K"
	case Code_K:
		s += "k"
	case Code_W:
		s += "w"
	case Code_MB:
		s += "B"
	case Code_MV:
		s += "V"
	case Code_NY:
		s += "Y"
	case Code_ND:
		s += "D"
	case Code_NZ:
		s += "Z"
	case Code_NGB:
		s += "Q"
	case Code_NG:
		s += "G"
	case Code_N:
		s += "n"
	case Code_MP:
		s += "P"
	case Code_M:
		s += "m"
	case Code_R:
		s += "r"
	default:
		return ""
	}
	switch getVowelCode(code) {
	case Code_ClosedA:
		s += "a"
	case Code_NasalA:
		s += "A"
	case Code_ClosedE:
		s += "e"
	case Code_NasalE:
		s += "E"
	case Code_ClosedI:
		s += "i"
	case Code_NasalI:
		s += "I"
	case Code_ClosedO:
		s += "o"
	case Code_NasalO:
		s += "O"
	case Code_OpenE:
		s += "x"
	case Code_OpenO:
		s += "c"
	case Code_ClosedU:
		s += "u"
	case Code_NasalU:
		s += "U"
	case Code_UnknownE:
		s += "X"
	case Code_UnknownO:
		s += "C"
	default:
		return ""
	}
	switch getPitchCode(code) {
	case Code_LowPitch:
		s += "_"
	case Code_MidPitch:
		s += ":"
	case Code_HighPitch:
		s += "^"
	}
	return s
}

var canonicalRE = regexp.MustCompile(`^(([ -]?)([~=]?)([hbvydzqgHpfltsKkwBVYDZQGnPmr])([aAeEiIoOxcuUXC])([_:^]?))$`)

func canonicalToSangoSyllableCode(s string) uint16 {
	// TODO: finish functionality and test on ../lexicon/lexicon.go
	return 0xFFFF
}

func extractBits(code uint16, lsb int, msb int) uint16 {
	if msb < lsb || lsb > 15 || msb < 0 {
		return uint16(0)
	}
	if lsb < 0 {
		lsb = 0
	}
	if msb > 15 {
		msb = 15
	}
	bits := uint16(code)
	bits >>= lsb
	shift := msb + 1 - lsb
	if shift < 15 {
		bits %= (1 << shift)
	}
	return bits
}
