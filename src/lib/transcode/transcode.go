package transcode

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"regexp"
	"slices"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

type SSE uint16

// Valid language codes that can be passed into encodeAsciiWord are "sg", "en", and "fr".
// Any other language code will encode as raw unicode runes rather than words.
func EncodeWord(languageCode string, word []byte) (sses []SSE) {
	return encodeWord(languageCode, word)
}

func Encode(out *bufio.Writer, in *bufio.Reader) error {
	defer out.Flush()
	phrase, err := io.ReadAll(norm.NFC.Reader(in))
	if err != nil {
		return err
	}
	sses := encode(phrase)
	fmt.Fprintf(out, "There are %v tokens:\n", len(sses))
	for k, sse := range sses {
		if sse&0b1000000000000000 != 0 {
			_, err = fmt.Fprintf(out, "#%v: 0b_%01b_%02b_%02b_%02b_%05b_%04b\n", k,
				sse&0b1000000000000000>>15,
				sse&0b0110000000000000>>13,
				sse&0b0001100000000000>>11,
				sse&0b0000011000000000>>9,
				sse&0b0000000111110000>>4,
				sse&0b0000000000001111>>0)
		} else if sse&0b0100000000000000 != 0 {
			_, err = fmt.Fprintf(out, "#%v: 0b_%02b_%01b_%05b_%08b\n", k,
				sse&0b1100000000000000>>14,
				sse&0b0010000000000000>>13,
				sse&0b0001111100000000>>8,
				sse&0b0000000011111111>>0)
		} else {
			_, err = fmt.Fprintf(out, "#%v: 0b_%02b_%014b\n", k,
				sse&0b1100000000000000>>14,
				sse&0b0011111111111111>>0)
		}
		if err != nil {
			return err
		}
	}
	return err
}

func Decode(out *bufio.Writer, in *bufio.Reader) error {
	return nil
}

////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

var (
	// Matches final syllable of a Sango word in NFD form.
	sangoSyllableRE = regexp.MustCompile(`(?i)(?:([-]?)((?:n(?:[dyz]?|gb?)|m[bv]?|kp?|gb?|[bdfhlprstvwyz]?)?)(?:([aeiou\xF8\x{0254}\x{0259}\x{025B}])([\x{0302}\x{0308}\x{0323}]?)(n?)))$`)

	// Matches whole Sango word + single rune on either side (if any) in NFC form.
	sangoWordRE = regexp.MustCompile(`(?i)(?:^|[^a-z\xE2\xE4\xEA\xEB\xEE\xEF\xF4\xF6\xF8\xFB\xFC\x{0254}\x{0259}\x{025B}\x{0302}\x{0308}\x{0323}\x{1EA1}\x{1EB9}\x{1ECB}\x{1ECD}\x{1EE5}-])((?:(?:[-]?)(?:(?:n(?:[dyz]?|gb?)-|m[bv]?|kp?|gb?|[bdfhlprstvwyz]?)?)(?:(?:[\xE2\xE4\xEA\xEB\xEE\xEF\xF4\xF6\xFB\xFC\x{1EA1}\x{1EB9}\x{1ECB}\x{1ECD}\x{1EE5}]|[aeiou\xF8\x{0254}\x{0259}\x{025B}][\x{0302}\x{0308}\x{0323}]?)(?:n?)))+)(?:$|[^a-z\xE2\xE4\xEA\xEB\xEE\xEF\xF4\xF6\xF8\xFB\xFC\x{0254}\x{0259}\x{025B}\x{0302}\x{0308}\x{0323}\x{1EA1}\x{1EB9}\x{1ECB}\x{1ECD}\x{1EE5}-])`)
)

const (
	// SSE 16-BIT ENCODING
	// TYPE            = 0b_T000_0000_0000_0000
	typeIsRune     = 0b_0000_0000_0000_0000
	typeIsSyllable = 0b_1000_0000_0000_0000
	// typeOnly        = 0b_1000_0000_0000_0000

	// UNICODE RUNE    = 0b_00UU_UUUU_UUUU_UUUU where
	// UUUUUUUUUUUUUU  = Unicode rune value (U+0000 - U+3FFF)
	runeOnly      = 0b_1100_0000_0000_0000
	runeIsUnicode = 0b_0000_0000_0000_0000
	// runeIsAscii     = 0b_0100_0000_0000_0000
	unicodeValueOnly  = 0b_0011_1111_1111_1111
	unicodeValueShift = 0

	// ASCII ENGLISH   = 0b_010L_LLLL_AAAA_AAAA
	// ASCII FRENCH    = 0b_011L_LLLL_AAAA_AAAA where
	//   LLLLL         = min(31,n), n  = # letters left
	//       AAAAAAAA  = ASCII letter value (U+00 - U+FF)
	// asciiOnly       = 0b_1110_0000_0000_0000
	asciiIsEnglish   = 0b_0100_0000_0000_0000
	asciiIsFrench    = 0b_0110_0000_0000_0000
	asciiLengthOnly  = 0b_0001_1111_0000_0000
	asciiLengthShift = 8
	asciiValueOnly   = 0b_0000_0000_1111_1111
	asciiValueShift  = 0

	// SANGO           = 0b_1SSX_XPPC_CCCC_VVVV where
	// SS              = min( 3,m), m = # syllables left
	//   XX            = Case : 00=Hidden , 01=lowercase, 10=-prefixed, 11=Uppercase
	//     PP          = Pitch: 00=Unknown, 01=LowTone  , 10=MidTone  , 11=HighTone
	//       CCCCC     = Consonant (first 3 on left below, last 2 on top)
	//            VVVV = Vowel     (first 2 on left below, last 2 on top)
	// where CCCCC and VVVV are set as follows (MSB on left, LSB on top):
	// | Bit | 00 | 01 | 10 | 11 |
	// +-----+----+----+----+----+
	// | 000 |    | f  | r  | k  |
	// | 001 | mv | v  | ng | g  |
	// | 010 | m  | p  | l  | kp |
	// | 011 | mb | b  | ngb| gb |
	// | 100 |    | s  | y  | h  |
	// | 101 | nz | z  | ny | w  |
	// | 110 | n  | t  | nd | d  |
	// |  00 |    | u  | ɔ  | ɛ  |
	// |  01 | a  | i  | o  | e  |
	// |  10 |    | uñ | ø  | ə  |
	// |  11 | añ | iñ | oñ | eñ |
	sangoLengthOnly  = 0b_0110_0000_0000_0000
	sangoLengthShift = 13
	// sangoCaseOnly   = 0b_0001_1000_0000_0000
	// sangoCaseShift  = 11
	sangoCaseHidden = 0b_0000_0000_0000_0000
	sangoCaseLower  = 0b_0000_1000_0000_0000
	sangoCaseHyphen = 0b_0001_0000_0000_0000
	sangoCaseUpper  = 0b_0001_1000_0000_0000
	// sangoPitchOnly  = 0b_0000_0110_0000_0000
	// sangoPitchShift = 9
	sangoPitchUnknown = 0b_0000_0000_0000_0000
	sangoPitchLow     = 0b_0000_0010_0000_0000
	sangoPitchMid     = 0b_0000_0100_0000_0000
	sangoPitchHigh    = 0b_0000_0110_0000_0000
	// sangoConsOnly   = 0b_0000_0001_1111_0000
	// sangoConsShift  = 4
	sangoConsInvalid = 0b_0000_0001_0000_0000
	// sangoVowelOnly  = 0b_0000_0000_0000_1111
	// sangoVowelShift = 0
	sangoVowelInvalid = 0b_0000_0000_0000_1000
)

var (
	consToSSE = map[string]SSE{
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
		"":   SSE(0b_0000_0000_0000_0000),
		"u":  SSE(0b_0000_0000_0000_0001),
		"ɔ":  SSE(0b_0000_0000_0000_0010),
		"ɛ":  SSE(0b_0000_0000_0000_0011),
		"a":  SSE(0b_0000_0000_0000_0100),
		"i":  SSE(0b_0000_0000_0000_0101),
		"o":  SSE(0b_0000_0000_0000_0110),
		"e":  SSE(0b_0000_0000_0000_0111),
		"un": SSE(0b_0000_0000_0000_1001),
		"ø":  SSE(0b_0000_0000_0000_1010),
		"ə":  SSE(0b_0000_0000_0000_1011),
		"an": SSE(0b_0000_0000_0000_1100),
		"in": SSE(0b_0000_0000_0000_1101),
		"on": SSE(0b_0000_0000_0000_1110),
		"øn": SSE(0b_0000_0000_0000_1110),
		"en": SSE(0b_0000_0000_0000_1111),
		"ən": SSE(0b_0000_0000_0000_1111),
	}
)

func encodeLastSyllable(word []byte, numSyllablesLeft int) (newWord []byte, sse SSE) {
	word = norm.NFD.Bytes(word)
	if len(word) == 0 {
		return
	}
	span := sangoSyllableRE.FindSubmatchIndex(word)
	if span == nil || len(span) != 12 || span[0] < 0 || span[1] <= span[0] {
		return
	}
	sse = typeIsSyllable | (sangoLengthOnly & SSE(min(3, numSyllablesLeft)<<sangoLengthShift))
	newWord = word[:span[0]]
	syllables := bytes.Runes(word[span[0]:span[1]])
	hyphen := string(word[span[2]:span[3]])

	cons := string(bytes.ToLower(word[span[4]:span[5]]))
	vowel := string(bytes.ToLower(word[span[6]:span[7]]))
	pitch := string(bytes.ToLower(word[span[8]:span[9]]))
	nasal := string(bytes.ToLower(word[span[10]:span[11]]))

	// Asciify vowel (and convert to lower case.
	// The `syllables[0]` rune preserves the overall case.

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
	// In the standard orthography, there is no dot below.
	// On output, the dot below should be stripped out.
	case "\u0323": // dot below, e.g. ọ
		sse |= sangoPitchLow
	case "\u0308": // diaeresis above, e.g. ö
		sse |= sangoPitchMid
	case "\u0302": // circumflex above, e.g. ô
		sse |= sangoPitchHigh
	default:
		sse |= sangoPitchUnknown
	}

	if sseCons, isFound := consToSSE[cons]; isFound {
		sse |= sseCons
	} else {
		sse |= sangoConsInvalid
	}

	if sseVowel, isFound := vowelToSSE[vowel+nasal]; isFound {
		sse |= sseVowel
	} else {
		sse |= sangoVowelInvalid
	}

	return
}

// Valid language codes that can be passed into encodeAsciiWord are "sg", "en", and "fr".
// Any other language code will encode as raw unicode runes rather than words.
func encodeWord(languageCode string, word []byte) (sses []SSE) {
	sses = []SSE{}
	if languageCode == "sg" {
		// Encode Sango by syllables, which is the fundamental phonemic unit.
		for numSyllablesLeft := 0; len(word) > 0; numSyllablesLeft++ {
			var sse SSE
			word, sse = encodeLastSyllable(word, numSyllablesLeft)
			sses = append(sses, sse)
		}
		slices.Reverse(sses)
		return
	}

	// For non-Sango text, encode each rune individually.
	// For English and French words, encode the language code and running length as well.
	// Also, replace common ASCII punctuation with fancier Unicode punctuation.
	word = bytes.ReplaceAll(word, []byte("..."), []byte("…"))
	word = bytes.ReplaceAll(word, []byte("<<"), []byte("«"))
	word = bytes.ReplaceAll(word, []byte(">>"), []byte("»"))
	word = bytes.ReplaceAll(word, []byte("``"), []byte("“"))
	word = bytes.ReplaceAll(word, []byte("''"), []byte("”"))
	word = bytes.ReplaceAll(word, []byte("---"), []byte("—"))
	word = bytes.ReplaceAll(word, []byte("--"), []byte("–"))
	ascii := bytes.Runes(word)
	n := len(ascii)
	for _, r := range ascii {
		n--
		if r > 0x3fff {
			r = 0x25a1 // use white square for runes that cannot be encoded
		}
		sse := SSE(typeIsRune)
		if r >= 0 && r <= 0xff && languageCode == "en" {
			sse |= asciiIsEnglish
			sse |= SSE(asciiLengthOnly & (min(31, n) << asciiLengthShift))
			sse |= SSE(asciiValueOnly & (r << asciiValueShift))
		} else if r >= 0 && r <= 0xff && languageCode == "fr" {
			sse |= asciiIsFrench
			sse |= SSE(asciiLengthOnly & (min(31, n) << asciiLengthShift))
			sse |= SSE(asciiValueOnly & (r << asciiValueShift))
		} else {
			sse |= runeIsUnicode
			sse |= SSE(unicodeValueOnly & (r << unicodeValueShift))
		}
		sses = append(sses, sse)
	}
	return
}

func encode(phrase []byte) (sses []SSE) {
	sses = []SSE{}
	spans := [][3]int{} // [sPre, sMid, eMid]
	for s, n := 0, len(phrase); s < n; {
		span := sangoWordRE.FindSubmatchIndex(phrase[s:n])
		if span == nil {
			spans = append(spans, [3]int{s, n, n})
			break
		} else if len(span) != 4 {
			panic("Bad span")
		} else {
			spans = append(spans, [3]int{s, s + span[2], s + span[3]})
			s += span[3]
		}
	}
	for j, span := range spans {
		if len(span) != 3 {
			log.Fatalf("Bad span (%v) from spans (%v)", span, spans)
		}
		log.Printf("=================== PRE #%v ===================", j)
		if s, e := span[0], span[1]; s < e {
			log.Printf("other = phrase[%v:%v] = %q\n", s, s, string(phrase[s:e]))
			// TODO: Use ../tokenize/wordlist_{en,fr,sg,sg_toneless}.cf to assign languageCode.
			// For now, just leave blank and encode it like punctuation.
			ssesNew := EncodeWord("", phrase[s:e])
			for k, sse := range ssesNew {
				log.Printf("sse[%v] = %04x = %016b\n", k, sse, sse)
			}
			sses = append(sses, ssesNew...)
		}

		log.Printf("=================== MID #%v ===================", j)
		if s, e := span[1], span[2]; s < e {
			log.Printf("sango = phrase[%v:%v] = %q\n", s, e, string(phrase[s:e]))
			ssesNew := EncodeWord("sg", phrase[s:e])
			for k, sse := range ssesNew {
				log.Printf("sse[%v] = %04x = %016b\n", k, sse, sse)
			}
			sses = append(sses, ssesNew...)
		}
	}
	return sses
}
