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

var utf8ToAsciiInput map[string]string = map[string]string{
	"-":  "",
	"A":  "qa",
	"Ä":  "qaj",
	"Â":  "qajj",
	"E":  "qe",
	"Ë":  "qej",
	"Ê":  "qejj",
	"Ɛ":  "qx",
	"Ɛ̈": "qxj",
	"Ɛ̂": "qxjj",
	"I":  "qi",
	"Ï":  "qij",
	"Î":  "qijj",
	"O":  "qo",
	"Ö":  "qoj",
	"Ô":  "qojj",
	"Ɔ":  "qc",
	"Ɔ̈": "qcj",
	"Ɔ̂": "qcjj",
	"U":  "qu",
	"Ü":  "quj",
	"Û":  "qujj",
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
	"ä":  "aj",
	"â":  "ajj",
	"e":  "e",
	"ë":  "ej",
	"ê":  "ejj",
	"ɛ":  "x",
	"ɛ̈": "xj",
	"ɛ̂": "xjj",
	"i":  "i",
	"ï":  "ij",
	"î":  "ijj",
	"o":  "o",
	"ö":  "oj",
	"ô":  "ojj",
	"ɔ":  "c",
	"ɔ̈": "cj",
	"ɔ̂": "cjj",
	"u":  "u",
	"ü":  "uj",
	"û":  "ujj",
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

var utf8VowelToAsciiOutput map[string]string = map[string]string{
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

var utf8ConsonantToAsciiOutput map[string]string = map[string]string{
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

var isHighPitch map[string]bool = map[string]bool{
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
	var boundaries int
	var word string
	for len(b) > 0 {
		c, b, boundaries, state = uniseg.Step(b, state)
		s := string(c)
		a, found := utf8ToAsciiInput[s]
		if found {
			word += a
		} else {
			fmt.Printf("%v", word)
			if s != "\n" {
				fmt.Printf("%v", s)
			}
			word = ""
		}
		if boundaries&uniseg.MaskSentence != 0 {
			fmt.Println("\nSENTENCE BREAK")
			if boundaries&uniseg.MaskLine == uniseg.LineMustBreak {
				fmt.Println("LINE BREAK")
			}
		} else if boundaries&uniseg.MaskLine == uniseg.LineMustBreak {
			fmt.Println("\nLINE BREAK")
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
			if consonant, found := utf8ConsonantToAsciiOutput[s]; found {
				consonant := strings.ToLower(consonant)
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
		if vowel, found := utf8VowelToAsciiOutput[s]; found {
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
			panic("\nConsonants {" + consonants + "} found not followed by a vowel {" + s + "}")
		} else {
			fmt.Printf("%s", s)
		}
	}
}

func decodeInput() {
	r := bufio.NewReader(os.Stdin)
	w := bufio.NewWriter(os.Stdout)
	o := norm.NFKC.Writer(w)
	b, err := io.ReadAll(r)
	check(err)
	_, err = o.Write(b)
	w.Flush()
	check(err)
}

func decodeOutput() {
	r := bufio.NewReader(os.Stdin)
	w := bufio.NewWriter(os.Stdout)
	o := norm.NFKC.Writer(w)
	b, err := io.ReadAll(r)
	check(err)
	_, err = o.Write(b)
	w.Flush()
	check(err)
}
