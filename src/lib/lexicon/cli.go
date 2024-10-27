package lexicon

import (
	"fmt"
	"log"
	"regexp"

	"github.com/spf13/cobra"
)

func Init(rootCmd *cobra.Command) {
	lookupCmd.Flags().StringVar(&tonelessFlagValue, "toneless", "", "Returns values only where this regexp partially matches toneless.")
	lookupCmd.Flags().StringVar(&lemmaFlagValue, "lemma", "", "Returns values only where this regexp partially matches lemma.")
	lookupCmd.Flags().StringVar(&udPosFlagValue, "ud_os", "", "Returns values only where this regexp partially matches uDPos.")
	lookupCmd.Flags().StringVar(&udFeatureFlagValue, "ud_feature", "", "Returns values only where this regexp partially matches uDFeature.")
	lookupCmd.Flags().StringVar(&categoryFlagValue, "category", "", "Returns values only where this regexp partially matches category.")
	lookupCmd.Flags().StringVar(&englishTranslationFlagValue, "english_translation", "", "Returns values only where this regexp partially matches english translation.")
	lookupCmd.Flags().StringVar(&englishDefinitionFlagValue, "english_definition", "", "Returns values only where this regexp partially matches english definition.")
	lookupCmd.Flags().IntVar(&frequencyMinFlagValue, "frequency_min", 1, "Returns values only where frequency_min <= row.frequency.")
	lookupCmd.Flags().IntVar(&frequencyMaxFlagValue, "frequency_max", 9, "Returns values only where frequency_max >= row.frequency.")
	lexiconCmd.AddCommand(lookupCmd)
	rootCmd.AddCommand(lexiconCmd)
}

var (
	tonelessFlagValue           string
	lemmaFlagValue              string
	udPosFlagValue              string
	udFeatureFlagValue          string
	categoryFlagValue           string
	englishTranslationFlagValue string
	englishDefinitionFlagValue  string
	frequencyMinFlagValue       int
	frequencyMaxFlagValue       int

	lexiconCmd = &cobra.Command{
		Use:   "lexicon",
		Short: "A CLI to interact with the Sango lexicon",
		Long:  "https://github.com/zokwezo/sango/blob/main/src/lib/lexicon/README.md",
	}

	lookupCmd = &cobra.Command{
		Use:   "lookup",
		Short: "Read from stdin, lookup UTF8 into NFKC, then write to stdout",
		Args:  cobra.MaximumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {

			f := DictRowRegexp{
				TonelessRE:           regexp.MustCompile(tonelessFlagValue),
				LemmaRE:              regexp.MustCompile(lemmaFlagValue),
				UDPosRE:              regexp.MustCompile(udPosFlagValue),
				UDFeatureRE:          regexp.MustCompile(udFeatureFlagValue),
				CategoryRE:           regexp.MustCompile(categoryFlagValue),
				EnglishTranslationRE: regexp.MustCompile(englishTranslationFlagValue),
				EnglishDefinitionRE:  regexp.MustCompile(englishDefinitionFlagValue),
				FrequencyMin:         frequencyMinFlagValue,
				FrequencyMax:         frequencyMaxFlagValue,
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
