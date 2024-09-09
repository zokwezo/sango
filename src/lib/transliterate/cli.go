package transliterate

import (
	"github.com/spf13/cobra"
)

func Init(rootCmd *cobra.Command) {
	rootCmd.AddCommand(transliterateCmd)

	transliterateCmd.AddCommand(normalizeCmd)

	transliterateCmd.AddCommand(encodeCmd)
	encodeCmd.AddCommand(encodeInputCmd)
	encodeCmd.AddCommand(encodeOutputCmd)

	transliterateCmd.AddCommand(decodeCmd)
	decodeCmd.AddCommand(decodeInputCmd)
	decodeCmd.AddCommand(decodeOutputCmd)
}

var (
	transliterateCmd = &cobra.Command{
		Use:   "transliterate",
		Short: "A CLI to transliterate Sango between UTF8 and ASCII",
		Long:  "https://github.com/zokwezo/sango/blob/main/src/lib/transliterate/README.md",
	}

	normalizeCmd = &cobra.Command{
		Use:   "normalize",
		Short: "XRead from stdin, normalize UTF8 into NFKC, then write to stdout",
		Args:  cobra.MaximumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			Normalize()
		},
	}

	encodeCmd = &cobra.Command{
		Use:   "encode",
		Short: "Read from stdin, encode UTF8 into ASCII format, then write to stdout",
	}

	encodeInputCmd = &cobra.Command{
		Use:   "input",
		Short: "Read from stdin, encode UTF8 into ASCII input format, then write to stdout",
		Run: func(cmd *cobra.Command, args []string) {
			EncodeInput()
		},
	}

	encodeOutputCmd = &cobra.Command{
		Use:   "output",
		Short: "Read from stdin, encode UTF8 into ASCII output format, then write to stdout",
		Run: func(cmd *cobra.Command, args []string) {
			EncodeOutput()
		},
	}

	decodeCmd = &cobra.Command{
		Use:   "decode",
		Short: "Read from stdin, decode ASCII into UTF8 format, then write to stdout",
		Args:  cobra.MaximumNArgs(0),
	}

	decodeInputCmd = &cobra.Command{
		Use:   "input",
		Short: "Read from stdin, decode ASCII into UTF8 input format, then write to stdout",
		Run: func(cmd *cobra.Command, args []string) {
			DecodeInput()
		},
	}

	decodeOutputCmd = &cobra.Command{
		Use:   "output",
		Short: "Read from stdin, decode ASCII into UTF8 output, then write to stdout",
		Run: func(cmd *cobra.Command, args []string) {
			DecodeOutput()
		},
	}
)
