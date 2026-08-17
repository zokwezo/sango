package sse

import (
	"strings"
	"testing"
)

func TestWriteAsUTF8(t *testing.T) {
	sses := [...]SSE{
		0x65E5_672C_8A9E_306F,
		0x0000_96E3_3057_3044,
		0x0021_0000_0020_00A7,
		0x8_062_BE5_451_320_FFF,
		0x9_062_BE5_451_320_FFF,
		0xA_062_BE5_451_320_FFF,
		0xB_062_BE5_451_320_FFF,
		0xC_062_BE5_451_320_FFF,
		0xD_062_BE5_451_320_FFF,
		0xE_062_BE5_451_320_FFF,
		0xF_062_BE5_451_320_FFF,
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
		0x8_062_BE5_451_320_FFF,
		0x9_062_BE5_451_320_FFF,
		0xA_062_BE5_451_320_FFF,
		0xB_062_BE5_451_320_FFF,
		0xC_062_BE5_451_320_FFF,
		0xD_062_BE5_451_320_FFF,
		0xE_062_BE5_451_320_FFF,
		0xF_062_BE5_451_320_FFF,
	}
	var s strings.Builder
	for _, sse := range sses {
		sse.WriteAsCanonicalTo(&s)
	}
	expect := "U+65E5U+672CU+8A9EU+306FU+96E3U+3057U+3044U+0021U+0020U+00A7" +
		"bx^-kc:Bi:tx_~bx^-kc:Bi:tx_=bx^-=kc:=Bi:=tx_ bx^-kc:Bi:tx_ ~bx^-kc:Bi:tx_ =bx^-=kc:=Bi:=tx_"
	actual := s.String()
	if actual != expect {
		t.Errorf("From CanonicalFromSSE(\n%#v\n),\nexpect: %#v\nactual: %#v\n\n",
			sses, expect, actual)
	}
}
