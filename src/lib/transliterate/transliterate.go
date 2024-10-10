package transliterate

import (
	"bufio"
	"io"
	"regexp"

	"golang.org/x/text/unicode/norm"
)

type replacer = struct {
	re   *regexp.Regexp
	repl []byte
}

var (
	left  string = "\ufff9"
	mid   string = "\ufffa"
	right string = "\ufffb"

	replacers = []replacer{
		{regexp.MustCompile(`__+`), []byte(left + `UNDERSCORES` + right)},
		{regexp.MustCompile(`_`), []byte(left + `UNDERSCORE` + right)},
		{regexp.MustCompile(`(\p{Pc})`), []byte(left + `PUNC_CONNECTOR` + mid + `$1` + right)},
		{regexp.MustCompile(`{`), []byte(left + `LEFT_BRACE` + right)},
		{regexp.MustCompile(`}`), []byte(left + `RIGHT_BRACE` + right)},
		{regexp.MustCompile(`\|`), []byte(left + `BAR` + right)},
		{regexp.MustCompile(`(\.{4,})`), []byte(left + `DOTS` + right)},
		{regexp.MustCompile(`\.{3}`), []byte(`…`)},
		{regexp.MustCompile(`\.{2}`), []byte(`‥`)},
		{regexp.MustCompile(`(\-{4,})`), []byte(left + `HLINE` + mid + `$1` + right)},
		{regexp.MustCompile(`\-{3}`), []byte(`—`)},
		{regexp.MustCompile(`\-{2}`), []byte(`–`)},
		{regexp.MustCompile(`(\*{2,})`), []byte(left + `STARS` + mid + `$1` + right)},
		{regexp.MustCompile(`\*`), []byte(left + `STAR}`)},
		{regexp.MustCompile(`\\`), []byte(left + `BACKSLASH` + right)},
		{regexp.MustCompile(`/`), []byte(left + `SLASH` + right)},
		{regexp.MustCompile(`<`), []byte(left + `LT` + right)},
		{regexp.MustCompile(`>`), []byte(left + `GT` + right)},
		{regexp.MustCompile(`\[`), []byte(left + `LEFT_BRACKET` + right)},
		{regexp.MustCompile(`\]`), []byte(left + `RIGHT_BRACKET` + right)},
		{regexp.MustCompile(`=`), []byte(left + `EQ` + right)},
		{regexp.MustCompile(`~`), []byte(left + `TILDE` + right)},
		{regexp.MustCompile(`\+`), []byte(left + `PLUS` + right)},
		{regexp.MustCompile(`\$`), []byte(left + `DOLLAR` + right)},
		{regexp.MustCompile(`\^`), []byte(left + `HAT` + right)},
		{regexp.MustCompile(`(\p{Nd}+(?:\.\p{Nd}+)?)`), []byte(left + `NUMBER` + mid + `$1}`)},
		{regexp.MustCompile(`(\p{Pi})`), []byte(left + `PUNC_INITIAL` + mid + `$1` + right)},
		{regexp.MustCompile(`(\p{Pf})`), []byte(left + `PUNC_FINAL` + mid + `$1` + right)},
		{regexp.MustCompile(`(\p{Pd})`), []byte(left + `PUNC_DASH` + mid + `$1` + right)},
		{regexp.MustCompile(`(\p{Po})`), []byte(left + `PUNC_OTHER` + mid + `$1` + right)},
		{regexp.MustCompile(`(\p{Ps})`), []byte(left + `PUNC_OPEN` + mid + `$1` + right)},
		{regexp.MustCompile(`(\p{Pe})`), []byte(left + `PUNC_CLOSE` + mid + `$1` + right)},
		{regexp.MustCompile(`(\p{Z}+)`), []byte(left + `SPACE` + mid + `$1}`)},
		{regexp.MustCompile(left), []byte(`{`)},
		{regexp.MustCompile(mid), []byte(`|`)},
		{regexp.MustCompile(right), []byte(`}`)},
	}
)

func Encode(out *bufio.Writer, in *bufio.Reader) error {
	defer out.Flush()
	b, err := io.ReadAll(norm.NFKD.Reader(in))
	if err != nil {
		return err
	}

	for _, r := range replacers {
		b = r.re.ReplaceAll(b, r.repl)
	}

	_, err = out.Write(b)
	return err
}

func Decode(out *bufio.Writer, in *bufio.Reader) error {
	defer out.Flush()
	b, err := io.ReadAll(in)
	if err != nil {
		return err
	}

	w := norm.NFKC.Writer(out)
	_, err = w.Write(b)
	w.Close()
	return err
}
