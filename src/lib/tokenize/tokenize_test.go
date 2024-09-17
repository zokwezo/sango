package tokenize

import "testing"

func TestTokenizeSango(t *testing.T) {
	check := func(s *string, expected []Token) {
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

	s := ""
	check(&s, []Token{})

	s = "mafüta tî ngû tî  Yikes! «Mon Dieu...» mɛ9 tî 《東京》 bâgara."
	check(&s, []Token{
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

func TestClassifySango(t *testing.T) {
	check := func(s *string, expected []Lemma) {
		actually := ClassifySango(s)
		if len(actually) == len(expected) {
			for k, l := range actually {
				r := expected[k]
				if l.Word != r.Word || l.Type != r.Type || l.Lang != r.Lang {
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

	s := ""
	check(&s, []Lemma{})

	s = "mafüta tî ngû tî  Yikes! «Mon Dieu...» mɛ9 tî 《東京》 bâgara."
	check(&s, []Lemma{
		{"mafüta", "WORD", "sg"},
		{" ", "SPACE", ""},
		{"tî", "WORD", "sg"},
		{" ", "SPACE", ""},
		{"ngû", "WORD", "sg"},
		{" ", "SPACE", ""},
		{"tî", "WORD", "sg"},
		{"  ", "SPACE", ""},
		{"Yikes", "WORD", "en"},
		{"!", "PUNC", "other"},
		{" ", "SPACE", ""},
		{"«", "PUNC", "open"},
		{"Mon", "WORD", "fr"},
		{" ", "SPACE", ""},
		{"Dieu", "WORD", "fr"},
		{"…", "PUNC", "other"},
		{"»", "PUNC", "close"},
		{" ", "SPACE", ""},
		{"mɛ", "WORD", "sg"},
		{"9", "NUM", ""},
		{" ", "SPACE", ""},
		{"tî", "WORD", "sg"},
		{" ", "SPACE", ""},
		{"《", "PUNC", "open"},
		{"東京", "OTHER", ""},
		{"》", "PUNC", "close"},
		{" ", "SPACE", ""},
		{"bâgara", "WORD", "sg"},
		{".", "PUNC", "other"},
	})
}
