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
	for k, r := range lex.LexiconRows() {
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
		sses, err := sse.CanonicalToSSEs(d)
		if err != nil {
			fmt.Printf("unexpected error returned from CanonicalToSSEs[%v]\nerr = %v", d, err)
			panic("bad Canonical")
		}
		// The lexicon.Canonical is missing hyphens. Loop through every infix combination until a match.
		word := sse.SSE(^uint64(0))
		if k == 0 && a == "" {
		} else if len(sses) == 1 {
			word = sses[0]
		} else {
			panic("bad len(sses)")
		}
		/*
				found := false
				heightless := ""
				if len(d)%3 != 0 {
					panic("bad len(d)")
				}
				numSyllables := len(d) / 3
				var InfixMask uint64 = 0b0_0_00_1_00000_0000_00_1_00000_0000_00_1_00000_0000_00_1_00000_0000_00_1_00000_0000_00
				InfixMask >>= 12 * (5 - numSyllables)
				InfixMask <<= 12 * (5 - numSyllables)
				var InfixPad uint64 = 0x0FFF_FFFF_FFFF_FFFF >> (12 * numSyllables)
				var InfixPass uint64 = ^InfixMask
			HyphenLoop:
				for s4 := range 2 {
					set4 := uint64(s4) << 59
					for s3 := range 2 {
						set3 := uint64(s3) << 47
						set3 |= set4
						for s2 := range 2 {
							set2 := uint64(s2) << 35
							set2 |= set3
							for s1 := range 2 {
								set1 := uint64(s1) << 23
								set1 |= set2
								for s0 := range 2 {
									set0 := uint64(s0) << 11
									set0 |= set1
									set0 &= InfixMask
									set0 |= InfixPad
									// fmt.Printf("word = %064b\n", uint64(word))
									// fmt.Printf("set0 = %064b\n", set0)
									x := sse.SSE(uint64(word)&InfixPass | set0)
									// fmt.Printf("x    = %064b\n", x)
									var s strings.Builder
									x.WriteAsHeightlessTo(&s)
									heightless = s.String()
									hhh := []byte(heightless)
									m := "\nutf8: "
									for _, z := range hhh {
										m += fmt.Sprintf(" %02X", z)
									}
									actual := []rune(heightless)
									m += "\nactual: "
									for _, z := range actual {
										m += fmt.Sprintf(" %04X", z)
									}
									expect := []rune(b)
									m += "\nexpect: "
									for _, z := range expect {
										m += fmt.Sprintf(" %04X", z)
									}
									// fmt.Printf("  %01b%01b%01b%01b%01b %q %q %v %v\n", s4, s3, s2, s1, s0, b, heightless, heightless == b, m)
									if heightless == b {
										found = true
										var sd strings.Builder
										x.WriteAsCanonicalTo(&sd)
										d = sd.String()
										break HyphenLoop
									}
								}
							}
						}
					}
				}
		*/
		{
			var sd strings.Builder
			word.WriteAsTonelessTo(&sd)
			if a != sd.String() {
				panic("bad Toneless")
			}
		}
		{
			var sd strings.Builder
			word.WriteAsHeightlessTo(&sd)
			if b != sd.String() {
				panic("bad Heightless")
			}
		}
		{
			var sd strings.Builder
			word.WriteAsLemmaTo(&sd)
			if c != sd.String() {
				panic("bad Lemma")
			}
		}
		{
			var sd strings.Builder
			word.WriteAsCanonicalTo(&sd)
			if d != sd.String() {
				panic("bad Canonical")
			}
		}
		fmt.Printf("\t\t{%q, %q, %q, %q, %q, %q, %q, %v, %q, %q},\n", a, b, c, d, e, f, g, h, i, j)
		/*
			if !found {
				fmt.Printf("heightless = %v\n", heightless)
				panic("TEST")
			}
		*/
	}
	/*
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
	*/
}
