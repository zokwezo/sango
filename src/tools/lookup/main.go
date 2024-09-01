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
			bytesField{"Toneless  ", &d.Toneless},
			bytesField{"SangoUTF8 ", &d.SangoUTF8},
			bytesField{"LexPos    ", &d.LexPos},
			bytesField{"UDPos     ", &d.UDPos},
			bytesField{"UDFeature ", &d.UDFeature},
			bytesField{"Category  ", &d.Category},
			bytesField{"English   ", &d.English},
		}
	}

	getRunesFields := func(d l.DictCols) []runesField {
		return []runesField{
			runesField{"Sango     ", &d.Sango},
		}
	}

	type Dict struct {
		name string
		rows l.DictRows
		cols l.DictCols
	}

	dicts := []Dict{
		{"AFFIXES", l.AffixesRows(), l.AffixesCols()},
		{"LEXICON", l.LexiconRows(), l.LexiconCols()},
	}

	for _, dict := range dicts {
		fmt.Printf("\n\n%v ROWS:\n%v\n", dict.name, dict.rows)
		fmt.Printf("\n\n%v COLS AS STRINGS:\n", dict.name)
		fmt.Printf("%v Cols[Frequency ] = %v\n", dict.name, dict.cols.Frequency)
		for _, bf := range getBytesFields(dict.cols) {
			for i, b := range *bf.field {
				fmt.Printf("%v Cols[%v][%v] = {%s}\n", dict.name, bf.name, i, string(b))
			}
		}
		for _, bf := range getRunesFields(dict.cols) {
			for i, b := range *bf.field {
				fmt.Printf("%v Cols[%v][%v] = {%s}\n", dict.name, bf.name, i, string(b))
			}
		}
	}

	// Look for first entry <= {Toneless: "butuma" Sango: "butuma"}
	cmpFunc := func(lhs, rhs l.DictRow) int {
		if cmp := strings.Compare(lhs.Toneless, rhs.Toneless); cmp != 0 {
			return cmp
		}
		return strings.Compare(lhs.Sango, rhs.Sango)
	}
	entry := l.DictRow{Toneless: "butuma", Sango: "butuma"}
	nBegin, found := slices.BinarySearchFunc(l.LexiconRows(), entry, cmpFunc)
	fmt.Printf("Looking for %v at entry[%v] (found = %v)\n", entry, nBegin, found)
	entry = l.DictRow{Toneless: "butuma", Sango: "butumb"}
	nEnd, _ := slices.BinarySearchFunc(l.LexiconRows(), entry, cmpFunc)
	for n := nBegin; n < nEnd; n++ {
		fmt.Printf("entry[%v] = %v\n", n, l.LexiconRows()[n])
	}
}
