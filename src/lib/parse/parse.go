// Implements parsing of Sango text constructing a top-down left-to-right parse tree using regexp.

package parse

import (
	"strings"
)

type ParseTree struct {
	Begin    int
	End      int
	Name     string
	Subtrees [][]ParseTree // [#][named]
}

type ParseMatcher int

const (
	SentenceFinal ParseMatcher = iota
	SangoWord
	SangoSyllable
	CanonicalWord
	CanonicalSyllable
	NumParseMatchers
)

func (pm ParseMatcher) String() string {
	return parseMatcherREs[pm].String()
}

// Creates a top-down left-to-right parse tree.
// The inputParseTree defines the substring of interest, e.g.
// inputParseTree = ParseTree{Begin: 0, End: len(s)}
func Parse(parseMatchers []ParseMatcher, s string, inputParseTree ParseTree) ParseTree {
	return parse(parseMatchers, s, inputParseTree)
}

// Creates a human-readable version of a parse tree.
func (parseTree ParseTree) Annotate(s string, alignSubstrings bool, namesToShow map[string]bool) strings.Builder {
	return annotate(parseTree, s, alignSubstrings, namesToShow)
}
