package parse

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Set global flags for all tests in this package
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// log.SetOutput(io.Discard)

	// Run the test suite
	os.Exit(m.Run())
}

func TestCanonicalToIndexes(t *testing.T) {
	log.Println("ENTER TestCanonicalToIndexes")
	s := ` ~ha_HO:-Do:ni^ =ha_HO:-Do:ni^ ha^Dx_ ba^ha_-mo_-tx_nx_ ~bx^-kc:Bi:tx_bx^-kc:Bi:tx_` +
		`bx^-=kc:=Bi:=tx_ bx^-kc:Bi:tx_ ~bx^-kc:Bi:tx_ ~bx^-=kc:=Bi:=tx_ ha_HO:Do:ni^`
	const M int = 10
	type E = [M]int
	expect := [...]E{
		E{1, 8, 5, 5, 5, 6, 6, 7, 7, 8},
		E{9, 15, 12, 12, 12, 13, 13, 14, 14, 15},
		E{16, 23, 20, 20, 20, 21, 21, 22, 22, 23},
		E{24, 30, 27, 27, 27, 28, 28, 29, 29, 30},
		E{31, 37, 34, 34, 34, 35, 35, 36, 36, 37},
		E{38, 44, 41, 41, 41, 42, 42, 43, 43, 44},
		E{45, 48, 45, 45, 45, 46, 46, 47, 47, 48},
		E{49, 55, 52, 52, 52, 53, 53, 54, 54, 55},
		E{56, 60, 56, 57, 57, 58, 58, 59, 59, 60},
		E{61, 73, 70, 70, 70, 71, 71, 72, 72, 73},
		E{74, 86, 83, 83, 83, 84, 84, 85, 85, 86},
		E{87, 99, 95, 96, 96, 97, 97, 98, 98, 99},
		E{100, 103, 100, 100, 100, 101, 101, 102, 102, 103},
		E{104, 113, 110, 110, 110, 111, 111, 112, 112, 113},
		E{114, 118, 114, 115, 115, 116, 116, 117, 117, 118},
		E{119, 128, 125, 125, 125, 126, 126, 127, 127, 128},
		E{129, 133, 129, 130, 130, 131, 131, 132, 132, 133},
		E{134, 146, 142, 143, 143, 144, 144, 145, 145, 146},
		E{147, 159, 156, 156, 156, 157, 157, 158, 158, 159},
	}
	actual := CanonicalToIndexes(s)
	log.Printf("actual = %#v\n", actual)
	na := len(actual)
	ne := len(expect)
	n := min(na, ne)
	if na != ne {
		t.Errorf("bad CanonicalToIndexes:\nfound %v matches but expected %v\n", na, ne)
	}
	for i := range n {
		a := actual[i]
		e := expect[i]
		m := len(a)
		if m != M {
			t.Errorf("bad CanonicalToIndexes:\nfound len(actual[%v]) = %v but expected %v\n", i, m, M)
		}
		for j := range m {
			if a[j] != e[j] {
				t.Errorf("bad CanonicalToIndexes:\nfound actual[%v][%v] = %v but expected %v\n", i, j, a[j], e[j])
			}
		}
	}
	log.Println("LEAVE TestCanonicalToIndexes")
}

func TestCanonicalToSpans(t *testing.T) {
	log.Println("ENTER TestCanonicalToSpans")
	s := ` ~ha_HO:-Do:ni^ =ha_HO:-Do:ni^ ha^Dx_ ba^ha_-mo_-tx_nx_ ~bx^-kc:Bi:tx_bx^-kc:Bi:tx_` +
		`bx^-=kc:=Bi:=tx_ bx^-kc:Bi:tx_ ~bx^-kc:Bi:tx_ ~bx^-=kc:=Bi:=tx_ ha_HO:Do:ni^`
	expect := [...]Span{
		Span{Begin: 0, End: 1, IsSango: false},
		Span{Begin: 1, End: 8, IsSango: true},
		Span{Begin: 8, End: 9, IsSango: false},
		Span{Begin: 9, End: 15, IsSango: true},
		Span{Begin: 15, End: 16, IsSango: false},
		Span{Begin: 16, End: 23, IsSango: true},
		Span{Begin: 23, End: 24, IsSango: false},
		Span{Begin: 24, End: 30, IsSango: true},
		Span{Begin: 30, End: 31, IsSango: false},
		Span{Begin: 31, End: 37, IsSango: true},
		Span{Begin: 37, End: 38, IsSango: false},
		Span{Begin: 38, End: 44, IsSango: true},
		Span{Begin: 44, End: 45, IsSango: false},
		Span{Begin: 45, End: 48, IsSango: true},
		Span{Begin: 48, End: 49, IsSango: false},
		Span{Begin: 49, End: 55, IsSango: true},
		Span{Begin: 55, End: 56, IsSango: false},
		Span{Begin: 56, End: 60, IsSango: true},
		Span{Begin: 60, End: 61, IsSango: false},
		Span{Begin: 61, End: 73, IsSango: true},
		Span{Begin: 73, End: 74, IsSango: false},
		Span{Begin: 74, End: 86, IsSango: true},
		Span{Begin: 86, End: 87, IsSango: false},
		Span{Begin: 87, End: 99, IsSango: true},
		Span{Begin: 99, End: 100, IsSango: false},
		Span{Begin: 100, End: 103, IsSango: true},
		Span{Begin: 103, End: 104, IsSango: false},
		Span{Begin: 104, End: 113, IsSango: true},
		Span{Begin: 113, End: 114, IsSango: false},
		Span{Begin: 114, End: 118, IsSango: true},
		Span{Begin: 118, End: 119, IsSango: false},
		Span{Begin: 119, End: 128, IsSango: true},
		Span{Begin: 128, End: 129, IsSango: false},
		Span{Begin: 129, End: 133, IsSango: true},
		Span{Begin: 133, End: 134, IsSango: false},
		Span{Begin: 134, End: 146, IsSango: true},
		Span{Begin: 146, End: 147, IsSango: false},
		Span{Begin: 147, End: 159, IsSango: true},
	}
	actual := CanonicalToSpans(s)
	log.Printf("actual = %#v\n", actual)
	na := len(actual)
	ne := len(expect)
	n := min(na, ne)
	if na != ne {
		t.Errorf("bad CanonicalToSpans:\nfound %v matches but expected %v\n", na, ne)
	}
	for i := range n {
		a := actual[i]
		e := expect[i]
		log.Printf("actual[%v] = s[%v:%v] = %q\n", i, a.Begin, a.End, s[a.Begin:a.End])
		if a.Begin != e.Begin || a.End != e.End || a.IsSango != e.IsSango {
			t.Errorf("bad CanonicalToSpans:\nfound span[%v] = %#v but expected %#v\n", i, a, e)
		}
	}
	log.Println("LEAVE TestCanonicalToSpans")
}

func TestSangoToIndexes(t *testing.T) {
	log.Println("ENTER TestSangoToIndexes")
	s := `日本語は難しい! § Ahöñ-ndönî AHÖÑ-NDÖNÎ ândɛ bâa-mo-tɛnɛ` +
		` BƐ̂-kɔ̈mbïtɛbɛ̂-kɔ̈mbïtɛ難BƐ̂-KƆ̈MBÏTƐ bɛ̂-kɔ̈mbïtɛ BƐ̂-kɔ̈mbïtɛ BƐ̂-KƆ̈MBÏTƐ ahöñndönî`
	const M int = 14
	type E = [M]int
	expect := [...]E{
		E{26, 32, 27, 28, 28, 32, -1, -1, 28, 32, 28, 30, 30, 32},
		E{33, 40, 38, 38, 38, 40, -1, -1, 38, 40, 38, 40, 40, 40},
		E{41, 47, 42, 43, 43, 47, -1, -1, 43, 47, 43, 45, 45, 47},
		E{48, 55, 53, 53, 53, 55, -1, -1, 53, 55, 53, 55, 55, 55},
		E{56, 62, 59, 60, 60, 62, 60, 62, 56, 59, 56, 58, 58, 59},
		E{63, 67, 66, 66, 66, 67, -1, -1, 66, 67, 66, 67, 67, 67},
		E{68, 70, 68, 69, 69, 70, -1, -1, 69, 70, 69, 70, 70, 70},
		E{71, 77, 74, 75, 75, 77, 75, 77, -1, -1, -1, -1, -1, -1},
		E{78, 83, 78, 79, 79, 83, 79, 83, -1, -1, -1, -1, -1, -1},
		E{84, 101, 96, 97, 97, 101, 97, 101, 91, 93, 91, 93, 93, 93},
		E{102, 114, 111, 112, 112, 114, 112, 114, 109, 111, 109, 111, 111, 111},
		E{117, 122, 117, 118, 118, 122, 118, 122, -1, -1, -1, -1, -1, -1},
		E{123, 135, 132, 133, 133, 135, 133, 135, 130, 132, 130, 132, 132, 132},
		E{136, 141, 136, 137, 137, 141, 137, 141, -1, -1, -1, -1, -1, -1},
		E{142, 154, 151, 152, 152, 154, 152, 154, 149, 151, 149, 151, 151, 151},
		E{155, 160, 155, 156, 156, 160, 156, 160, -1, -1, -1, -1, -1, -1},
		E{161, 173, 170, 171, 171, 173, 171, 173, 168, 170, 168, 170, 170, 170},
		E{174, 179, 174, 175, 175, 179, 175, 179, -1, -1, -1, -1, -1, -1},
		E{180, 192, 189, 190, 190, 192, 190, 192, 187, 189, 187, 189, 189, 189},
		E{193, 206, 204, 204, 204, 206, -1, -1, 204, 206, 204, 206, 206, 206},
	}
	actual := SangoToIndexes(s)
	log.Printf("actual = %#v\n", actual)
	na := len(actual)
	ne := len(expect)
	n := min(na, ne)
	if na != ne {
		t.Errorf("bad SangoToIndexes:\nfound %v matches but expected %v\n", na, ne)
	}
	for i := range n {
		a := actual[i]
		e := expect[i]
		m := len(a)
		if m != M {
			t.Errorf("bad SangoToIndexes:\nfound len(actual[%v]) = %v but expected %v\n", i, m, M)
		}
		for j := range m {
			if a[j] != e[j] {
				t.Errorf("bad SangoToIndexes:\nfound actual[%v][%v] = %v but expected %v\n", i, j, a[j], e[j])
			}
		}
	}
	log.Println("LEAVE TestSangoToIndexes")
}

func TestSangoToSpans(t *testing.T) {
	log.Println("ENTER TestSangoToSpans")
	s := `日本語は難しい! § Ahöñ-ndönî AHÖÑ-NDÖNÎ ândɛ bâa-mo-tɛnɛ` +
		` BƐ̂-kɔ̈mbïtɛbɛ̂-kɔ̈mbïtɛ難BƐ̂-KƆ̈MBÏTƐ bɛ̂-kɔ̈mbïtɛ BƐ̂-kɔ̈mbïtɛ BƐ̂-KƆ̈MBÏTƐ ahöñndönî`
	expect := [...]Span{
		Span{Begin: 0, End: 26, IsSango: false},
		Span{Begin: 26, End: 32, IsSango: true},
		Span{Begin: 32, End: 33, IsSango: false},
		Span{Begin: 33, End: 40, IsSango: true},
		Span{Begin: 40, End: 41, IsSango: false},
		Span{Begin: 41, End: 47, IsSango: true},
		Span{Begin: 47, End: 48, IsSango: false},
		Span{Begin: 48, End: 55, IsSango: true},
		Span{Begin: 55, End: 56, IsSango: false},
		Span{Begin: 56, End: 62, IsSango: true},
		Span{Begin: 62, End: 63, IsSango: false},
		Span{Begin: 63, End: 67, IsSango: true},
		Span{Begin: 67, End: 68, IsSango: false},
		Span{Begin: 68, End: 70, IsSango: true},
		Span{Begin: 70, End: 71, IsSango: false},
		Span{Begin: 71, End: 77, IsSango: true},
		Span{Begin: 77, End: 78, IsSango: false},
		Span{Begin: 78, End: 83, IsSango: true},
		Span{Begin: 83, End: 84, IsSango: false},
		Span{Begin: 84, End: 101, IsSango: true},
		Span{Begin: 101, End: 102, IsSango: false},
		Span{Begin: 102, End: 114, IsSango: true},
		Span{Begin: 114, End: 117, IsSango: false},
		Span{Begin: 117, End: 122, IsSango: true},
		Span{Begin: 122, End: 123, IsSango: false},
		Span{Begin: 123, End: 135, IsSango: true},
		Span{Begin: 135, End: 136, IsSango: false},
		Span{Begin: 136, End: 141, IsSango: true},
		Span{Begin: 141, End: 142, IsSango: false},
		Span{Begin: 142, End: 154, IsSango: true},
		Span{Begin: 154, End: 155, IsSango: false},
		Span{Begin: 155, End: 160, IsSango: true},
		Span{Begin: 160, End: 161, IsSango: false},
		Span{Begin: 161, End: 173, IsSango: true},
		Span{Begin: 173, End: 174, IsSango: false},
		Span{Begin: 174, End: 179, IsSango: true},
		Span{Begin: 179, End: 180, IsSango: false},
		Span{Begin: 180, End: 192, IsSango: true},
		Span{Begin: 192, End: 193, IsSango: false},
		Span{Begin: 193, End: 206, IsSango: true},
	}
	actual := SangoToSpans(s)
	na := len(actual)
	ne := len(expect)
	n := min(na, ne)
	if na != ne {
		t.Errorf("bad SangoToSpans:\nfound %v matches but expected %v\n", na, ne)
	}
	for i := range n {
		aa := actual[i]
		ee := expect[i]
		if aa.Begin != ee.Begin || aa.End != ee.End || aa.IsSango != ee.IsSango {
			t.Errorf("bad SangoToSpans:\nfound span[%v] = %#v but expected %#v\n", i, aa, ee)
		}
		if aa.IsSango {
			GetSangoSyllables(s, aa)
		} else {
			unicodeRunes := []rune(s[aa.Begin:aa.End])
			for k, rune := range unicodeRunes {
				fmt.Printf("rune[%v] = 0x%08X = %q\n", k, rune, string(rune))
			}
		}
	}
	log.Println("LEAVE TestSangoToSpans")
}
