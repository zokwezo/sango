package sse

import (
	"log"
	"testing"
	"unicode/utf8"
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
				for _, v := range []uint16{0, 1, 3, 7, 13} {
					for _, p := range []uint16{0, 1, 2, 3} {
						q := uint16(1) << 15
						q += uint16(s) << 13
						q += uint16(x) << 11
						q += uint16(c) << 6
						q += uint16(v) << 2
						q += uint16(p)
						expectToken := SSEtoken(q)
						var cc, vv []byte
						for _, r := range sseTokenToCons[SSEtoken(c<<6)] {
							cc = utf8.AppendRune(cc, r)
						}
						for _, r := range sseTokenToVowel[SSEtoken(v<<2)] {
							vv = utf8.AppendRune(vv, r)
						}
						actualSse := MakeSangoSSE(x, p, string(cc), string(vv), s)
						actualToken := actualSse.SSEToken()
						expectSse := actualToken.AsSangoSSE()
						actualSseStr := actualSse.FullString()
						expectSseStr := expectSse.FullString()
						if actualToken != expectToken || actualSseStr != expectSseStr {
							t.Errorf("for n=%02b x=%02b c=%05b v=%04b p=%02b\nactual = %016b\nexpect = %016b\n",
								s, x, c, v, p, actualToken, expectToken)
							t.Errorf("for n=%02b x=%02b c=%05b v=%04b p=%02b\nactual = %s\nexpect = %s\n",
								s, x, c, v, p, actualSseStr, expectSseStr)
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
	assertEqual(t, MakeSangoSSE(SangoCaseTitle, SangoPitchHigh, "b", "ɛ", 1).FullString(),
		`{token=0b_1_01_01_01101_0011_11 UTF8=[66 201 155 204 130] UTF8/glyph=[[66] [201 155 204 130]] runes=[66 603 770] string="Bɛ̂" lang=sg numTokensLeft=1}`)
	assertEqual(t, MakeSangoSSE(SangoCaseHyphen, SangoPitchMid, "b", "in", 0).FullString(),
		`{token=0b_1_00_10_01101_1101_10 UTF8=[45 98 195 175 110] UTF8/glyph=[[45] [98] [195 175] [110]] runes=[45 98 239 110] string="-bïn" lang=sg numTokensLeft=0}`)
	assertEqual(t, MakeSangoSSE(SangoCaseUpper, SangoPitchUnknown, "b", "ə", 1).FullString(),
		`{token=0b_1_01_11_01101_1011_00 UTF8=[66 198 143 204 163] UTF8/glyph=[[66] [198 143 204 163]] runes=[66 399 803] string="BƏ̣" lang=sg numTokensLeft=1}`)
	assertEqual(t, MakeSangoSSE(SangoCaseUpper, SangoPitchUnknown, "b", "i", 0).FullString(),
		`{token=0b_1_00_11_01101_0101_00 UTF8=[66 225 187 138] UTF8/glyph=[[66] [225 187 138]] runes=[66 7882] string="BỊ" lang=sg numTokensLeft=0}`)
}

func assertEqual(t *testing.T, a string, b string) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}
