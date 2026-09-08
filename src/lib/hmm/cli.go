package hmm

import (
	"github.com/spf13/cobra"
)

func Init(rootCmd *cobra.Command) {
	rootCmd.AddCommand(hmmCmd)

	hmmCmd.AddCommand(trainCmd)
	trainCmd.Flags().StringP("in", "i", "",
		"input filename of the training data")
	trainCmd.Flags().StringP("out", "o", "",
		"output filename of the HMM counts")
	trainCmd.MarkFlagRequired("out")

	hmmCmd.AddCommand(predictCmd)
	predictCmd.Flags().StringP("in", "i", "",
		"input filename of the HMM model")
	// predictCmd.MarkFlagRequired("in")

	hmmCmd.AddCommand(evaluateCmd)
	evaluateCmd.Flags().StringP("actual", "a", "",
		"predicted text which is evaluated against the expected text")
	evaluateCmd.Flags().StringP("expect", "e", "", "expected human-corrected text")
	evaluateCmd.MarkFlagRequired("actual")
	evaluateCmd.MarkFlagRequired("expect")
}

var (
	hmmCmd = &cobra.Command{
		Use:   "hmm",
		Short: "A CLI to generate, use, and evaluate a Hidden Markov Model (HMM) for diacritic restoration of Sango text",
		Long:  "https://github.com/zokwezo/sango/blob/main/src/hmm/README.md",
	}

	trainCmd = &cobra.Command{
		Use:   "train",
		Short: "Trains an HMM model on Sango text having (hopefully accurate) diacritics.",
		Long:  "Read Sango training text from stdin and train HMM counts of its words and their diacritics.",
		Args:  cobra.MaximumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			trainingTextFilename, err := cmd.Flags().GetString("in")
			if err != nil {
				return err
			}
			modelOutputFilename, err := cmd.Flags().GetString("out")
			if err != nil {
				return err
			}
			return MainTrain(trainingTextFilename, modelOutputFilename)
		},
	}

	predictCmd = &cobra.Command{
		Use:   "predict",
		Short: "Predicts an HMM model on Sango text having (hopefully accurate) diacritics.",
		Long:  "Read Sango predicting text from stdin and restore their vowels' diacritics.",
		Args:  cobra.MaximumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			modelInputFilename, err := cmd.Flags().GetString("in")
			if err != nil {
				return err
			}
			return MainPredict(modelInputFilename)
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
