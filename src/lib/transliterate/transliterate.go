package transliterate

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/rivo/uniseg"
	"golang.org/x/text/unicode/norm"
)

func EncodeInput()  { encodeInput() }
func EncodeOutput() { encodeOutput() }
func DecodeInput()  { decodeInput() }
func DecodeOutput() { decodeOutput() }

////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var utf8ToAsciiInput = map[string]string{
	"-":  "",
	"A":  "A",
	"Ä":  "jA",
	"Â":  "JA",
	"E":  "E",
	"Ë":  "jE",
	"Ê":  "JE",
	"Ɛ":  "X",
	"Ɛ̈": "jX",
	"Ɛ̂": "JX",
	"I":  "I",
	"Ï":  "jI",
	"Î":  "JI",
	"O":  "O",
	"Ö":  "jO",
	"Ô":  "JO",
	"Ɔ":  "C",
	"Ɔ̈": "jC",
	"Ɔ̂": "JC",
	"U":  "U",
	"Ü":  "jU",
	"Û":  "JU",
	"B":  "B",
	"D":  "D",
	"F":  "F",
	"G":  "G",
	"H":  "H",
	"K":  "K",
	"L":  "L",
	"M":  "M",
	"N":  "N",
	"P":  "P",
	"R":  "R",
	"S":  "S",
	"T":  "T",
	"V":  "V",
	"W":  "W",
	"Y":  "Y",
	"Z":  "Z",
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

var utf8VowelToAsciiOutput = map[string]string{
	"A":  "a",
	"Ä":  "a",
	"Â":  "a",
	"E":  "e",
	"Ë":  "e",
	"Ê":  "e",
	"Ɛ":  "x",
	"Ɛ̈": "x",
	"Ɛ̂": "x",
	"I":  "i",
	"Ï":  "i",
	"Î":  "i",
	"O":  "o",
	"Ö":  "o",
	"Ô":  "o",
	"Ɔ":  "c",
	"Ɔ̈": "c",
	"Ɔ̂": "c",
	"U":  "u",
	"Ü":  "u",
	"Û":  "u",
	"a":  "a",
	"ä":  "a",
	"â":  "a",
	"e":  "e",
	"ë":  "e",
	"ê":  "e",
	"ɛ":  "x",
	"ɛ̈": "x",
	"ɛ̂": "x",
	"i":  "i",
	"ï":  "i",
	"î":  "i",
	"o":  "o",
	"ö":  "o",
	"ô":  "o",
	"ɔ":  "c",
	"ɔ̈": "c",
	"ɔ̂": "c",
	"u":  "u",
	"ü":  "u",
	"û":  "u",
}

var utf8ConsonantToAsciiOutput = map[string]string{
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

var isHighPitch = map[string]bool{
	"Ä":  false,
	"Â":  true,
	"Ë":  false,
	"Ê":  true,
	"Ɛ̈": false,
	"Ɛ̂": true,
	"Ï":  false,
	"Î":  true,
	"Ö":  false,
	"Ô":  true,
	"Ɔ̈": false,
	"Ɔ̂": true,
	"Ü":  false,
	"Û":  true,
	"ä":  false,
	"â":  true,
	"ë":  false,
	"ê":  true,
	"ɛ̈": false,
	"ɛ̂": true,
	"ï":  false,
	"î":  true,
	"ö":  false,
	"ô":  true,
	"ɔ̈": false,
	"ɔ̂": true,
	"ü":  false,
	"û":  true,
}

func encodeInput() {
	r := norm.NFKC.Reader(bufio.NewReader(os.Stdin))
	b, err := io.ReadAll(r)
	check(err)

	state := -1
	var c []byte
	var word string
	for len(b) > 0 {
		c, b, _, state = uniseg.Step(b, state)
		s := string(c)
		a, isSangoUTF8 := utf8ToAsciiInput[s]
		if isSangoUTF8 {
			word += a
		} else {
			fmt.Printf("%v", word)
			if s != "\n" {
				fmt.Printf("%v", s)
			}
			word = ""
		}
	}
}

func encodeOutput() {
	r := norm.NFKC.Reader(bufio.NewReader(os.Stdin))
	b, err := io.ReadAll(r)
	check(err)

	state := -1
	var c []byte
	for len(b) > 0 {
		consonants := ""
		s := ""
		for len(b) > 0 {
			c, b, _, state = uniseg.Step(b, state)
			s = string(c)
			if consonant, isConsonant := utf8ConsonantToAsciiOutput[s]; isConsonant {
				if consonants == "n" && consonant != "d" && consonant != "g" && consonant != "y" && consonant != "z" {
					fmt.Print("N")
					consonants = ""
					continue
				}
				consonants += consonant
			} else {
				break
			}
		}
		if vowel, isVowel := utf8VowelToAsciiOutput[s]; isVowel {
			if isHigh, isMedOrHigh := isHighPitch[s]; isMedOrHigh {
				vowel = strings.ToUpper(vowel)
				if isHigh {
					consonants = strings.ToUpper(consonants)
				}
			}
			fmt.Printf("%s%s", consonants, vowel)
		} else if consonants == "n" {
			fmt.Print("N")
		} else if consonants != "" {
			panic("\nConsonants {" + consonants + "} not followed by a vowel {" + s + "}")
		} else {
			fmt.Printf("%s", s)
		}
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

func decodeInput() {
	r := bufio.NewReader(os.Stdin)
	isUpperCaseLetter := false
	isUpperCaseWord := false
	pitch := 0
	for c, _, err := r.ReadRune(); err == nil; c, _, err = r.ReadRune() {
		s := string(c)
		switch s {
		case "j":
			pitch = 1
			continue
		case "J":
			pitch = 2
			continue
		case "q":
			isUpperCaseLetter = true
			continue
		case "Q":
			isUpperCaseWord = true
			continue
		}
		isUpperCase := isUpperCaseLetter || isUpperCaseWord
		if m, isVowel := asciiInputToUtf8[isUpperCase][pitch][s]; isVowel {
			fmt.Print(m)
		} else if _, isConsonant := utf8ConsonantToAsciiOutput[s]; isConsonant {
			if isUpperCase {
				fmt.Print(strings.ToUpper(s))
			} else {
				fmt.Print(s)
			}
		} else {
			fmt.Print(s)
			isUpperCaseWord = false
		}
		pitch = 0
		isUpperCaseLetter = false
	}
}

func decodeOutput() {
  panic("'transliterate decode output' is not implemented")
}
