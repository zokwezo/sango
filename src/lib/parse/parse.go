// Implements parsing of Sango text constructing a top-down left-to-right parse tree using regexp.

package parse

import (
	"cmp"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"golang.org/x/text/width"
)

type ParseTree struct {
	Begin    int
	End      int
	Subtrees []map[string]ParseTree
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
func (parseTree ParseTree) Annotate(s string, alignSubstrings bool, namesNotToShow map[string]bool) string {
	return annotate(parseTree, s, alignSubstrings, namesNotToShow)
}

//////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

var parseMatcherREs = func() [NumParseMatchers]*regexp.Regexp {
	const (
		canonicalSyllablePattern = `(?<syllable>(?<shift>[~=#]?)(?<consonant>[hHbBqQdDfgGklrmnpKPstvVwyYzZ])(?<vowelNasal>[aAeEiIoOxcuUXC])(?<pitch>[_:^]?))`
		sangoSyllablePattern     = `(?<syllable>(?i)(?<consonant>(?i)b|d|f|gb|g|h|kp|k|l|mb|mp|mv|m|nd|ngb|ng|ny|nz|n|p|r|s|t|v|w|y|z|)(?<vowelNasal>(?<openVowel>ɛ̈|ɛ̂|ɛ̣|ɛ|ɔ̈|ɔ̂|ɔ̣|ɔ|ẍ|x̂|x̣|x|c̈|ĉ|c̣|c)|(?<closeVowelNasal>(?<closeVowel>ä|â|ạ|a|ë|ê|ẹ|e|ï|î|i|ị|ö|ô|ọ|o|ü|û|ụ|u)(?<nasal>ñ|n|))))`
		canonicalWordPattern     = `(?<word>(?:(?:(?:[~=#]?)(?:[hHbBqQdDfgGklrmnpKPstvVwyYzZ])(?:[aAeEiIoOxcuUXC])(?:[_:^]?)))+)`
		sangoWordPattern         = `(?<word>(?:(?:(?i)(?:(?i)b|d|f|gb|g|h|kp|k|l|mb|mp|mv|m|nd|ngb|ng|ny|nz|n|p|r|s|t|v|w|y|z|)(?:(?:ɛ̈|ɛ̂|ɛ̣|ɛ|ɔ̈|ɔ̂|ɔ̣|ɔ|ẍ|x̂|x̣|x|c̈|ĉ|c̣|c)|(?:(?:ä|â|ạ|a|ë|ê|ẹ|e|ï|î|i|ị|ö|ô|ọ|o|ü|û|ụ|u)(?:ñ|n|)))))+)`
		sentencePattern          = `(?<sentence>.*?(?:[.!?](?:’|”|›|»)*[ \n]+|\n\n+|$))`
	)

	return [NumParseMatchers]*regexp.Regexp{
		regexp.MustCompile(sentencePattern),
		regexp.MustCompile(sangoWordPattern),
		regexp.MustCompile(sangoSyllablePattern),
		regexp.MustCompile(canonicalWordPattern),
		regexp.MustCompile(canonicalSyllablePattern),
	}
}()

func parse(parseMatchers []ParseMatcher, s string, inputParseTree ParseTree) ParseTree {
	if inputParseTree.Begin < 0 || inputParseTree.Begin > inputParseTree.End || inputParseTree.End > len(s) {
		panic("Bad inputParseTree")
	}
	if inputParseTree.Begin == inputParseTree.End {
		return inputParseTree
	}
	if len(parseMatchers) == 0 {
		return inputParseTree
	}
	re := parseMatcherREs[parseMatchers[0]]
	indexFor := map[string]int{}
	for _, name := range re.SubexpNames() {
		if re.SubexpIndex(name) > 0 {
			if i := re.SubexpIndex(name); i >= 0 {
				indexFor[name] = i
			}
		}
	}
	b := inputParseTree.Begin
	e := inputParseTree.End
	inputParseTree.Subtrees = []map[string]ParseTree{}
	indexes := re.FindAllStringSubmatchIndex(s[b:e], -1)
	for i, index := range indexes {
		for j := range index {
			indexes[i][j] += b
		}
	}
	for _, index := range indexes {
		c := index[0]
		if c > b {
			parseSubtrees := map[string]ParseTree{}
			parseSubtrees["unmatched"] = ParseTree{Begin: b, End: c}
			inputParseTree.Subtrees = append(inputParseTree.Subtrees, parseSubtrees)
		}
		d := index[1]
		if d > c {
			parseSubtrees := map[string]ParseTree{}
			for name, j := range indexFor {
				if x, y := index[2*j], index[2*j+1]; x >= 0 && y >= 0 {
					childParseTree := ParseTree{Begin: x, End: y}
					parseSubtrees[name] = parse(parseMatchers[1:], s, childParseTree)
				}
			}
			inputParseTree.Subtrees = append(inputParseTree.Subtrees, parseSubtrees)
		}
		b = d
	}
	if e > b {
		parseSubtrees := map[string]ParseTree{}
		parseSubtrees["unmatched"] = ParseTree{Begin: b, End: e}
		inputParseTree.Subtrees = append(inputParseTree.Subtrees, parseSubtrees)
	}
	return inputParseTree
}

// getStringWidth calculates how many terminal columns a string occupies.
func getStringWidth(s string) int {
	totalWidth := 0
	for _, r := range s {
		p := width.LookupRune(r)
		switch p.Kind() {
		case width.EastAsianWide, width.EastAsianFullwidth:
			totalWidth += 2
		case width.EastAsianNarrow, width.EastAsianHalfwidth, width.Neutral:
			totalWidth += 1
		// Ambiguous or zero-width characters (like combining accents) add 0 or 1
		default:
			totalWidth += 1
		}
	}
	return totalWidth
}

func annotate(parseTree ParseTree, s string, alignSubstrings bool, namesNotToShow map[string]bool) string {
	type Row struct {
		begin int
		end   int
		value string
	}
	rowCompare := func(l, r Row) int {
		if c := cmp.Compare(l.begin, r.begin); c != 0 {
			return c
		}
		if c := -cmp.Compare(l.end, r.end); c != 0 {
			return c
		}
		return cmp.Compare(l.value, r.value)
	}
	rows := []Row{}
	var f func(string, ParseTree)
	f = func(n string, c ParseTree) {
		if 0 <= c.Begin && c.Begin < c.End && c.End <= len(s) {
			if !namesNotToShow[n] {
				rowValue := ""
				if alignSubstrings {
					rowValue = fmt.Sprintf("%-20v [%3v:%3v] = %v%v\n", n, c.Begin, c.End,
						strings.Repeat(" ", getStringWidth(s[0:c.Begin])), s[c.Begin:c.End])
				} else {
					rowValue = fmt.Sprintf("%-20v [%3v:%3v] = %v\n", n, c.Begin, c.End, s[c.Begin:c.End])
				}
				rows = append(rows, Row{begin: c.Begin, end: c.End, value: rowValue})
			}
			for _, parseSubtree := range c.Subtrees {
				for name, d := range parseSubtree {
					f(name, d)
				}
			}
		}
	}
	f("whole", parseTree)
	slices.SortStableFunc(rows, rowCompare)
	var o strings.Builder
	o.WriteString("\n")
	for _, row := range rows {
		o.WriteString(row.value)
	}
	return o.String()
}
