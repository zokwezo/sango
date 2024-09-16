// Tokenizes a string based on regular expressions.
// Provides a specialization for Sango (allowing for code switching with other languages).

package tokenize

import (
	"regexp"
)

var SangoTokenizerRegexps = []*regexp.Regexp{
	regexp.MustCompile(`\pZ+`),                    // whitespace
	regexp.MustCompile(`\p{Nd}+(?:[.,]\p{Nd}*)*`), // numbers
	regexp.MustCompile(`\pP`),                     // punctuation
	regexp.MustCompile(`^(?:(?i)` +
		`(?:n(?:[dyz]?|gb?)|m[bv]?|kp?|gb?|[bdfhlprstvwyz]?)?` +
		`(?:(?:ä|ë|ï|ö|ü|â|ê|î|ô|û|a|e|i|o|u)n?|ɛ̂|ɛ̈|ɛ|ɔ̂|ɔ̈|ɔ))+$`), // Sango words
	regexp.MustCompile(`^\p{Latin}+$`), // words in other languages with a Latin script
} // implicitly, everything else

type Token = struct {
	begin, end int
	reIndex    int
}

func Tokenize(s *string, regexps []*regexp.Regexp) (*string, []Token) {
	if s == nil {
		return s, nil
	}
	nRE := len(regexps)
	var tokenize func(Token) []Token
	tokenize = func(candidate Token) (tokens []Token) {
		tokens = []Token{}
		b := candidate.begin
		e := candidate.end
		r := candidate.reIndex
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
