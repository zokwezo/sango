// Implements parsing of Sango text constructing a top-down left-to-right parse tree using regexp.

package parse

import (
	"fmt"
	"strings"

	"github.com/dlclark/regexp2/v2"
	"github.com/dlclark/regexp2/v2/compat"
	"golang.org/x/text/width"
)

const lettersPattern = `(?:(?:[A-Za-z]|Ɛ̈|Ɛ̂|Ɛ̣|Ɛ|Ɔ̈|Ɔ̂|Ɔ̣|Ɔ|Ẍ|X̂|X̣|C̈|Ĉ|C̣|Ä|Â|Ạ|Ë|Ê|Ẹ|Ï|Î|Ị|Ö|Ô|Ọ|Ü|Û|Ụ|Ñ|ɛ̈|ɛ̂|ɛ̣|ɛ|ɔ̈|ɔ̂|ɔ̣|ɔ|ẍ|x̂|x̣|c̈|ĉ|c̣|ä|â|ạ|ë|ê|ẹ|ï|î|ị|ö|ô|ọ|ü|û|ụ|ñ)\p{M}*)+`

var (
	validGroupNameRE = regexp2.MustCompile(`[A-Za-z][A-Za-z0-9_]*`, regexp2.RE2, regexp2.OptionMaintainCaptureOrder())
	parseMatcherREs  = func() [NumParseMatchers]*regexp2.Regexp {
		const (
			canonicalSyllablePattern = `(?<syllable>(?<consonant>[hHbBqQdDfgGklrmnpKPstvVwyYzZ])(?<vowelNasal>[aAeEiIoOxcuUXC])(?<pitch>[%_:^]?))`
			sangoSyllablePattern     = `(?i)(?<syllable>(?<consonant>b|d|f|gb|g|h|kp|k|l|mb|mp|mv|m|nd|ngb|ng|ny|nz|n|p|r|s|t|v|w|y|z|)(?<vowelNasal>(?<openVowel>ɛ̈|ɛ̂|ɛ̣|ɛ|ɔ̈|ɔ̂|ɔ̣|ɔ|ẍ|x̂|x̣|x|c̈|ĉ|c̣|c)|(?<closeVowelNasal>(?<closeVowel>ä|â|ạ|a|ë|ê|ẹ|e|ï|î|i|ị|ö|ô|ọ|o|ü|û|ụ|u)(?<nasal>(?:ñ|n)(?!(?:d|g|gb|y|z)?(?:ɛ̈|ɛ̂|ɛ̣|ɛ|ɔ̈|ɔ̂|ɔ̣|ɔ|ẍ|x̂|x̣|x|c̈|ĉ|c̣|c|ä|â|ạ|a|ë|ê|ẹ|e|ï|î|i|ị|ö|ô|ọ|o|ü|û|ụ|u))|))))`
			canonicalWordPattern     = `(?<canonical>(?<word>(?<shift>[~=#@])(?<syllables>(?:[hHbBqQdDfgGklrmnpKPstvVwyYzZ][aAeEiIoOxcuUXC][%_:^])+))|(?<unicode>(U[+](?<hexcode>[0-9A-F]{4}))|(?<symbol>.)))`
			sangoWordPattern         = `(?i)(?<sango>(?<!` + lettersPattern + `)(?<word>(?<syllables>(?:(?:b|d|f|gb|g|h|kp|k|l|mb|mp|mv|m|nd|ngb|ng|ny|nz|n|p|r|s|t|v|w|y|z|)(?:(?:ɛ̈|ɛ̂|ɛ̣|ɛ|ɔ̈|ɔ̂|ɔ̣|ɔ|ẍ|x̂|x̣|x|c̈|ĉ|c̣|c)|(?:(?:ä|â|ạ|a|ë|ê|ẹ|e|ï|î|i|ị|ö|ô|ọ|o|ü|û|ụ|u)(?:ñ|n|))))+))(?!` + lettersPattern + `)|(?<symbol>.))`
			sentencePattern          = `(?s)(?<sentence>(?:(?![.!?]|U\+00(?:2E|21|3F)|\n\n).)*(?:[.!? \n’”›»]|U\+(?:00(?:2E|21|3F|20|0A)|20(?:19|1D|3A)))*)`
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
)

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
	r := parseMatcherREs[parseMatcher]
	if r == nil {
		return parseTree
	}
	m := r.GroupNameFromNumber(1)
	if m == "" || m == "1" {
		m = g + "Child"
	}
	unmatched := "NON" + m
	parseTree.Subtrees = [][]ParseTree{}
	b := parseTree.Begin
	e := parseTree.End
	if e <= b {
		return parseTree
	}
	indexes := compat.Wrap(r).FindAllStringSubmatchIndex(s[b:e], -1)
	for i, index := range indexes {
		for j := range index {
			indexes[i][j] += b
		}
	}
	for _, index := range indexes {
		c := index[0]
		if c > b {
			if _, found := parseMatcherForGroupName[unmatched]; found {
				parseSubtrees := []ParseTree{ParseTree{Begin: b, End: c, Name: unmatched}}
				parseTree.Subtrees = append(parseTree.Subtrees, parseSubtrees)
			}
		}
		d := index[1]
		if d > c {
			parseSubtrees := []ParseTree{}
			numMatches := len(index) / 2
			for k := range numMatches {
				if n := r.GroupNameFromNumber(k); n != "" {
					// Ensure n is also not a stringified number
					if matched, err := validGroupNameRE.MatchString(n); matched && err == nil {
						if x, y := index[2*k], index[2*k+1]; x >= 0 && y >= 0 {
							parseSubtree := parse(parseMatcherForGroupName, s, ParseTree{Begin: x, End: y, Name: n})
							if parseSubtree.Name != "" {
								parseSubtrees = append(parseSubtrees, parseSubtree)
							}
						}
					}
				}
			}
			parseTree.Subtrees = append(parseTree.Subtrees, parseSubtrees)
		}
		b = d
	}
	if e > b {
		if _, found := parseMatcherForGroupName[unmatched]; found {
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
