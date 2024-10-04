package tokenize

import (
	"bufio"
	"os"

	"github.com/spf13/cobra"
)

func Init(rootCmd *cobra.Command) {
	rootCmd.AddCommand(tokenizeCmd)
}

var (
	tokenizeCmd = &cobra.Command{
		Use:   "tokenize",
		Short: "A CLI to tokenize text into Sango, English, French, punctuation, and whitespace",
		Long:  "https://github.com/zokwezo/sango/blob/main/src/lib/tokenize/README.md",
		Run: func(cmd *cobra.Command, args []string) {
			in := bufio.NewReader(os.Stdin)
			out := bufio.NewWriter(os.Stdout)
			defer out.Flush()
			lemmas := ClassifySango(in)
			for _, lemma := range lemmas {
				if _, err := out.WriteString("{"); err != nil {
					panic(err)
				}
				if _, err := out.WriteString(lemma.Type); err != nil {
					panic(err)
				}
				if lemma.Lang != "" {
					if _, err := out.WriteString(":"); err != nil {
						panic(err)
					}
					if _, err := out.WriteString(lemma.Lang); err != nil {
						panic(err)
					}
				}
				if _, err := out.WriteString("|"); err != nil {
					panic(err)
				}
				if _, err := out.WriteString(lemma.Sango); err != nil {
					panic(err)
				}
				if _, err := out.WriteString("}\n"); err != nil {
					panic(err)
				}
			}
		},
	}
)
