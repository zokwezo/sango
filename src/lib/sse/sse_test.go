package sse

import (
	"fmt"
	"strings"
	"testing"
)

func TestCanonicalToCodes(t *testing.T) {
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
}

func TestCodesToSSEs(t *testing.T) {
	u := func(v uint16) sseCode { return sseCode{value: v, isSango: false} }
	s := func(v uint16) sseCode { return sseCode{value: v, isSango: true} }
	codes := []sseCode{
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
	expect := []SSE{
		0x65E5_672C_8A9E_306F, 0x0000_96E3_3057_3044, 0x0021_0020_00A7_0000,
		0xE_089_0F6_A72_463_000, 0xF_089_0F6_A72_463_000, 0xD_08B_255_000_000_000,
		0xD_10B_089_C31_D95_455, 0xE_117_B6E_162_595_117, 0x9_B6E_162_595_000_000,
		0x0000_96E3_0000_0000, 0xB_117_B6E_162_595_000, 0xD_117_B6E_162_595_000,
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
	actual := codesToSSEs(codes)
	actualHex := dumpSSEs(actual)
	expectHex := dumpSSEs(expect)
	if actualHex != expectHex {
		t.Errorf("bad codesToSSEs\nexpect: %v\nactual: %v\n", expectHex, actualHex)
	}
}

func TestWriteAsUTF8MixedKind(t *testing.T) {
	sses := []SSE{
		0x65E5_672C_8A9E_306F, 0x0000_96E3_3057_3044, 0x0021_0020_00A7_0000,
		0xE_089_0F6_A72_463_000, 0xF_089_0F6_A72_463_000, 0xD_08B_255_000_000_000,
		0xD_10B_089_C31_D95_455, 0xE_117_B6E_162_595_117, 0x9_B6E_162_595_000_000,
		0x0000_96E3_0000_0000, 0xB_117_B6E_162_595_000, 0xD_117_B6E_162_595_000,
		0xE_117_B6E_162_595_000, 0xF_117_B6E_162_595_000, 0xD_089_0F6_272_463_000,
	}
	var s strings.Builder
	for _, sse := range sses {
		s.WriteString("|")
		sse.WriteAsUTF8To(&s)
	}
	s.WriteString("|")
	expect := `|日本語は|難しい|! §| Ahöñ-ndönî| AHÖÑ-NDÖNÎ` +
		`| ândɛ| bâa-mo-tɛnɛ| BƐ̂-kɔ̈mbïtɛbɛ̂|-kɔ̈mbïtɛ|難|BƐ̂-KƆ̈MBÏTƐ` +
		`| bɛ̂-kɔ̈mbïtɛ| BƐ̂-kɔ̈mbïtɛ| BƐ̂-KƆ̈MBÏTƐ| ahöñndönî|`
	actual := s.String()
	if actual != expect {
		t.Errorf("in TestWriteAsUTF8MixedKindTo(\n%#v\n),\nexpect: %#v\nactual: %#v\n\n",
			sses, expect, actual)
	}
}

func TestWriteAsTonelessMixedKind(t *testing.T) {
	sses := []SSE{
		0x65E5_672C_8A9E_306F, 0x0000_96E3_3057_3044, 0x0021_0020_00A7_0000,
		0xE_089_0F6_A72_463_000, 0xF_089_0F6_A72_463_000, 0xD_08B_255_000_000_000,
		0xD_10B_089_C31_D95_455, 0xE_117_B6E_162_595_117, 0x9_B6E_162_595_000_000,
		0x0000_96E3_0000_0000, 0xB_117_B6E_162_595_000, 0xD_117_B6E_162_595_000,
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
}

func TestWriteAsToneless(t *testing.T) {
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
}

func TestWriteAsHeightless(t *testing.T) {
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
}

func TestWriteAsLemma(t *testing.T) {
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
}

func TestWriteAsUTF8(t *testing.T) {
	sses := []SSE{
		0xE_089_0F6_A72_463_000, 0xF_089_0F6_A72_463_000, 0xD_08B_255_000_000_000,
		0xD_10B_089_C31_D95_455, 0xE_117_B6E_162_595_117, 0x9_B6E_162_595_000_000,
		0xB_117_B6E_162_595_000, 0xD_117_B6E_162_595_000, 0xE_117_B6E_162_595_000,
		0xF_117_B6E_162_595_000, 0xD_089_0F6_272_463_000,
	}
	var s strings.Builder
	for _, sse := range sses {
		s.WriteString("|")
		sse.WriteAsUTF8To(&s)
	}
	s.WriteString("|")
	expect := `| Ahöñ-ndönî| AHÖÑ-NDÖNÎ| ândɛ| bâa-mo-tɛnɛ| BƐ̂-kɔ̈mbïtɛbɛ̂` +
		`|-kɔ̈mbïtɛ|BƐ̂-KƆ̈MBÏTƐ| bɛ̂-kɔ̈mbïtɛ| BƐ̂-kɔ̈mbïtɛ| BƐ̂-KƆ̈MBÏTƐ| ahöñndönî|`
	actual := s.String()
	if actual != expect {
		t.Errorf("in TestWriteAsLemmaTo(\n%#v\n),\nexpect: %#v\nactual: %#v\n\n",
			sses, expect, actual)
	}
}

func TestWriteEmptyAsToneless(t *testing.T) {
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
}

func TestWriteEmptyAsHeightless(t *testing.T) {
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
}

func TestWriteEmptyAsLemma(t *testing.T) {
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
}

func TestWriteEmptyAsUTF8(t *testing.T) {
	sses := []SSE{}
	var s strings.Builder
	for _, sse := range sses {
		sse.WriteAsUTF8To(&s)
	}
	expect := ""
	actual := s.String()
	if actual != expect {
		t.Errorf("in TestWriteEmptyAsLemmaTo(\n%#v\n),\nexpect: %#v\nactual: %#v\n\n",
			sses, expect, actual)
	}
}

func TestEmptyCanonicalToCodes(t *testing.T) {
	actual, leftOffAt := canonicalToCodes("")
	if leftOffAt != 0 {
		t.Errorf("bad leftOffAt")
	}
	if len(actual) != 0 {
		t.Errorf("found nonempty codes")
	}
}

func TestEmptyCodeToSSEs(t *testing.T) {
	actual := codesToSSEs([]sseCode{})
	if len(actual) != 0 {
		t.Errorf("bad codesToSSEs, expected empty array")
	}
}

func TestGoodCanonicalToSSEs(t *testing.T) {
	c := `U+65E5U+672CU+8A9EU+306FU+96E3U+3057U+3044U+0021U+0020U+00A7 ~ha_HO:-Do:ni^` +
		` =ha_HO:-Do:ni^ ha^Dx_ ba^ha_-mo_-tx_nx_ ~bx^-kc:Bi:tx_bx^-kc:Bi:tx_U+96E3` +
		`=bx^-=kc:=Bi:=tx_ bx^-kc:Bi:tx_ ~bx^-kc:Bi:tx_ =bx^-=kc:=Bi:=tx_ ha_HO:Do:ni^`
	expect := []SSE{
		0x65E5_672C_8A9E_306F, 0x0000_96E3_3057_3044, 0x0021_0020_00A7_0000,
		0xE_089_0F6_A72_463_000, 0xF_089_0F6_A72_463_000, 0xD_08B_255_000_000_000,
		0xD_10B_089_C31_D95_455, 0xE_117_B6E_162_595_117, 0x9_B6E_162_595_000_000,
		0x0_000_96E_300_000_000, 0xB_117_B6E_162_595_000, 0xD_117_B6E_162_595_000,
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
}

func TestWriteGoodCanonicalAsLemma(t *testing.T) {
	sses := []SSE{
		0xE_089_0F6_A72_463_000, 0xF_089_0F6_A72_463_000, 0xD_08B_255_000_000_000,
		0xD_10B_089_C31_D95_455, 0xE_117_B6E_162_595_117, 0x9_B6E_162_595_000_000,
		0xB_117_B6E_162_595_000, 0xD_117_B6E_162_595_000,
		0xE_117_B6E_162_595_000, 0xF_117_B6E_162_595_000, 0xD_089_0F6_272_463_000,
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
}

func TestBadCanonicalToSSEs(t *testing.T) {
	c := `U+65E5U+672CU+8A9EU+306FU+96E3U+3057U+3044U+0021U+0020U+00A7` +
		` ~ha_HO:-Do:ni^ =ha_HO:-jo:ni^ ha^Dx_ ba^ha_-mo_-tx_nx_ ~bx^-kc:Bi:tx_bx^-kc:Bi:tx_` +
		`U+96E3=bx^-=kc:=Bi:=tx_ bx^-kc:Bi:tx_ ~bx^-kc:Bi:tx_ =bx^-=kc:=Bi:=tx_ ha_HO:Do:ni^`
	expect := []SSE{
		0x65E5_672C_8A9E_306F, 0x0000_96E3_3057_3044, 0x0021_0020_00A7_0000,
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
	expectErr := fmt.Errorf("cannot parse Canonical string starting at s[83:] = %q", "-jo:ni^ ha...")
	actual, actualErr := CanonicalToSSEs(c)
	if actualErr.Error() != expectErr.Error() {
		t.Errorf("expected error not returned from CanonicalToSSEs\nactualErr = %v\nexpectErr = %v", actualErr, expectErr)
	}
	actualHex := dumpSSEs(actual)
	expectHex := dumpSSEs(expect)
	if actualHex != expectHex {
		t.Errorf("bad BadCanonicalToSSEs\nexpect: %v\nactual: %v\n", expectHex, actualHex)
	}
}

func TestCanonicalToSSEsForUnknownPitch(t *testing.T) {
	c := `~haHO-Doni =haHO-Doni haDx baha-mo-txnx ~bx-kcBitxbx-kcBitxU+96E3` +
		`=bx-=kc=Bi=tx bx-kcBitx ~bx-kcBitx =bx-=kc=Bi=tx haHODoni`
	expect := []SSE{
		0xA_088_0F4_A70_460_000, 0xF_088_0F4_A70_460_000, 0xD_088_254_000_000_000,
		0xD_108_088_C30_D94_454, 0xE_114_B6C_160_594_114, 0x9_B6C_160_594_000_000,
		0x0_000_96E_300_000_000, 0xB_114_B6C_160_594_000, 0xD_114_B6C_160_594_000,
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
}

func TestWriteAsLemmaForUnknownPitch(t *testing.T) {
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
}

func TestWriteAsCanonicalForUnknownPitch(t *testing.T) {
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
}

func TestCanonicalToSSEsToCanonical(t *testing.T) {
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
	// expect := c, but with Titlecase suppressed for nonleading syllables.
	expect := `~haHO-Doni =ha=HO-=Do=ni haDx baha-mo-txnx ~bx-kcBitx` +
		`bx-kcBitxU+96E3=bx-=kc=Bi=tx bx-kcBitx ~bx-kcBitx =bx-=kc=Bi=tx haHODoni`
	actual := s.String()
	if actual != expect {
		t.Errorf("bad BadCanonicalToSSEs\nexpect: %v\nactual: %v\n", expect, actual)
	}
}
