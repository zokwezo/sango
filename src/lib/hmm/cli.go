package hmm

import (
	"github.com/spf13/cobra"
)

func Init(rootCmd *cobra.Command) {
	rootCmd.AddCommand(hmmCmd)

	hmmCmd.AddCommand(accumulateCmd)
	accumulateCmd.Flags().StringP("out", "o", "",
		"output filename of the HMM counts accumulated from the text in stdin.")
	accumulateCmd.MarkFlagRequired("out")

	hmmCmd.AddCommand(mergeCmd)
	mergeCmd.Flags().StringArrayP("in", "i", []string{},
		"input filenames of the HMM counts to be merged into a single model count.")
	mergeCmd.Flags().StringP("out", "o", "",
		"output filename of the merged HMM counts.")
	mergeCmd.MarkFlagRequired("in")
	mergeCmd.MarkFlagRequired("out")

	hmmCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringP("in", "i", "",
		"input filename of the HMM counts.")
	generateCmd.Flags().StringP("out", "o", "",
		"output filename of the HMM model generated from the input model.")
	generateCmd.MarkFlagRequired("in")
	generateCmd.MarkFlagRequired("out")

	hmmCmd.AddCommand(predictCmd)
	predictCmd.Flags().StringP("in", "i", "",
		"input filename of the HMM model.")
	predictCmd.MarkFlagRequired("in")

	hmmCmd.AddCommand(evaluateCmd)
	evaluateCmd.Flags().StringP("actual", "a", "",
		"predicted text which is evaluated against the expected text.")
	evaluateCmd.Flags().StringP("expect", "e", "", "expected human-corrected text.")
	evaluateCmd.MarkFlagRequired("actual")
	evaluateCmd.MarkFlagRequired("expect")
}

var (
	hmmCmd = &cobra.Command{
		Use:   "hmm",
		Short: "A CLI to generate, use, and evaluate a Hidden Markov Model (HMM) for diacritic restoration of Sango text",
		Long:  "https://github.com/zokwezo/sango/blob/main/src/hmm/README.md",
	}

	accumulateCmd = &cobra.Command{
		Use:   "accumulate",
		Short: "Trains an HMM model on Sango text having (hopefully accurate) diacritics.",
		Long:  "Read Sango training text from stdin and accumulate HMM counts of its words and their diacritics.",
		Args:  cobra.MaximumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			outFilename, err := cmd.Flags().GetString("out")
			if err != nil {
				return err
			}
			return MainAccumulate(outFilename)
		},
	}

	mergeCmd = &cobra.Command{
		Use:   "merge",
		Short: "Merges multiple accumulated HMM counts into one.",
		Long:  "The utility of this is to improve a model by adding additional corpora without having to retrain on the entire corpus.",
		Args:  cobra.MaximumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			inFilenames, err := cmd.Flags().GetStringArray("in")
			if err != nil {
				return err
			}
			outFilename, err := cmd.Flags().GetString("out")
			if err != nil {
				return err
			}
			return MainMerge(inFilenames, outFilename)
		},
	}

	generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Convert HMM counts into an HMM model.",
		Long:  "After generating a model, additional corpora can no longer be added to it.",
		Args:  cobra.MaximumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			inFilename, err := cmd.Flags().GetString("in")
			if err != nil {
				return err
			}
			outFilename, err := cmd.Flags().GetString("out")
			if err != nil {
				return err
			}
			return MainGenerate(inFilename, outFilename)
		},
	}

	predictCmd = &cobra.Command{
		Use:   "predict",
		Short: "Predicts Sango diacritics for input Sango text.",
		Long:  "Read Sango input text from stdin, strips any diacritics, adds new diacritics predicted with the supplied HMM model, and outputs the result to stdout. Any non-Sango text is passed through unchanged.",
		Args:  cobra.MaximumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			inFilename, err := cmd.Flags().GetString("in")
			if err != nil {
				return err
			}
			return MainPredict(inFilename)
		},
	}

	evaluateCmd = &cobra.Command{
		Use:   "evaluate",
		Short: "Evaluates the quality of predicted Sango diacritics.",
		Long:  "Read test Sango text and human-verified Sango text and outputs statistics on the precision, recall, and F1 score.",
		Args:  cobra.MaximumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			actualFilename, err := cmd.Flags().GetString("actual")
			if err != nil {
				return err
			}
			expectFilename, err := cmd.Flags().GetString("expect")
			if err != nil {
				return err
			}
			return MainEvaluate(actualFilename, expectFilename)
		},
	}
)
