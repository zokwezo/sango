// Restores vowel height and pitch to Sango input.

package restore

import (
	"fmt"
	"slices"
	"strings"

	"github.com/zokwezo/sango/src/lib/lexicon"
	"github.com/zokwezo/sango/src/lib/tokenize"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type lexiconRowsFromTonelessMap = map[string][]*lexicon.DictRow

var lexiconRowsFromToneless = func() lexiconRowsFromTonelessMap {
	out := lexiconRowsFromTonelessMap{}
	for _, row := range lexicon.LexiconRows() {
		out[row.Toneless] = append(out[row.Toneless], &row)
	}
	return out
}()

var titleCaseOf cases.Caser = cases.Title(language.Und)
var upperCaseOf cases.Caser = cases.Upper(language.Und)
var lowerCaseOf cases.Caser = cases.Lower(language.Und)

func RestoreSangoVowels(s string) string {
	out := ""
	lemmas := tokenize.ClassifySango(&s)
	for k, lemma := range lemmas {
		w := s[lemma.Source.Begin:lemma.Source.End]
		var c cases.Caser
		if w == titleCaseOf.String(w) {
			c = titleCaseOf
		} else if w == upperCaseOf.String(w) {
			c = upperCaseOf
		} else if w == lowerCaseOf.String(w) {
			c = lowerCaseOf
		} else {
			panic("Mixed case unsupported")
		}
		o := []string{}
		fmt.Printf("Lemma[%v] = %v\n", k, lemma)
		switch lemma.Lang {
		case "sg":
			for r, row := range lexiconRowsFromToneless[lemma.Toneless] {
				if row.Frequency <= 6 && row.Sango == lowerCaseOf.String(w) {
					o = append(o, c.String(row.Sango))
					fmt.Printf("  Lex[%v] = {%q, %q, %q, %q, %v, %q}\n", r, row.Sango, row.UDPos, row.UDFeature, row.Category, row.Frequency, row.EnglishTranslation)
				}
			}
		case "SG":
			for r, row := range lexiconRowsFromToneless[lemma.Toneless] {
				if row.Frequency <= 6 {
					o = append(o, row.Sango)
					fmt.Printf("  Lex[%v] = {%q, %q, %q, %q, %v, %q}\n", r, row.Sango, row.UDPos, row.UDFeature, row.Category, row.Frequency, row.EnglishTranslation)
				}
			}
		default:
			o = append(o, c.String(lemma.Sango))
		}
		slices.Sort(o)
		o = slices.Compact(o)
		out += strings.Join(o, "|")
	}
	return out
}
