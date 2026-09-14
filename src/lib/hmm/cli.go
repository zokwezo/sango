package hmm

import (
	"github.com/spf13/cobra"
)

func Init(rootCmd *cobra.Command) {
	rootCmd.AddCommand(hmmCmd)

	hmmCmd.AddCommand(trainCmd)
	trainCmd.Flags().StringP("in", "i", "",
		"filename of the input training text")
	trainCmd.Flags().StringP("model", "m", "",
		"filename of the generated HMM")
	trainCmd.MarkFlagRequired("in")
	trainCmd.MarkFlagRequired("model")

	hmmCmd.AddCommand(predictCmd)
	predictCmd.Flags().StringP("in", "i", "",
		"filename of the input text")
	predictCmd.Flags().StringP("model", "m", "",
		"input filename of the HMM")
	predictCmd.Flags().StringP("orig_codes", "o", "",
		"filename of the encoded original input text")
	predictCmd.Flags().StringP("pred_codes", "p", "",
		"filename of the encoded predicted output text")
	predictCmd.MarkFlagRequired("in")
	predictCmd.MarkFlagRequired("model")
	predictCmd.MarkFlagRequired("orig_codes")
	predictCmd.MarkFlagRequired("pred_codes")

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
			inputTextFilename, err := cmd.Flags().GetString("in")
			if err != nil {
				return err
			}
			modelOutputFilename, err := cmd.Flags().GetString("model")
			if err != nil {
				return err
			}
			return MainTrain(inputTextFilename, modelOutputFilename)
		},
	}

	predictCmd = &cobra.Command{
		Use:   "predict",
		Short: "Predicts an HMM model on Sango text having (hopefully accurate) diacritics.",
		Long:  "Read Sango predicting text from stdin and restore their vowels' diacritics.",
		Args:  cobra.MaximumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			inputTextFilename, err := cmd.Flags().GetString("in")
			if err != nil {
				return err
			}
			modelInputFilename, err := cmd.Flags().GetString("model")
			if err != nil {
				return err
			}
			origCodesFilename, err := cmd.Flags().GetString("orig_codes")
			if err != nil {
				return err
			}
			predCodesFilename, err := cmd.Flags().GetString("pred_codes")
			if err != nil {
				return err
			}
			return MainPredict(inputTextFilename, modelInputFilename, origCodesFilename, predCodesFilename)
		},
	}

	evaluateCmd = &cobra.Command{
		Use:   "evaluate",
		Short: "Evaluates the quality of predicted Sango diacritics.",
		Long:  "Read test Sango text and human-verified Sango text and outputs statistics on the precision, recall, and F1 score.",
		Args:  cobra.MaximumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			actualCodesFilename, err := cmd.Flags().GetString("actual")
			if err != nil {
				return err
			}
			expectCodesFilename, err := cmd.Flags().GetString("expect")
			if err != nil {
				return err
			}
			return MainEvaluate(actualCodesFilename, expectCodesFilename)
		},
	}
)
