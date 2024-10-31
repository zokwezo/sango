// Example that parses and then serializes a ConLLU file.

package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/brewingweasel/go-conllu"
)

func SerializeConLLU(sentences []conllu.Sentence, out io.Writer) {
	for _, sentence := range sentences {
		fmt.Fprintf(out, "# text = %s\n", sentence.Text)
		for _, t := range sentence.Tokens {
			form := t.Form
			if form == "" {
				form = "_"
			}
			lemma := t.Lemma
			if lemma == "" {
				lemma = "_"
			}
			upos := t.UPOS
			if upos == "" {
				upos = "_"
			}
			xpos := t.XPOS
			if xpos == "" {
				xpos = "_"
			}
			feats := ""
			for _, m := range t.Feats {
				if feats != "" {
					feats += "|"
				}
				feats += m.Feature + "=" + m.Value
			}
			if feats == "" {
				feats = "_"
			}
			deps := ""
			for _, d := range t.Deps {
				if deps != "" {
					deps += "|"
				}
				deps += fmt.Sprintf("%v:%s", d.Head, d.Deprel)
			}
			if deps == "" {
				deps = "_"
			}
			misc := strings.Join(t.Misc, "|")
			if misc == "" {
				misc = "_"
			}
			fmt.Fprintf(out, "%v\t%v\t%v\t%v\t%v\t%v\t%v\t_\t_\t%v\n",
				t.ID, form, lemma, upos, xpos, feats, deps, misc)
		}
		fmt.Fprintln(out)
	}
}

func main() {
	url := "https://raw.githubusercontent.com/zokwezo/sango/refs/heads/main/corpora/les_ruses_de_tere/tere_na_nguru.conllu"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching URL:", err)
		return
	}
	defer resp.Body.Close()

	sentences, errs := conllu.Parse(resp.Body)
	err = errors.Join(errs...)
	if err != nil {
		log.Fatal(err)
	}

	SerializeConLLU(sentences, io.Writer(bufio.NewWriter(os.Stdout)))
}
