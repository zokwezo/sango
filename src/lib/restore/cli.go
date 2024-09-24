package restore

import (
	"bufio"
	"io"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/text/unicode/norm"
)

func Init(rootCmd *cobra.Command) {
	rootCmd.AddCommand(restoreCmd)
}

var (
	restoreCmd = &cobra.Command{
		Use:   "restore",
		Short: "A CLI to restore vowel height and pitch to Sango text",
		Run: func(cmd *cobra.Command, args []string) {
			in := bufio.NewReader(os.Stdin)
			r := norm.NFKC.Reader(in)
			b, err := io.ReadAll(r)
			if err != nil {
				panic(err)
			}
			out := bufio.NewWriter(os.Stdout)
			defer out.Flush()
			s := string(b)

			restored := RestoreSangoVowels(s)
			if _, err := out.WriteString(restored); err != nil {
				panic(err)
			}
		},
	}
)
