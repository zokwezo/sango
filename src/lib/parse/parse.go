// Implements parsing of Sango text constructing a top-down left-to-right parse tree using regexp.

package parse

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dlclark/regexp2/v2"
	"github.com/dlclark/regexp2/v2/compat"
	"golang.org/x/text/width"
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

//////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

var parseMatcherREs = func() [NumParseMatchers]*regexp2.Regexp {
	const (
		canonicalSyllablePattern = `(?<syllable>(?<shift>[~=#]?)(?<consonant>[hHbBqQdDfgGklrmnpKPstvVwyYzZ])(?<vowelNasal>[aAeEiIoOxcuUXC])(?<pitch>[_:^]?))`
		sangoSyllablePattern     = `(?<syllable>(?i)(?<consonant>(?i)b|d|f|gb|g|h|kp|k|l|mb|mp|mv|m|nd|ngb|ng|ny|nz|n|p|r|s|t|v|w|y|z|)(?<vowelNasal>(?<openVowel>ɛ̈|ɛ̂|ɛ̣|ɛ|ɔ̈|ɔ̂|ɔ̣|ɔ|ẍ|x̂|x̣|x|c̈|ĉ|c̣|c)|(?<closeVowelNasal>(?<closeVowel>ä|â|ạ|a|ë|ê|ẹ|e|ï|î|i|ị|ö|ô|ọ|o|ü|û|ụ|u)(?<nasal>(?:ñ|n)(?!(?i)(?:d|g|gb|y|z)?(?:ɛ̈|ɛ̂|ɛ̣|ɛ|ɔ̈|ɔ̂|ɔ̣|ɔ|ẍ|x̂|x̣|x|c̈|ĉ|c̣|c|ä|â|ạ|a|ë|ê|ẹ|e|ï|î|i|ị|ö|ô|ọ|o|ü|û|ụ|u))|))))`
		canonicalWordPattern     = `(?<word>(?<![\p{L}\p{M}_:^~])(?:(?:(?:[~=#]?)(?:[hHbBqQdDfgGklrmnpKPstvVwyYzZ])(?:[aAeEiIoOxcuUXC])(?:[_:^]?)))+(?![\p{L}\p{M}_:^~]))`
		sangoWordPattern         = `(?<word>(?<![\p{L}\p{M}])(?:(?:(?i)(?:(?i)b|d|f|gb|g|h|kp|k|l|mb|mp|mv|m|nd|ngb|ng|ny|nz|n|p|r|s|t|v|w|y|z|)(?:(?:ɛ̈|ɛ̂|ɛ̣|ɛ|ɔ̈|ɔ̂|ɔ̣|ɔ|ẍ|x̂|x̣|x|c̈|ĉ|c̣|c)|(?:(?:ä|â|ạ|a|ë|ê|ẹ|e|ï|î|i|ị|ö|ô|ọ|o|ü|û|ụ|u)(?:ñ|n|)))))+(?![\p{L}\p{M}]))`
		sentencePattern          = `(?<sentence>.*?(?:[.!?](?:’|”|›|»)*[ \n]+|\n\n+|$))`
	)

	return [NumParseMatchers]*regexp2.Regexp{
		regexp2.MustCompile(sentencePattern, regexp2.RE2, regexp2.OptionMaintainCaptureOrder()),
		regexp2.MustCompile(sangoWordPattern, regexp2.RE2, regexp2.OptionMaintainCaptureOrder()),
		regexp2.MustCompile(sangoSyllablePattern, regexp2.RE2, regexp2.OptionMaintainCaptureOrder()),
		regexp2.MustCompile(canonicalWordPattern, regexp2.RE2, regexp2.OptionMaintainCaptureOrder()),
		regexp2.MustCompile(canonicalSyllablePattern, regexp2.RE2, regexp2.OptionMaintainCaptureOrder()),
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
	type IndexFor struct {
		Name  string
		Index int
	}
	indexFor := []IndexFor{}
	for _, name := range re.GetGroupNames() {
		// Annoyingly, regexp2.Regexp.GetGroupNames names unnamed capture groups
		// as strconv.Itoa(i) for the ith capture group. Filter these out.
		if i := re.GroupNumberFromName(name); i >= 0 {
			if number := strconv.Itoa(i); name != number {
				indexFor = append(indexFor, IndexFor{Name: name, Index: i})
			}
		}
	}
	b := inputParseTree.Begin
	e := inputParseTree.End
	inputParseTree.Subtrees = [][]ParseTree{}
	indexes := compat.Wrap(re).FindAllStringSubmatchIndex(s[b:e], -1)
	for i, index := range indexes {
		for j := range index {
			indexes[i][j] += b
		}
	}
	for _, index := range indexes {
		c := index[0]
		if parseMatchers[0] == SentenceFinal {
			// Move sentence start back to the end of the previous sentence.
			delta := c - b
			for _, i := range indexFor {
				j := i.Index
				index[2*j] -= delta
			}
			c = b
		}
		if c > b {
			parseSubtrees := []ParseTree{ParseTree{Begin: b, End: c, Name: "unmatched"}}
			inputParseTree.Subtrees = append(inputParseTree.Subtrees, parseSubtrees)
		}
		d := index[1]
		if d > c {
			parseSubtrees := []ParseTree{}
			for _, i := range indexFor {
				name, j := i.Name, i.Index
				if x, y := index[2*j], index[2*j+1]; x >= 0 && y >= 0 {
					childParseTree := ParseTree{Begin: x, End: y, Name: name}
					parseSubtrees = append(parseSubtrees, parse(parseMatchers[1:], s, childParseTree))
				}
			}
			inputParseTree.Subtrees = append(inputParseTree.Subtrees, parseSubtrees)
		}
		b = d
	}
	if e > b {
		parseSubtrees := []ParseTree{ParseTree{Begin: b, End: e, Name: "unmatched"}}
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

func annotate(parseTree ParseTree, s string, alignSubstrings bool, namesToShow map[string]bool) strings.Builder {
	o := strings.Builder{}
	o.WriteString("\n")
	var f func(ParseTree)
	f = func(c ParseTree) {
		if 0 <= c.Begin && c.Begin < c.End && c.End <= len(s) {
			n := c.Name
			if namesToShow[n] {
				if alignSubstrings {
					o.WriteString(fmt.Sprintf("%-20v [%3v:%3v] = %v%v\n", n, c.Begin, c.End,
						strings.Repeat(" ", getStringWidth(s[0:c.Begin])), s[c.Begin:c.End]))
				} else {
					o.WriteString(fmt.Sprintf("%-20v [%3v:%3v] = %v\n", n, c.Begin, c.End, s[c.Begin:c.End]))
				}
			}
			for _, parseSubtree := range c.Subtrees {
				for _, parseTree := range parseSubtree {
					f(parseTree)
				}
			}
		}
	}
	f(parseTree)
	return o
}
