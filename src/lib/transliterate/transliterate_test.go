package transliterate

import (
	"bufio"
	"bytes"
	"testing"
)

func fromString(in string, transliterate func(*bufio.Writer, *bufio.Reader) error) string {
	var out bytes.Buffer
	w := bufio.NewWriter(&out)
	r := bufio.NewReader(bytes.NewBufferString(in))
	if err := transliterate(w, r); err != nil {
		panic(err)
	}
	w.Flush()
	return out.String()
}

func TestTransliterate(t *testing.T) {
	asciiExpect := "``He said `What...' a-b c--d e---f g==>h h<==g xJc.ej.oq XJC.Ej.Oq'' <<,C'est vrai?>>"
	sangoExpect := "“He said ‘What…’ a-b ɔ–d e—f g⟹h h⟸g ɛ̣ɔə̂ø̈ Ɛ̣ƆƏ̂Ø̈” «Ç'est vrai?»"
	sangoActual := fromString(asciiExpect, Decode)
	if sangoActual != sangoExpect {
		t.Errorf("Sango actual = %s", sangoActual)
		t.Errorf("Sango expect = %s", sangoExpect)
	}
	asciiActual := fromString(sangoExpect, Encode)
	if asciiActual != asciiExpect {
		t.Errorf("ASCII actual = %s", asciiActual)
		t.Errorf("ASCII expect = %s", asciiExpect)
	}
}
