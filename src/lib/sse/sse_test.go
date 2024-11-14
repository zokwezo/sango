package sse

import (
	"fmt"
	"log"
	"testing"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func TestUnicodeSSE(t *testing.T) {
	for _, u := range []rune{-1, 65, 127, 255, 256, 50000} {
		expectToken := SSEtoken(max(0, u) % 16384)
		actualSse := MakeUnicodeSSE(u)
		actualToken := actualSse.SSEToken()
		if actualToken != expectToken {
			t.Errorf("for rune %v (%c) found %016b but expected %016b\n",
				u, u, actualToken, expectToken)
		}
		expectSse := actualToken.AsUnicodeSSE()
		actualSseStr := actualSse.FullString()
		expectSseStr := expectSse.FullString()
		if actualSseStr != expectSseStr {
			t.Errorf("for rune %v (%c):\n  actual = %s\n  expect = %s\n",
				u, u, actualSseStr, expectSseStr)
		}
	}
}

func TestAsciiSSE(t *testing.T) {
	for _, isFrench := range []bool{false, true} {
		for _, n := range []int{-1, 0, 1, 2, 7, 31, 32} {
			for _, a := range []rune{-1, 65, 127, 255, 256, 50000} {
				q := uint16(16384)
				if isFrench {
					q += 8192
				}
				q += uint16(min(31, max(0, n))) * 256
				q += uint16(max(0, a)) % 256
				expectToken := SSEtoken(q)
				actualSse := MakeAsciiSSE(a, isFrench, n)
				actualToken := actualSse.SSEToken()
				if actualToken != expectToken {
					t.Errorf("for isFrench=%v n=%v a=%v found %016b but expected %016b\n",
						isFrench, n, a, actualToken, expectToken)
				}
				expectSse := actualToken.AsAsciiSSE()
				actualSseStr := actualSse.FullString()
				expectSseStr := expectSse.FullString()
				if actualSseStr != expectSseStr {
					t.Errorf("for isFrench=%v n=%v a=%v:\n  actual = %s\n  expect = %s\n",
						isFrench, n, a, actualSseStr, expectSseStr)
				}
			}
		}
	}
}

func TestSangoSSE(t *testing.T) {
	for _, s := range []int{0, 1, 2, 3} {
		for _, x := range []uint16{0, 1, 2, 3} {
			for _, c := range []uint16{0, 1, 3, 11, 25} {
				// Don't test uppercase if there is no consonant, since that looks just like titlecase
				if x == 3 && c == 0 {
					continue
				}
				for _, v := range []uint16{1, 3, 7, 13} {
					for _, p := range []uint16{0, 1, 2, 3} {
						q := uint16(1) << 15
						q += uint16(s) << 13
						q += uint16(x) << 11
						q += uint16(c) << 6
						q += uint16(v) << 2
						q += uint16(p)
						expectToken := SSEtoken(q)
						expectSse := expectToken.AsSangoSSE()
						expectSyllable := expectSse.String()
						actualSse, _ := MakeSangoSSE(expectSyllable, s)
						actualToken := actualSse.SSEToken()
						actualSyllable := actualSse.String()
						actualSseStr := actualSse.FullString()
						expectSseStr := expectSse.FullString()
						if actualToken != expectToken {
							t.Errorf("for n=%02b x=%02b c=%05b v=%04b p=%02b\nactual = %016b\nexpect = %016b\n",
								s, x, c, v, p, actualToken, expectToken)
							panic("DONE 1")
						} else if actualSyllable != expectSyllable {
							t.Errorf("for n=%02b x=%02b c=%05b v=%04b p=%02b\nactual = %s\nexpect = %s\n",
								s, x, c, v, p, actualSyllable, expectSyllable)
							panic("DONE 2")
						} else if actualSseStr != expectSseStr {
							t.Errorf("for n=%02b x=%02b c=%05b v=%04b p=%02b\nactual = %s\nexpect = %s\n",
								s, x, c, v, p, actualSseStr, expectSseStr)
							panic("DONE 3")
						}
					}
				}
			}
		}
	}
}

func TestVariousSSE(t *testing.T) {
	assertEqual(t, MakeAsciiSSE('H', false, 2).FullString(),
		`{token=0b_01_0_00010_01001000 UTF8=[72] UTF8/glyph=[[72]] runes=[72] string="H" lang=en numTokensLeft=2}`)
	assertEqual(t, MakeAsciiSSE('i', false, 1).FullString(),
		`{token=0b_01_0_00001_01101001 UTF8=[105] UTF8/glyph=[[105]] runes=[105] string="i" lang=en numTokensLeft=1}`)
	assertEqual(t, MakeAsciiSSE('!', false, 0).FullString(),
		`{token=0b_01_0_00000_00100001 UTF8=[33] UTF8/glyph=[[33]] runes=[33] string="!" lang=en numTokensLeft=0}`)

	sse, w := MakeSangoSSE("Bɛ̂", 1)
	assertEqual(t, w, "")
	assertEqual(t, sse.FullString(),
		`{token=0b_1_01_01_00001_0100_11 UTF8=[66 201 155 204 130] UTF8/glyph=[[66] [201 155 204 130]] runes=[66 603 770] string="Bɛ̂" lang=sg numTokensLeft=1}`)

	sse, w = MakeSangoSSE("-bïn", 0)
	assertEqual(t, w, "")
	assertEqual(t, sse.FullString(),
		`{token=0b_1_00_10_00001_1000_10 UTF8=[45 98 195 175 110] UTF8/glyph=[[45] [98] [195 175] [110]] runes=[45 98 239 110] string="-bïn" lang=sg numTokensLeft=0}`)

	sse, w = MakeSangoSSE("Bə̣", 2)
	assertEqual(t, w, "")
	assertEqual(t, sse.FullString(),
		`{token=0b_1_10_01_00001_0011_00 UTF8=[66 201 153 204 163] UTF8/glyph=[[66] [201 153 204 163]] runes=[66 601 803] string="Bə̣" lang=sg numTokensLeft=2}`)

	sse, w = MakeSangoSSE("BI", 3)
	assertEqual(t, w, "")
	assertEqual(t, sse.FullString(),
		`{token=0b_1_11_11_00001_0111_01 UTF8=[66 73] UTF8/glyph=[[66] [73]] runes=[66 73] string="BI" lang=sg numTokensLeft=3}`)
}

func TestEncodeSangoWord(t *testing.T) {
	s := "\n"
	for k, a := range EncodeSangoWord("Bïkua-ɔ̂kɔ") {
		s += fmt.Sprintf("sse[%v] = %s\n", k, a.FullString())
	}
	assertEqual(t, s, `
sse[0] = {token=0b_1_11_01_00001_0111_10 UTF8=[66 195 175] UTF8/glyph=[[66] [195 175]] runes=[66 239] string="Bï" lang=sg numTokensLeft=3}
sse[1] = {token=0b_1_11_00_00111_1101_01 UTF8=[107 117] UTF8/glyph=[[107] [117]] runes=[107 117] string="ku" lang=sg numTokensLeft=3}
sse[2] = {token=0b_1_10_00_00000_0001_01 UTF8=[97] UTF8/glyph=[[97]] runes=[97] string="a" lang=sg numTokensLeft=2}
sse[3] = {token=0b_1_01_10_00000_1010_11 UTF8=[45 201 148 204 130] UTF8/glyph=[[45] [201 148 204 130]] runes=[45 596 770] string="-ɔ̂" lang=sg numTokensLeft=1}
sse[4] = {token=0b_1_00_00_00111_1010_01 UTF8=[107 201 148] UTF8/glyph=[[107] [201 148]] runes=[107 596] string="kɔ" lang=sg numTokensLeft=0}
`)
}

func assertEqual(t *testing.T, a string, b string) {
	if a != b {
		t.Fatalf("\nactual %s !=\nexpect %s", a, b)
	}
}
