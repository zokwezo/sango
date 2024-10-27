package transliterate

import (
	"bufio"
	"bytes"
	"io"

	"golang.org/x/text/unicode/norm"
)

func Normalize(out *bufio.Writer, in *bufio.Reader) error {
	defer out.Flush()
	b, err := io.ReadAll(norm.NFC.Reader(in))
	if err != nil {
		return err
	}
	_, err = out.Write(b)
	return err
}

func Unnormalize(out *bufio.Writer, in *bufio.Reader) error {
	defer out.Flush()
	b, err := io.ReadAll(norm.NFD.Reader(in))
	if err != nil {
		return err
	}
	_, err = out.Write(b)
	return err
}

func Encode(out *bufio.Writer, in *bufio.Reader) error {
	defer out.Flush()
	b, err := io.ReadAll(norm.NFD.Reader(in))
	if err != nil {
		return err
	}
	b = bytes.ReplaceAll(b, []byte("\u0308"), []byte("q"))
	b = bytes.ReplaceAll(b, []byte("\u0302"), []byte("j"))
	b = bytes.ReplaceAll(b, []byte("\u0323"), []byte("J"))
	b = bytes.ReplaceAll(b, []byte("ø"), []byte(".o"))
	b = bytes.ReplaceAll(b, []byte("Ø"), []byte(".O"))
	b = bytes.ReplaceAll(b, []byte("ə"), []byte(".e"))
	b = bytes.ReplaceAll(b, []byte("Ə"), []byte(".E"))
	b = bytes.ReplaceAll(b, []byte("ç"), []byte(",ɔ"))
	b = bytes.ReplaceAll(b, []byte("Ç"), []byte(",Ɔ"))
	b = bytes.ReplaceAll(b, []byte("ɔ"), []byte("c"))
	b = bytes.ReplaceAll(b, []byte("Ɔ"), []byte("C"))
	b = bytes.ReplaceAll(b, []byte("ɛ"), []byte("x"))
	b = bytes.ReplaceAll(b, []byte("Ɛ"), []byte("X"))
	b = bytes.ReplaceAll(b, []byte("⟹"), []byte("==>"))
	b = bytes.ReplaceAll(b, []byte("⟸"), []byte("<=="))
	b = bytes.ReplaceAll(b, []byte("–"), []byte("--"))
	b = bytes.ReplaceAll(b, []byte("—"), []byte("---"))
	b = bytes.ReplaceAll(b, []byte("…"), []byte("..."))
	b = bytes.ReplaceAll(b, []byte("»"), []byte(">>"))
	b = bytes.ReplaceAll(b, []byte("«"), []byte("<<"))
	b = bytes.ReplaceAll(b, []byte("’ "), []byte("' "))
	b = bytes.ReplaceAll(b, []byte(" ‘"), []byte(" \x60"))
	b = bytes.ReplaceAll(b, []byte("”"), []byte("''"))
	b = bytes.ReplaceAll(b, []byte("“"), []byte("\x60\x60"))
	w := norm.NFC.Writer(out)
	_, err = w.Write(b)
	w.Close()
	return err
}

func Decode(out *bufio.Writer, in *bufio.Reader) error {
	defer out.Flush()
	b, err := io.ReadAll(norm.NFD.Reader(in))
	if err != nil {
		return err
	}
	b = bytes.ReplaceAll(b, []byte("\x60\x60"), []byte("“"))
	b = bytes.ReplaceAll(b, []byte("''"), []byte("”"))
	b = bytes.ReplaceAll(b, []byte(" \x60"), []byte(" ‘"))
	b = bytes.ReplaceAll(b, []byte(" '"), []byte(" ‘"))
	b = bytes.ReplaceAll(b, []byte("' "), []byte("’ "))
	b = bytes.ReplaceAll(b, []byte("<<"), []byte("«"))
	b = bytes.ReplaceAll(b, []byte(">>"), []byte("»"))
	b = bytes.ReplaceAll(b, []byte("..."), []byte("…"))
	b = bytes.ReplaceAll(b, []byte("---"), []byte("—"))
	b = bytes.ReplaceAll(b, []byte("--"), []byte("–"))
	b = bytes.ReplaceAll(b, []byte("<=="), []byte("⟸"))
	b = bytes.ReplaceAll(b, []byte("==>"), []byte("⟹"))
	b = bytes.ReplaceAll(b, []byte("X"), []byte("Ɛ"))
	b = bytes.ReplaceAll(b, []byte("x"), []byte("ɛ"))
	b = bytes.ReplaceAll(b, []byte("C"), []byte("Ɔ"))
	b = bytes.ReplaceAll(b, []byte("c"), []byte("ɔ"))
	b = bytes.ReplaceAll(b, []byte(",Ɔ"), []byte("Ç"))
	b = bytes.ReplaceAll(b, []byte(",ɔ"), []byte("ç"))
	b = bytes.ReplaceAll(b, []byte(".E"), []byte("Ə"))
	b = bytes.ReplaceAll(b, []byte(".e"), []byte("ə"))
	b = bytes.ReplaceAll(b, []byte(".O"), []byte("Ø"))
	b = bytes.ReplaceAll(b, []byte(".o"), []byte("ø"))
	b = bytes.ReplaceAll(b, []byte("J"), []byte("\u0323"))
	b = bytes.ReplaceAll(b, []byte("j"), []byte("\u0302"))
	b = bytes.ReplaceAll(b, []byte("q"), []byte("\u0308"))
	_, err = out.Write(b)
	return err
}
