package tokenize

import (
	"fmt"
	"strings"
	"testing"
)

func checkTokenize(t *testing.T, s *string, expected []Token) {
	in := strings.NewReader(*s)
	o, actually := TokenizeSango(in)
	if *o != *s {
		t.Errorf("o (%v) != s (%v)", *o, *s)
	} else if len(actually) != len(expected) {
		t.Errorf("len(actually) = %v but len(expected) = %v\n", len(actually), len(expected))
	} else {
		for k, l := range actually {
			r := expected[k]
			if l.Begin != r.Begin {
				t.Errorf("actually[%v].Begin = %v\n", k, l.Begin)
				t.Errorf("expected[%v].Begin = %v\n", k, r.Begin)
				break
			}
			if l.End != r.End {
				t.Errorf("actually[%v].End = %v\n", k, l.End)
				t.Errorf("expected[%v].End = %v\n", k, r.End)
				break
			}
			if l.REindex != r.REindex {
				t.Errorf("actually[%v].REindex = %v\n", k, l.REindex)
				t.Errorf("expected[%v].REindex = %v\n", k, r.REindex)
				break
			}
			if l.Begin != r.Begin || l.End != r.End || l.REindex != r.REindex {
				t.Errorf("actually %v (%v) RE#%v != expected %v (%v) RE#%v",
					l, (*s)[l.Begin:l.End], l.REindex, r, (*s)[r.Begin:r.End], r.REindex)
				break
			}
		}
		return
	}
}

func TestTokenizeEmpty(t *testing.T) {
	s := "kcjliqngbaj hoqnndoq tijnli hojntiq"
	checkTokenize(t, &s, []Token{
		{0, 11, 3},
		{11, 12, 0},
		{12, 20, 3},
		{20, 21, 0},
		{21, 27, 3},
		{27, 28, 0},
		{28, 35, 3},
	})
}

func TestTokenizeSango(t *testing.T) {
	s := "mafuqta tij nguj tij  Yikes! «Mon Dieu...» mx9 tij 《東京》 bajgara."
	checkTokenize(t, &s, []Token{
		{0, 7, 3},
		{7, 8, 0},
		{8, 11, 3},
		{11, 12, 0},
		{12, 16, 3},
		{16, 17, 0},
		{17, 20, 3},
		{20, 22, 0},
		{22, 27, 4},
		{27, 28, 2},
		{28, 29, 0},
		{29, 31, 2},
		{31, 34, 3},
		{34, 35, 0},
		{35, 39, 3},
		{39, 42, 2},
		{42, 44, 2},
		{44, 45, 0},
		{45, 47, 3},
		{47, 48, 1},
		{48, 49, 0},
		{49, 52, 3},
		{52, 53, 0},
		{53, 56, 2},
		{56, 62, 5},
		{62, 65, 2},
		{65, 66, 0},
		{66, 73, 3},
		{73, 74, 2},
	})
}

func checkClassify(t *testing.T, s *string, expected []Lemma) {
	in := strings.NewReader(*s)
	actually := ClassifySango(in)
	if len(actually) == len(expected) {
		for k, l := range actually {
			r := expected[k]
			if l.Toneless != r.Toneless || l.Sango != r.Sango || l.Type != r.Type || l.Lang != r.Lang {
				t.Errorf("actually[%v] = %v != expected[%v] = %v", k, l, k, r)
				break
			}
		}
		return
	}
	t.Errorf("s = %v", *s)
	if len(expected) == 0 {
		t.Errorf("expected = %v", expected)
	} else {
		for k, lemma := range expected {
			t.Errorf("expected[%v] = %v", k, lemma)
		}
	}
	if len(actually) == 0 {
		t.Errorf("actually = %v", actually)
	} else {
		for k, lemma := range actually {
			t.Errorf("actually[%v] = %v", k, lemma)
		}
		for _, l := range actually {
			fmt.Printf("    {Token{%v, %v, %v}, %q, %q, %q, %q\n", l.Source.Begin, l.Source.End, l.Source.REindex, l.Toneless, l.Sango, l.Type, l.Lang)
		}
	}
}

func TestClassifyEmpty(t *testing.T) {
	s := "kcjliqngbaj hoqnndoq tijnli hojntiq"
	checkClassify(t, &s, []Lemma{
		{Token{0, 11, 3}, "kolingba", "kcjliqngbaj", "WORD", "XX"},
		{Token{11, 12, 0}, " ", " ", "SPACE", ""},
		{Token{12, 20, 3}, "honndo", "hoqnndoq", "WORD", "XX"},
		{Token{20, 21, 0}, " ", " ", "SPACE", ""},
		{Token{21, 27, 3}, "tinli", "tijnli", "WORD", "XX"},
		{Token{27, 28, 0}, " ", " ", "SPACE", ""},
		{Token{28, 35, 3}, "honti", "hojntiq", "WORD", "sg"},
	})
}

func TestClassifySango(t *testing.T) {
	s := "mafuqta tij nguj niJ atiq  Yikes! «Mon Dieu...» mx9 tij tenx txqnxqngcq txqnxq tij asdfgk《東京》."
	checkClassify(t, &s, []Lemma{
		{Token{0, 7, 3}, "mafuta", "mafuqta", "WORD", "sg"},
		{Token{7, 8, 0}, " ", " ", "SPACE", ""},
		{Token{8, 11, 3}, "ti", "tij", "WORD", "sg"},
		{Token{11, 12, 0}, " ", " ", "SPACE", ""},
		{Token{12, 16, 3}, "ngu", "nguj", "WORD", "sg"},
		{Token{16, 17, 0}, " ", " ", "SPACE", ""},
		{Token{17, 20, 3}, "ni", "nij", "WORD", "sg"},
		{Token{20, 21, 0}, " ", " ", "SPACE", ""},
		{Token{21, 25, 3}, "ati", "atiq", "WORD", "sg"},
		{Token{25, 27, 0}, "  ", "  ", "SPACE", ""},
		{Token{27, 32, 4}, "yikes", "yikes", "WORD", "en"},
		{Token{32, 33, 2}, "!", "!", "PUNC", "other"},
		{Token{33, 34, 0}, " ", " ", "SPACE", ""},
		{Token{34, 36, 2}, "«", "«", "PUNC", "open"},
		{Token{36, 39, 3}, "mon", "mon", "WORD", "fr"},
		{Token{39, 40, 0}, " ", " ", "SPACE", ""},
		{Token{40, 44, 3}, "dieu", "dieu", "WORD", "fr"},
		{Token{44, 47, 2}, "...", "...", "PUNC", "other"},
		{Token{47, 49, 2}, "»", "»", "PUNC", "close"},
		{Token{49, 50, 0}, " ", " ", "SPACE", ""},
		{Token{50, 52, 3}, "me", "mx", "WORD", "sg"},
		{Token{52, 53, 1}, "9", "9", "NUM", ""},
		{Token{53, 54, 0}, " ", " ", "SPACE", ""},
		{Token{54, 57, 3}, "ti", "tij", "WORD", "sg"},
		{Token{57, 58, 0}, " ", " ", "SPACE", ""},
		{Token{58, 62, 3}, "tene", "tenx", "WORD", "SG"},
		{Token{62, 63, 0}, " ", " ", "SPACE", ""},
		{Token{63, 73, 3}, "tenengo", "txqnxqngcq", "WORD", "sg"},
		{Token{73, 74, 0}, " ", " ", "SPACE", ""},
		{Token{74, 79, 3}, "tene", "txqnxq", "WORD", "sg"},
		{Token{79, 80, 0}, " ", " ", "SPACE", ""},
		{Token{80, 83, 3}, "ti", "tij", "WORD", "sg"},
		{Token{83, 84, 0}, " ", " ", "SPACE", ""},
		{Token{84, 91, 3}, "asdfgk", "asdfgk", "WORD", "XX"},
		{Token{91, 94, 2}, "《", "《", "PUNC", "open"},
		{Token{94, 100, 5}, "東京", "東京", "OTHER", ""},
		{Token{100, 103, 2}, "》", "》", "PUNC", "close"},
		{Token{103, 104, 2}, ".", ".", "PUNC", "other"},
	})
}
