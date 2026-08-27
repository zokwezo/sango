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

func TestUnpadUnicodeWords(t *testing.T) {
	log.Println("ENTER TestUnpadUnicodeWords")
	words := []uint64{
		0x65E5_672C_8A9E_306F,
		0x0000_96E3_3057_3044,
		0x0010_96E3_3057_3044,
		0x0021_0020_00A7_0000,
		0x0000_96E3_3057_0000,
		0x0010_96E3_3057_0000,
		0x0021_0020_0000_0000,
		0x0000_96E3_0000_0000,
		0x0010_96E3_0000_0000,
		0x0021_0000_0000_0000,
		0x0000_0000_0000_0000,
	}
	expectsUnpadded := []uint64{
		0x65E5_672C_8A9E_306F,
		0x0000_96E3_3057_3044,
		0x0010_96E3_3057_3044,
		0x0000_0021_0020_00A7,
		0x0000_0000_96E3_3057,
		0x0000_0010_96E3_3057,
		0x0000_0000_0021_0020,
		0x0000_0000_0000_96E3,
		0x0000_0000_0010_96E3,
		0x0000_0000_0000_0021,
		0x0000_0000_0000_0000,
	}
	for k, word := range words {
		expect := expectsUnpadded[k]
		actual := unpadRight(word)
		if actual != expect {
			t.Errorf("bad UnpadUnicodeWords[%v](%#016X)\nexpect: %#016X\nactual: %#016X\n", k, word, expect, actual)
		}
	}
	log.Println("LEAVE TestUnpadUnicodeWords")
}

func TestPadUnicodeWords(t *testing.T) {
	log.Println("ENTER TestPadUnicodeWords")
	words := []uint64{
		0x65E5_672C_8A9E_306F,
		0x0010_96E3_3057_3044,
		0x0000_0021_0020_00A7,
		0x0000_0010_96E3_3057,
		0x0000_0000_0021_0020,
		0x0000_0000_0010_96E3,
		0x0000_0000_0000_0021,
		0x0000_0000_0000_0000,
	}
	expectsPadded := []uint64{
		0x65E5_672C_8A9E_306F,
		0x0010_96E3_3057_3044,
		0x0021_0020_00A7_0000,
		0x0010_96E3_3057_0000,
		0x0021_0020_0000_0000,
		0x0010_96E3_0000_0000,
		0x0021_0000_0000_0000,
		0x0000_0000_0000_0000,
	}
	for k, word := range words {
		expect := expectsPadded[k]
		actual := padRight(word)
		if actual != expect {
			t.Errorf("bad PadUnicodeWords[%v](%#016X)\nexpect: %#016X\nactual: %#016X\n", k, word, expect, actual)
		}
	}
	log.Println("LEAVE TestPadUnicodeWords")
}

func TestUnpadSangoWords(t *testing.T) {
	log.Println("ENTER TestUnpadSangoWords")
	words := []uint64{
		0xD_10B_089_C31_D95_455,
		0xF_089_0F6_A72_463_000,
		0x9_B6E_162_595_000_000,
		0xD_08B_255_000_000_000,
		0xE_08B_000_000_000_000,
		0xF_000_000_000_000_000,
	}
	expectsUnpadded := []uint64{
		0xD_10B_089_C31_D95_455,
		0x000_F_089_0F6_A72_463,
		0x000_000_9_B6E_162_595,
		0x000_000_000_D_08B_255,
		0x000_000_000_000_E_08B,
		0x000_000_000_000_000_F,
	}
	for k, word := range words {
		expect := expectsUnpadded[k]
		actual := unpadRight(word)
		if actual != expect {
			t.Errorf("bad UnpadSangoWords[%v](%#016X)\nexpect: %#016X\nactual: %#016X\n", k, word, expect, actual)
		}
	}
	log.Println("LEAVE TestUnpadSangoWords")
}

func TestPadSangoWords(t *testing.T) {
	log.Println("ENTER TestPadSangoWords")
	words := []uint64{
		0xD_10B_089_C31_D95_455,
		0x000_F_089_0F6_A72_463,
		0x000_000_9_B6E_162_595,
		0x000_000_000_D_08B_255,
		0x000_000_000_000_E_08B,
		0x000_000_000_000_000_F,
	}
	expectsPadded := []uint64{
		0xD_10B_089_C31_D95_455,
		0xF_089_0F6_A72_463_000,
		0x9_B6E_162_595_000_000,
		0xD_08B_255_000_000_000,
		0xE_08B_000_000_000_000,
		0xF_000_000_000_000_000,
	}
	for k, word := range words {
		expect := expectsPadded[k]
		actual := padRight(word)
		if actual != expect {
			t.Errorf("bad PadSangoWords[%v](%#016X)\nexpect: %#016X\nactual: %#016X\n", k, word, expect, actual)
		}
	}
	log.Println("LEAVE TestPadSangoWords")
}

func TestUnpadUnicodeFollowedByPad(t *testing.T) {
	log.Println("ENTER TestUnpadUnicodeFollowedByPad")
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
	log.Println("LEAVE TestUnpadUnicodeFollowedByPad")
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
		s(0xE089), s(0x90F6), s(0x9A72), s(0x9463), s(0xF089),
		s(0x90F6), s(0x9A72), s(0x9463), s(0xD08B), s(0x9255),
		s(0xD10B), s(0x9089), s(0x9C31), s(0x9D95), s(0x9455),
		s(0xE117), s(0x9B6E), s(0x9162), s(0x9595), s(0x9117),
		s(0x9B6E), s(0x9162), s(0x9595),
		u(0x96E3),
		s(0xB117), s(0xBB6E), s(0xB162), s(0xB595), s(0xD117),
		s(0x9B6E), s(0x9162), s(0x9595), s(0xE117), s(0x9B6E),
		s(0x9162), s(0x9595), s(0xF117), s(0xBB6E), s(0xB162),
		s(0xB595), s(0xD089), s(0x90F6), s(0x9272), s(0x9463),
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
	expect := []SSE{
		0x65E5_672C_8A9E_306F, 0x0010_96E3_3057_3044, 0x0021_0020_00A7_0000,
		0xE_089_0F6_A72_463_000, 0xF_089_0F6_A72_463_000, 0xD_08B_255_000_000_000,
		0xD_10B_089_C31_D95_455, 0xE_117_B6E_162_595_117, 0x9_B6E_162_595_000_000,
		0x0_010_96E_300_000_000, 0xB_117_B6E_162_595_000, 0xD_117_B6E_162_595_000,
		0xE_117_B6E_162_595_000, 0xF_117_B6E_162_595_000, 0xD_089_0F6_272_463_000,
	}
	dumpSSEs := func(sses []SSE) string {
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
	expect := []SSE{
		0x65E5_672C_8A9E_306F, 0x0010_96E3_3057_3044, 0x0021_0020_00A7_0000,
		0xE_089_0F6_A72_463_000, 0xF_089_0F6_000_000_000,
	} // results up to broken parse
	dumpSSEs := func(sses []SSE) string {
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
	expect := []SSE{
		0xA_088_0F4_A70_460_000, 0xF_088_0F4_A70_460_000, 0xD_088_254_000_000_000,
		0xD_108_088_C30_D94_454, 0xE_114_B6C_160_594_114, 0x9_B6C_160_594_000_000,
		0x0_010_96E_300_000_000, 0xB_114_B6C_160_594_000, 0xD_114_B6C_160_594_000,
		0xE_114_B6C_160_594_000, 0xF_114_B6C_160_594_000, 0xD_088_0F4_270_460_000,
	}
	dumpSSEs := func(sses []SSE) string {
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
		s(0xE089), s(0x90F6), s(0x9A72), s(0x9463), // ` Ahöñ-ndönî`
		s(0xE089), s(0xB0F6), s(0xBA72), s(0xB463), // ` AHÖÑ-NDÖNÎ`
		s(0xD08B), s(0x9255), // ` ândɛ`
		s(0xD10B), s(0x9089), s(0x9C31), s(0x9D95), s(0x9455), // ` bâa-mo-tɛnɛ`
		s(0xF117), s(0x9B6E), s(0x9162), s(0x9595), // ` BƐ̂-kɔ̈mbïtɛb`
		s(0x9117), s(0x9B6E), s(0x9162), s(0x9595), //  `bɛ̂-kɔ̈mbïtɛ`
		u(0x96E3),                                  //  `難`
		s(0xB117), s(0xBB6E), s(0xB162), s(0xB595), //  `BƐ̂-KƆ̈MBÏTƐ`
		s(0xD117), s(0x9B6E), s(0x9162), s(0x9595), // ` bɛ̂-kɔ̈mbïtɛ`
		s(0xF117), s(0x9B6E), s(0x9162), s(0x9595), // ` BƐ̂-kɔ̈mbïtɛ`
		s(0xF117), s(0xBB6E), s(0xB162), s(0xB595), // ` BƐ̂-KƆ̈MBÏTƐ`
		s(0xD089), s(0x90F6), s(0x9272), s(0x9463), // ` ahöñndönî`
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
	expect := []SSE{
		0x65E5_672C_8A9E_306F, 0x0010_96E3_3057_3044, 0x0021_0020_00A7_0000, // `日本語は難しい! §`
		0xE_089_0F6_A72_463_000, // ` Ahöñ-ndönî`
		0xF_089_0F6_A72_463_000, // ` AHÖÑ-NDÖNÎ`
		0xD_08B_255_000_000_000, // ` ândɛ`
		0xD_10B_089_C31_D95_455, // ` bâa-mo-tɛnɛ`
		0xF_117_B6E_162_595_117, // ` BƐ̂-kɔ̈mbïtɛb`  // if any syllable is UPPER, the whole word is
		0x9_B6E_162_595_000_000, //  `bɛ̂-kɔ̈mbïtɛ`
		0x0010_96E3_0000_0000,   //  `難`
		0xB_117_B6E_162_595_000, //  `BƐ̂-KƆ̈MBÏTƐ`
		0xD_117_B6E_162_595_000, // ` bɛ̂-kɔ̈mbïtɛ`
		0xF_117_B6E_162_595_000, // ` BƐ̂-kɔ̈mbïtɛ`
		0xF_117_B6E_162_595_000, // ` BƐ̂-KƆ̈MBÏTƐ`
		0xD_089_0F6_272_463_000, // ` ahöñndönî`
	}
	dumpSSEs := func(sses []SSE) string {
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
	expect := []SSE{0x65E5_672C_8A9E_0000} // results up to broken parse
	dumpSSEs := func(sses []SSE) string {
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
	expect := []SSE{
		0xA_088_0F4_270_460_000, 0xF_088_0F4_270_460_000, 0xD_088_258_000_000_000,
		0xD_108_088_430_598_458, 0xF_118_B70_160_598_118, 0x9_B70_160_598_000_000,
		0x0_010_96E_300_000_000, 0xB_117_B6E_162_594_000, 0xD_117_B6E_162_594_000,
		0xF_117_B6E_162_594_000, 0xF_117_B6E_162_594_000, 0xD_088_0F6_272_463_000,
	}
	dumpSSEs := func(sses []SSE) string {
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
	// TODO: Set UPPER case only if a syllable has more than one letter, otherwise set Title case.
	//       If any syllable in a word has UPPER, set all to UPPER.
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
		s(0xE089), s(0x90F6), s(0x9A72), s(0x9463), // ` Ahöñ-ndönî`
		s(0xE089), s(0xB0F6), s(0xBA72), s(0xB463), // ` AHÖÑ-NDÖNÎ`
		s(0xD08B), s(0x9255), // ` ândɛ`
		s(0xD10B), s(0x9089), s(0x9C31), s(0x9D95), s(0x9455), // ` bâa-mo-tɛnɛ`
		s(0xF117), s(0x9B6E), s(0x9162), s(0x9595), // ` BƐ̂-kɔ̈mbïtɛb`
		s(0x9117), s(0x9B6E), s(0x9162), s(0x9595), //  `bɛ̂-kɔ̈mbïtɛ`
		u(0x96E3),                                  //  `難`
		s(0xB117), s(0xBB6E), s(0xB162), s(0xB595), //  `BƐ̂-KƆ̈MBÏTƐ`
		s(0xD117), s(0x9B6E), s(0x9162), s(0x9595), // ` bɛ̂-kɔ̈mbïtɛ`
		s(0xF117), s(0x9B6E), s(0x9162), s(0x9595), // ` BƐ̂-kɔ̈mbïtɛ`
		s(0xF117), s(0xBB6E), s(0xB162), s(0xB595), // ` BƐ̂-KƆ̈MBÏTƐ`
		s(0xD089), s(0x90F6), s(0x9272), s(0x9463), // ` ahöñndönî`
	}
	expect := []SSE{
		0x65E5_672C_8A9E_306F, 0x0010_96E3_3057_3044, 0x0021_0020_00A7_0000, // `日本語は難しい! §`
		0xE_089_0F6_A72_463_000, // ` Ahöñ-ndönî`
		0xF_089_0F6_A72_463_000, // ` AHÖÑ-NDÖNÎ`
		0xD_08B_255_000_000_000, // ` ândɛ`
		0xD_10B_089_C31_D95_455, // ` bâa-mo-tɛnɛ`
		0xF_117_B6E_162_595_117, // ` BƐ̂-kɔ̈mbïtɛb`  // if any syllable is UPPER, the whole word is
		0x9_B6E_162_595_000_000, //  `bɛ̂-kɔ̈mbïtɛ`
		0x0010_96E3_0000_0000,   //  `難`
		0xB_117_B6E_162_595_000, //  `BƐ̂-KƆ̈MBÏTƐ`
		0xD_117_B6E_162_595_000, // ` bɛ̂-kɔ̈mbïtɛ`
		0xF_117_B6E_162_595_000, // ` BƐ̂-kɔ̈mbïtɛ`
		0xF_117_B6E_162_595_000, // ` BƐ̂-KƆ̈MBÏTƐ`
		0xD_089_0F6_272_463_000, // ` ahöñndönî`
	}
	dumpSSEs := func(sses []SSE) string {
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
	sses := []SSE{
		0x65E5_672C_8A9E_306F, 0x0010_96E3_3057_3044, 0x0021_0020_00A7_0000,
		0xE_089_0F6_A72_463_000, 0xF_089_0F6_A72_463_000, 0xD_08B_255_000_000_000,
		0xD_10B_089_C31_D95_455, 0xE_117_B6E_162_595_117, 0x9_B6E_162_595_000_000,
		0x0010_96E3_0000_0000, 0xB_117_B6E_162_595_000, 0xD_117_B6E_162_595_000,
		0xE_117_B6E_162_595_000, 0xF_117_B6E_162_595_000, 0xD_089_0F6_272_463_000,
	}
	var s strings.Builder
	for _, sse := range sses {
		s.WriteString("|")
		sse.WriteAsUtf8To(&s)
	}
	s.WriteString("|")
	expect := `|日本語は|難しい|! §| Ahöñ-ndönî| AHÖÑ-NDÖNÎ` +
		`| ândɛ| bâa-mo-tɛnɛ| BƐ̂-kɔ̈mbïtɛbɛ̂|-kɔ̈mbïtɛ|難|BƐ̂-KƆ̈MBÏTƐ` +
		`| bɛ̂-kɔ̈mbïtɛ| BƐ̂-kɔ̈mbïtɛ| BƐ̂-KƆ̈MBÏTƐ| ahöñndönî|`
	actual := s.String()
	if actual != expect {
		t.Errorf("in TestWriteAsUtf8MixedKindTo(\n%#v\n),\nexpect: %#v\nactual: %#v\n\n",
			sses, expect, actual)
	}
	log.Println("LEAVE TestWriteAsUtf8MixedKind")
}

func TestWriteAsTonelessMixedKind(t *testing.T) {
	log.Println("ENTER TestWriteAsTonelessMixedKind")
	sses := []SSE{
		0x65E5_672C_8A9E_306F, 0x0010_96E3_3057_3044, 0x0021_0020_00A7_0000,
		0xE_089_0F6_A72_463_000, 0xF_089_0F6_A72_463_000, 0xD_08B_255_000_000_000,
		0xD_10B_089_C31_D95_455, 0xE_117_B6E_162_595_117, 0x9_B6E_162_595_000_000,
		0x0010_96E3_0000_0000, 0xB_117_B6E_162_595_000, 0xD_117_B6E_162_595_000,
		0xE_117_B6E_162_595_000, 0xF_117_B6E_162_595_000, 0xD_089_0F6_272_463_000,
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
	sses := []SSE{
		0xE_089_0F6_A72_463_000, 0xF_089_0F6_A72_463_000, 0xD_08B_255_000_000_000,
		0xD_10B_089_C31_D95_455, 0xE_117_B6E_162_595_117, 0x9_B6E_162_595_000_000,
		0xB_117_B6E_162_595_000, 0xD_117_B6E_162_595_000, 0xE_117_B6E_162_595_000,
		0xF_117_B6E_162_595_000, 0xD_089_0F6_272_463_000,
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
	sses := []SSE{
		0xE_089_0F6_A72_463_000, 0xF_089_0F6_A72_463_000, 0xD_08B_255_000_000_000,
		0xD_10B_089_C31_D95_455, 0xE_117_B6E_162_595_117, 0x9_B6E_162_595_000_000,
		0xB_117_B6E_162_595_000, 0xD_117_B6E_162_595_000, 0xE_117_B6E_162_595_000,
		0xF_117_B6E_162_595_000, 0xD_089_0F6_272_463_000,
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
	sses := []SSE{
		0xE_089_0F6_A72_463_000, 0xF_089_0F6_A72_463_000, 0xD_08B_255_000_000_000,
		0xD_10B_089_C31_D95_455, 0xE_117_B6E_162_595_117, 0x9_B6E_162_595_000_000,
		0xB_117_B6E_162_595_000, 0xD_117_B6E_162_595_000, 0xE_117_B6E_162_595_000,
		0xF_117_B6E_162_595_000, 0xD_089_0F6_272_463_000,
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
	sses := []SSE{
		0xA_088_0F4_A70_460_000, 0xF_088_0F4_A70_460_000, 0xD_088_254_000_000_000,
		0xD_108_088_C30_D94_454, 0xE_114_B6C_160_594_114, 0x9_B6C_160_594_000_000,
		0x0_000_96E_300_000_000, 0xB_114_B6C_160_594_000, 0xD_114_B6C_160_594_000,
		0xE_114_B6C_160_594_000, 0xF_114_B6C_160_594_000, 0xD_088_0F4_270_460_000,
	}
	var s strings.Builder
	for _, sse := range sses {
		s.WriteString("|")
		sse.WriteAsLemmaTo(&s)
	}
	s.WriteString("|")
	expect := `|Ạhọn-ndọnị| ẠHỌÑ-NDỌNỊ| ạndɛ̣| bạạ-mọ-tɛ̣nɛ̣| BƐ̣-kɔ̣mbịtɛ̣bɛ̣|` +
		`-kɔ̣mbịtɛ̣|難|BƐ̣-KƆ̣MBỊTƐ̣| bɛ̣-kɔ̣mbịtɛ̣| BƐ̣-kɔ̣mbịtɛ̣| BƐ̣-KƆ̣MBỊTƐ̣| ạhọnndọnị|`
	actual := s.String()
	if actual != expect {
		t.Errorf("in TestWriteAsLemmaTo(\n%#v\n),\nexpect: %#v\nactual: %#v\n\n",
			sses, expect, actual)
	}
	log.Println("LEAVE TestWriteAsLemmaForUnknownPitch")
}

func TestWriteAsCanonicalForUnknownPitch(t *testing.T) {
	log.Println("ENTER TestWriteAsCanonicalForUnknownPitch")
	sses := []SSE{
		0xA_088_0F4_A70_460_000, 0xF_088_0F4_A70_460_000, 0xD_088_254_000_000_000,
		0xD_108_088_C30_D94_454, 0xE_114_B6C_160_594_114, 0x9_B6C_160_594_000_000,
		0x0_000_96E_300_000_000, 0xB_114_B6C_160_594_000, 0xD_114_B6C_160_594_000,
		0xE_114_B6C_160_594_000, 0xF_114_B6C_160_594_000, 0xD_088_0F4_270_460_000,
	}
	var s strings.Builder
	for _, sse := range sses {
		s.WriteString("|")
		sse.WriteAsCanonicalTo(&s)
	}
	s.WriteString("|")
	expect := `|~haHO-Doni| =ha=HO-=Do=ni| haDx| baha-mo-txnx| ~bx-kcBitxbx|-kcBitx` +
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
	sses := []SSE{}
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
	sses := []SSE{}
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
	sses := []SSE{}
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
	sses := []SSE{}
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
	sses := []SSE{
		0xE_089_0F6_A72_463_000, 0xF_089_0F6_A72_463_000, 0xD_08B_255_000_000_000,
		0xD_10B_089_C31_D95_455, 0xE_117_B6E_162_595_117, 0x9_B6E_162_595_000_000,
		0xB_117_B6E_162_595_000, 0xD_117_B6E_162_595_000, 0xE_117_B6E_162_595_000,
		0xF_117_B6E_162_595_000, 0xD_089_0F6_272_463_000,
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
