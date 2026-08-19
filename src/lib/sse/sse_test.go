package sse

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func TestWriteAsUTF8(t *testing.T) {
	sses := [...]SSE{
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
	sses := [...]SSE{
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

func TestCanonicalToSyllables(t *testing.T) {
	s := "U+65E5U+672CU+8A9EU+306FU+96E3U+3057U+3044U+0021U+0020U+00A7" +
		"bx^-kc:Bi:tx_~bx^-kc:Bi:tx_=bx^-=kc:=Bi:=tx_ bx^-kc:Bi:tx_ ~bx^-kc:Bi:tx_ =bx^-=kc:=Bi:=tx_"
	expect, expectLen := []uint16{
		0x65E5, 0x672C, 0x8A9E, 0x306F, 0x96E3, 0x3057, 0x3044, 0x0021, 0x0020, 0x00A7,
		0x008E, 0x0AE5, 0x00D9, 0x050C,
		0x108E, 0x0AE5, 0x00D9, 0x050C,
		0x208E, 0x2AE5, 0x20D9, 0x250C,
		0x408E, 0x0AE5, 0x00D9, 0x050C,
		0x508E, 0x0AE5, 0x00D9, 0x050C,
		0x608E, 0x2AE5, 0x20D9, 0x250C,
	}, len(s)
	actual, actualLen := canonicalToSyllables(s)
	if actualLen != expectLen {
		fmt.Printf("error found at s[%v](good: %q bad: %q\n", actualLen, s[0:actualLen], s[actualLen:])
	}
	actualHex := "{"
	for _, x := range actual {
		actualHex += ", 0x" + strconv.FormatUint(uint64(x), 16)
	}
	actualHex += "}"
	expectHex := "{"
	for _, x := range expect {
		expectHex += ", 0x" + strconv.FormatUint(uint64(x), 16)
	}
	expectHex += "}"
	if actualHex != expectHex {
		fmt.Printf("bad canonicalToSyllables\nexpect: %v\nactual: %v\n", expectHex, actualHex)
	}
}

func TestSyllablesToSSEs(t *testing.T) {
	syllables := []uint16{
		0x65E5, 0x672C, 0x8A9E, 0x306F, 0x96E3, 0x3057, 0x3044, 0x0021, 0x0020, 0x00A7,
		0x008E, 0x0AE5, 0x00D9, 0x050C,
		0x108E, 0x0AE5, 0x00D9, 0x050C,
		0x208E, 0x2AE5, 0x20D9, 0x250C,
		0x408E, 0x0AE5, 0x00D9, 0x050C,
		0x508E, 0x0AE5, 0x00D9, 0x050C,
		0x608E, 0x2AE5, 0x20D9, 0x250C,
	}
	expect := []SSE{} // TODO: Update with real result after SyllablesToSSEs is implemented
	actual := SyllablesToSSEs(syllables)
	actualHex := "{"
	for _, x := range actual {
		actualHex += ", 0x" + strconv.FormatUint(uint64(x), 16)
	}
	actualHex += "}"
	expectHex := "{"
	for _, x := range expect {
		expectHex += ", 0x" + strconv.FormatUint(uint64(x), 16)
	}
	expectHex += "}"
	if actualHex != expectHex {
		fmt.Printf("bad SyllablesToSSEs\nexpect: %v\nactual: %v\n", expectHex, actualHex)
	}
}
