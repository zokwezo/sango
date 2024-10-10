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
	for sango, ascii := range sangoAscii {
		if out := fromString(sango, Encode); out != ascii {
			t.Errorf("Encode: %s", sango)
			t.Errorf("Expect: %s", ascii)
			t.Errorf("Actual: %s", out)
		}
	}
}

var sangoAscii = map[string]string{
	`.......I ~ you...said|wrote 3.25, {{Hello}}! \\@/.`: `{DOTS}I` +
		`{SPACE| }{TILDE}{SPACE| }you{PUNC_OTHER|…}said{BAR}wrote{SPACE| }` +
		`{NUMBER|3{PUNC_OTHER|.}25{PUNC_CLOSE|}}{PUNC_OTHER|,}{SPACE| }` +
		`{LEFT_BRACE}{LEFT_BRACE}Hello{RIGHT_BRACE}{RIGHT_BRACE}{PUNC_OTHER|!}` +
		`{SPACE| }{BACKSLASH}{BACKSLASH}{PUNC_OTHER|@}{SLASH}{PUNC_OTHER|.}`,
}
