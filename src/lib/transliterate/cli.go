package transliterate

import (
	"bufio"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func Init(rootCmd *cobra.Command) {
	rootCmd.AddCommand(transliterateCmd)
	transliterateCmd.AddCommand(encodeCmd)
	transliterateCmd.AddCommand(decodeCmd)
	transliterateCmd.AddCommand(normalizeCmd)
	transliterateCmd.AddCommand(unnormalizeCmd)
}

var (
	transliterateCmd = &cobra.Command{
		Use:   "transliterate",
		Short: "A CLI to transliterate Sango between UTF8 and ASCII",
		Long:  "https://github.com/zokwezo/sango/blob/main/src/lib/transliterate/README.md",
	}

	encodeCmd = &cobra.Command{
		Use:   "encode",
		Short: "Read from stdin, encode ASCII into UTF8 format, then write to stdout",
		Run: func(cmd *cobra.Command, args []string) {
			in := bufio.NewReader(os.Stdin)
			out := bufio.NewWriter(os.Stdout)
			if err := Encode(out, in); err != nil {
				log.Fatal(err)
			}
			out.Flush()
		},
	}

	decodeCmd = &cobra.Command{
		Use:   "decode",
		Short: "Read from stdin, decode ASCII into UTF8 format, then write to stdout",
		Run: func(cmd *cobra.Command, args []string) {
			in := bufio.NewReader(os.Stdin)
			out := bufio.NewWriter(os.Stdout)
			if err := Decode(out, in); err != nil {
				log.Fatal(err)
			}
			out.Flush()
		},
	}

	normalizeCmd = &cobra.Command{
		Use:   "normalize",
		Short: "Read from stdin, normalize ASCII into UTF8 NFC format, then write to stdout",
		Run: func(cmd *cobra.Command, args []string) {
			in := bufio.NewReader(os.Stdin)
			out := bufio.NewWriter(os.Stdout)
			if err := Normalize(out, in); err != nil {
				log.Fatal(err)
			}
			out.Flush()
		},
	}

	unnormalizeCmd = &cobra.Command{
		Use:   "unnormalize",
		Short: "Read from stdin, unnormalize ASCII into UTF8 NFD format, then write to stdout",
		Run: func(cmd *cobra.Command, args []string) {
			in := bufio.NewReader(os.Stdin)
			out := bufio.NewWriter(os.Stdout)
			if err := Unnormalize(out, in); err != nil {
				log.Fatal(err)
			}
			out.Flush()
		},
	}
)
