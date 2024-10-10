package transcode

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"regexp"
	"slices"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

func Encode(out *bufio.Writer, in *bufio.Reader) error {
	defer out.Flush()
	encode(in)
	return nil
}
func Decode(out *bufio.Writer, in *bufio.Reader) error {
	return nil
}

////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

// The `\b` word boundary matches the middle of
// `(?=[0-9A-Za-z_])()(?=[^0-9A-Za-z_]|\z)` which unfortunately also
// matches an apostrophe, and look ahead/behind is not supported in
// RE2 syntax, so it cannot be added into the above character classes.
// This means that e.g. the French "c'est" is mistakenly split into a
// Sango word "c" and and separate ASCII word "est".
// Instead, use the logic in ../tokenize/tokenize.go with recursive
// regexps to better isolate components.
var (
	sangoREPattern = `(?i)(?:([-]?)((?:n(?:[dyz]?|gb?)|m[bv]?|kp?|gb?|[bdfhlprstvwyz]?)?)` +
		`(?:([aeiouxc\x{0186}\x{0190}\x{0254}\x{025B}])([\x60jq\x{0302}\x{0306}\x{0308}]?)(n?)))`
	sangoSyllableRE = regexp.MustCompile(sangoREPattern + `$`)
	sangoWordRE     = regexp.MustCompile(`\b` + sangoREPattern + `+\b`)
)

type SSE uint16

const (
	// payloadOnly       = 0b_1000_0000_0000_0000
	payloadIsRune     = 0b_0000_0000_0000_0000
	payloadIsSyllable = 0b_1000_0000_0000_0000

	runeOnly      = 0b_1100_0000_0000_0000
	runeIsUnicode = 0b_0000_0000_0000_0000
	// runeIsAscii   = 0b_0100_0000_0000_0000

	unicodeValueOnly  = 0b_0011_1111_1111_1111
	unicodeValueShift = 0

	// asciiOnly      = 0b_1110_0000_0000_0000
	asciiIsEnglish = 0b_0100_0000_0000_0000
	// asciiIsFrench  = 0b_0110_0000_0000_0000

	asciiLengthOnly  = 0b_0001_1111_0000_0000
	asciiLengthShift = 8
	asciiValueOnly   = 0b_0000_0000_1111_1111
	asciiValueShift  = 0

	sangoLengthOnly  = 0b_0110_0000_0000_0000
	sangoLengthShift = 13

	// sangoCaseOnly   = 0b_0001_1000_0000_0000
	// sangoCaseShift  = 11
	sangoCaseHidden = 0b_0000_0000_0000_0000
	sangoCaseLower  = 0b_0000_1000_0000_0000
	sangoCaseHyphen = 0b_0001_0000_0000_0000
	sangoCaseUpper  = 0b_0001_1000_0000_0000

	// sangoPitchOnly    = 0b_0000_0110_0000_0000
	// sangoPitchShift   = 9
	sangoPitchUnknown = 0b_0000_0000_0000_0000
	sangoPitchLow     = 0b_0000_0010_0000_0000
	sangoPitchMid     = 0b_0000_0100_0000_0000
	sangoPitchHigh    = 0b_0000_0110_0000_0000

	// sangoConsonantOnly    = 0b_0000_0001_1111_0000
	// sangoConsonantShift   = 4
	sangoConsonantInvalid = 0b_0000_0001_0000_0000

	// sangoVowelOnly    = 0b_0000_0000_0000_1111
	// sangoVowelShift   = 0
	sangoVowelInvalid = 0b_0000_0000_0000_1000
)

var (
	consonantToSSE = map[string]SSE{
		"":    SSE(0b_0000_0000_0000_0000),
		"f":   SSE(0b_0000_0000_0001_0000),
		"r":   SSE(0b_0000_0000_0010_0000),
		"k":   SSE(0b_0000_0000_0011_0000),
		"mv":  SSE(0b_0000_0000_0100_0000),
		"v":   SSE(0b_0000_0000_0101_0000),
		"ng":  SSE(0b_0000_0000_0110_0000),
		"g":   SSE(0b_0000_0000_0111_0000),
		"m":   SSE(0b_0000_0000_1000_0000),
		"p":   SSE(0b_0000_0000_1001_0000),
		"l":   SSE(0b_0000_0000_1010_0000),
		"kp":  SSE(0b_0000_0000_1011_0000),
		"mb":  SSE(0b_0000_0000_1100_0000),
		"b":   SSE(0b_0000_0000_1101_0000),
		"ngb": SSE(0b_0000_0000_1110_0000),
		"gb":  SSE(0b_0000_0000_1111_0000),
		"s":   SSE(0b_0000_0001_0001_0000),
		"y":   SSE(0b_0000_0001_0010_0000),
		"h":   SSE(0b_0000_0001_0011_0000),
		"nz":  SSE(0b_0000_0001_0100_0000),
		"z":   SSE(0b_0000_0001_0101_0000),
		"ny":  SSE(0b_0000_0001_0110_0000),
		"w":   SSE(0b_0000_0001_0111_0000),
		"n":   SSE(0b_0000_0001_1000_0000),
		"t":   SSE(0b_0000_0001_1001_0000),
		"nd":  SSE(0b_0000_0001_1010_0000),
		"d":   SSE(0b_0000_0001_1011_0000),
	}

	vowelToSSE = map[string]SSE{
		"":    SSE(0b_0000_0000_0000_0000),
		"u":   SSE(0b_0000_0000_0000_0001),
		"c":   SSE(0b_0000_0000_0000_0010),
		"x":   SSE(0b_0000_0000_0000_0011),
		"a":   SSE(0b_0000_0000_0000_0100),
		"i":   SSE(0b_0000_0000_0000_0101),
		"o":   SSE(0b_0000_0000_0000_0110),
		"e":   SSE(0b_0000_0000_0000_0111),
		"un":  SSE(0b_0000_0000_0000_1001),
		"co":  SSE(0b_0000_0000_0000_1010),
		"xe":  SSE(0b_0000_0000_0000_1011),
		"an":  SSE(0b_0000_0000_0000_1100),
		"in":  SSE(0b_0000_0000_0000_1101),
		"on":  SSE(0b_0000_0000_0000_1110),
		"con": SSE(0b_0000_0000_0000_1110),
		"en":  SSE(0b_0000_0000_0000_1111),
		"xen": SSE(0b_0000_0000_0000_1111),
	}
)

func encodeLastSyllable(word []byte, numSyllablesLeft int) (newWord []byte, sse SSE) {
	if len(word) == 0 {
		return
	}
	span := sangoSyllableRE.FindSubmatchIndex(word)
	if span == nil || len(span) != 12 || span[0] < 0 || span[1] <= span[0] {
		return
	}
	sse = payloadIsSyllable | (sangoLengthOnly & SSE(min(3, numSyllablesLeft)<<sangoLengthShift))
	newWord = word[:span[0]]
	syllables := bytes.Runes(word[span[0]:span[1]])
	hyphen := string(word[span[2]:span[3]])

	consonant := string(bytes.ToLower(word[span[4]:span[5]]))
	vowel := string(bytes.ToLower(word[span[6]:span[7]]))
	pitch := string(bytes.ToLower(word[span[8]:span[9]]))
	nasal := string(bytes.ToLower(word[span[10]:span[11]]))

	// Asciify pitch
	if pitch == "" {
		pitch = "q" // unknown
	} else {
		pitch = strings.ReplaceAll(pitch, "\u0300", "j") // low
		pitch = strings.ReplaceAll(pitch, "\u0302", "J") // high
		pitch = strings.ReplaceAll(pitch, "\u0306", "q") // unknown
		pitch = strings.ReplaceAll(pitch, "\u0308", "Q") // mid
	}

	// Asciify vowel (and convert to lower case.
	// The `syllables[0]` rune preserves the overall case.
	vowel = strings.ReplaceAll(vowel, "Ɛ", "x")
	vowel = strings.ReplaceAll(vowel, "Ɔ", "c")
	vowel = strings.ReplaceAll(vowel, "ɛ", "x")
	vowel = strings.ReplaceAll(vowel, "ɔ", "c")

	// If pitch is unknown, assume that height is also unknown.
	if pitch == "q" {
		switch vowel {
		case "e":
			vowel = "xe"
		case "o":
			vowel = "co"
		}
	}

	if len(syllables) == 0 {
		return
	}
	if hyphen == "-" {
		sse |= sangoCaseHyphen
	} else if hyphen != "" {
		sse |= sangoCaseHidden
	} else if unicode.IsUpper(syllables[0]) {
		sse |= sangoCaseUpper
	} else {
		sse |= sangoCaseLower
	}

	switch pitch {
	case "q":
		sse |= sangoPitchUnknown
	case "j":
		sse |= sangoPitchLow
	case "Q":
		sse |= sangoPitchMid
	case "J":
		sse |= sangoPitchHigh
	default:
		panic("Bad pitch")
	}

	if sseConsonant, isFound := consonantToSSE[consonant]; isFound {
		sse |= sseConsonant
	} else {
		sse |= sangoConsonantInvalid
	}

	if sseVowel, isFound := vowelToSSE[vowel+nasal]; isFound {
		sse |= sseVowel
	} else {
		sse |= sangoVowelInvalid
	}

	return
}

func encodeSangoWord(word []byte) (sses []SSE) {
	if sses != nil {
		panic("sses starts out nonnil")
	}
	sses = []SSE{}
	for numSyllablesLeft := 0; len(word) > 0; numSyllablesLeft++ {
		var sse SSE
		word, sse = encodeLastSyllable(word, numSyllablesLeft)
		sses = append(sses, sse)
	}
	slices.Reverse(sses)
	return
}

func encode(in io.Reader) (sses []SSE) {
	if sses != nil {
		panic("sses starts out nonnil")
	}
	phrase, err := io.ReadAll(norm.NFKD.Reader(in))
	if err != nil {
		return
	}
	sses = []SSE{}
	ePrev := 0
	sangoSpans := sangoWordRE.FindAllIndex(phrase, -1)
	for j, wordSpan := range sangoSpans {
		log.Println("=================== PRE ===================")
		if len(wordSpan) != 2 {
			panic("Bad wordSpan")
		}
		s, e := wordSpan[0], wordSpan[1]
		log.Printf("word[%v] = phrase[%v:%v] = %q\n", j, ePrev, s, string(phrase[ePrev:s]))
		for k, c := range string(phrase[ePrev:s]) {
			log.Printf("  byte[%v] = '%c' = %02x\n", k, c, c)
		}
		ascii := bytes.Runes(phrase[ePrev:s])
		n := len(ascii)
		for k, r := range ascii {
			numAsciisLeft := (n - 1) - k
			if r > 0x3fff {
				r = 0x25a1 // use white square for runes that cannot be encoded
			}
			sse := SSE(payloadIsRune)
			if unicode.IsLetter(r) {
				sse |= asciiIsEnglish
				sse |= SSE(asciiLengthOnly & (min(31, numAsciisLeft) << asciiLengthShift))
				sse |= SSE(asciiValueOnly & (r << asciiValueShift))
			} else {
				sse |= runeIsUnicode
				sse |= SSE(unicodeValueOnly & (r << unicodeValueShift))
			}
			sses = append(sses, sse)
		}
		log.Printf("sses[%v] = %v\n", j, sses)
		ePrev = e
		log.Println("=================== MID ===================")
		// Process Sango word
		word := phrase[s:e]
		log.Printf("word[%v] = phrase[%v:%v] = %q\n", j, s, e, string(word))
		sses = append(sses, encodeSangoWord(word)...)
		for k, sse := range sses {
			log.Printf("sse[%v] = %04x = %016b\n", k, sse, sse)
		}
	}
	log.Println("=================== POST ==================")
	log.Printf("word[final] = phrase[%v:%v] = %q\n", ePrev, len(phrase), string(phrase[ePrev:]))
	ascii := bytes.Runes(phrase[ePrev:])
	n := len(ascii)
	for k, r := range ascii {
		numAsciisLeft := (n - 1) - k
		if r > 0x3fff {
			r = 0x25a1 // use white square for runes that cannot be encoded
		}
		sse := SSE(payloadIsRune)
		if unicode.IsLetter(r) {
			sse |= asciiIsEnglish
			sse |= SSE(asciiLengthOnly & (min(31, numAsciisLeft) << asciiLengthShift))
			sse |= SSE(asciiValueOnly & (r << asciiValueShift))
		} else {
			sse |= runeIsUnicode
			sse |= SSE(unicodeValueOnly & (r << unicodeValueShift))
		}
		sses = append(sses, sse)
	}
	log.Printf("sses[final] = %v\n", sses)
	for k, sse := range sses {
		log.Printf("sse[%v] = %04x = %016b\n", k, sse, sse)
	}
	return sses
}
