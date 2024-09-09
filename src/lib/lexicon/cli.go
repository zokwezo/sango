package lexicon

import (
	"fmt"
	"log"
	"regexp"

	"github.com/spf13/cobra"
)

func Init(rootCmd *cobra.Command) {
	lookupCmd.Flags().StringVar(&tonelessFlagValue, "toneless", "", "Returns values only where this regexp partially matches toneless.")
	lookupCmd.Flags().StringVar(&sangoFlagValue, "sango", "", "Returns values only where this regexp partially matches sango.")
	lookupCmd.Flags().StringVar(&lexPosFlagValue, "lex_pos", "", "Returns values only where this regexp partially matches lexPos.")
	lookupCmd.Flags().StringVar(&udPosFlagValue, "ud_os", "", "Returns values only where this regexp partially matches uDPos.")
	lookupCmd.Flags().StringVar(&udFeatureFlagValue, "ud_feature", "", "Returns values only where this regexp partially matches uDFeature.")
	lookupCmd.Flags().StringVar(&categoryFlagValue, "category", "", "Returns values only where this regexp partially matches category.")
	lookupCmd.Flags().StringVar(&englishFlagValue, "english", "", "Returns values only where this regexp partially matches english.")
	lookupCmd.Flags().IntVar(&frequencyMinFlagValue, "frequency_min", 1, "Returns values only where frequency_min <= row.frequency.")
	lookupCmd.Flags().IntVar(&frequencyMaxFlagValue, "frequency_max", 9, "Returns values only where frequency_max >= row.frequency.")
	lexiconCmd.AddCommand(lookupCmd)
	rootCmd.AddCommand(lexiconCmd)
}

var (
	tonelessFlagValue     string
	sangoFlagValue        string
	lexPosFlagValue       string
	udPosFlagValue        string
	udFeatureFlagValue    string
	categoryFlagValue     string
	englishFlagValue      string
	frequencyMinFlagValue int
	frequencyMaxFlagValue int

	lexiconCmd = &cobra.Command{
		Use:   "lexicon",
		Short: "A CLI to lexicon Sango between UTF8 and ASCII",
		Long:  "https://github.com/zokwezo/sango/blob/main/src/lib/lexicon/README.md",
	}

	lookupCmd = &cobra.Command{
		Use:   "lookup",
		Short: "Read from stdin, lookup UTF8 into NFKC, then write to stdout",
		Args:  cobra.MaximumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {

			f := DictRowRegexp{
				TonelessRE:   regexp.MustCompile(tonelessFlagValue),
				SangoRE:      regexp.MustCompile(sangoFlagValue),
				LexPosRE:     regexp.MustCompile(lexPosFlagValue),
				UDPosRE:      regexp.MustCompile(udPosFlagValue),
				UDFeatureRE:  regexp.MustCompile(udFeatureFlagValue),
				CategoryRE:   regexp.MustCompile(categoryFlagValue),
				EnglishRE:    regexp.MustCompile(englishFlagValue),
				FrequencyMin: frequencyMinFlagValue,
				FrequencyMax: frequencyMaxFlagValue,
			}

			dictRows, err := Lookup(LexiconRows(), f)
			if err != nil {
				log.Fatal(err)
			}
			for k, row := range dictRows {
				fmt.Printf("row[%v] = %v\n", k, row)
			}
		},
	}
)
