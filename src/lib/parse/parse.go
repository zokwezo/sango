// Implements operations on 64-bit word(s)

package parse

import (
	"fmt"
	"regexp"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Span struct {
	Begin   int
	End     int
	IsSango bool
}

func CanonicalToIndexes(s string) [][]int {
	return toIndexes(canonicalWordRE, s)
}

func CanonicalToSpans(s string) []Span {
	return toSpans(canonicalWordRE, s)
}

func SangoToIndexes(s string) [][]int {
	return toIndexes(sangoWordRE, s)
}

func SangoToSpans(s string) []Span {
	return toSpans(sangoWordRE, s)
}

func GetSangoSyllables(s string, span Span) {
	getSangoSyllables(s, span)
}

//////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

const (
	canonicalSyllablePattern = `(?<shift>[~=#]?)(?<consonant>[hHbBqQdDfgGklrmnpKPstvVwyYzZ])(?<vowelNasal>[aAeEiIoOxcuUXC])(?<pitch>[_:^]?)`
	sangoSyllablePattern     = `(?i)(?<consonant>(?i)b|d|f|gb|g|h|kp|k|l|mb|mp|mv|m|nd|ngb|ng|ny|nz|n|p|r|s|t|v|w|y|z|)` +
		`(?<vowelNasal>(?<openVowel>ɛ̈|ɛ̂|ɛ̣|ɛ|ɔ̈|ɔ̂|ɔ̣|ɔ|ẍ|x̂|x̣|x|c̈|ĉ|c̣|c)|(?<closeVowelNasal>(?<closeVowel>ä|â|ạ|a|ë|ê|ẹ|e|ï|î|i|ị|ö|ô|ọ|o|ü|û|ụ|u)(?<nasal>ñ|n|)))`
	canonicalWordPattern = `(?:` + canonicalSyllablePattern + `)+`
	sangoWordPattern     = `(?:` + sangoSyllablePattern + `)+`
)

var (
	canonicalSyllableRE = regexp.MustCompile(canonicalSyllablePattern)
	canonicalWordRE     = regexp.MustCompile(canonicalWordPattern) // 14 capturing groups
	sangoSyllableRE     = regexp.MustCompile(sangoSyllablePattern)
	sangoWordRE         = regexp.MustCompile(sangoWordPattern)
)

// Returns the accumulated codes and the index where processing stopped.
// The latter will equal len(s) on success, else the index of the first error.
func toIndexes(re *regexp.Regexp, s string) [][]int {
	return re.FindAllStringSubmatchIndex(s, -1)
}

// Returns the accumulated codes and the index where processing stopped.
// The latter will equal len(s) on success, else the index of the first error.
func toSpans(re *regexp.Regexp, s string) []Span {
	indexes := re.FindAllStringSubmatchIndex(s, -1)
	fmt.Printf("indexes = %#v\n", indexes)
	spans := make([]Span, 0, len(indexes)*2+1)
	n := len(s)
	a := 0
	for _, index := range indexes {
		b := index[0]
		if b > a {
			spans = append(spans, Span{IsSango: false, Begin: a, End: b})
		}
		c := index[1]
		if c > b {
			spans = append(spans, Span{IsSango: true, Begin: b, End: c})
		}
		a = c
	}
	if n > a {
		spans = append(spans, Span{IsSango: false, Begin: a, End: n})
	}
	return spans
}

func getSangoSyllables(s string, span Span) {
	if !span.IsSango {
		panic("word is not Sango")
	}
	ss := s[span.Begin:span.End]
	ii := toIndexes(sangoSyllableRE, ss)
	numSyllables := len(ii)
	for j := range numSyllables {
		numIndexes := len(ii[j])
		if numIndexes != 14 {
			panic("numIndexes != 14")
		}
		for k := range numIndexes {
			ii[j][k] += span.Begin
		}
	}
	// If consonant N was unjustly stolen to be the nasal of the previous syllable, put it back.
	// In case of ambiguity, assume consonant N (user can use an apostrophe or hyphen separator to force nasal N).
	for jCurr := 1; jCurr < numSyllables; jCurr++ {
		jPrev := jCurr - 1
		bPrevNasal := ii[jPrev][12]
		ePrevNasal := ii[jPrev][13]
		bCurrConsonant := ii[jCurr][2]
		eCurrConsonant := ii[jCurr][3]
		if d := ePrevNasal - bPrevNasal; 0 < bPrevNasal && bPrevNasal < ePrevNasal && ePrevNasal <= bCurrConsonant && bCurrConsonant <= eCurrConsonant {
			switch cases.Lower(language.English).String(s[bCurrConsonant:min(bCurrConsonant+1, eCurrConsonant)]) {
			case "d":
				fallthrough
			case "g":
				fallthrough
			case "y":
				fallthrough
			case "z":
				ii[jPrev][11] -= d
				fallthrough
			case "":
				// Return stolen nasal from previous syllable to current one.
				ii[jPrev][1] -= d
				ii[jPrev][5] -= d
				ii[jPrev][9] -= d
				ii[jPrev][13] -= d
				if ii[jPrev][12] != ii[jPrev][13] {
					panic("after returning nasal, prev nasal is not empty")
				}
				ii[jCurr][0] -= d
				ii[jCurr][2] -= d
			}
		}
	}
	for j := range numSyllables {
		fmt.Printf("word[%v] ", j)
		bw := ii[j][0]
		ew := ii[j][1]
		if ew < bw || bw < 0 {
			panic("bad extracted word")
		}
		fmt.Printf("= s[%v:%v] = %q =>", bw, ew, s[bw:ew])
		bc := ii[j][2]
		ec := ii[j][3]
		if ec < bc || bc < 0 {
			panic("bad extracted consonant")
		}
		fmt.Printf(" (s[%v:%v] = %q", bc, ec, s[bc:ec])
		bv := ii[j][6]
		ev := ii[j][7]
		if bv < 0 || ev < bv {
			bv = ii[j][10]
			ev = ii[j][11]
		}
		if ev < bv || bv < 0 {
			panic("bad extracted vowel")
		}
		fmt.Printf(", s[%v:%v] = %q", bv, ev, s[bv:ev])
		bn := ii[j][12]
		en := ii[j][13]
		if en < bn || bn < 0 {
			panic("bad extracted nasal")
		}
		fmt.Printf(", s[%v:%v] = %q)\n", bn, en, s[bn:en])
	}
}
