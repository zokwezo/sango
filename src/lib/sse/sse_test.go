package sse

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	// Set global flags for all tests in this package
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(io.Discard)

	// Run the test suite
	os.Exit(m.Run())
}

// TODO: Add test to compare SSEs.
func TestSyllableCompare(t *testing.T) {
	canonicals := [5]string{" =ha_", " ~ha_", "-Do:", "ni^", "HO:"}
	expect := [5][5]int{
		{0, 0, -1, -1, -1},
		{0, 0, -1, -1, -1},
		{1, 1, 0, -1, -1},
		{1, 1, 1, 0, -1},
		{1, 1, 1, 1, 0},
	}
	for l, lhs := range canonicals {
		for r, rhs := range canonicals {
			if actual := CanonicalCompare(lhs, rhs); actual != expect[l][r] {
				t.Errorf("badCanonicalCompare[%v][%v](%q, %q)\nexpected %v but found %v\n",
					l, r, lhs, rhs, expect[l][r], actual)
			}
		}
	}
}

func TestUnpadRight(t *testing.T) {
	log.Println("ENTER TestUnpadRight")
	testCases := [][2]uint64{
		{0x65E5_672C_8A9E_306F, 0x65E5_672C_8A9E_306F},
		{0x0010_96E3_3057_3044, 0x0010_96E3_3057_3044},
		{0x0021_0020_00A7_0000, 0x0000_0021_0020_00A7},
		{0x0010_96E3_3057_0000, 0x0000_0010_96E3_3057},
		{0x0021_0020_0000_0000, 0x0000_0000_0021_0020},
		{0x0010_96E3_0000_0000, 0x0000_0000_0010_96E3},
		{0x0021_0000_0000_0000, 0x0000_0000_0000_0021},
		{0x0000_0000_0000_0000, 0x0000_0000_0000_0000},
		{0xD_10B_089_C31_D95_455, 0xD_10B_089_C31_D95_455},
		{0xF_089_0F6_A72_463_000, 0x000_F_089_0F6_A72_463},
		{0x9_B6E_162_595_000_000, 0x000_000_9_B6E_162_595},
		{0xD_08B_255_000_000_000, 0x000_000_000_D_08B_255},
		{0xE_08B_000_000_000_000, 0x000_000_000_000_E_08B},
		{0xF_000_000_000_000_000, 0x000_000_000_000_000_F},
	}
	for k, testCase := range testCases {
		word := testCase[0]
		expect := testCase[1]
		actual := unpadRight(word)
		if actual != expect {
			t.Errorf("bad UnpadRight[%v](%#016X)\nexpect: %#016X\nactual: %#016X\n", k, word, expect, actual)
		}
		// check idempotence
		actual2 := unpadRight(actual)
		if actual2 != actual {
			t.Errorf("bad UnpadRight[%v](%#016X)\nactual2: %#016X\nexpect: %#016X\n", k, actual, actual2, expect)
		}
		// check roundtrip
		word2 := padRight(expect)
		if word2 != word {
			t.Errorf("bad PadRight[%v](%#016X)\nword2: %#016X\nword: %#016X\n", k, expect, word2, word)
		}
	}
	log.Println("LEAVE TestUnpadRight")
}

func TestPadRight(t *testing.T) {
	log.Println("ENTER TestPadRight")
	testCases := [][2]uint64{
		{0x65E5_672C_8A9E_306F, 0x65E5_672C_8A9E_306F},
		{0x0010_96E3_3057_3044, 0x0010_96E3_3057_3044},
		{0x0000_0021_0020_00A7, 0x0021_0020_00A7_0000},
		{0x0000_0010_96E3_3057, 0x0010_96E3_3057_0000},
		{0x0000_0000_0021_0020, 0x0021_0020_0000_0000},
		{0x0000_0000_0010_96E3, 0x0010_96E3_0000_0000},
		{0x0000_0000_0000_0021, 0x0021_0000_0000_0000},
		{0x0000_0000_0000_0000, 0x0000_0000_0000_0000},
		{0xD_10B_089_C31_D95_455, 0xD_10B_089_C31_D95_455},
		{0x000_F_089_0F6_A72_463, 0xF_089_0F6_A72_463_000},
		{0x000_000_9_B6E_162_595, 0x9_B6E_162_595_000_000},
		{0x000_000_000_D_08B_255, 0xD_08B_255_000_000_000},
		{0x000_000_000_000_E_08B, 0xE_08B_000_000_000_000},
		{0x000_000_000_000_000_F, 0xF_000_000_000_000_000},
	}
	for k, testCase := range testCases {
		word := testCase[0]
		expect := testCase[1]
		actual := padRight(word)
		if actual != expect {
			t.Errorf("bad PadRight[%v](%#016X)\nexpect: %#016X\nactual: %#016X\n", k, word, expect, actual)
		}
		// check idempotence
		actual2 := padRight(actual)
		if actual2 != actual {
			t.Errorf("bad PadRight[%v](%#016X)\nactual2: %#016X\nexpect: %#016X\n", k, actual, actual2, expect)
		}
		// check roundtrip
		word2 := unpadRight(expect)
		if word2 != word {
			t.Errorf("bad UnpadRight[%v](%#016X)\nword2: %#016X\nword: %#016X\n", k, expect, word2, word)
		}
	}
	log.Println("LEAVE TestPadRight")
}

func TestGetShortCode(t *testing.T) {
	log.Println("ENTER TestGetShortCode")
	testCases := [][2]uint64{
		{0x65E5_672C_8A9E_306F, 0x65E5_672C_8A9E_306F},
		{0x0010_96E3_3057_3044, 0x0010_96E3_3057_3044},
		{0x0021_0020_00A7_0000, 0x0000_0021_0020_00A7},
		{0x0010_96E3_3057_0000, 0x0000_0010_96E3_3057},
		{0x0021_0020_0000_0000, 0x0000_0000_0021_0020},
		{0x0010_96E3_0000_0000, 0x0000_0000_0010_96E3},
		{0x0021_0000_0000_0000, 0x0000_0000_0000_0021},
		{0x0000_0000_0000_0000, 0x0000_0000_0000_0000},
		{0xD_10B_089_C31_D95_455, 0xD_10B_089_C31_D95_455},
		{0xF_089_0F6_A72_463_000, 0x000_F_089_0F6_A72_463},
		{0x9_B6E_162_595_000_000, 0x000_000_9_B6E_162_595},
		{0xD_08B_255_000_000_000, 0x000_000_000_D_08B_255},
		{0xE_08B_000_000_000_000, 0x000_000_000_000_E_08B},
		{0xF_000_000_000_000_000, 0x000_000_000_000_000_F},
	}
	for k, testCase := range testCases {
		sse := SSE(testCase[0])
		expect := testCase[1]
		actual := sse.GetShortCode()
		if actual != expect {
			t.Errorf("bad GetShortCode[%v](%#016X)\nexpect: %#016X\nactual: %#016X\n", k, sse, expect, actual)
		}
		// check roundtrip
		sse2 := FromShortCode(expect)
		if sse2 != sse {
			t.Errorf("bad PadRight[%v](%#016X)\nsse2: %#016X\nsse: %#016X\n", k, expect, sse2, sse)
		}
	}
	log.Println("LEAVE TestGetShortCode")
}

func TestFromShortCode(t *testing.T) {
	log.Println("ENTER TestFromShortCode")
	testCases := [][2]uint64{
		{0x65E5_672C_8A9E_306F, 0x65E5_672C_8A9E_306F},
		{0x0010_96E3_3057_3044, 0x0010_96E3_3057_3044},
		{0x0000_0021_0020_00A7, 0x0021_0020_00A7_0000},
		{0x0000_0010_96E3_3057, 0x0010_96E3_3057_0000},
		{0x0000_0000_0021_0020, 0x0021_0020_0000_0000},
		{0x0000_0000_0010_96E3, 0x0010_96E3_0000_0000},
		{0x0000_0000_0000_0021, 0x0021_0000_0000_0000},
		{0x0000_0000_0000_0000, 0x0000_0000_0000_0000},
		{0xD_10B_089_C31_D95_455, 0xD_10B_089_C31_D95_455},
		{0x000_F_089_0F6_A72_463, 0xF_089_0F6_A72_463_000},
		{0x000_000_9_B6E_162_595, 0x9_B6E_162_595_000_000},
		{0x000_000_000_D_08B_255, 0xD_08B_255_000_000_000},
		{0x000_000_000_000_E_08B, 0xE_08B_000_000_000_000},
		{0x000_000_000_000_000_F, 0xF_000_000_000_000_000},
	}
	for k, testCase := range testCases {
		word := testCase[0]
		expect := SSE(testCase[1])
		actual := FromShortCode(word)
		if actual != expect {
			t.Errorf("bad PadRight[%v](%#016X)\nexpect: %#016X\nactual: %#016X\n", k, word, expect, actual)
		}
		// check roundtrip
		word2 := expect.GetShortCode()
		if word2 != word {
			t.Errorf("bad UnpadRight[%v](%#016X)\nword2: %#016X\nword: %#016X\n", k, expect, word2, word)
		}
	}
	log.Println("LEAVE TestFromShortCode")
}

func TestUnpadFollowedByPad(t *testing.T) {
	log.Println("ENTER TestUnpadFollowedByPad")
	msbs := []uint64{
		0x0010_9ED2_0000_0000,
		0x002F_0000_0000_0000,
		0x0115_0000_0000_0000,
		0x5143_0000_0000_0000,
		0x002F_002F_0000_0000,
		0x0115_0115_0000_0000,
		0x5143_5143_0000_0000,
		0x5143_9ED2_0000_0000,
	}
	lsbs := []uint64{
		0x0000_0000_0000_0000,
		0x0000_0000_0000_002F,
		0x0000_0000_0000_0115,
		0x0000_0000_0000_5143,
		0x0000_0000_0000_9ED2,
		0x0000_0000_0010_9ED2,
		0x0000_0000_002F_0000,
		0x0000_0000_0115_0000,
		0x0000_0000_5143_0000,
		0x0000_0000_9ED2_0000,
		0x0000_0000_002F_002F,
		0x0000_0000_0115_0115,
		0x0000_0000_5143_5143,
		0x0000_0000_5143_9ED2,
		0x0000_0000_9ED2_9ED2,
	}
	for _, msb := range msbs {
		for _, lsb := range lsbs {
			word := msb | lsb
			unpadded := unpadRight(word)
			padded := padRight(unpadded)
			reunpadded := unpadRight(padded)
			if padded != word {
				t.Errorf("\n    word = %#016X  unpadded = %#016X\n  padded = %#016X\n", word, unpadded, padded)
				return
			}
			if reunpadded != unpadded {
				t.Errorf("\n    word = %#016X  padded = %#016X\n  reunpadded = %#016X\n", unpadded, padded, reunpadded)
				return
			}
		}
	}
	log.Println("LEAVE TestUnpadFollowedByPad")
}

//////////////////////////////////////////////////////////////////////////////

func TestEmptyCanonicalToCodes(t *testing.T) {
	log.Println("ENTER TestEmptyCanonicalToCodes")
	actual, leftOffAt := canonicalToCodes("")
	if leftOffAt != 0 {
		t.Errorf("bad leftOffAt")
	}
	if len(actual) != 0 {
		t.Errorf("found nonempty codes")
	}
	log.Println("LEAVE TestEmptyCanonicalToCodes")
}

func TestCanonicalToCodes(t *testing.T) {
	log.Println("ENTER TestCanonicalToCodes")
	c := `U+65E5U+672CU+8A9EU+306FU+96E3U+3057U+3044U+0021U+0020U+00A7 ~ha_HO:-Do:ni^ ` +
		`=ha_HO:-Do:ni^ ha^Dx_ ba^ha_-mo_-tx_nx_ ~bx^-kc:Bi:tx_bx^-kc:Bi:tx_U+96E3` +
		`=bx^-=kc:=Bi:=tx_ bx^-kc:Bi:tx_ ~bx^-kc:Bi:tx_ =bx^-=kc:=Bi:=tx_ ha_HO:Do:ni^`
	u := func(v uint16) sseCode { return sseCode{value: v, isSango: false} }
	s := func(v uint16) sseCode { return sseCode{value: v, isSango: true} }
	expect := []sseCode{
		u(0x65E5), u(0x672C), u(0x8A9E), u(0x306F), u(0x96E3),
		u(0x3057), u(0x3044), u(0x0021), u(0x0020), u(0x00A7),
		s(0xE089), s(0x9236), s(0x9C72), s(0x9423), s(0xF089),
		s(0x9236), s(0x9C72), s(0x9423), s(0xD08B), s(0x9455),
		s(0xD0CB), s(0x9089), s(0x9B31), s(0x9E55), s(0x9415),
		s(0xE0D7), s(0x9A6E), s(0x9362), s(0x9655), s(0x90D7),
		s(0x9A6E), s(0x9362), s(0x9655),
		u(0x96E3),
		s(0xB0D7), s(0xBA6E), s(0xB362), s(0xB655), s(0xD0D7),
		s(0x9A6E), s(0x9362), s(0x9655), s(0xE0D7), s(0x9A6E),
		s(0x9362), s(0x9655), s(0xF0D7), s(0xBA6E), s(0xB362),
		s(0xB655), s(0xD089), s(0x9236), s(0x9472), s(0x9423),
	}
	expectLen := len(c)
	actual, actualLen := canonicalToCodes(c)
	if actualLen != expectLen {
		t.Errorf("in TestCanonicalToCodes: error found at c[%v](good: %q bad: %q\n", actualLen, c[0:actualLen], c[actualLen:])
	}
	var actualHex string
	for _, x := range actual {
		if x.isSango {
			actualHex += fmt.Sprintf("S+%04X ", x.value)
		} else {
			actualHex += fmt.Sprintf("%U ", x.value)
		}
	}
	var expectHex string
	for _, x := range expect {
		if x.isSango {
			expectHex += fmt.Sprintf("S+%04X ", x.value)
		} else {
			expectHex += fmt.Sprintf("%U ", x.value)
		}
	}
	if actualHex != expectHex {
		t.Errorf("in TestCanonicalToCodes: bad canonicalToCodes\nexpect: %v\nactual: %v\n", expectHex, actualHex)
	}
	log.Println("LEAVE TestCanonicalToCodes")
}

func TestGoodCanonicalToSSEs(t *testing.T) {
	log.Println("ENTER TestGoodCanonicalToSSEs")
	c := `U+65E5U+672CU+8A9EU+306FU+96E3U+3057U+3044U+0021U+0020U+00A7 ~ha_HO:-Do:ni^` +
		` =ha_HO:-Do:ni^ ha^Dx_ ba^ha_-mo_-tx_nx_ ~bx^-kc:Bi:tx_bx^-kc:Bi:tx_U+96E3` +
		`bx^-=kc:=Bi:=tx_ bx^-kc:Bi:tx_ ~bx^-kc:Bi:tx_ ~bx^-=kc:=Bi:=tx_ ha_HO:Do:ni^`
	expect := SSEs{
		0x65E5_672C_8A9E_306F, 0x0010_96E3_3057_3044, 0x0021_0020_00A7_0000,
		0xE_089_236_C72_423_000, 0xF_089_236_C72_423_000, 0xD_08B_455_000_000_000,
		0xD_0CB_089_B31_E55_415, 0xE_0D7_A6E_362_655_0D7, 0x9_A6E_362_655_000_000,
		0x0_010_96E_300_000_000, 0xB_0D7_A6E_362_655_000, 0xD_0D7_A6E_362_655_000,
		0xE_0D7_A6E_362_655_000, 0xF_0D7_A6E_362_655_000, 0xD_089_236_472_423_000,
	}
	dumpSSEs := func(sses SSEs) string {
		s := "{"
		for _, sse := range sses {
			s += fmt.Sprintf(" %016X", sse)
		}
		s += " }"
		return s
	}
	actual, err := CanonicalToSSEs(c)
	if err != nil {
		t.Errorf("unexpected error returned from CanonicalToSSEs\nerr = %v", err)
	}
	actualHex := dumpSSEs(actual)
	expectHex := dumpSSEs(expect)
	if actualHex != expectHex {
		t.Errorf("bad GoodCanonicalToSSEs\nexpect: %v\nactual: %v\n", expectHex, actualHex)
	}
	log.Println("LEAVE TestGoodCanonicalToSSEs")
}

func TestBadCanonicalToSSEs(t *testing.T) {
	log.Println("ENTER TestBadCanonicalToSSEs")
	c := `U+65E5U+672CU+8A9EU+306FU+96E3U+3057U+3044U+0021U+0020U+00A7` +
		` ~ha_HO:-Do:ni^ =ha_HO:-jo:ni^ ha^Dx_ ba^ha_-mo_-tx_nx_ ~bx^-kc:Bi:tx_bx^-kc:Bi:tx_` +
		`U+96E3=bx^-=kc:=Bi:=tx_ bx^-kc:Bi:tx_ ~bx^-kc:Bi:tx_ =bx^-=kc:=Bi:=tx_ ha_HO:Do:ni^`
	expect := SSEs{
		0x65E5_672C_8A9E_306F, 0x0010_96E3_3057_3044, 0x0021_0020_00A7_0000,
		0xE_089_236_C72_423_000, 0xF_089_236_000_000_000,
	} // results up to broken parse
	dumpSSEs := func(sses SSEs) string {
		s := "{"
		for _, sse := range sses {
			s += fmt.Sprintf(" %016X", sse)
		}
		s += " }"
		return s
	}
	expectErr := fmt.Errorf("cannot parse string starting at s[83:] = %q", "-jo:ni^ ha...")
	actual, actualErr := CanonicalToSSEs(c)
	if actualErr.Error() != expectErr.Error() {
		t.Errorf("expected error not returned from CanonicalToSSEs\nactualErr = %v\nexpectErr = %v", actualErr, expectErr)
	}
	actualHex := dumpSSEs(actual)
	expectHex := dumpSSEs(expect)
	if actualHex != expectHex {
		t.Errorf("bad BadCanonicalToSSEs\nexpect: %v\nactual: %v\n", expectHex, actualHex)
	}
	log.Println("LEAVE TestBadCanonicalToSSEs")
}

func TestCanonicalToSSEsForUnknownPitch(t *testing.T) {
	log.Println("ENTER TestCanonicalToSSEsForUnknownPitch")
	c := `~haHO-Doni =haHO-Doni haDx baha-mo-txnx ~bx-kcBitxbx-kcBitxU+96E3` +
		`=bx-=kc=Bi=tx bx-kcBitx ~bx-kcBitx =bx-=kc=Bi=tx haHODoni`
	expect := SSEs{
		0xA_088_234_C70_420_000, 0xF_088_234_C70_420_000, 0xD_088_454_000_000_000,
		0xD_0C8_088_B30_E54_414, 0xE_0D4_A6C_360_654_0D4, 0x9_A6C_360_654_000_000,
		0x0_010_96E_300_000_000, 0xB_0D4_A6C_360_654_000, 0xD_0D4_A6C_360_654_000,
		0xE_0D4_A6C_360_654_000, 0xF_0D4_A6C_360_654_000, 0xD_088_234_470_420_000,
	}
	dumpSSEs := func(sses SSEs) string {
		s := "{"
		for _, sse := range sses {
			s += fmt.Sprintf(" %016X", sse)
		}
		s += " }"
		return s
	}
	actual, err := CanonicalToSSEs(c)
	if err != nil {
		t.Errorf("unexpected error returned from CanonicalToSSEs\nerr = %v", err)
	}
	actualHex := dumpSSEs(actual)
	expectHex := dumpSSEs(expect)
	if actualHex != expectHex {
		t.Errorf("bad GoodCanonicalToSSEs\nexpect: %v\nactual: %v\n", expectHex, actualHex)
	}
	log.Println("LEAVE TestCanonicalToSSEsForUnknownPitch")
}

func TestCanonicalToSSEsToCanonical(t *testing.T) {
	log.Println("ENTER TestCanonicalToSSEsToCanonical")
	c := `~ha~HO-Doni =haHO-Doni haDx baha-mo-txnx ~bx-kcBitxbx-kcBitxU+96E3` +
		`=bx-=kc=Bi=tx bx-kcBitx ~bx-kcBitx =bx-=kc=Bi=tx haHODoni`
	sses, err := CanonicalToSSEs(c)
	if err != nil {
		t.Errorf("unexpected error returned from CanonicalToSSEs\nerr = %v", err)
		return
	}
	var s strings.Builder
	for _, sse := range sses {
		sse.WriteAsCanonicalTo(&s)
	}
	// expect c, but with Titlecase suppressed for nonleading syllables.
	expect := `~haHO-Doni =ha=HO-=Do=ni haDx baha-mo-txnx ~bx-kcBitx` +
		`bx-kcBitxU+96E3=bx-=kc=Bi=tx bx-kcBitx ~bx-kcBitx =bx-=kc=Bi=tx haHODoni`
	actual := s.String()
	if actual != expect {
		t.Errorf("bad BadCanonicalToSSEs\nexpect: %v\nactual: %v\n", expect, actual)
	}
	log.Println("LEAVE TestCanonicalToSSEsToCanonical")
}

//////////////////////////////////////////////////////////////////////////////

func TestEmptyUtf8ToCodes(t *testing.T) {
	log.Println("ENTER TestEmptyUtf8ToCodes")
	actual, leftOffAt := utf8ToCodes("", FromLemma)
	if leftOffAt != 0 {
		t.Errorf("bad leftOffAt")
	}
	if len(actual) != 0 {
		t.Errorf("found nonempty codes")
	}
	log.Println("LEAVE TestEmptyUtf8ToCodes")
}

func TestUtf8ToCodes(t *testing.T) {
	log.Println("ENTER TestUtf8ToCodes")
	c := `日本語は難しい! § Ahöñ-ndönî AHÖÑ-NDÖNÎ ândɛ bâa-mo-tɛnɛ` +
		` BƐ̂-kɔ̈mbïtɛbɛ̂-kɔ̈mbïtɛ難BƐ̂-KƆ̈MBÏTƐ bɛ̂-kɔ̈mbïtɛ BƐ̂-kɔ̈mbïtɛ BƐ̂-KƆ̈MBÏTƐ ahöñndönî`
	u := func(v uint16) sseCode { return sseCode{value: v, isSango: false} }
	s := func(v uint16) sseCode { return sseCode{value: v, isSango: true} }
	expect := []sseCode{
		u(0x65E5), u(0x672C), u(0x8A9E), u(0x306F), u(0x96E3), // `日本語は難`
		u(0x3057), u(0x3044), u(0x0021), u(0x0020), u(0x00A7), // `しい! §`
		s(0xE089), s(0x9236), s(0x9C72), s(0x9423), // ` Ahöñ-ndönî`
		s(0xE089), s(0xB236), s(0xBC72), s(0xB423), // ` AHÖÑ-NDÖNÎ`
		s(0xD08B), s(0x9455), // ` ândɛ`
		s(0xD0CB), s(0x9089), s(0x9B31), s(0x9E55), s(0x9415), // ` bâa-mo-tɛnɛ`
		s(0xF0D7), s(0x9A6E), s(0x9362), s(0x9655), // ` BƐ̂-kɔ̈mbïtɛb`
		s(0x90D7), s(0x9A6E), s(0x9362), s(0x9655), //  `bɛ̂-kɔ̈mbïtɛ`
		u(0x96E3),                                  //  `難`
		s(0xB0D7), s(0xBA6E), s(0xB362), s(0xB655), //  `BƐ̂-KƆ̈MBÏTƐ`
		s(0xD0D7), s(0x9A6E), s(0x9362), s(0x9655), // ` bɛ̂-kɔ̈mbïtɛ`
		s(0xF0D7), s(0x9A6E), s(0x9362), s(0x9655), // ` BƐ̂-kɔ̈mbïtɛ`
		s(0xF0D7), s(0xBA6E), s(0xB362), s(0xB655), // ` BƐ̂-KƆ̈MBÏTƐ`
		s(0xD089), s(0x9236), s(0x9472), s(0x9423), // ` ahöñndönî`
	}
	expectLen := len(c)
	actual, actualLen := utf8ToCodes(c, FromLemma)
	if actualLen != expectLen {
		t.Errorf("in TestUtf8ToCodes: error found at c[%v](good: %q bad: %q\n", actualLen, c[0:actualLen], c[actualLen:])
	}
	var actualHex string
	for _, x := range actual {
		if x.isSango {
			actualHex += fmt.Sprintf("S+%04X ", x.value)
		} else {
			actualHex += fmt.Sprintf("%U ", x.value)
		}
	}
	var expectHex string
	for _, x := range expect {
		if x.isSango {
			expectHex += fmt.Sprintf("S+%04X ", x.value)
		} else {
			expectHex += fmt.Sprintf("%U ", x.value)
		}
	}
	if actualHex != expectHex {
		t.Errorf("in TestUtf8ToCodes: bad utf8ToCodes\nexpect: %v\nactual: %v\n", expectHex, actualHex)
	}
	log.Println("LEAVE TestUtf8ToCodes")
}

func TestGoodUtf8ToSSEs(t *testing.T) {
	log.Println("ENTER TestGoodUtf8ToSSEs")
	c := `日本語は難しい! § Ahöñ-ndönî AHÖÑ-NDÖNÎ ândɛ bâa-mo-tɛnɛ` +
		` BƐ̂-kɔ̈mbïtɛbɛ̂-kɔ̈mbïtɛ難BƐ̂-KƆ̈MBÏTƐ bɛ̂-kɔ̈mbïtɛ BƐ̂-kɔ̈mbïtɛ BƐ̂-KƆ̈MBÏTƐ ahöñndönî`
	expect := SSEs{
		0x65E5_672C_8A9E_306F, 0x0010_96E3_3057_3044, 0x0021_0020_00A7_0000, // `日本語は難しい! §`
		0xE_089_236_C72_423_000, // ` Ahöñ-ndönî`
		0xF_089_236_C72_423_000, // ` AHÖÑ-NDÖNÎ`
		0xD_08B_455_000_000_000, // ` ândɛ`
		0xD_0CB_089_B31_E55_415, // ` bâa-mo-tɛnɛ`
		0xF_0D7_A6E_362_655_0D7, // ` BƐ̂-kɔ̈mbïtɛbɛ̂`  [if any syllable is UPPER, the whole word is]
		0x9_A6E_362_655_000_000, //    `-kɔ̈mbïtɛ`
		0x0010_96E3_0000_0000,   //  `難`
		0xB_0D7_A6E_362_655_000, //  `BƐ̂-KƆ̈MBÏTƐ`
		0xD_0D7_A6E_362_655_000, // ` bɛ̂-kɔ̈mbïtɛ`
		0xF_0D7_A6E_362_655_000, // ` BƐ̂-kɔ̈mbïtɛ`
		0xF_0D7_A6E_362_655_000, // ` BƐ̂-KƆ̈MBÏTƐ`
		0xD_089_236_472_423_000, // ` ahöñndönî`
	}
	dumpSSEs := func(sses SSEs) string {
		s := "{"
		for _, sse := range sses {
			s += fmt.Sprintf(" %016X", sse)
		}
		s += " }"
		return s
	}
	actual, err := Utf8ToSSEs(c, FromLemma)
	if err != nil {
		t.Errorf("unexpected error returned from Utf8ToSSEs\nerr = %v", err)
	}
	actualHex := dumpSSEs(actual)
	expectHex := dumpSSEs(expect)
	if actualHex != expectHex {
		t.Errorf("bad GoodUtf8ToSSEs\nexpect: %v\nactual: %v\n", expectHex, actualHex)
	}
	log.Println("LEAVE TestGoodUtf8ToSSEs")
}

func TestBadUtf8ToSSEs(t *testing.T) {
	log.Println("ENTER TestBadUtf8ToSSEs")
	c := `日本語𓋹はしい!`
	expect := SSEs{0x65E5_672C_8A9E_0000} // results up to broken parse
	dumpSSEs := func(sses SSEs) string {
		s := "{"
		for _, sse := range sses {
			s += fmt.Sprintf(" %016X", sse)
		}
		s += " }"
		return s
	}
	expectErr := fmt.Errorf(`cannot parse string starting at s[9:] = "𓋹はし..."`)
	actual, actualErr := Utf8ToSSEs(c, FromLemma)
	if actualErr == nil || actualErr.Error() != expectErr.Error() {
		t.Errorf("expected error not returned from Utf8ToSSEs\nactualErr = %v\nexpectErr = %v", actualErr, expectErr)
	}
	actualHex := dumpSSEs(actual)
	expectHex := dumpSSEs(expect)
	if actualHex != expectHex {
		t.Errorf("bad BadUtf8ToSSEs\nexpect: %v\nactual: %v\n", expectHex, actualHex)
	}
	log.Println("LEAVE TestBadUtf8ToSSEs")
}

func TestUtf8ToSSEsFromToneless(t *testing.T) {
	log.Println("ENTER TestUtf8ToSSEsFromToneless")
	c := `Ahonndoni AHONNDONI ande baamotene` +
		` BE-kombitebe-kombite難BƐ̂-KƆ̈MBÏTƐ bɛ̂-kɔ̈mbïtɛ BƐ̂-kɔ̈mbïtɛ BƐ̂-KƆ̈MBÏTƐ ahöñndönî`
	expect := SSEs{
		0xA_088_234_470_420_000, 0xF_088_234_470_420_000, 0xD_088_458_000_000_000,
		0xD_0C8_088_330_658_418, 0xF_0D8_A70_360_658_0D8, 0x9_A70_360_658_000_000,
		0x0_010_96E_300_000_000, 0xB_0D7_A6E_362_654_000, 0xD_0D7_A6E_362_654_000,
		0xF_0D7_A6E_362_654_000, 0xF_0D7_A6E_362_654_000, 0xD_088_236_472_423_000,
	}
	dumpSSEs := func(sses SSEs) string {
		s := "{"
		for _, sse := range sses {
			s += fmt.Sprintf(" %016X", sse)
		}
		s += " }"
		return s
	}
	actual, err := Utf8ToSSEs(c, FromToneless)
	if err != nil {
		t.Errorf("unexpected error returned from Utf8ToSSEs\nerr = %v", err)
	}
	actualHex := dumpSSEs(actual)
	expectHex := dumpSSEs(expect)
	if actualHex != expectHex {
		log.Println("bad GoodUtf8ToSSEs")
		t.Errorf("bad GoodUtf8ToSSEs\nexpect: %v\nactual: %v\n", expectHex, actualHex)
	} else {
		log.Println("good GoodUtf8ToSSEs")
	}
	log.Println("LEAVE TestUtf8ToSSEsFromToneless")
}

func TestUtf8ToSSEsToUtf8(t *testing.T) {
	log.Println("ENTER TestUtf8ToSSEsToUtf8")
	c := `Ahöñ-ndönî AHÖÑ-NDÖNÎ ândɛ bâa-mo-tɛnɛ` +
		` BƐ̂-kɔ̈mbïtɛbɛ̂-kɔ̈mbïtɛ難BƐ̂-KƆ̈MBÏTƐ bɛ̂-kɔ̈mbïtɛ BƐ̂-kɔ̈mbïtɛ BƐ̂-KƆ̈MBÏTƐ ahöñndönî`
	sses, err := Utf8ToSSEs(c, FromLemma)
	if err != nil {
		t.Errorf("unexpected error returned from Utf8ToSSEs\nerr = %v", err)
		return
	}
	var s strings.Builder
	for _, sse := range sses {
		sse.WriteAsUtf8To(&s)
	}
	expect := `Ahöñ-ndönî AHÖÑ-NDÖNÎ ândɛ bâa-mo-tɛnɛ` +
		` BƐ̂-KƆ̈MBÏTƐBƐ̂-kɔ̈mbïtɛ難BƐ̂-KƆ̈MBÏTƐ bɛ̂-kɔ̈mbïtɛ BƐ̂-KƆ̈MBÏTƐ BƐ̂-KƆ̈MBÏTƐ ahöñndönî`
	actual := s.String()
	if actual != expect {
		t.Errorf("bad BadUtf8ToSSEs\nexpect: %v\nactual: %v\n", expect, actual)
	}
	log.Println("LEAVE TestUtf8ToSSEsToUtf8")
}

//////////////////////////////////////////////////////////////////////////////

func TestCodesToSSEs(t *testing.T) {
	log.Println("ENTER TestCodesToSSEs")
	u := func(v uint16) sseCode { return sseCode{value: v, isSango: false} }
	s := func(v uint16) sseCode { return sseCode{value: v, isSango: true} }
	codes := []sseCode{
		u(0x65E5), u(0x672C), u(0x8A9E), u(0x306F), u(0x96E3), // `日本語は難`
		u(0x3057), u(0x3044), u(0x0021), u(0x0020), u(0x00A7), // `しい! §`
		s(0xE089), s(0x9236), s(0x9C72), s(0x9423), // ` Ahöñ-ndönî`
		s(0xE089), s(0xB236), s(0xBC72), s(0xB423), // ` AHÖÑ-NDÖNÎ`
		s(0xD08B), s(0x9455), // ` ândɛ`
		s(0xD0CB), s(0x9089), s(0x9B31), s(0x9E55), s(0x9415), // ` bâa-mo-tɛnɛ`
		s(0xF0D7), s(0x9A6E), s(0x9362), s(0x9655), // ` BƐ̂-kɔ̈mbïtɛb`
		s(0x90D7), s(0x9A6E), s(0x9362), s(0x9655), //  `bɛ̂-kɔ̈mbïtɛ`
		u(0x96E3),                                  //  `難`
		s(0xB0D7), s(0xBA6E), s(0xB362), s(0xB655), //  `BƐ̂-KƆ̈MBÏTƐ`
		s(0xD0D7), s(0x9A6E), s(0x9362), s(0x9655), // ` bɛ̂-kɔ̈mbïtɛ`
		s(0xF0D7), s(0x9A6E), s(0x9362), s(0x9655), // ` BƐ̂-kɔ̈mbïtɛ`
		s(0xF0D7), s(0xBA6E), s(0xB362), s(0xB655), // ` BƐ̂-KƆ̈MBÏTƐ`
		s(0xD089), s(0x9236), s(0x9472), s(0x9423), // ` ahöñndönî`
	}
	expect := SSEs{
		0x65E5_672C_8A9E_306F, 0x0010_96E3_3057_3044, 0x0021_0020_00A7_0000, // `日本語は難しい! §`
		0xE_089_236_C72_423_000, // ` Ahöñ-ndönî`
		0xF_089_236_C72_423_000, // ` AHÖÑ-NDÖNÎ`
		0xD_08B_455_000_000_000, // ` ândɛ`
		0xD_0CB_089_B31_E55_415, // ` bâa-mo-tɛnɛ`
		0xF_0D7_A6E_362_655_0D7, // ` BƐ̂-kɔ̈mbïtɛbɛ̂`  [if any syllable is UPPER, the whole word is]
		0x9_A6E_362_655_000_000, //    `-kɔ̈mbïtɛ`
		0x0010_96E3_0000_0000,   //  `難`
		0xB_0D7_A6E_362_655_000, //  `BƐ̂-KƆ̈MBÏTƐ`
		0xD_0D7_A6E_362_655_000, // ` bɛ̂-kɔ̈mbïtɛ`
		0xF_0D7_A6E_362_655_000, // ` BƐ̂-kɔ̈mbïtɛ`
		0xF_0D7_A6E_362_655_000, // ` BƐ̂-KƆ̈MBÏTƐ`
		0xD_089_236_472_423_000, // ` ahöñndönî`
	}
	dumpSSEs := func(sses SSEs) string {
		s := "{"
		for _, sse := range sses {
			s += fmt.Sprintf(" %016X", sse)
		}
		s += " }"
		return s
	}
	actual := codesToSSEs(codes)
	actualHex := dumpSSEs(actual)
	expectHex := dumpSSEs(expect)
	if actualHex != expectHex {
		t.Errorf("bad codesToSSEs\nexpect: %v\nactual: %v\n", expectHex, actualHex)
	}
	log.Println("LEAVE TestCodesToSSEs")
}

func TestWriteAsUtf8MixedKind(t *testing.T) {
	log.Println("ENTER TestWriteAsUtf8MixedKind")
	sses := SSEs{
		0x65E5_672C_8A9E_306F, 0x0010_96E3_3057_3044, 0x0021_0020_00A7_0000, // `日本語は難しい! §`
		0xE_089_236_C72_423_000, // ` Ahöñ-ndönî`
		0xF_089_236_C72_423_000, // ` AHÖÑ-NDÖNÎ`
		0xD_08B_455_000_000_000, // ` ândɛ`
		0xD_0CB_089_B31_E55_415, // ` bâa-mo-tɛnɛ`
		0xF_0D7_A6E_362_655_0D7, // ` BƐ̂-kɔ̈mbïtɛbɛ̂`  [if any syllable is UPPER, the whole word is]
		0x9_A6E_362_655_000_000, //    `-kɔ̈mbïtɛ`
		0x0010_96E3_0000_0000,   //  `難`
		0xB_0D7_A6E_362_655_000, //  `BƐ̂-KƆ̈MBÏTƐ`
		0xD_0D7_A6E_362_655_000, // ` bɛ̂-kɔ̈mbïtɛ`
		0xF_0D7_A6E_362_655_000, // ` BƐ̂-kɔ̈mbïtɛ`
		0xF_0D7_A6E_362_655_000, // ` BƐ̂-KƆ̈MBÏTƐ`
		0xD_089_236_472_423_000, // ` ahöñndönî`
	}
	var s strings.Builder
	for _, sse := range sses {
		s.WriteString("|")
		sse.WriteAsUtf8To(&s)
	}
	s.WriteString("|")
	expect := `|日本語は|難しい|! §| Ahöñ-ndönî| AHÖÑ-NDÖNÎ` +
		`| ândɛ| bâa-mo-tɛnɛ| BƐ̂-KƆ̈MBÏTƐBƐ̂|-kɔ̈mbïtɛ|難|BƐ̂-KƆ̈MBÏTƐ` +
		`| bɛ̂-kɔ̈mbïtɛ| BƐ̂-KƆ̈MBÏTƐ| BƐ̂-KƆ̈MBÏTƐ| ahöñndönî|`
	actual := s.String()
	if actual != expect {
		t.Errorf("in TestWriteAsUtf8MixedKindTo(\n%#v\n),\nexpect: %#v\nactual: %#v\n\n",
			sses, expect, actual)
	}
	log.Println("LEAVE TestWriteAsUtf8MixedKind")
}

func TestWriteAsTonelessMixedKind(t *testing.T) {
	log.Println("ENTER TestWriteAsTonelessMixedKind")
	sses := SSEs{
		0x65E5_672C_8A9E_306F, 0x0010_96E3_3057_3044, 0x0021_0020_00A7_0000, // `日本語は難しい! §`
		0xE_089_236_C72_423_000, // ` Ahöñ-ndönî`
		0xF_089_236_C72_423_000, // ` AHÖÑ-NDÖNÎ`
		0xD_08B_455_000_000_000, // ` ândɛ`
		0xD_0CB_089_B31_E55_415, // ` bâa-mo-tɛnɛ`
		0xE_0D7_A6E_362_655_0D7, // ` Bɛ̂-kɔ̈mbïtɛbɛ̂`  [if any syllable is UPPER, the whole word is]
		0x9_A6E_362_655_000_000, //    `-kɔ̈mbïtɛ`
		0x0010_96E3_0000_0000,   //  `難`
		0xB_0D7_A6E_362_655_000, //  `BƐ̂-KƆ̈MBÏTƐ`
		0xD_0D7_A6E_362_655_000, // ` bɛ̂-kɔ̈mbïtɛ`
		0xF_0D7_A6E_362_655_000, // ` BƐ̂-kɔ̈mbïtɛ`
		0xF_0D7_A6E_362_655_000, // ` BƐ̂-KƆ̈MBÏTƐ`
		0xD_089_236_472_423_000, // ` ahöñndönî`
	}
	var s strings.Builder
	for _, sse := range sses {
		s.WriteString("|")
		sse.WriteAsTonelessTo(&s)
	}
	s.WriteString("|")
	expect := `|日本語は|難しい|! §|ahonndoni|ahonndoni|ande|baamotene|bekombitebe` +
		`|kombite|難|bekombite|bekombite|bekombite|bekombite|ahonndoni|`
	actual := s.String()
	if actual != expect {
		t.Errorf("in TestWriteAsTonelessMixedKindTo(\n%#v\n),\nexpect: %#v\nactual: %#v\n\n",
			sses, expect, actual)
	}
	log.Println("LEAVE TestWriteAsTonelessMixedKind")
}

func TestWriteAsToneless(t *testing.T) {
	log.Println("ENTER TestWriteAsToneless")
	sses := SSEs{
		0xE_089_236_C72_423_000, // ` Ahöñ-ndönî`
		0xF_089_236_C72_423_000, // ` AHÖÑ-NDÖNÎ`
		0xD_08B_455_000_000_000, // ` ândɛ`
		0xD_0CB_089_B31_E55_415, // ` bâa-mo-tɛnɛ`
		0xE_0D7_A6E_362_655_0D7, // ` Bɛ̂-kɔ̈mbïtɛbɛ̂`  [if any syllable is UPPER, the whole word is]
		0x9_A6E_362_655_000_000, //    `-kɔ̈mbïtɛ`
		0xB_0D7_A6E_362_655_000, //  `BƐ̂-KƆ̈MBÏTƐ`
		0xD_0D7_A6E_362_655_000, // ` bɛ̂-kɔ̈mbïtɛ`
		0xF_0D7_A6E_362_655_000, // ` BƐ̂-kɔ̈mbïtɛ`
		0xF_0D7_A6E_362_655_000, // ` BƐ̂-KƆ̈MBÏTƐ`
		0xD_089_236_472_423_000, // ` ahöñndönî`
	}
	var s strings.Builder
	for _, sse := range sses {
		s.WriteString("|")
		sse.WriteAsTonelessTo(&s)
	}
	s.WriteString("|")
	expect := `|ahonndoni|ahonndoni|ande|baamotene|bekombitebe|kombite` +
		`|bekombite|bekombite|bekombite|bekombite|ahonndoni|`
	actual := s.String()
	if actual != expect {
		t.Errorf("in TestWriteAsTonelessTo(\n%#v\n),\nexpect: %#v\nactual: %#v\n\n",
			sses, expect, actual)
	}
	log.Println("LEAVE TestWriteAsToneless")
}

func TestWriteAsHeightless(t *testing.T) {
	log.Println("ENTER TestWriteAsHeightless")
	sses := SSEs{
		0xE_089_236_C72_423_000, // ` Ahöñ-ndönî`
		0xF_089_236_C72_423_000, // ` AHÖÑ-NDÖNÎ`
		0xD_08B_455_000_000_000, // ` ândɛ`
		0xD_0CB_089_B31_E55_415, // ` bâa-mo-tɛnɛ`
		0xE_0D7_A6E_362_655_0D7, // ` Bɛ̂-kɔ̈mbïtɛbɛ̂`  [if any syllable is UPPER, the whole word is]
		0x9_A6E_362_655_000_000, //    `-kɔ̈mbïtɛ`
		0xB_0D7_A6E_362_655_000, //  `BƐ̂-KƆ̈MBÏTƐ`
		0xD_0D7_A6E_362_655_000, // ` bɛ̂-kɔ̈mbïtɛ`
		0xE_0D7_A6E_362_655_000, // ` Bɛ̂-kɔ̈mbïtɛ`
		0xF_0D7_A6E_362_655_000, // ` BƐ̂-KƆ̈MBÏTƐ`
		0xD_089_236_472_423_000, // ` ahöñndönî`
	}
	var s strings.Builder
	for _, sse := range sses {
		s.WriteString("|")
		sse.WriteAsHeightlessTo(&s)
	}
	s.WriteString("|")
	expect := `| Ahön-ndönî| AHÖÑ-NDÖNÎ| ânde| bâa-mo-tene| BƐ̂-kömbïtebê|-kömbïte` +
		`|BƐ̂-KƆ̈MBÏTƐ| bê-kömbïte| BƐ̂-kömbïte| BƐ̂-KƆ̈MBÏTƐ| ahönndönî|`
	actual := s.String()
	if actual != expect {
		t.Errorf("in TestWriteAsHeightlessTo(\n%#v\n),\nexpect: %#v\nactual: %#v\n\n",
			sses, expect, actual)
	}
	log.Println("LEAVE TestWriteAsHeightless")
}

func TestWriteAsLemma(t *testing.T) {
	log.Println("ENTER TestWriteAsLemma")
	sses := SSEs{
		0xE_089_236_C72_423_000, // ` Ahöñ-ndönî`
		0xF_089_236_C72_423_000, // ` AHÖÑ-NDÖNÎ`
		0xD_08B_455_000_000_000, // ` ândɛ`
		0xD_0CB_089_B31_E55_415, // ` bâa-mo-tɛnɛ`
		0xE_0D7_A6E_362_655_0D7, // ` Bɛ̂-kɔ̈mbïtɛbɛ̂`  [if any syllable is UPPER, the whole word is]
		0x9_A6E_362_655_000_000, //    `-kɔ̈mbïtɛ`
		0xB_0D7_A6E_362_655_000, //  `BƐ̂-KƆ̈MBÏTƐ`
		0xD_0D7_A6E_362_655_000, // ` bɛ̂-kɔ̈mbïtɛ`
		0xE_0D7_A6E_362_655_000, // ` Bɛ̂-kɔ̈mbïtɛ`
		0xF_0D7_A6E_362_655_000, // ` BƐ̂-KƆ̈MBÏTƐ`
		0xD_089_236_472_423_000, // ` ahöñndönî`
	}
	var s strings.Builder
	for _, sse := range sses {
		s.WriteString("|")
		sse.WriteAsLemmaTo(&s)
	}
	s.WriteString("|")
	expect := `| Ahön-ndönî| AHÖÑ-NDÖNÎ| ândɛ| bâa-mo-tɛnɛ| BƐ̂-kɔ̈mbïtɛbɛ̂` +
		`|-kɔ̈mbïtɛ|BƐ̂-KƆ̈MBÏTƐ| bɛ̂-kɔ̈mbïtɛ| BƐ̂-kɔ̈mbïtɛ| BƐ̂-KƆ̈MBÏTƐ| ahönndönî|`
	actual := s.String()
	if actual != expect {
		t.Errorf("in TestWriteAsLemmaTo(\n%#v\n),\nexpect: %#v\nactual: %#v\n\n",
			sses, expect, actual)
	}
	log.Println("LEAVE TestWriteAsLemma")
}

func TestWriteAsLemmaForUnknownPitch(t *testing.T) {
	log.Println("ENTER TestWriteAsLemmaForUnknownPitch")
	sses := SSEs{
		0xE_088_234_C70_420_000, // ` Ahöñ-ndönî`
		0xF_088_234_C70_420_000, // ` AHÖÑ-NDÖNÎ`
		0xD_088_454_000_000_000, // ` ândɛ`
		0xD_0C8_088_B30_E54_414, // ` bâa-mo-tɛnɛ`
		0xE_0D4_A6C_360_654_0D4, // ` Bɛ̂-kɔ̈mbïtɛbɛ̂`  [if any syllable is UPPER, the whole word is]
		0x9_A6C_360_654_000_000, //    `-kɔ̈mbïtɛ`
		0xB_0D4_A6C_360_654_000, //  `BƐ̂-KƆ̈MBÏTƐ`
		0xD_0D4_A6C_360_654_000, // ` bɛ̂-kɔ̈mbïtɛ`
		0xE_0D4_A6C_360_654_000, // ` Bɛ̂-kɔ̈mbïtɛ`
		0xF_0D4_A6C_360_654_000, // ` BƐ̂-KƆ̈MBÏTƐ`
		0xD_088_234_470_420_000, // ` ahöñndönî`
	}
	var s strings.Builder
	for _, sse := range sses {
		s.WriteString("|")
		sse.WriteAsLemmaTo(&s)
	}
	s.WriteString("|")
	expect := `| Ạhọn-ndọnị| ẠHỌÑ-NDỌNỊ| ạndɛ̣| bạạ-mọ-tɛ̣nɛ̣| BƐ̣-kɔ̣mbịtɛ̣bɛ̣|` +
		`-kɔ̣mbịtɛ̣|BƐ̣-KƆ̣MBỊTƐ̣| bɛ̣-kɔ̣mbịtɛ̣| BƐ̣-kɔ̣mbịtɛ̣| BƐ̣-KƆ̣MBỊTƐ̣| ạhọnndọnị|`
	actual := s.String()
	if actual != expect {
		t.Errorf("in TestWriteAsLemmaTo(\n%#v\n),\nexpect: %#v\nactual: %#v\n\n",
			sses, expect, actual)
	}
	log.Println("LEAVE TestWriteAsLemmaForUnknownPitch")
}

func TestWriteAsCanonicalForUnknownPitch(t *testing.T) {
	log.Println("ENTER TestWriteAsCanonicalForUnknownPitch")
	sses := SSEs{
		0xE_088_234_C70_420_000, // ` Ahöñ-ndönî`
		0xF_088_234_C70_420_000, // ` AHÖÑ-NDÖNÎ`
		0xD_088_454_000_000_000, // ` ândɛ`
		0xD_0C8_088_B30_E54_414, // ` bâa-mo-tɛnɛ`
		0xE_0D4_A6C_360_654_0D4, // ` Bɛ̂-kɔ̈mbïtɛbɛ̂`  [if any syllable is UPPER, the whole word is]
		0x9_A6C_360_654_000_000, //    `-kɔ̈mbïtɛ`
		0x0010_96E3_0000_0000,   //  `難`
		0xB_0D4_A6C_360_654_000, //  `BƐ̂-KƆ̈MBÏTƐ`
		0xD_0D4_A6C_360_654_000, // ` bɛ̂-kɔ̈mbïtɛ`
		0xE_0D4_A6C_360_654_000, // ` Bɛ̂-kɔ̈mbïtɛ`
		0xF_0D4_A6C_360_654_000, // ` BƐ̂-KƆ̈MBÏTƐ`
		0xD_088_234_470_420_000, // ` ahöñndönî`
	}
	var s strings.Builder
	for _, sse := range sses {
		s.WriteString("|")
		sse.WriteAsCanonicalTo(&s)
	}
	s.WriteString("|")
	expect := `| ~haHO-Doni| =ha=HO-=Do=ni| haDx| baha-mo-txnx| ~bx-kcBitxbx|-kcBitx` +
		`|U+96E3|=bx-=kc=Bi=tx| bx-kcBitx| ~bx-kcBitx| =bx-=kc=Bi=tx| haHODoni|`
	actual := s.String()
	if actual != expect {
		t.Errorf("in TestWriteAsLemmaTo(\n%#v\n),\nexpect: %#v\nactual: %#v\n\n",
			sses, expect, actual)
	}
	log.Println("LEAVE TestWriteAsCanonicalForUnknownPitch")
}

func TestWriteEmptyAsToneless(t *testing.T) {
	log.Println("ENTER TestWriteEmptyAsToneless")
	sses := SSEs{}
	var s strings.Builder
	for _, sse := range sses {
		sse.WriteAsTonelessTo(&s)
	}
	expect := ""
	actual := s.String()
	if actual != expect {
		t.Errorf("in TestWriteEmptyAsTonelessTo(\n%#v\n),\nexpect: %#v\nactual: %#v\n\n",
			sses, expect, actual)
	}
	log.Println("LEAVE TestWriteEmptyAsToneless")
}

func TestWriteEmptyAsHeightless(t *testing.T) {
	log.Println("ENTER TestWriteEmptyAsHeightless")
	sses := SSEs{}
	var s strings.Builder
	for _, sse := range sses {
		sse.WriteAsHeightlessTo(&s)
	}
	expect := ""
	actual := s.String()
	if actual != expect {
		t.Errorf("in TestWriteEmptyAsHeightlessTo(\n%#v\n),\nexpect: %#v\nactual: %#v\n\n",
			sses, expect, actual)
	}
	log.Println("LEAVE TestWriteEmptyAsHeightless")
}

func TestWriteEmptyAsLemma(t *testing.T) {
	log.Println("ENTER TestWriteEmptyAsLemma")
	sses := SSEs{}
	var s strings.Builder
	for _, sse := range sses {
		sse.WriteAsLemmaTo(&s)
	}
	expect := ""
	actual := s.String()
	if actual != expect {
		t.Errorf("in TestWriteEmptyAsLemmaTo(\n%#v\n),\nexpect: %#v\nactual: %#v\n\n",
			sses, expect, actual)
	}
	log.Println("LEAVE TestWriteEmptyAsLemma")
}

func TestWriteEmptyAsUtf8(t *testing.T) {
	log.Println("ENTER TestWriteEmptyAsUtf8")
	sses := SSEs{}
	var s strings.Builder
	for _, sse := range sses {
		sse.WriteAsUtf8To(&s)
	}
	expect := ""
	actual := s.String()
	if actual != expect {
		t.Errorf("in TestWriteEmptyAsLemmaTo(\n%#v\n),\nexpect: %#v\nactual: %#v\n\n",
			sses, expect, actual)
	}
	log.Println("LEAVE TestWriteEmptyAsUtf8")
}

//////////////////////////////////////////////////////////////////////////////

func TestWriteAsUtf8(t *testing.T) {
	log.Println("ENTER TestWriteAsUtf8")
	sses := SSEs{
		0xE_089_236_C72_423_000, // ` Ahöñ-ndönî`
		0xF_089_236_C72_423_000, // ` AHÖÑ-NDÖNÎ`
		0xD_08B_455_000_000_000, // ` ândɛ`
		0xD_0CB_089_B31_E55_415, // ` bâa-mo-tɛnɛ`
		0xE_0D7_A6E_362_655_0D7, // ` Bɛ̂-kɔ̈mbïtɛbɛ̂`  [if any syllable is UPPER, the whole word is]
		0x9_A6E_362_655_000_000, //    `-kɔ̈mbïtɛ`
		0xB_0D7_A6E_362_655_000, //  `BƐ̂-KƆ̈MBÏTƐ`
		0xD_0D7_A6E_362_655_000, // ` bɛ̂-kɔ̈mbïtɛ`
		0xE_0D7_A6E_362_655_000, // ` Bɛ̂-kɔ̈mbïtɛ`
		0xF_0D7_A6E_362_655_000, // ` BƐ̂-KƆ̈MBÏTƐ`
		0xD_089_236_472_423_000, // ` ahöñndönî`
	}
	var s strings.Builder
	for _, sse := range sses {
		s.WriteString("|")
		sse.WriteAsUtf8To(&s)
	}
	s.WriteString("|")
	expect := `| Ahöñ-ndönî| AHÖÑ-NDÖNÎ| ândɛ| bâa-mo-tɛnɛ| BƐ̂-kɔ̈mbïtɛbɛ̂` +
		`|-kɔ̈mbïtɛ|BƐ̂-KƆ̈MBÏTƐ| bɛ̂-kɔ̈mbïtɛ| BƐ̂-kɔ̈mbïtɛ| BƐ̂-KƆ̈MBÏTƐ| ahöñndönî|`
	actual := s.String()
	if actual != expect {
		t.Errorf("in TestWriteAsLemmaTo(\n%#v\n),\nexpect: %#v\nactual: %#v\n\n",
			sses, expect, actual)
	}
	log.Println("LEAVE TestWriteAsUtf8")
}
