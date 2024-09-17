// Tokenizes a string based on regular expressions.
// Provides a specialization for Sango (allowing for code switching with other languages).

package tokenize

import (
	"regexp"
	"strings"
)

var SangoTokenizerRegexps = []*regexp.Regexp{
	regexp.MustCompile(`(?:\p{Zl}|\p{Zp}|\p{Zs}|\p{Z}|\s)+`),                     // sep/whitespace
	regexp.MustCompile(`\p{Nd}+(?:[.,]\p{Nd}*)*`),                                // numbers
	regexp.MustCompile(`\p{Pi}|\p{Pf}|\p{Ps}|\p{Pe}|\p{Pd}|\p{Pc}|\p{Po}|\p{P}`), // punctuation
	regexp.MustCompile(`^(?:(?i)` +
		`(?:n(?:[dyz]?|gb?)|m[bv]?|kp?|gb?|[bdfhlprstvwyz]?)?` +
		`(?:(?:ä|ë|ï|ö|ü|â|ê|î|ô|û|a|e|i|o|u)n?|ɛ̂|ɛ̈|ɛ|ɔ̂|ɔ̈|ɔ))+$`), // Sango
	regexp.MustCompile(`^\p{Latin}+$`), // English/French
} // everything else

type Token = struct {
	Begin   int
	End     int
	REindex int
}

type Lemma struct {
	Word string
	Type string
	Lang string
}

func ClassifySango(s *string) []Lemma {
	return classify(TokenizeSango(s))
}

func TokenizeSango(s *string) (*string, []Token) {
	return tokenize(s, SangoTokenizerRegexps)
}

//////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

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
				w = "…"
				wLC = "…"
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
			if _, isSg := sgWords[wLC]; isSg {
				l = "sg"
				break
			}
			fallthrough
		case 4:
			t = "WORD"
			if _, isFr := frWords[wLC]; isFr {
				l = "fr"
			} else if _, isEn := enWords[wLC]; isEn {
				l = "en"
			} else {
				l = "XX"
			}
		}
		lemmas = append(lemmas, Lemma{w, t, l})
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
		if spans == nil {
			tokens = append(tokens, Token{b, e, nRE})
			return tokens
		}
		aa := b // start of nonmatching span
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
