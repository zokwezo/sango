package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/zokwezo/sango/src/lib/lexicon"
	"github.com/zokwezo/sango/src/lib/restore"
	"github.com/zokwezo/sango/src/lib/tokenize"
	"github.com/zokwezo/sango/src/lib/transcode"
	"github.com/zokwezo/sango/src/lib/transliterate"
)

var (
	sangoCmd = &cobra.Command{
		Use:   "sango",
		Short: "A CLI to run Sango language tools and server",
		Long:  "https://github.com/zokwezo/sango/blob/main/src/README.md",
	}
)

func init() {
	lexicon.Init(sangoCmd)
	restore.Init(sangoCmd)
	tokenize.Init(sangoCmd)
	transcode.Init(sangoCmd)
	transliterate.Init(sangoCmd)
}

func main() {
	log.SetFlags(log.Lshortfile)
	if err := sangoCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
