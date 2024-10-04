// Tokenizes a string based on regular expressions.
// Provides a specialization for Sango (allowing for code switching with other languages).

package tokenize

import (
	_ "embed"
	"io"
	"log"
	"regexp"
	"strings"

	cuckoo "github.com/panmari/cuckoofilter"
	"golang.org/x/text/unicode/norm"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

type Token = struct {
	Begin   int
	End     int
	REindex int
}

func TokenizeSango(in io.Reader) (*string, []Token) {
	r := norm.NFKC.Reader(in)
	b, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}
	s := string(b)
	return tokenize(&s, sangoTokenizerRegexps)
}

type Lemma struct {
	Source   Token
	Toneless string
	Sango    string
	Type     string
	Lang     string
}

func ClassifySango(in io.Reader) []Lemma {
	return classify(TokenizeSango(in))
}

//////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

var sangoTokenizerRegexps = []*regexp.Regexp{
	regexp.MustCompile(`\p{Z}+`),                  // whitespace
	regexp.MustCompile(`\p{Nd}+(?:[.,]\p{Nd}*)*`), // numbers
	regexp.MustCompile(`\.{3}|\p{P}`),             // punctuation
	regexp.MustCompile(`^(?:(?i)` +
		`(?:n(?:[dyz]?|gb?)|m[bv]?|kp?|gb?|[bdfhlprstvwyz]?)?` +
		`(?:(?:ä|ë|ï|ö|ü|â|ê|î|ô|û|a|e|i|o|u)n?|ɛ̂|ɛ̈|ɛ|ɔ̂|ɔ̈|ɔ))+$`), // Sango
	regexp.MustCompile(`^\p{Latin}+$`), // English/French
} // everything else

//go:embed wordlist_en.cf
var enWordListEncodedCuckooFilter []byte

//go:embed wordlist_fr.cf
var frWordListEncodedCuckooFilter []byte

//go:embed wordlist_sg.cf
var sgWordListEncodedCuckooFilter []byte

//go:embed wordlist_sg_toneless.cf
var sgTonelessWordListEncodedCuckooFilter []byte

var enWords = getWordListFromEncodedCuckooFilter(enWordListEncodedCuckooFilter)
var frWords = getWordListFromEncodedCuckooFilter(frWordListEncodedCuckooFilter)
var sgWords = getWordListFromEncodedCuckooFilter(sgWordListEncodedCuckooFilter)
var sgTonelessWords = getWordListFromEncodedCuckooFilter(sgTonelessWordListEncodedCuckooFilter)

func getWordListFromEncodedCuckooFilter(b []byte) *cuckoo.Filter {
	cf, err := cuckoo.Decode(b)
	if err != nil {
		panic(err)
	}
	return cf
}

var toLowPitch = func() func(r rune) rune {
	m := map[rune]rune{
		'ä': 'a', 'â': 'a', 'ë': 'e', 'ê': 'e', 'ï': 'i',
		'î': 'i', 'ö': 'o', 'ô': 'o', 'ü': 'u', 'û': 'u'}
	return func(r rune) rune {
		if c, found := m[r]; found {
			return c
		}
		return r
	}
}()

func classify(s *string, tokens []Token) []Lemma {
	Pi := regexp.MustCompile(`\p{Pi}`)
	Pf := regexp.MustCompile(`\p{Pf}`)
	Ps := regexp.MustCompile(`\p{Ps}`)
	Pe := regexp.MustCompile(`\p{Pe}`)
	Pd := regexp.MustCompile(`\p{Pd}`)
	Pc := regexp.MustCompile(`\p{Pc}`)
	Po := regexp.MustCompile(`\p{Po}`)
	P := regexp.MustCompile(`\p{P}`)

	if s == nil || tokens == nil {
		return nil
	}
	lemmas := []Lemma{}
	for _, token := range tokens {
		b := token.Begin
		e := token.End
		r := token.REindex
		w := (*s)[b:e]
		wLC := strings.ToLower(w)
		wToneless := wLC
		for _, p := range [...][2]string{{"ɛ̂", "e"}, {"ɛ̈", "e"}, {"ɛ", "e"}, {"ɔ̂", "o"}, {"ɔ̈", "o"}, {"ɔ", "o"}} {
			wToneless = strings.ReplaceAll(wToneless, p[0], p[1])
		}
		wToneless = strings.Map(toLowPitch, wToneless)
		l := ""
		t := "OTHER"
		switch r {
		case 0:
			t = "SPACE"
		case 1:
			t = "NUM"
		case 2:
			t = "PUNC"
			if w == "..." {
				wLC = "…"
				wToneless = "…"
			}
			if Pi.MatchString(wLC) || Ps.MatchString(wLC) {
				l = "open"
			} else if Pf.MatchString(wLC) || Pe.MatchString(wLC) {
				l = "close"
			} else if Pd.MatchString(wLC) {
				l = "dash"
			} else if Pc.MatchString(wLC) {
				l = "connector"
			} else if Po.MatchString(wLC) {
				l = "other"
			} else if P.MatchString(wLC) {
				l = "punc"
			}
		case 3:
			t = "WORD"
			if sgWords.Lookup([]byte(wLC)) {
				l = "sg"
				// TODO: Create and use sgHeightlessWords.Lookup to automatically keep
				// extant pitch accent when restoring vowel height. This is the most
				// common case since the standard orthography does not encode height
				// but does encode pitch.
				if sgTonelessWords.Lookup([]byte(wLC)) {
					// There is ambiguity when matching a Sango word with no pitch or height accents.
					// This word (possibly by coincidence) matches a lexicon entry, but absence of
					// accent marks does not prove there shouldn't be any. It would be much clearer
					// if low pitch and close vowels (e and o) had explicit markings, or if conversely
					// unaccented words were marked explicitly as unknown.
					// Indicate this ambiguity with a mixed-case language code.
					l = "Sg"
				}
				break
			}
			if sgTonelessWords.Lookup([]byte(wToneless)) {
				// Lexeme does not match as is, but would with different accents.
				// Indicate this with uppercase language code.
				l = "SG"
				break
			}
			fallthrough
		case 4:
			t = "WORD"
			if frWords.Lookup([]byte(wLC)) {
				l = "fr"
			} else if enWords.Lookup([]byte(wLC)) {
				l = "en"
			} else if l == "" {
				l = "XX"
			}
		}
		lemmas = append(lemmas, Lemma{token, wToneless, wLC, t, l})
	}
	return lemmas
}

func tokenize(s *string, regexps []*regexp.Regexp) (*string, []Token) {
	if s == nil {
		return s, nil
	}
	nRE := len(regexps)
	var tokenize func(Token) []Token
	tokenize = func(candidate Token) (tokens []Token) {
		tokens = []Token{}
		b := candidate.Begin
		e := candidate.End
		r := candidate.REindex
		if b >= e {
			return tokens
		}
		for ; r < nRE && regexps[r] == nil; r++ {
		}
		if r >= nRE {
			tokens = append(tokens, Token{b, e, nRE})
			return tokens
		}
		var f int = max(0, e-b)
		spans := regexps[r].FindAllStringSubmatchIndex((*s)[b:e], -1)
		spans = append(spans, []int{f, f}) // harmless, but makes logic easier
		aa := b                            // start of nonmatching span
		for _, span := range spans {
			bb := b + span[0] // end of nonmatching span = start of matching span
			ee := b + span[1] // end of matching span
			if aa >= bb {
				aa = max(bb, ee)
			} else {
				token := Token{aa, bb, r + 1}
				aa = max(bb, ee)
				subtokens := tokenize(token)
				tokens = append(tokens, subtokens...)
			}
			if bb < ee {
				tokens = append(tokens, Token{bb, ee, r})
			}
		}
		return tokens
	}
	return s, tokenize(Token{0, len(*s), 0})
}
