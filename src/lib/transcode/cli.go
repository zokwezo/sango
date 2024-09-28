package transcode

import (
	"bufio"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func Init(rootCmd *cobra.Command) {
	encodeCmd.Flags().BoolVar(&useJForPitchFlagValue, "useJForPitch", false, "Use the letter 'j' instead of uppercase for pitch accent.")
	rootCmd.AddCommand(transcodeCmd)
	transcodeCmd.AddCommand(normalizeCmd)
	transcodeCmd.AddCommand(encodeCmd)
	transcodeCmd.AddCommand(decodeCmd)
}

var (
	useJForPitchFlagValue bool

	transcodeCmd = &cobra.Command{
		Use:   "transcode",
		Short: "A CLI to transcode Sango between UTF8 and ASCII",
		Long:  "https://github.com/zokwezo/sango/blob/main/src/lib/transcode/README.md",
	}

	normalizeCmd = &cobra.Command{
		Use:   "normalize",
		Short: "Read from stdin, normalize UTF8 into NFKC, then write to stdout",
		Args:  cobra.MaximumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			if err := Normalize(bufio.NewWriter(os.Stdout), bufio.NewReader(os.Stdin)); err != nil {
				log.Fatal(err)
			}
		},
	}

	encodeCmd = &cobra.Command{
		Use:   "encode",
		Short: "Read from stdin, encode UTF8 into ASCII format, then write to stdout",
		Run: func(cmd *cobra.Command, args []string) {
			if err := Encode(bufio.NewWriter(os.Stdout), bufio.NewReader(os.Stdin), useJForPitchFlagValue); err != nil {
				log.Fatal(err)
			}
		},
	}

	decodeCmd = &cobra.Command{
		Use:   "decode",
		Short: "Read from stdin, decode ASCII into UTF8 format, then write to stdout",
		Run: func(cmd *cobra.Command, args []string) {
			if err := Decode(bufio.NewWriter(os.Stdout), bufio.NewReader(os.Stdin)); err != nil {
				log.Fatal(err)
			}
		},
	}
)
