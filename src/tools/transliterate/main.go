package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/zokwezo/sango/src/lib/transliterate"
)

var (
	transliterateCmd = &cobra.Command{
		Use:   "transliterate",
		Short: "A CLI to transliterate Sango between UTF8 and ASCII",
		Long:  "https://github.com/zokwezo/sango/blob/main/src/lib/transliterate/README.md",
		Args:  cobra.MaximumNArgs(0),
	}

	encodeCmd = &cobra.Command{
		Use:   "encode",
		Short: "Read from stdin, encode UTF8 into ASCII, then write to stdout",
		Long:  "https://github.com/zokwezo/sango/blob/main/src/lib/transliterate/README.md",
		Args:  cobra.MaximumNArgs(0),
	}

	encodeInputCmd = &cobra.Command{
		Use:   "input",
		Short: "Read from stdin, encode UTF8 into ASCII input, then write to stdout",
		Long:  "https://github.com/zokwezo/sango/blob/main/src/lib/transliterate/README.md",
		Run: func(cmd *cobra.Command, args []string) {
			transliterate.EncodeInput()
		},
	}

	encodeOutputCmd = &cobra.Command{
		Use:   "output",
		Short: "Read from stdin, encode UTF8 into ASCII output, then write to stdout",
		Long:  "https://github.com/zokwezo/sango/blob/main/src/lib/transliterate/README.md",
		Run: func(cmd *cobra.Command, args []string) {
			transliterate.EncodeOutput()
		},
	}

	decodeCmd = &cobra.Command{
		Use:   "decode",
		Short: "Read from stdin, decode ASCII into UTF8, then write to stdout",
		Long:  "https://github.com/zokwezo/sango/blob/main/src/lib/transliterate/README.md",
		Args:  cobra.MaximumNArgs(0),
	}

	decodeInputCmd = &cobra.Command{
		Use:   "input",
		Short: "Read from stdin, decode ASCII into UTF8 input, then write to stdout",
		Long:  "https://github.com/zokwezo/sango/blob/main/src/lib/transliterate/README.md",
		Run: func(cmd *cobra.Command, args []string) {
			transliterate.DecodeInput()
		},
	}

	decodeOutputCmd = &cobra.Command{
		Use:   "output",
		Short: "Read from stdin, decode ASCII into UTF8 output, then write to stdout",
		Long:  "https://github.com/zokwezo/sango/blob/main/src/lib/transliterate/README.md",
		Run: func(cmd *cobra.Command, args []string) {
			transliterate.DecodeOutput()
		},
	}
)

func init() {
	transliterateCmd.AddCommand(encodeCmd)
	encodeCmd.AddCommand(encodeInputCmd)
	encodeCmd.AddCommand(encodeOutputCmd)

	transliterateCmd.AddCommand(decodeCmd)
	decodeCmd.AddCommand(decodeInputCmd)
	decodeCmd.AddCommand(decodeOutputCmd)
}

func main() {
	log.SetFlags(log.Lshortfile)
	if err := transliterateCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
