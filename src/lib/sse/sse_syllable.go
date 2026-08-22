// SSE Sango Syllable
//
// Encoding for a phonemically-valid Sango syllable.

package sse

import (
	"fmt"
	"regexp"
	"strings"
)

func IsSango(code uint16) bool { return isSango(code) }
func IsValid(code uint16) bool { return isValid(code) }

type PrefixCode uint16
type ShiftCode uint16
type InfixCode uint16
type ConsonantCode uint16
type VowelCode uint16
type PitchCode uint16

const (
	IsSango_MASK        uint16        = 0b1_0_00_0_00000_0000_00 // Assumed to be Unicode rune if not set
	PrefixCode_None     PrefixCode    = 0b0_0_00_0_00000_0000_00
	PrefixCode_Space    PrefixCode    = 0b0_1_00_0_00000_0000_00
	PrefixCode_MASK     uint16        = 0b0_1_00_0_00000_0000_00
	ShiftCode_Invisible ShiftCode     = 0b0_0_00_0_00000_0000_00
	ShiftCode_lower     ShiftCode     = 0b0_0_01_0_00000_0000_00
	ShiftCode_Title     ShiftCode     = 0b0_0_10_0_00000_0000_00
	ShiftCode_UPPER     ShiftCode     = 0b0_0_11_0_00000_0000_00
	ShiftCode_MASK      uint16        = 0b0_0_11_0_00000_0000_00
	InfixCode_None      InfixCode     = 0b0_0_00_0_00000_0000_00
	InfixCode_Hyphen    InfixCode     = 0b0_0_00_1_00000_0000_00
	InfixCode_MASK      uint16        = 0b0_0_00_1_00000_0000_00
	ConsonantCode_None  ConsonantCode = 0b0_0_00_0_00000_0000_00 // Not found in valid Sango
	ConsonantCode_h     ConsonantCode = 0b0_0_00_0_00010_0000_00 // omitted from UTF8
	ConsonantCode_H     ConsonantCode = 0b0_0_00_0_00011_0000_00 // h
	ConsonantCode_b     ConsonantCode = 0b0_0_00_0_00100_0000_00 // b
	ConsonantCode_B     ConsonantCode = 0b0_0_00_0_00101_0000_00 // mb
	ConsonantCode_q     ConsonantCode = 0b0_0_00_0_00110_0000_00 // gb
	ConsonantCode_Q     ConsonantCode = 0b0_0_00_0_00111_0000_00 // ngb
	ConsonantCode_d     ConsonantCode = 0b0_0_00_0_01000_0000_00 // d
	ConsonantCode_D     ConsonantCode = 0b0_0_00_0_01001_0000_00 // nd
	ConsonantCode_f     ConsonantCode = 0b0_0_00_0_01010_0000_00 // f
	ConsonantCode_g     ConsonantCode = 0b0_0_00_0_01011_0000_00 // g
	ConsonantCode_G     ConsonantCode = 0b0_0_00_0_01100_0000_00 // ng
	ConsonantCode_k     ConsonantCode = 0b0_0_00_0_01101_0000_00 // k
	ConsonantCode_l     ConsonantCode = 0b0_0_00_0_01110_0000_00 // l
	ConsonantCode_r     ConsonantCode = 0b0_0_00_0_01111_0000_00 // r
	ConsonantCode_m     ConsonantCode = 0b0_0_00_0_10000_0000_00 // m
	ConsonantCode_n     ConsonantCode = 0b0_0_00_0_10001_0000_00 // n
	ConsonantCode_p     ConsonantCode = 0b0_0_00_0_10010_0000_00 // p
	ConsonantCode_K     ConsonantCode = 0b0_0_00_0_10011_0000_00 // kp
	ConsonantCode_P     ConsonantCode = 0b0_0_00_0_10100_0000_00 // mp
	ConsonantCode_s     ConsonantCode = 0b0_0_00_0_10101_0000_00 // s
	ConsonantCode_t     ConsonantCode = 0b0_0_00_0_10110_0000_00 // t
	ConsonantCode_v     ConsonantCode = 0b0_0_00_0_10111_0000_00 // v
	ConsonantCode_V     ConsonantCode = 0b0_0_00_0_11000_0000_00 // mv
	ConsonantCode_w     ConsonantCode = 0b0_0_00_0_11001_0000_00 // w
	ConsonantCode_y     ConsonantCode = 0b0_0_00_0_11010_0000_00 // y
	ConsonantCode_Y     ConsonantCode = 0b0_0_00_0_11011_0000_00 // ny
	ConsonantCode_z     ConsonantCode = 0b0_0_00_0_11100_0000_00 // z
	ConsonantCode_Z     ConsonantCode = 0b0_0_00_0_11101_0000_00 // nz
	ConsonantCode_MASK  uint16        = 0b0_0_00_0_11111_0000_00
	VowelCode_None      VowelCode     = 0b0_0_00_0_00000_0000_00 // Not found in valid Sango
	VowelCode_a         VowelCode     = 0b0_0_00_0_00000_0010_00 // a
	VowelCode_A         VowelCode     = 0b0_0_00_0_00000_0011_00 // añ
	VowelCode_X         VowelCode     = 0b0_0_00_0_00000_0100_00 // x (e with unknown height)
	VowelCode_x         VowelCode     = 0b0_0_00_0_00000_0101_00 // ɛ
	VowelCode_e         VowelCode     = 0b0_0_00_0_00000_0110_00 // e
	VowelCode_E         VowelCode     = 0b0_0_00_0_00000_0111_00 // eñ
	VowelCode_i         VowelCode     = 0b0_0_00_0_00000_1000_00 // i
	VowelCode_I         VowelCode     = 0b0_0_00_0_00000_1001_00 // iñ
	VowelCode_C         VowelCode     = 0b0_0_00_0_00000_1010_00 // c (o with unknown height)
	VowelCode_c         VowelCode     = 0b0_0_00_0_00000_1011_00 // ɔ
	VowelCode_o         VowelCode     = 0b0_0_00_0_00000_1100_00 // o
	VowelCode_O         VowelCode     = 0b0_0_00_0_00000_1101_00 // oñ
	VowelCode_u         VowelCode     = 0b0_0_00_0_00000_1110_00 // u
	VowelCode_U         VowelCode     = 0b0_0_00_0_00000_1111_00 // uñ
	VowelCode_MASK      uint16        = 0b0_0_00_0_00000_1111_00
	PitchCode_Unknown   PitchCode     = 0b0_0_00_0_00000_0000_00
	PitchCode_Low       PitchCode     = 0b0_0_00_0_00000_0000_01
	PitchCode_Mid       PitchCode     = 0b0_0_00_0_00000_0000_10
	PitchCode_High      PitchCode     = 0b0_0_00_0_00000_0000_11
	PitchCode_MASK      uint16        = 0b0_0_00_0_00000_0000_11
)

func getPrefixCode(code uint16) PrefixCode       { return PrefixCode(code & PrefixCode_MASK) }
func getShiftCode(code uint16) ShiftCode         { return ShiftCode(code & ShiftCode_MASK) }
func getInfixCode(code uint16) InfixCode         { return InfixCode(code & InfixCode_MASK) }
func getConsonantCode(code uint16) ConsonantCode { return ConsonantCode(code & ConsonantCode_MASK) }
func getVowelCode(code uint16) VowelCode         { return VowelCode(code & VowelCode_MASK) }
func getPitchCode(code uint16) PitchCode         { return PitchCode(code & PitchCode_MASK) }

func isSango(code uint16) bool {
	switch code & IsSango_MASK {
	case IsSango_MASK:
	default:
		return false
	}
	return true
}

func hasValidPrefix(code uint16) bool {
	switch getPrefixCode(code) {
	case PrefixCode_None:
	case PrefixCode_Space:
	default:
		return false
	}
	return true
}

func hasValidShift(code uint16) bool {
	switch getShiftCode(code) {
	case ShiftCode_Invisible:
	case ShiftCode_lower:
	case ShiftCode_Title:
	case ShiftCode_UPPER:
	default:
		return false
	}
	return true
}

func hasValidInfix(code uint16) bool {
	switch getInfixCode(code) {
	case InfixCode_None:
	case InfixCode_Hyphen:
	default:
		return false
	}
	return true
}

func hasValidConsonant(code uint16) bool {
	switch getConsonantCode(code) {
	case ConsonantCode_h:
	case ConsonantCode_H:
	case ConsonantCode_b:
	case ConsonantCode_B:
	case ConsonantCode_q:
	case ConsonantCode_Q:
	case ConsonantCode_d:
	case ConsonantCode_D:
	case ConsonantCode_f:
	case ConsonantCode_g:
	case ConsonantCode_G:
	case ConsonantCode_k:
	case ConsonantCode_l:
	case ConsonantCode_r:
	case ConsonantCode_m:
	case ConsonantCode_n:
	case ConsonantCode_p:
	case ConsonantCode_K:
	case ConsonantCode_P:
	case ConsonantCode_s:
	case ConsonantCode_t:
	case ConsonantCode_v:
	case ConsonantCode_V:
	case ConsonantCode_w:
	case ConsonantCode_y:
	case ConsonantCode_Y:
	case ConsonantCode_z:
	case ConsonantCode_Z:
	default:
		return false
	}
	return true
}

func hasValidVowel(code uint16) bool {
	switch getVowelCode(code) {
	case VowelCode_a:
	case VowelCode_A:
	case VowelCode_X:
	case VowelCode_x:
	case VowelCode_e:
	case VowelCode_E:
	case VowelCode_i:
	case VowelCode_I:
	case VowelCode_C:
	case VowelCode_c:
	case VowelCode_o:
	case VowelCode_O:
	case VowelCode_u:
	case VowelCode_U:
	default:
		return false
	}
	return true
}

func hasValidPitch(code uint16) bool {
	switch getPitchCode(code) {
	case PitchCode_Unknown:
	case PitchCode_Low:
	case PitchCode_Mid:
	case PitchCode_High:
	default:
		return false
	}
	return true
}

func isValid(code uint16) bool {
	if code == 0 {
		return false
	}
	if !isSango(code) {
		return true
	}
	return hasValidPrefix(code) &&
		hasValidShift(code) &&
		hasValidInfix(code) &&
		hasValidConsonant(code) &&
		hasValidVowel(code) &&
		hasValidPitch(code)
}

func utf8FromSangoCodeValue(code uint16, options WriteUTF8Options) string {
	s := ""
	shiftCode := getShiftCode(code)
	if shiftCode == ShiftCode_Invisible {
		return ""
	}
	pitch := PitchCode_Low
	if options.WithPitch {
		pitch = getPitchCode(code)
	}
	switch getConsonantCode(code) {
	case ConsonantCode_h:
		// omit consonant for unaspirated h
	case ConsonantCode_H:
		s += "h"
	case ConsonantCode_b:
		s += "b"
	case ConsonantCode_B:
		s += "mb"
	case ConsonantCode_q:
		s += "gb"
	case ConsonantCode_Q:
		s += "ngb"
	case ConsonantCode_d:
		s += "d"
	case ConsonantCode_D:
		s += "nd"
	case ConsonantCode_f:
		s += "f"
	case ConsonantCode_g:
		s += "g"
	case ConsonantCode_G:
		s += "ng"
	case ConsonantCode_k:
		s += "k"
	case ConsonantCode_l:
		s += "l"
	case ConsonantCode_r:
		s += "r"
	case ConsonantCode_m:
		s += "m"
	case ConsonantCode_n:
		s += "n"
	case ConsonantCode_p:
		s += "p"
	case ConsonantCode_K:
		s += "kp"
	case ConsonantCode_P:
		s += "mp"
	case ConsonantCode_s:
		s += "s"
	case ConsonantCode_t:
		s += "t"
	case ConsonantCode_v:
		s += "v"
	case ConsonantCode_V:
		s += "mv"
	case ConsonantCode_w:
		s += "w"
	case ConsonantCode_y:
		s += "y"
	case ConsonantCode_Y:
		s += "ny"
	case ConsonantCode_z:
		s += "z"
	case ConsonantCode_Z:
		s += "nz"
	default:
		return ""
	}
	switch pitch {
	case PitchCode_Low:
		switch getVowelCode(code) {
		case VowelCode_a:
			s += "a"
		case VowelCode_A:
			s += "añ"
		case VowelCode_X:
			s += "x"
		case VowelCode_x:
			s += "ɛ"
		case VowelCode_e:
			s += "e"
		case VowelCode_E:
			s += "eñ"
		case VowelCode_i:
			s += "i"
		case VowelCode_I:
			s += "iñ"
		case VowelCode_C:
			s += "c"
		case VowelCode_c:
			s += "ɔ"
		case VowelCode_o:
			s += "o"
		case VowelCode_O:
			s += "oñ"
		case VowelCode_u:
			s += "u"
		case VowelCode_U:
			s += "uñ"
		default:
			return ""
		}
	case PitchCode_Mid:
		switch getVowelCode(code) {
		case VowelCode_a:
			s += "ä"
		case VowelCode_A:
			s += "äñ"
		case VowelCode_X:
			s += "ẍ"
		case VowelCode_x:
			s += "ɛ̈"
		case VowelCode_e:
			s += "ë"
		case VowelCode_E:
			s += "ëñ"
		case VowelCode_i:
			s += "ï"
		case VowelCode_I:
			s += "ïñ"
		case VowelCode_C:
			s += "c̈"
		case VowelCode_c:
			s += "ɔ̈"
		case VowelCode_o:
			s += "ö"
		case VowelCode_O:
			s += "öñ"
		case VowelCode_u:
			s += "ü"
		case VowelCode_U:
			s += "üñ"
		default:
			return ""
		}
	case PitchCode_High:
		switch getVowelCode(code) {
		case VowelCode_a:
			s += "â"
		case VowelCode_A:
			s += "âñ"
		case VowelCode_X:
			s += "x̂"
		case VowelCode_x:
			s += "ɛ̂"
		case VowelCode_e:
			s += "ê"
		case VowelCode_E:
			s += "êñ"
		case VowelCode_i:
			s += "î"
		case VowelCode_I:
			s += "îñ"
		case VowelCode_C:
			s += "ĉ"
		case VowelCode_c:
			s += "ɔ̂"
		case VowelCode_o:
			s += "ô"
		case VowelCode_O:
			s += "ôñ"
		case VowelCode_u:
			s += "û"
		case VowelCode_U:
			s += "ûñ"
		default:
			return ""
		}
	case PitchCode_Unknown:
		switch getVowelCode(code) {
		case VowelCode_a:
			s += "ạ"
		case VowelCode_A:
			s += "ạñ"
		case VowelCode_X:
			s += "x̣"
		case VowelCode_x:
			s += "ɛ̣"
		case VowelCode_e:
			s += "ẹ"
		case VowelCode_E:
			s += "ẹñ"
		case VowelCode_i:
			s += "ị"
		case VowelCode_I:
			s += "ịñ"
		case VowelCode_C:
			s += "c̣"
		case VowelCode_c:
			s += "ɔ̣"
		case VowelCode_o:
			s += "ọ"
		case VowelCode_O:
			s += "ọñ"
		case VowelCode_u:
			s += "ụ"
		case VowelCode_U:
			s += "ụñ"
		default:
			return ""
		}
	}
	if options.WithShift {
		switch shiftCode {
		case ShiftCode_lower:
			s = strings.ToLower(s)
		case ShiftCode_Title:
			s = strings.ToTitle(s)
		case ShiftCode_UPPER:
			s = strings.ToUpper(s)
		}
	}
	if !options.WithHeight {
		s = strings.ReplaceAll(s, "ɛ̈", "ë")
		s = strings.ReplaceAll(s, "ɔ̈", "ö")
		s = strings.ReplaceAll(s, "ɛ̂", "ê")
		s = strings.ReplaceAll(s, "ɔ̂", "ô")
		s = strings.ReplaceAll(s, "ɛ̣", "ẹ")
		s = strings.ReplaceAll(s, "ɔ̣", "ọ")
		s = strings.ReplaceAll(s, "ɛ", "e")
		s = strings.ReplaceAll(s, "ɔ", "o")
		s = strings.ReplaceAll(s, "ẍ", "ë")
		s = strings.ReplaceAll(s, "c̈", "ö")
		s = strings.ReplaceAll(s, "x̂", "ê")
		s = strings.ReplaceAll(s, "ĉ", "ô")
		s = strings.ReplaceAll(s, "x̣", "ẹ")
		s = strings.ReplaceAll(s, "c̣", "ọ")
		s = strings.ReplaceAll(s, "x", "e")
		s = strings.ReplaceAll(s, "c", "o")
	}
	if !options.WithNTilde {
		s = strings.ReplaceAll(s, "ñ", "n")
	}
	if getPrefixCode(code) == PrefixCode_Space {
		return options.ForSpaceUse + s
	}
	if getInfixCode(code) == InfixCode_Hyphen {
		return options.ForHyphenUse + s
	}
	return s
}

func canonicalFromSangoCodeValue(code uint16) string {
	s := ""
	if getPrefixCode(code) == PrefixCode_Space {
		s += " "
	} else if getInfixCode(code) == InfixCode_Hyphen {
		s += "-"
	}
	switch getShiftCode(code) {
	case ShiftCode_lower:
		// No case prefix
	case ShiftCode_Title:
		s += "~"
	case ShiftCode_UPPER:
		s += "="
	case ShiftCode_Invisible:
		s += "#"
	}
	switch getConsonantCode(code) {
	case ConsonantCode_h:
		s += "h"
	case ConsonantCode_H:
		s += "H"
	case ConsonantCode_b:
		s += "b"
	case ConsonantCode_d:
		s += "d"
	case ConsonantCode_f:
		s += "f"
	case ConsonantCode_g:
		s += "g"
	case ConsonantCode_q:
		s += "q"
	case ConsonantCode_k:
		s += "k"
	case ConsonantCode_K:
		s += "K"
	case ConsonantCode_l:
		s += "l"
	case ConsonantCode_r:
		s += "r"
	case ConsonantCode_m:
		s += "m"
	case ConsonantCode_B:
		s += "B"
	case ConsonantCode_P:
		s += "P"
	case ConsonantCode_V:
		s += "V"
	case ConsonantCode_n:
		s += "n"
	case ConsonantCode_D:
		s += "D"
	case ConsonantCode_G:
		s += "G"
	case ConsonantCode_Q:
		s += "Q"
	case ConsonantCode_Y:
		s += "Y"
	case ConsonantCode_Z:
		s += "Z"
	case ConsonantCode_p:
		s += "p"
	case ConsonantCode_s:
		s += "s"
	case ConsonantCode_t:
		s += "t"
	case ConsonantCode_v:
		s += "v"
	case ConsonantCode_w:
		s += "w"
	case ConsonantCode_y:
		s += "y"
	case ConsonantCode_z:
		s += "z"
	default:
		return ""
	}
	switch getVowelCode(code) {
	case VowelCode_a:
		s += "a"
	case VowelCode_A:
		s += "A"
	case VowelCode_X:
		s += "X"
	case VowelCode_x:
		s += "x"
	case VowelCode_e:
		s += "e"
	case VowelCode_E:
		s += "E"
	case VowelCode_i:
		s += "i"
	case VowelCode_I:
		s += "I"
	case VowelCode_C:
		s += "C"
	case VowelCode_c:
		s += "c"
	case VowelCode_o:
		s += "o"
	case VowelCode_O:
		s += "O"
	case VowelCode_u:
		s += "u"
	case VowelCode_U:
		s += "U"
	default:
		return ""
	}
	switch getPitchCode(code) {
	case PitchCode_Low:
		s += "_"
	case PitchCode_Mid:
		s += ":"
	case PitchCode_High:
		s += "^"
	}
	return s
}

var canonicalRE = regexp.MustCompile(`U[+]([0-9A-F]{4})|([ -]?)([~=#]?)([hHbBqQdDfgGklrmnpKPstvVwyYzZ])([aAeEiIoOxcuUXC])([_:^]?)`)

func canonicalToSangoCodeValue(affix, shift, consonant, vowel, pitch string) (uint16, error) {
	var code uint16 = 0x8000
	switch affix {
	case "":
		// do nothing
	case " ":
		code |= uint16(PrefixCode_Space)
	case "-":
		code |= uint16(InfixCode_Hyphen)
	default:
		return IsSango_MASK, fmt.Errorf("bad affix %q", affix)
	}
	switch shift {
	case "":
		code |= uint16(ShiftCode_lower)
	case "~":
		code |= uint16(ShiftCode_Title)
	case "=":
		code |= uint16(ShiftCode_UPPER)
	case "#":
		code |= uint16(ShiftCode_Invisible)
	default:
		return IsSango_MASK, fmt.Errorf("bad shift %q", shift)
	}
	switch consonant {
	case "h":
		code |= uint16(ConsonantCode_h)
	case "H":
		code |= uint16(ConsonantCode_H)
	case "b":
		code |= uint16(ConsonantCode_b)
	case "B":
		code |= uint16(ConsonantCode_B)
	case "q":
		code |= uint16(ConsonantCode_q)
	case "Q":
		code |= uint16(ConsonantCode_Q)
	case "d":
		code |= uint16(ConsonantCode_d)
	case "D":
		code |= uint16(ConsonantCode_D)
	case "f":
		code |= uint16(ConsonantCode_f)
	case "g":
		code |= uint16(ConsonantCode_g)
	case "G":
		code |= uint16(ConsonantCode_G)
	case "k":
		code |= uint16(ConsonantCode_k)
	case "l":
		code |= uint16(ConsonantCode_l)
	case "r":
		code |= uint16(ConsonantCode_r)
	case "m":
		code |= uint16(ConsonantCode_m)
	case "n":
		code |= uint16(ConsonantCode_n)
	case "p":
		code |= uint16(ConsonantCode_p)
	case "K":
		code |= uint16(ConsonantCode_K)
	case "P":
		code |= uint16(ConsonantCode_P)
	case "s":
		code |= uint16(ConsonantCode_s)
	case "t":
		code |= uint16(ConsonantCode_t)
	case "v":
		code |= uint16(ConsonantCode_v)
	case "V":
		code |= uint16(ConsonantCode_V)
	case "w":
		code |= uint16(ConsonantCode_w)
	case "y":
		code |= uint16(ConsonantCode_y)
	case "Y":
		code |= uint16(ConsonantCode_Y)
	case "z":
		code |= uint16(ConsonantCode_z)
	case "Z":
		code |= uint16(ConsonantCode_Z)
	default:
		return IsSango_MASK, fmt.Errorf("bad consonant %q", consonant)
	}
	switch vowel {
	case "a":
		code |= uint16(VowelCode_a)
	case "A":
		code |= uint16(VowelCode_A)
	case "e":
		code |= uint16(VowelCode_e)
	case "E":
		code |= uint16(VowelCode_E)
	case "i":
		code |= uint16(VowelCode_i)
	case "I":
		code |= uint16(VowelCode_I)
	case "o":
		code |= uint16(VowelCode_o)
	case "O":
		code |= uint16(VowelCode_O)
	case "x":
		code |= uint16(VowelCode_x)
	case "c":
		code |= uint16(VowelCode_c)
	case "u":
		code |= uint16(VowelCode_u)
	case "U":
		code |= uint16(VowelCode_U)
	case "X":
		code |= uint16(VowelCode_X)
	case "C":
		code |= uint16(VowelCode_C)
	default:
		return IsSango_MASK, fmt.Errorf("bad vowel %q", vowel)
	}
	switch pitch {
	case "":
		code |= uint16(PitchCode_Unknown)
	case "_":
		code |= uint16(PitchCode_Low)
	case ":":
		code |= uint16(PitchCode_Mid)
	case "^":
		code |= uint16(PitchCode_High)
	default:
		return IsSango_MASK, fmt.Errorf("bad pitch %q", pitch)
	}
	return code, nil
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
