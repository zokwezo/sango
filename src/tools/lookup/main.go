package main

import (
	"fmt"

	l "github.com/zokwezo/sango/src/lib/lexicon"
)

type BytesField struct {
	name  string
	field *[][]byte
}

type RunesField struct {
	name  string
	field *[][]rune
}

func GetBytesFields(d l.Dict) []BytesField {
	return []BytesField{
		BytesField{"Toneless  ", &d.Toneless},
		BytesField{"SangoUTF8 ", &d.SangoUTF8},
		BytesField{"LexPos    ", &d.LexPos},
		BytesField{"UDPos     ", &d.UDPos},
		BytesField{"UDFeature ", &d.UDFeature},
		BytesField{"Category  ", &d.Category},
		BytesField{"English   ", &d.English},
	}
}

func GetRunesFields(d l.Dict) []RunesField {
	return []RunesField{
		RunesField{"Sango     ", &d.Sango},
	}
}

func main() {
	fmt.Printf("\n\nAFFIXES:\n")
	for k, v := range Affixes.Level {
		fmt.Printf("AffixesLevel[%v] = {%v}\n", k, v)
	}
	for _, bf := range GetBytesFields(Affixes) {
		for i, b := range *bf.field {
			fmt.Printf("AffixesBytes[%v][%v] = {%s}\n", bf.name, i, string(b))
		}
	}
	for _, bf := range GetRunesFields(Affixes) {
		for i, b := range *bf.field {
			fmt.Printf("AffixesRunes[%v][%v] = {%s}\n", bf.name, i, string(b))
		}
	}

	fmt.Printf("\n\nLEXICON:\n")
	for k, v := range Lexicon.Level {
		fmt.Printf("LexiconLevel[%v] = {%v}\n", k, v)
	}
	for _, bf := range GetBytesFields(Lexicon) {
		for i, b := range *bf.field {
			fmt.Printf("LexiconBytes[%v][%v] = {%s}\n", bf.name, i, string(b))
		}
	}
	for _, bf := range GetRunesFields(Lexicon) {
		for i, b := range *bf.field {
			fmt.Printf("LexiconRunes[%v][%v] = {%s}\n", bf.name, i, string(b))
		}
	}
}
