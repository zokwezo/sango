// Converts UTF8 to normalized UTF8.
package main

import (
	"bufio"
	"io"
	"os"

	"golang.org/x/text/unicode/norm"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	r := norm.NFC.Reader(bufio.NewReader(os.Stdin))
	w := bufio.NewWriter(os.Stdout)
	dat, err := io.ReadAll(r)
	s := string(dat)
	check(err)
	_, err = io.WriteString(w, s)
	w.Flush()
	check(err)
}
