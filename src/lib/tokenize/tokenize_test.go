package tokenize

import (
	"testing"
)

func TestTokenize(t *testing.T) {
	check := func(s *string, expected []Token) {
		o, actually := Tokenize(s, SangoTokenizerRegexps)
		if o != s {
			t.Errorf("o (%v) != s (%v)", o, s)
			t.Errorf("o (%v) != s (%v)", *o, *s)
		} else if len(actually) == len(expected) {
			for k, l := range actually {
				r := expected[k]
				if l.begin != r.begin || l.end != r.end {
					t.Errorf("actually %v (%v) != expected %v (%v)", l, (*s)[l.begin:l.end], r, (*s)[r.begin:r.end])
					break
				}
			}
			return
		}
		if len(expected) == 0 {
			t.Errorf("expected = %v", expected)
		} else {
			for k, token := range expected {
				t.Errorf("RE #%v expected[%v]%v = %v", token.reIndex, k, token, (*s)[token.begin:token.end])
			}
		}
		if len(actually) == 0 {
			t.Errorf("actually = %v", actually)
		} else {
			for k, token := range actually {
				t.Errorf("RE #%v actually[%v]%v = %v", token.reIndex, k, token, (*s)[token.begin:token.end])
			}
		}
		t.Errorf("s = %v", *s)
	}

	s := ""
	check(&s, []Token{})

	s = "mafüta tî ngû tî Golly! mɛ9 tî 東京 bâgara."
	check(&s, []Token{
		{0, 7, 3},
		{7, 8, 0},
		{8, 11, 3},
		{11, 12, 0},
		{12, 16, 3},
		{16, 17, 0},
		{17, 20, 3},
		{20, 21, 0},
		{21, 26, 4},
		{26, 27, 2},
		{27, 28, 0},
		{28, 31, 3},
		{31, 32, 1},
		{32, 33, 0},
		{33, 36, 3},
		{36, 37, 0},
		{37, 43, 5},
		{43, 44, 0},
		{44, 51, 3},
		{51, 52, 2},
	})
}
