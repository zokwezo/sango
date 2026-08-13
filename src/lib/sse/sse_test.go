package sse

import (
	"testing"
)

func TestExtractBits(t *testing.T) {
	var xx = []struct {
		src    uint64
		lsb    int
		msb    int
		expect uint64
	}{
		{src: 0b0100100011, lsb: -3, msb: 888, expect: 0b0100100011},
		{src: 0b0100100011, lsb: 888, msb: -3, expect: 0},
		{src: 0b0100100011, lsb: 1, msb: 8, expect: 0b10010001},
		{src: 0b0100100011, lsb: 8, msb: 1, expect: 0},
		{src: 0b0100100011, lsb: 0, msb: 5, expect: 0b100011},
		{src: 0b0100100011, lsb: 5, msb: 0, expect: 0},
		{src: 10548167345049833471, lsb: 48, msb: 59, expect: 0b001001100010},
	}
	for k, x := range xx {
		actual := extractBits(x.src, x.lsb, x.msb)
		if actual != x.expect {
			t.Errorf("In x #%v, from extractBits(%b, %v, %v),\nexpect: %b\nactual: %b\n",
				k, x.src, x.lsb, x.msb, x.expect, actual)
		}
	}
}

func TestUTF8FromUnicodeSSEs(t *testing.T) {
	sses := [...]SSE{
		0x00e697a5e69cac00,
		0x00e8aa9ee381af00,
		0x00e99ba3e3819700,
		0x00e3818421000000,
	}
	expect := "日本語は難しい!"
	actual := UTF8FromSSEs(sses[:])
	if actual != expect {
		t.Errorf("From UTF8FromSSE(\n%#v\n),\nexpect: %#v\nactual: %#v\n\n",
			sses, expect, actual)
	}
}

func TestUTF8FromSangoSSEs(t *testing.T) {
	sses := [...]SSE{0x9_062_BE5_451_320_FFF}
	expect := "bɛ̂-kɔ̈mbïtɛ"
	actual := UTF8FromSSEs(sses[:])
	if actual != expect {
		t.Errorf("From UTF8FromSSE(\n%#v\n),\nexpect: %#v\nactual: %#v\n\n",
			sses, expect, actual)
	}
}

func TestCanonicalFromUnicodeSSEs(t *testing.T) {
	sses := [...]SSE{
		0x00e697a5e69cac00,
		0x00e8aa9ee381af00,
		0x00e99ba3e3819700,
		0x00e3818421000000,
	}
	expectUTF8 := "E697A5E69CACE8AA9EE381AFE99BA3E38197E3818421"
	actualUTF8 := CanonicalFromSSEs(sses[:])
	if actualUTF8 != expectUTF8 {
		t.Errorf("From CanonicalFromSSE(\n%#v\n),\nexpect: %#v\nactual: %#v\n\n",
			sses, expectUTF8, actualUTF8)
	}
}

func TestCanonicalFromSangoSSEs(t *testing.T) {
	sses := [...]SSE{0x9_062_BE5_451_320_FFF}
	//expectUTF8 := "Bɛ̂-kɔ̈mbïtɛ"
	expect := "bx^-kc:Bi:tx_"
	actual := CanonicalFromSSEs(sses[:])
	if actual != expect {
		t.Errorf("From CanonicalFromSSE(\n%#v\n),\nexpect: %#v\nactual: %#v\n\n",
			sses, expect, actual)
	}
}
