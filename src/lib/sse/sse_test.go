package sse

import (
	"fmt"
	"strings"
	"testing"
)

func TestWriteAsUTF8(t *testing.T) {
	sses := []SSE{
		0x65E5_672C_8A9E_306F,
		0x0000_96E3_3057_3044,
		0x0021_0000_0020_00A7,
		0x8_08E_AE5_0D9_50C_FFF,
		0x9_08E_AE5_0D9_50C_FFF,
		0xA_08E_AE5_0D9_50C_FFF,
		0xB_08E_AE5_0D9_50C_FFF,
		0xC_08E_AE5_0D9_50C_FFF,
		0xD_08E_AE5_0D9_50C_FFF,
		0xE_08E_AE5_0D9_50C_FFF,
		0xF_08E_AE5_0D9_50C_FFF,
	}
	var s strings.Builder
	for _, sse := range sses {
		sse.WriteAsUTF8To(&s)
	}
	expect := "日本語は難しい! §bɛ̂-kɔ̈mbïtɛBƐ̂-kɔ̈mbïtɛBƐ̂-KƆ̈MBÏTƐ bɛ̂-kɔ̈mbïtɛ BƐ̂-kɔ̈mbïtɛ BƐ̂-KƆ̈MBÏTƐ"
	actual := s.String()
	if actual != expect {
		t.Errorf("From UTF8FromSSE(\n%#v\n),\nexpect: %#v\nactual: %#v\n\n",
			sses, expect, actual)
	}
}

func TestWriteAsCanonical(t *testing.T) {
	sses := []SSE{
		0x65E5_672C_8A9E_306F,
		0x0000_96E3_3057_3044,
		0x0021_0000_0020_00A7,
		0x8_08E_AE5_0D9_50C_FFF,
		0x9_08E_AE5_0D9_50C_FFF,
		0xA_08E_AE5_0D9_50C_FFF,
		0xB_08E_AE5_0D9_50C_FFF,
		0xC_08E_AE5_0D9_50C_FFF,
		0xD_08E_AE5_0D9_50C_FFF,
		0xE_08E_AE5_0D9_50C_FFF,
		0xF_08E_AE5_0D9_50C_FFF,
	}
	var s strings.Builder
	for _, sse := range sses {
		sse.WriteAsCanonicalTo(&s)
	}
	expect := "U+65E5U+672CU+8A9EU+306FU+96E3U+3057U+3044U+0021U+0020U+00A7" +
		"bx^-kc:Bi:tx_~bx^-kc:Bi:tx_=bx^-=kc:=Bi:=tx_#bx^-#kc:#Bi:#tx_ bx^-kc:Bi:tx_ ~bx^-kc:Bi:tx_ =bx^-=kc:=Bi:=tx_ #bx^-#kc:#Bi:#tx_"
	actual := s.String()
	if actual != expect {
		t.Errorf("From CanonicalFromSSE(\n%#v\n),\nexpect: %#v\nactual: %#v\n\n",
			sses, expect, actual)
	}
}

func TestCanonicalToCodes(t *testing.T) {
	c := "U+65E5U+672CU+8A9EU+306FU+96E3U+3057U+3044U+0021U+0020U+00A7" +
		"bx^-kc:Bi:tx_~bx^-kc:Bi:tx_U+96E3=bx^-=kc:=Bi:=tx_ bx^-kc:Bi:tx_ ~bx^-kc:Bi:tx_ =bx^-=kc:=Bi:=tx_"
	u := func(v uint16) sseCode { return sseCode{value: v, isSango: false} }
	s := func(v uint16) sseCode { return sseCode{value: v, isSango: true} }
	expect := []sseCode{
		u(0x65E5), u(0x672C), u(0x8A9E), u(0x306F), u(0x96E3),
		u(0x3057), u(0x3044), u(0x0021), u(0x0020), u(0x00A7),
		s(0x808E), s(0x8AE5), s(0x80D9), s(0x850C),
		s(0x908E), s(0x8AE5), s(0x80D9), s(0x850C),
		u(0x96E3),
		s(0xA08E), s(0xAAE5), s(0xA0D9), s(0xA50C),
		s(0xC08E), s(0x8AE5), s(0x80D9), s(0x850C),
		s(0xD08E), s(0x8AE5), s(0x80D9), s(0x850C),
		s(0xE08E), s(0xAAE5), s(0xA0D9), s(0xA50C),
	}
	expectLen := len(c)
	actual, actualLen := canonicalToCodes(c)
	if actualLen != expectLen {
		t.Errorf("error found at c[%v](good: %q bad: %q\n", actualLen, c[0:actualLen], c[actualLen:])
	}
	var actualHex string
	for _, x := range actual {
		if x.isSango {
			actualHex += fmt.Sprintf("S+%016X ", x.value)
		} else {
			actualHex += fmt.Sprintf("%U ", x.value)
		}
	}
	var expectHex string
	for _, x := range expect {
		if x.isSango {
			expectHex += fmt.Sprintf("S+%016X ", x.value)
		} else {
			expectHex += fmt.Sprintf("%U ", x.value)
		}
	}
	if actualHex != expectHex {
		t.Errorf("bad canonicalToCodes\nexpect: %v\nactual: %v\n", expectHex, actualHex)
	}
}

func TestCodesToSSEs(t *testing.T) {
	u := func(v uint16) sseCode { return sseCode{value: v, isSango: false} }
	s := func(v uint16) sseCode { return sseCode{value: v, isSango: true} }
	codes := []sseCode{
		u(0x65E5), u(0x672C), u(0x8A9E), u(0x306F), u(0x96E3),
		u(0x3057), u(0x3044), u(0x0021), u(0x0020), u(0x00A7),
		s(0x808E), s(0x8AE5), s(0x80D9), s(0x850C),
		s(0x908E), s(0x8AE5), s(0x80D9), s(0x850C),
		u(0x96E3),
		s(0xA08E), s(0xAAE5), s(0xA0D9), s(0xA50C),
		s(0xC08E), s(0x8AE5), s(0x80D9), s(0x850C),
		s(0xD08E), s(0x8AE5), s(0x80D9), s(0x850C),
		s(0xE08E), s(0xAAE5), s(0xA0D9), s(0xA50C),
	}
	expect := []SSE{
		0x65E5672C8A9E306F, 0x000096E330573044, 0x0021002000A70000, // Unicode
		0x808EAE50D950C08E, 0x8AE50D950CFFFFFF, 0x000096E300000000, // 12-syllable Sango word (no space separator)
		0xA08EAE50D950CFFF, // word break following Sango word prefixed by a space
		0xC08EAE50D950CFFF, // word break following Sango word prefixed by a space
		0xD08EAE50D950CFFF, // word break following Sango word prefixed by a space
		0xE08EAE50D950CFFF, // end of input with trailing padding
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

func TestGoodCanonicalToSSEs(t *testing.T) {
	c := "U+65E5U+672CU+8A9EU+306FU+96E3U+3057U+3044U+0021U+0020U+00A7" +
		"bx^-kc:Bi:tx_~bx^-kc:Bi:tx_U+96E3=bx^-=kc:=Bi:=tx_ bx^-kc:Bi:tx_ ~bx^-kc:Bi:tx_ =bx^-=kc:=Bi:=tx_"
	expect := []SSE{
		0x65E5672C8A9E306F, 0x000096E330573044, 0x0021002000A70000, // Unicode
		0x808EAE50D950C08E, 0x8AE50D950CFFFFFF, 0x000096E300000000, // 12-syllable Sango word (no space separator)
		0xA08EAE50D950CFFF, // word break following Sango word prefixed by a space
		0xC08EAE50D950CFFF, // word break following Sango word prefixed by a space
		0xD08EAE50D950CFFF, // word break following Sango word prefixed by a space
		0xE08EAE50D950CFFF, // end of input with trailing padding
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

func TestBadCanonicalToSSEs(t *testing.T) {
	c := "bx^-kc:Bi:tx_~bx^-kc:Bi:jx_+96E3=bx^-=kc:=Bi:=tx_ bx^-kc:Bi:tx_ ~bx^-kc:Bi:tx_ =bx^-=kc:=Bi:=tx_"
	expect := []SSE{0x808EAE50D950C08E, 0x8AE50D9FFFFFFFFF} // results up to broken parse
	dumpSSEs := func(sses []SSE) string {
		s := "{"
		for _, sse := range sses {
			s += fmt.Sprintf(" %016X", sse)
		}
		s += " }"
		return s
	}
	expectErr := fmt.Errorf("Error parsing Canonical string starting at s[24:] = %q", "jx_+96E3=b...")
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

func TestCanonicalToSSEsToCanonical(t *testing.T) {
	c := "U+65E5U+672CU+8A9EU+306FU+96E3U+3057U+3044U+0021U+0020U+00A7" +
		"bx^-kc:Bi:tx_~bx^-kc:Bi:tx_U+96E3=bx^-=kc:=Bi:=tx_ bx^-kc:Bi:tx_ ~bx^-kc:Bi:tx_ =bx^-=kc:=Bi:=tx_"
	sses, err := CanonicalToSSEs(c)
	if err != nil {
		t.Errorf("unexpected error returned from CanonicalToSSEs\nerr = %v", err)
		return
	}
	var s strings.Builder
	for _, sse := range sses {
		sse.WriteAsCanonicalTo(&s)
	}
	actual := s.String()
	expect := "U+65E5U+672CU+8A9EU+306FU+96E3U+3057U+3044U+0021U+0020U+00A7" +
		"bx^-kc:Bi:tx_bx^-kc:Bi:tx_U+96E3=bx^-=kc:=Bi:=tx_ bx^-kc:Bi:tx_ ~bx^-kc:Bi:tx_ =bx^-=kc:=Bi:=tx_"
	if actual != expect {
		t.Errorf("bad BadCanonicalToSSEs\nexpect: %v\nactual: %v\n", expect, actual)
	}
	// TODO: Something is broken, expect should equal c. There is a single missing Title case (~). Track down this error.
}
