package transliterate

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/rivo/uniseg"
	"golang.org/x/text/unicode/norm"
)

func Normalize(out *bufio.Writer, in *bufio.Reader) error    { return normalize(out, in) }
func EncodeInput(out *bufio.Writer, in *bufio.Reader) error  { return encodeInput(out, in) }
func EncodeOutput(out *bufio.Writer, in *bufio.Reader) error { return encodeOutput(out, in) }
func DecodeInput(out *bufio.Writer, in *bufio.Reader) error  { return decode(out, in, false) }
func DecodeOutput(out *bufio.Writer, in *bufio.Reader) error { return decode(out, in, true) }

////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

// TODO: Greatly simplify the logic by replacing with substring tokenization
// over the fixed set of 7650 tokens comprised by the outer product of consonants
// {"", "b", "d", "f", "g", "gb", "h", "k", "kp", "l", "m", "mb", "mv", "n",
//  "nd", "ng", "ngb", "ny", "nz", "p", "r", "s", "t", "v", "w", "y", "z"},
// vowels {a, an, e, en, ɛ, i, in, o, on, ɔ, u, un}, pitch {low, mid, high},
// and case {lower, upper}^letter, trying longest token to shortest token.
// This can be done efficiently with maps.

func normalize(out *bufio.Writer, in *bufio.Reader) error {
	defer out.Flush()
	r := norm.NFKC.Reader(in)
	b, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	_, err = out.Write(b)
	return err
}

func encodeInput(out *bufio.Writer, in *bufio.Reader) error {
	defer out.Flush()
	src, err := io.ReadAll(norm.NFKC.Reader(in))
	if err != nil {
		return err
	}
	state := -1
	var dst []byte
	for len(src) > 0 {
		dst, src, _, state = uniseg.Step(src, state)
		if a, isSangoUTF8 := utf8ToAsciiInput[string(dst)]; isSangoUTF8 {
			if _, err = out.WriteString(a); err != nil {
				return err
			}
		} else {
			if _, err = out.Write(dst); err != nil {
				return err
			}
		}
	}
	return err
}

const (
	lowPitch = iota
	midPitch
	highPitch
)

type asciiAndPitch = struct {
	ascii string
	pitch int
}

func encodeOutput(out *bufio.Writer, in *bufio.Reader) error {
	defer out.Flush()
	src, err := io.ReadAll(norm.NFKC.Reader(in))
	if err != nil {
		return err
	}
	state := -1
	var dst []byte
	consonantsWithQ := ""
	consonantsWithoutQ := ""
	for len(src) > 0 {
		dst, src, _, state = uniseg.Step(src, state)
		dstStr := string(dst)
		if consonant, isConsonant := lowercaseAsciiConsonant[dstStr]; isConsonant {
			if consonantsWithoutQ == "n" && consonant != "d" && consonant != "g" && consonant != "y" && consonant != "z" {
				if _, err = out.WriteRune('N'); err != nil {
					return err
				}
				consonantsWithoutQ = consonant
				consonantsWithQ = consonant
				continue
			}
			if len(dstStr) > 0 && dstStr == strings.ToUpper(dstStr) {
				consonantsWithQ += "q"
			}
			consonantsWithoutQ += consonant
			consonantsWithQ += consonant
			continue
		}
		if asciiPitch, isVowel := asciiAndPitchFromUTF8Vowel[dstStr]; isVowel {
			if consonantsWithoutQ != "" && asciiPitch.pitch == highPitch {
				consonantsWithQ = strings.ReplaceAll(strings.ReplaceAll(strings.ToUpper(consonantsWithQ), "Q", "q"), "J", "j")
			} else if consonantsWithoutQ == "" && asciiPitch.pitch == midPitch {
				consonantsWithQ += "j"
			}
			if len(dstStr) > 0 && dstStr == strings.ToUpper(dstStr) {
				consonantsWithQ += "q"
			}
			if _, err = out.WriteString(consonantsWithQ); err != nil {
				return err
			}
			if _, err = out.WriteString(asciiPitch.ascii); err != nil {
				return err
			}
		} else if consonantsWithoutQ == "n" {
			if _, err = out.WriteRune('N'); err != nil {
				return err
			}
			if _, err = out.Write(dst); err != nil {
				return err
			}
		} else {
			// TODO: This is a word break. If every other glyph output is a 'q', replace with single initial 'Q'.
			if _, err = out.WriteString(consonantsWithQ); err != nil {
				return err
			}
			if _, err = out.Write(dst); err != nil {
				return err
			}
		}
		consonantsWithQ = ""
		consonantsWithoutQ = ""
	}
	if consonantsWithoutQ == "n" {
		_, err = out.WriteRune('N')
	} else {
		_, err = out.WriteString(consonantsWithQ)
	}
	return err
}

func decode(out *bufio.Writer, in *bufio.Reader, isOutputFormat bool) error {
	defer out.Flush()
	isUpperCaseLetter := false
	isUpperCaseWord := false
	numLowerCaseConsonants := 0
	numUpperCaseConsonants := 0
	pitch := 0
	word := ""
	priorLetter := ' '
	for srcRune, _, err := in.ReadRune(); err == nil; srcRune, _, err = in.ReadRune() {
		if srcRune == 'j' {
			pitch = 1
			continue
		}
		if srcRune == 'J' {
			pitch = 2
			continue
		}
		if srcRune == 'q' {
			isUpperCaseLetter = true
			continue
		}
		if srcRune == 'Q' {
			isUpperCaseWord = true
			continue
		}
		isUpperCase := isUpperCaseLetter || isUpperCaseWord
		isUpperCaseLetter = false
		srcStr := string(srcRune)
		if consonant, isConsonant := lowercaseAsciiConsonant[srcStr]; isConsonant {
			// If the prior letter was an 'n' or 'N', it was syllable ending (vowel nasal), not a consonant.
			// In this case, restore the consonant count used to determine whether the vowel is high or mid.
			if consonant != "d" && consonant != "g" && consonant != "y" && consonant != "z" {
				if priorLetter == 'n' {
					numLowerCaseConsonants--
				} else if priorLetter == 'N' {
					numUpperCaseConsonants--
				}
			}
			if consonant == srcStr {
				numLowerCaseConsonants++
			} else {
				numUpperCaseConsonants++
			}
			if isUpperCase {
				word += strings.ToUpper(consonant)
			} else {
				word += consonant
			}
			priorLetter = srcRune
			continue
		}
		if asciiVowel, isAsciiVowel := lowercaseAsciiVowel[srcStr]; isAsciiVowel {
			if isOutputFormat && pitch == 0 {
				if numUpperCaseConsonants > 0 {
					pitch = 2
				} else if asciiVowel == srcStr {
					// Leave pitch unchanged
				} else if numLowerCaseConsonants == 0 {
					pitch = 2
				} else {
					pitch = 1
				}
			}
			vowel, isVowel := asciiInputToUtf8[isUpperCase][pitch][asciiVowel]
			if !isVowel {
				return fmt.Errorf("asciiInputToUtf8[%v][%v][%v] does not map to a UTF8 vowel", isUpperCase, pitch, asciiVowel)
			}
			word += vowel
			pitch = 0
			numLowerCaseConsonants = 0
			numUpperCaseConsonants = 0
			priorLetter = srcRune
			continue
		}
		if isOutputFormat {
			// Autocorrect words starting with high pitch vowel that should actually be middle pitch.
		autocorrect:
			switch word {
			// RARE: allowlist of known words that start with a middle pitch vowel
			case "âpɛ":
				word = "äpɛ"
			case "âpɛ̈":
				word = "äpɛ̈"
			case "ɛ̂":
				word = "ɛ̈"
			case "êkälïtïse":
				word = "ëkälïtïse"
			case "êpätîte":
				word = "ëpätîte"
			case "î":
				word = "ï"
			case "îrï":
				word = "ïrï"

			case "Âpɛ":
				word = "Äpɛ"
			case "Âpɛ̈":
				word = "Äpɛ̈"
			case "Ɛ̂":
				word = "Ɛ̈"
			case "Êkälïtïse":
				word = "Ëkälïtïse"
			case "Êpätîte":
				word = "Ëpätîte"
			case "Î":
				word = "Ï"
			case "Îrï":
				word = "Ïrï"

			case "ÂPƐ":
				word = "ÄPƐ"
			case "ÂPƐ̈":
				word = "ÄPƐ̈"
			case "ÊKÄLÏTÏSE":
				word = "ËKÄLÏTÏSE"
			case "ÊPÄTÎTE":
				word = "ËPÄTÎTE"
			case "ÎRÏ":
				word = "ÏRÏ"

				// False friends ending in "-ngɔ̈" that are not gerunds
			case "îngɔ̈":
				word = "îngɔ̈"
			case "Îngɔ̈":
				word = "Îngɔ̈"
			case "ÎNGƆ̈":
				word = "ÎNGƆ̈"

			default:
				// COMMON: any word that ends in "ngɔ̈" or "NGƆ̈" (verbal gerund form)
				// Also include "ngö" and "ngö" in case the vowel height is incorrect.
				if strings.HasSuffix(word, "ngɔ̈") ||
					strings.HasSuffix(word, "ngƆ̈") ||
					strings.HasSuffix(word, "nGɔ̈") ||
					strings.HasSuffix(word, "nGƆ̈") ||
					strings.HasSuffix(word, "Ngɔ̈") ||
					strings.HasSuffix(word, "NgƆ̈") ||
					strings.HasSuffix(word, "NGɔ̈") ||
					strings.HasSuffix(word, "NGƆ̈") ||
					strings.HasSuffix(word, "ngö") ||
					strings.HasSuffix(word, "ngÖ") ||
					strings.HasSuffix(word, "nGö") ||
					strings.HasSuffix(word, "nGÖ") ||
					strings.HasSuffix(word, "Ngö") ||
					strings.HasSuffix(word, "NgÖ") ||
					strings.HasSuffix(word, "NGö") ||
					strings.HasSuffix(word, "NGÖ") {
					for _, m := range asciiInputToUtf8 {
						for v, mid := range m[1] {
							hi := m[2][v]
							if suffix, found := strings.CutPrefix(word, hi); found {
								word = mid + suffix
								break autocorrect
							}
						}
					}
				}
			}
			priorLetter = srcRune
		}
		if _, err = out.WriteString(word); err != nil {
			return err
		}
		if _, err = out.WriteString(srcStr); err != nil {
			return err
		}
		word = ""
		pitch = 0
		isUpperCaseWord = false
		numLowerCaseConsonants = 0
		numUpperCaseConsonants = 0
	}
	_, err := out.WriteString(word)
	return err
}

var utf8ToAsciiInput = map[string]string{
	"A":  "qa",
	"Ä":  "jqa",
	"Â":  "Jqa",
	"E":  "qe",
	"Ë":  "jqe",
	"Ê":  "Jqe",
	"Ɛ":  "qx",
	"Ɛ̈": "jqx",
	"Ɛ̂": "Jqx",
	"I":  "qi",
	"Ï":  "jqi",
	"Î":  "Jqi",
	"O":  "qo",
	"Ö":  "jqo",
	"Ô":  "Jqo",
	"Ɔ":  "qc",
	"Ɔ̈": "jqc",
	"Ɔ̂": "Jqc",
	"U":  "qu",
	"Ü":  "jqu",
	"Û":  "Jqu",
	"B":  "qb",
	"D":  "qd",
	"F":  "qf",
	"G":  "qg",
	"H":  "qh",
	"K":  "qk",
	"L":  "ql",
	"M":  "qm",
	"N":  "qn",
	"P":  "qp",
	"R":  "qr",
	"S":  "qs",
	"T":  "qt",
	"V":  "qv",
	"W":  "qw",
	"Y":  "qy",
	"Z":  "qz",
	"a":  "a",
	"ä":  "ja",
	"â":  "Ja",
	"e":  "e",
	"ë":  "je",
	"ê":  "Je",
	"ɛ":  "x",
	"ɛ̈": "jx",
	"ɛ̂": "Jx",
	"i":  "i",
	"ï":  "ji",
	"î":  "Ji",
	"o":  "o",
	"ö":  "jo",
	"ô":  "Jo",
	"ɔ":  "c",
	"ɔ̈": "jc",
	"ɔ̂": "Jc",
	"u":  "u",
	"ü":  "ju",
	"û":  "Ju",
	"b":  "b",
	"d":  "d",
	"f":  "f",
	"g":  "g",
	"h":  "h",
	"k":  "k",
	"l":  "l",
	"m":  "m",
	"n":  "n",
	"p":  "p",
	"r":  "r",
	"s":  "s",
	"t":  "t",
	"v":  "v",
	"w":  "w",
	"y":  "y",
	"z":  "z",
}

var lowercaseAsciiVowel = map[string]string{
	"A": "a",
	"E": "e",
	"X": "x",
	"I": "i",
	"O": "o",
	"C": "c",
	"U": "u",
	"a": "a",
	"e": "e",
	"x": "x",
	"i": "i",
	"o": "o",
	"c": "c",
	"u": "u",
}

var lowercaseAsciiConsonant = map[string]string{
	"B": "b",
	"D": "d",
	"F": "f",
	"G": "g",
	"H": "h",
	"K": "k",
	"L": "l",
	"M": "m",
	"N": "n",
	"P": "p",
	"R": "r",
	"S": "s",
	"T": "t",
	"V": "v",
	"W": "w",
	"Y": "y",
	"Z": "z",
	"b": "b",
	"d": "d",
	"f": "f",
	"g": "g",
	"h": "h",
	"k": "k",
	"l": "l",
	"m": "m",
	"n": "n",
	"p": "p",
	"r": "r",
	"s": "s",
	"t": "t",
	"v": "v",
	"w": "w",
	"y": "y",
	"z": "z",
}

var asciiAndPitchFromUTF8Vowel = map[string]asciiAndPitch{
	"A":  {"a", lowPitch},
	"Ä":  {"A", midPitch},
	"Â":  {"A", highPitch},
	"E":  {"e", lowPitch},
	"Ë":  {"E", midPitch},
	"Ê":  {"E", highPitch},
	"Ɛ":  {"x", lowPitch},
	"Ɛ̈": {"X", midPitch},
	"Ɛ̂": {"X", highPitch},
	"I":  {"i", lowPitch},
	"Ï":  {"I", midPitch},
	"Î":  {"I", highPitch},
	"O":  {"o", lowPitch},
	"Ö":  {"O", midPitch},
	"Ô":  {"O", highPitch},
	"Ɔ":  {"c", lowPitch},
	"Ɔ̈": {"C", midPitch},
	"Ɔ̂": {"C", highPitch},
	"U":  {"u", lowPitch},
	"Ü":  {"U", midPitch},
	"Û":  {"U", highPitch},
	"a":  {"a", lowPitch},
	"ä":  {"A", midPitch},
	"â":  {"A", highPitch},
	"e":  {"e", lowPitch},
	"ë":  {"E", midPitch},
	"ê":  {"E", highPitch},
	"ɛ":  {"x", lowPitch},
	"ɛ̈": {"X", midPitch},
	"ɛ̂": {"X", highPitch},
	"i":  {"i", lowPitch},
	"ï":  {"I", midPitch},
	"î":  {"I", highPitch},
	"o":  {"o", lowPitch},
	"ö":  {"O", midPitch},
	"ô":  {"O", highPitch},
	"ɔ":  {"c", lowPitch},
	"ɔ̈": {"C", midPitch},
	"ɔ̂": {"C", highPitch},
	"u":  {"u", lowPitch},
	"ü":  {"U", midPitch},
	"û":  {"U", highPitch},
}

// asciiInputToUtf8[isUpperCase][pitch] UTF8
var asciiInputToUtf8 = map[bool]map[int]map[string]string{
	false: {
		0: {
			"a": "a",
			"e": "e",
			"x": "ɛ",
			"i": "i",
			"o": "o",
			"c": "ɔ",
			"u": "u",
		},
		1: {
			"a": "ä",
			"e": "ë",
			"x": "ɛ̈",
			"i": "ï",
			"o": "ö",
			"c": "ɔ̈",
			"u": "ü",
		},
		2: {
			"a": "â",
			"e": "ê",
			"x": "ɛ̂",
			"i": "î",
			"o": "ô",
			"c": "ɔ̂",
			"u": "û",
		},
	},
	true: {
		0: {
			"a": "A",
			"e": "E",
			"x": "Ɛ",
			"i": "I",
			"o": "O",
			"c": "Ɔ",
			"u": "U",
		},
		1: {
			"a": "Ä",
			"e": "Ë",
			"x": "Ɛ̈",
			"i": "Ï",
			"o": "Ö",
			"c": "Ɔ̈",
			"u": "Ü",
		},
		2: {
			"a": "Â",
			"e": "Ê",
			"x": "Ɛ̂",
			"i": "Î",
			"o": "Ô",
			"c": "Ɔ̂",
			"u": "Û",
		},
	},
}
