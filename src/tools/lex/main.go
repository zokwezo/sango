package main

import (
	"fmt"
	"strings"

	lex "github.com/zokwezo/sango/src/lib/lexicon"
	sse "github.com/zokwezo/sango/src/lib/sse"
)

func main() {
	posCount := make(map[string]int)
	categoryCount := make(map[string]int)
	frequencyCount := make(map[int]int)
	for _, r := range lex.LexiconRows() {
		a := strings.ToLower(r.Toneless)
		b := strings.ToLower(r.Heightless)
		c := strings.ToLower(r.Lemma)
		d := r.Canonical
		e := strings.ToUpper(r.UDPos)
		f := r.UDFeature
		g := strings.ToUpper(r.Category)
		h := r.Frequency
		i := r.EnglishTranslation
		j := r.EnglishDefinition
		posCount[e]++
		categoryCount[g]++
		frequencyCount[h]++

		// Flag lexicon entries with bad Canonical
		words, err := sse.CanonicalToSSEs(d)
		if err != nil {
			fmt.Printf("unexpected error returned from CanonicalToSSEs[%v]\nerr = %v", d, err)
			panic("bad Canonical")
		}
		var sa strings.Builder
		for _, word := range words {
			word.WriteAsTonelessTo(&sa)
		}
		a = strings.ReplaceAll(a, "-", " ")
		aa := strings.ReplaceAll(sa.String(), "-", " ")
		if a != aa {
			fmt.Printf("bad Toneless: expected %q but found %q\n", a, aa)
			panic("bad Toneless")
		}
		a = sa.String()
		var sb strings.Builder
		for _, word := range words {
			word.WriteAsHeightlessTo(&sb)
		}
		b = strings.ReplaceAll(b, "-", " ")
		bb := strings.ReplaceAll(sb.String(), "-", " ")
		if b != bb {
			fmt.Printf("bad Heightless: expected %q but found %q\n", b, bb)
			panic("bad Heightless")
		}
		b = sb.String()
		var sc strings.Builder
		for _, word := range words {
			word.WriteAsLemmaTo(&sc)
		}
		c = strings.ReplaceAll(c, "-", " ")
		cc := strings.ReplaceAll(sc.String(), "-", " ")
		if c != cc {
			fmt.Printf("bad Lemma: expected %q but found %q\n", c, cc)
			panic("bad Lemma")
		}
		c = sc.String()
		var sd strings.Builder
		for _, word := range words {
			word.WriteAsCanonicalTo(&sd)
		}
		if d != sd.String() {
			fmt.Printf("bad Canonical: expected %q but found %q\n", d, sd.String())
			panic("bad Canonical")
		}
		if d != sd.String() {
			fmt.Printf("%q != %q\n", d, sd.String())
			panic("bad Canonical")
		}
		fmt.Printf("\t\t{%q, %q, %q, %q, %q, %q, %q, %v, %q, %q},\n", a, b, c, d, e, f, g, h, i, j)
	}

	// Output some statistics.
	fmt.Println("")
	fmt.Println("POS")
	posTotal := 0
	for _, count := range posCount {
		posTotal += count
	}
	for pos, count := range posCount {
		fmt.Printf("%s\t%v\t%v\n", pos, count, (200*count+posTotal)/(posTotal*2))
	}
	fmt.Println("")
	fmt.Println("CATEGORY")
	categoryTotal := 0
	for _, count := range categoryCount {
		categoryTotal += count
	}
	for category, count := range categoryCount {
		fmt.Printf("%s\t%v\t%v\n", category, count, (200*count+categoryTotal)/(categoryTotal*2))
	}
	fmt.Println("")
	fmt.Println("FREQUENCY")
	frequencyTotal := 0
	for _, count := range frequencyCount {
		frequencyTotal += count
	}
	for frequency, count := range frequencyCount {
		fmt.Printf("%v\t%v\t%v\n", frequency, count, (200*count+frequencyTotal)/(frequencyTotal*2))
	}
}
