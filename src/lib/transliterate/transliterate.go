package transliterate

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/rivo/uniseg"
	"golang.org/x/text/unicode/norm"
)

func EncodeInput()  { encodeInput() }
func EncodeOutput() { encodeOutput() }
func DecodeInput()  { decode(false) }
func DecodeOutput() { decode(true) }

////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var utf8ToAsciiInput = map[string]string{
	"-":  "",
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

type AsciiIsHighPitch = struct {
	ascii       string
	isHighPitch bool
}

var asciiIsHighPitchFromUTF8Vowel = map[string]AsciiIsHighPitch{
	"A":  {"a", false},
	"Ä":  {"A", false},
	"Â":  {"A", true},
	"E":  {"e", false},
	"Ë":  {"E", false},
	"Ê":  {"E", true},
	"Ɛ":  {"x", false},
	"Ɛ̈": {"X", false},
	"Ɛ̂": {"X", true},
	"I":  {"i", false},
	"Ï":  {"I", false},
	"Î":  {"I", true},
	"O":  {"o", false},
	"Ö":  {"O", false},
	"Ô":  {"O", true},
	"Ɔ":  {"c", false},
	"Ɔ̈": {"C", false},
	"Ɔ̂": {"C", true},
	"U":  {"u", false},
	"Ü":  {"U", false},
	"Û":  {"U", true},
	"a":  {"a", false},
	"ä":  {"A", false},
	"â":  {"A", true},
	"e":  {"e", false},
	"ë":  {"E", false},
	"ê":  {"E", true},
	"ɛ":  {"x", false},
	"ɛ̈": {"X", false},
	"ɛ̂": {"X", true},
	"i":  {"i", false},
	"ï":  {"I", false},
	"î":  {"I", true},
	"o":  {"o", false},
	"ö":  {"O", false},
	"ô":  {"O", true},
	"ɔ":  {"c", false},
	"ɔ̈": {"C", false},
	"ɔ̂": {"C", true},
	"u":  {"u", false},
	"ü":  {"U", false},
	"û":  {"U", true},
}

func encodeInput() {
	r := norm.NFKC.Reader(bufio.NewReader(os.Stdin))
	b, err := io.ReadAll(r)
	check(err)

	state := -1
	var ccc []byte
	for len(b) > 0 {
		ccc, b, _, state = uniseg.Step(b, state)
		s := string(ccc)
		if a, isSangoUTF8 := utf8ToAsciiInput[s]; isSangoUTF8 {
			fmt.Printf("%s", a)
		} else {
			fmt.Printf("%s", s)
		}
	}
}

func encodeOutput() {
	r := norm.NFKC.Reader(bufio.NewReader(os.Stdin))
	b, err := io.ReadAll(r)
	check(err)

	state := -1
	var c []byte
	consonantsWithQ := ""
	consonantsWithoutQ := ""
	for len(b) > 0 {
		c, b, _, state = uniseg.Step(b, state)
		s := string(c)
		if consonant, isConsonant := lowercaseAsciiConsonant[s]; isConsonant {
			if consonantsWithoutQ == "n" && consonant != "d" && consonant != "g" && consonant != "y" && consonant != "z" {
				fmt.Print("N")
				consonantsWithoutQ = consonant
				consonantsWithQ = consonant
				continue
			}
			if len(s) > 0 && s == strings.ToUpper(s) {
				consonantsWithQ += "q"
			}
			consonantsWithoutQ += consonant
			consonantsWithQ += consonant
			continue
		}
		if asciiIsHighPitch, isVowel := asciiIsHighPitchFromUTF8Vowel[s]; isVowel {
			if asciiIsHighPitch.isHighPitch {
				consonantsWithQ = strings.ReplaceAll(strings.ToUpper(consonantsWithQ), "Q", "q")
			}
			if len(s) > 0 && s == strings.ToUpper(s) {
				consonantsWithQ += "q"
			}
			fmt.Printf("%s%s", consonantsWithQ, asciiIsHighPitch.ascii)
		} else if consonantsWithoutQ == "n" {
			fmt.Printf("%s%s", "N", s)
		} else {
			// TODO: This is a word break. If every other glyph output is a 'q', replace with single initial 'Q'.
			fmt.Printf("%s%s", consonantsWithQ, s)
		}
		consonantsWithQ = ""
		consonantsWithoutQ = ""
	}
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
			"u": "o",
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

func decode(isOutputFormat bool) {
	r := bufio.NewReader(os.Stdin)
	isUpperCaseLetter := false
	isUpperCaseWord := false
	numLowerCaseConsonants := 0
	numUpperCaseConsonants := 0
	pitch := 0
	word := ""
	for c, _, err := r.ReadRune(); err == nil; c, _, err = r.ReadRune() {
		if c == 'j' {
			pitch = 1
			continue
		}
		if c == 'J' {
			pitch = 2
			continue
		}
		if c == 'q' {
			isUpperCaseLetter = true
			continue
		}
		if c == 'Q' {
			isUpperCaseWord = true
			continue
		}
		isUpperCase := isUpperCaseLetter || isUpperCaseWord
		isUpperCaseLetter = false
		s := string(c)
		if consonant, isConsonant := lowercaseAsciiConsonant[s]; isConsonant {
			if consonant == s {
				numLowerCaseConsonants++
			} else {
				numUpperCaseConsonants++
			}
			if isUpperCase {
				word += strings.ToUpper(consonant)
			} else {
				word += consonant
			}
			continue
		}
		if asciiVowel, isAsciiVowel := lowercaseAsciiVowel[s]; isAsciiVowel {
			if isOutputFormat && pitch == 0 {
				if numUpperCaseConsonants > 0 {
					pitch = 2
				} else if asciiVowel == s {
					// Leave pitch unchanged
				} else if numLowerCaseConsonants == 0 {
					pitch = 2
				} else {
					pitch = 1
				}
			}
			vowel, isVowel := asciiInputToUtf8[isUpperCase][pitch][asciiVowel]
			if !isVowel {
				panic("asciiInputToUtf8[" + strconv.FormatBool(isUpperCase) + "][" + strconv.Itoa(pitch) + "][" + asciiVowel +
					"] does not map to a UTF8 vowel")
			}
			word += vowel
			pitch = 0
			numLowerCaseConsonants = 0
			numUpperCaseConsonants = 0
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

			default:
				// COMMON: any word that ends in "ngɔ̈" or "NGƆ̈" (verbal gerund form)
				// Also include "ngö" and "ngö" in case the vowel height is incorrect.
				if strings.HasSuffix(word, "ngɔ̈") ||
					strings.HasSuffix(word, "NGƆ̈") ||
					strings.HasSuffix(word, "ngö") ||
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
		}

		fmt.Printf("%s%s", word, s)
		word = ""
		pitch = 0
		isUpperCaseWord = false
		numLowerCaseConsonants = 0
		numUpperCaseConsonants = 0
	}
}
