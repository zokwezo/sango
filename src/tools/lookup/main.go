package main

import (
	"fmt"
	"slices"
	"strings"

	l "github.com/zokwezo/sango/src/lib/lexicon"
)

func main() {
	type bytesField struct {
		name  string
		field *[][]byte
	}

	type runesField struct {
		name  string
		field *[][]rune
	}

	getBytesFields := func(d l.DictCols) []bytesField {
		return []bytesField{
			bytesField{"Toneless           ", &d.Toneless},
			bytesField{"Heightless         ", &d.Heightless},
			bytesField{"LemmaUTF8          ", &d.LemmaUTF8},
			bytesField{"UDPos              ", &d.UDPos},
			bytesField{"UDFeature          ", &d.UDFeature},
			bytesField{"Category           ", &d.Category},
			bytesField{"EnglishTranslation ", &d.EnglishTranslation},
			bytesField{"EnglishDefinition  ", &d.EnglishDefinition},
		}
	}

	getRunesFields := func(d l.DictCols) []runesField {
		return []runesField{
			runesField{"Lemma              ", &d.Lemma},
		}
	}

	rows := l.LexiconRows()
	cols := l.LexiconCols()
	fmt.Printf("\n\nROWS:\n")
	for _, r := range rows {
		fmt.Printf("\t\t{%q, %q, %q, %q, %q, %q, %v, %q, %q},\n",
			r.Toneless,
			r.Heightless,
			r.Lemma,
			r.UDPos,
			r.UDFeature,
			r.Category,
			r.Frequency,
			r.EnglishTranslation,
			r.EnglishDefinition)
	}
	fmt.Printf("\n\nCOLS AS STRINGS:\n")
	fmt.Printf("Cols[Frequency ] = %v\n", cols.Frequency)
	for _, bf := range getBytesFields(cols) {
		for i, b := range *bf.field {
			fmt.Printf("Cols[%v][%v] = {%s}\n", bf.name, i, string(b))
		}
	}
	for _, bf := range getRunesFields(cols) {
		for i, b := range *bf.field {
			fmt.Printf("Cols[%v][%v] = {%s}\n", bf.name, i, string(b))
		}
	}

	cmpFunc := func(lhs, rhs l.DictRow) int {
		return strings.Compare(lhs.Toneless, rhs.Toneless)
	}
	entry := l.DictRow{Toneless: "nzoni"}
	nBegin, found := slices.BinarySearchFunc(l.LexiconRows(), entry, cmpFunc)
	fmt.Printf("Looking for %v at entry[%v] (found = %v)\n", entry, nBegin, found)
	entryEnd := l.DictRow{Toneless: "nzoni ", Lemma: ""}
	nEnd, _ := slices.BinarySearchFunc(l.LexiconRows(), entryEnd, cmpFunc)
	for n := nBegin; n < nEnd; n++ {
		fmt.Printf("entry[%v] = %v\n", n, l.LexiconRows()[n])
	}
}
