package transcode

import (
	"bufio"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func Init(rootCmd *cobra.Command) {
	rootCmd.AddCommand(transcodeCmd)
	transcodeCmd.AddCommand(encodeCmd)
	transcodeCmd.AddCommand(decodeCmd)
}

var (
	transcodeCmd = &cobra.Command{
		Use:   "transcode",
		Short: "A CLI to transcode Sango between UTF8 and ASCII",
		Long:  "https://github.com/zokwezo/sango/blob/main/src/lib/transcode/README.md",
	}

	encodeCmd = &cobra.Command{
		Use:   "encode",
		Short: "Read from stdin, encode UTF8 into SSE tokens, then write to stdout",
		Run: func(cmd *cobra.Command, args []string) {
			if err := Encode(bufio.NewWriter(os.Stdout), bufio.NewReader(os.Stdin)); err != nil {
				log.Fatal(err)
			}
		},
	}

	decodeCmd = &cobra.Command{
		Use:   "decode",
		Short: "Read from stdin, decode SSE tokens into UTF8, then write to stdout",
		Run: func(cmd *cobra.Command, args []string) {
			if err := Decode(bufio.NewWriter(os.Stdout), bufio.NewReader(os.Stdin)); err != nil {
				log.Fatal(err)
			}
		},
	}
)
