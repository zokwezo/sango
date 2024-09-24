package tokenize

import "testing"

func checkTokenize(t *testing.T, s *string, expected []Token) {
	o, actually := TokenizeSango(s)
	if o != s {
		t.Errorf("o (%v) != s (%v)", o, s)
		t.Errorf("o (%v) != s (%v)", *o, *s)
	} else if len(actually) == len(expected) {
		for k, l := range actually {
			r := expected[k]
			if l.Begin != r.Begin || l.End != r.End || l.REindex != r.REindex {
				t.Errorf("actually %v (%v) RE#%v != expected %v (%v) RE#%v",
					l, (*s)[l.Begin:l.End], l.REindex, r, (*s)[r.Begin:r.End], r.REindex)
				break
			}
		}
		return
	}
	t.Errorf("s = %v", *s)
	if len(expected) == 0 {
		t.Errorf("expected = %v", expected)
	} else {
		for k, token := range expected {
			t.Errorf("RE #%v expected[%v]%v = %v", token.REindex, k, token, (*s)[token.Begin:token.End])
		}
	}
	if len(actually) == 0 {
		t.Errorf("actually = %v", actually)
	} else {
		for k, token := range actually {
			t.Errorf("RE #%v actually[%v]%v = %v", token.REindex, k, token, (*s)[token.Begin:token.End])
		}
	}
}

func TestTokenizeEmpty(t *testing.T) {
	s := ""
	checkTokenize(t, &s, []Token{})
}

func TestTokenizeSango(t *testing.T) {
	s := "mafüta tî ngû tî  Yikes! «Mon Dieu...» mɛ9 tî 《東京》 bâgara."
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
		{45, 48, 3},
		{48, 49, 1},
		{49, 50, 0},
		{50, 53, 3},
		{53, 54, 0},
		{54, 57, 2},
		{57, 63, 5},
		{63, 66, 2},
		{66, 67, 0},
		{67, 74, 3},
		{74, 75, 2},
	})
}

func checkClassify(t *testing.T, s *string, expected []Lemma) {
	actually := ClassifySango(s)
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
	}
}

func TestClassifyEmpty(t *testing.T) {
	s := ""
	checkClassify(t, &s, []Lemma{})
}

func TestClassifySango(t *testing.T) {
	s := "mafüta tî ngû tî  Yikes! «Mon Dieu...» mɛ9 tî 《東京》 bâgara bânga."
	checkClassify(t, &s, []Lemma{
		{Token{0, 7, 3}, "mafuta", "mafüta", "WORD", "sg"},
		{Token{7, 8, 0}, " ", " ", "SPACE", ""},
		{Token{8, 11, 3}, "ti", "tî", "WORD", "sg"},
		{Token{11, 12, 0}, " ", " ", "SPACE", ""},
		{Token{12, 16, 3}, "ngu", "ngû", "WORD", "sg"},
		{Token{16, 17, 0}, " ", " ", "SPACE", ""},
		{Token{17, 20, 3}, "ti", "tî", "WORD", "sg"},
		{Token{20, 22, 0}, "  ", "  ", "SPACE", ""},
		{Token{22, 27, 4}, "yikes", "yikes", "WORD", "en"},
		{Token{27, 28, 2}, "!", "!", "PUNC", "other"},
		{Token{28, 29, 0}, " ", " ", "SPACE", ""},
		{Token{29, 31, 2}, "«", "«", "PUNC", "open"},
		{Token{31, 34, 3}, "mon", "mon", "WORD", "fr"},
		{Token{34, 35, 0}, " ", " ", "SPACE", ""},
		{Token{35, 39, 3}, "dieu", "dieu", "WORD", "fr"},
		{Token{39, 42, 2}, "…", "…", "PUNC", "other"},
		{Token{42, 44, 2}, "»", "»", "PUNC", "close"},
		{Token{44, 45, 0}, " ", " ", "SPACE", ""},
		{Token{45, 48, 3}, "me", "mɛ", "WORD", "sg"},
		{Token{48, 49, 1}, "9", "9", "NUM", ""},
		{Token{49, 50, 0}, " ", " ", "SPACE", ""},
		{Token{50, 53, 3}, "ti", "tî", "WORD", "sg"},
		{Token{53, 54, 0}, " ", " ", "SPACE", ""},
		{Token{54, 57, 2}, "《", "《", "PUNC", "open"},
		{Token{57, 63, 5}, "東京", "東京", "OTHER", ""},
		{Token{63, 66, 2}, "》", "》", "PUNC", "close"},
		{Token{66, 67, 0}, " ", " ", "SPACE", ""},
		{Token{67, 74, 3}, "bagara", "bâgara", "WORD", "sg"},
		{Token{74, 75, 0}, " ", " ", "SPACE", ""},
		{Token{75, 81, 3}, "banga", "bânga", "WORD", "sg"},
		{Token{81, 82, 2}, ".", ".", "PUNC", "other"},
	})
}

func TestClassifyTonelessSango(t *testing.T) {
	s := "mafuta ti ngu tî  Yikes! «Mon Dieu...» me9 ti 《東京》 bagara banga."
	checkClassify(t, &s, []Lemma{
		{Token{0, 7, 3}, "mafuta", "mafuta", "WORD", "SG"},
		{Token{7, 8, 0}, " ", " ", "SPACE", ""},
		{Token{8, 11, 3}, "ti", "ti", "WORD", "SG"},
		{Token{11, 12, 0}, " ", " ", "SPACE", ""},
		{Token{12, 16, 3}, "ngu", "ngu", "WORD", "SG"},
		{Token{16, 17, 0}, " ", " ", "SPACE", ""},
		{Token{17, 20, 3}, "ti", "tî", "WORD", "sg"},
		{Token{20, 22, 0}, "  ", "  ", "SPACE", ""},
		{Token{22, 27, 4}, "yikes", "yikes", "WORD", "en"},
		{Token{27, 28, 2}, "!", "!", "PUNC", "other"},
		{Token{28, 29, 0}, " ", " ", "SPACE", ""},
		{Token{29, 31, 2}, "«", "«", "PUNC", "open"},
		{Token{31, 34, 3}, "mon", "mon", "WORD", "fr"},
		{Token{34, 35, 0}, " ", " ", "SPACE", ""},
		{Token{35, 39, 3}, "dieu", "dieu", "WORD", "fr"},
		{Token{39, 42, 2}, "…", "…", "PUNC", "other"},
		{Token{42, 44, 2}, "»", "»", "PUNC", "close"},
		{Token{44, 45, 0}, " ", " ", "SPACE", ""},
		{Token{45, 48, 3}, "me", "me", "WORD", "SG"},
		{Token{48, 49, 1}, "9", "9", "NUM", ""},
		{Token{49, 50, 0}, " ", " ", "SPACE", ""},
		{Token{50, 53, 3}, "ti", "ti", "WORD", "SG"},
		{Token{53, 54, 0}, " ", " ", "SPACE", ""},
		{Token{54, 57, 2}, "《", "《", "PUNC", "open"},
		{Token{57, 63, 5}, "東京", "東京", "OTHER", ""},
		{Token{63, 66, 2}, "》", "》", "PUNC", "close"},
		{Token{66, 67, 0}, " ", " ", "SPACE", ""},
		{Token{67, 74, 3}, "bagara", "bagara", "WORD", "SG"},
		{Token{74, 75, 0}, " ", " ", "SPACE", ""},
		{Token{75, 81, 3}, "banga", "banga", "WORD", "Sg"},
		{Token{81, 82, 2}, ".", ".", "PUNC", "other"},
	})
}
