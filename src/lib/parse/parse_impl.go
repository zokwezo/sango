// Implements parsing of Sango text constructing a top-down left-to-right parse tree using regexp.

package parse

import (
	"fmt"
	"strings"

	"github.com/dlclark/regexp2/v2"
	"github.com/dlclark/regexp2/v2/compat"
	"golang.org/x/text/width"
)

var parseMatcherREs = func() [NumParseMatchers]*regexp2.Regexp {
	const (
		canonicalSyllablePattern = `(?<syllable>(?<consonant>[hHbBqQdDfgGklrmnpKPstvVwyYzZ])(?<vowelNasal>[aAeEiIoOxcuUXC])(?<pitch>[%_:^]?))`
		sangoSyllablePattern     = `(?i)(?<syllable>(?<consonant>b|d|f|gb|g|h|kp|k|l|mb|mp|mv|m|nd|ngb|ng|ny|nz|n|p|r|s|t|v|w|y|z|)(?<vowelNasal>(?<openVowel>ɛ̈|ɛ̂|ɛ̣|ɛ|ɔ̈|ɔ̂|ɔ̣|ɔ|ẍ|x̂|x̣|x|c̈|ĉ|c̣|c)|(?<closeVowelNasal>(?<closeVowel>ä|â|ạ|a|ë|ê|ẹ|e|ï|î|i|ị|ö|ô|ọ|o|ü|û|ụ|u)(?<nasal>(?:ñ|n)(?!(?:d|g|gb|y|z)?(?:ɛ̈|ɛ̂|ɛ̣|ɛ|ɔ̈|ɔ̂|ɔ̣|ɔ|ẍ|x̂|x̣|x|c̈|ĉ|c̣|c|ä|â|ạ|a|ë|ê|ẹ|e|ï|î|i|ị|ö|ô|ọ|o|ü|û|ụ|u))|))))`
		canonicalWordPattern     = `(?<word>(?<shift>[@~=#])(?<syllables>(?:[hHbBqQdDfgGklrmnpKPstvVwyYzZ][aAeEiIoOxcuUXC][%_:^])+))`
		sangoWordPattern         = `(?i)(?<![\p{L}\p{M}])(?<word>(?<syllables>(?:(?:b|d|f|gb|g|h|kp|k|l|mb|mp|mv|m|nd|ngb|ng|ny|nz|n|p|r|s|t|v|w|y|z|)(?:(?:ɛ̈|ɛ̂|ɛ̣|ɛ|ɔ̈|ɔ̂|ɔ̣|ɔ|ẍ|x̂|x̣|x|c̈|ĉ|c̣|c)|(?:(?:ä|â|ạ|a|ë|ê|ẹ|e|ï|î|i|ị|ö|ô|ọ|o|ü|û|ụ|u)(?:ñ|n|))))+))(?![\p{L}\p{M}])`
		sentencePattern          = `(?<sentence>.*?(?:[.!?](?:’|”|›|»)*[ \n]+|\n\n+|$))`
	)
	return [NumParseMatchers]*regexp2.Regexp{
		nil,
		regexp2.MustCompile(sentencePattern, regexp2.RE2, regexp2.OptionMaintainCaptureOrder()),
		regexp2.MustCompile(sangoWordPattern, regexp2.RE2, regexp2.OptionMaintainCaptureOrder()),
		regexp2.MustCompile(sangoSyllablePattern, regexp2.RE2, regexp2.OptionMaintainCaptureOrder()),
		regexp2.MustCompile(canonicalWordPattern, regexp2.RE2, regexp2.OptionMaintainCaptureOrder()),
		regexp2.MustCompile(canonicalSyllablePattern, regexp2.RE2, regexp2.OptionMaintainCaptureOrder()),
	}
}()

func parse(parseMatcherForGroupName map[string]ParseMatcher, s string, parseTree ParseTree) ParseTree {
	if parseTree.Begin < 0 || parseTree.Begin > parseTree.End || parseTree.End > len(s) {
		panic("Bad parseTree")
	}
	g := parseTree.Name
	parseMatcher, found := parseMatcherForGroupName[g]
	if !found {
		return ParseTree{}
	}
	if parseMatcher == Literal {
		return parseTree
	}
	re := parseMatcherREs[parseMatcher]
	if re == nil {
		return parseTree
	}
	parseTree.Subtrees = [][]ParseTree{}
	b := parseTree.Begin
	e := parseTree.End
	indexes := compat.Wrap(re).FindAllStringSubmatchIndex(s[b:e], -1)
	for i, index := range indexes {
		for j := range index {
			indexes[i][j] += b
		}
	}
	unmatched := g + "Unmatched"
	for _, index := range indexes {
		c := index[0]
		if parseMatcher == Sentence {
			// Move sentence start back to the end of the previous sentence.
			delta := c - b
			for j := range len(index) / 2 {
				index[2*j] -= delta
			}
			c = b
		}
		if c > b {
			if p, found := parseMatcherForGroupName[unmatched]; found && p == Literal {
				parseSubtrees := []ParseTree{ParseTree{Begin: b, End: c, Name: unmatched}}
				parseTree.Subtrees = append(parseTree.Subtrees, parseSubtrees)
			}
		}
		d := index[1]
		if d > c {
			parseSubtrees := []ParseTree{}
			numMatches := len(index) / 2
			for k := range numMatches {
				if n := re.GroupNameFromNumber(k); n != "" {
					if x, y := index[2*k], index[2*k+1]; x >= 0 && y >= 0 {
						parseSubtrees = append(parseSubtrees, parse(parseMatcherForGroupName, s, ParseTree{Begin: x, End: y, Name: n}))
					}
				}
			}
			parseTree.Subtrees = append(parseTree.Subtrees, parseSubtrees)
		}
		b = d
	}
	if e > b {
		if p, found := parseMatcherForGroupName[unmatched]; found && p == Literal {
			parseSubtrees := []ParseTree{ParseTree{Begin: b, End: e, Name: unmatched}}
			parseTree.Subtrees = append(parseTree.Subtrees, parseSubtrees)
		}
	}
	return parseTree
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
					o.WriteString(fmt.Sprintf("%-20v [%3v:%3v] = %q\n", n, c.Begin, c.End, s[c.Begin:c.End]))
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
