package main

import (
	"fmt"

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
}
