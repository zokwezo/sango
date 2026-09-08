// Hidden Markov Model (HMM) with Kneser-Ney smoothing

// Used for diacritic restoration of Sango text.

package hmm

import (
	"fmt"
	"io/ioutil"

	"google.golang.org/protobuf/proto"
)

func NewHMM() *HMM {
	return &HMM{
		States:              make(map[string]bool),
		Vocab:               make(map[string]bool),
		TransitionFromState: make(map[string]*Then),
		EmissionFromState:   make(map[string]*Then),
		StateCount:          make(map[string]float64),
	}
}

// TrainingInstance represents a sequence of tokens and their true state/label sequence.
// In this context, a Token is a Sango word without diacritics and a State is the
// same word with diacritics.
type TrainingInstance struct {
	Tokens []string
	States []string
}

func MainTrain(modelOutputFilename string) error {
	// Tiny low-resource corpus
	// Notice how "Francisco" ALWAYS appears with "NOUN", making it highly rigid.
	// "Apple" appears with both "NOUN" and "PROPN" (the company), making it versatile.
	trainingData := []TrainingInstance{
		{Tokens: []string{"San", "Francisco"}, States: []string{"PROPN", "PROPN"}},
		{Tokens: []string{"Eat", "an", "apple"}, States: []string{"VERB", "DET", "NOUN"}},
		{Tokens: []string{"Apple", "announces", "iPhone"}, States: []string{"PROPN", "VERB", "PROPN"}},
	}

	h := NewHMM()
	err := h.Train(trainingData)
	if err != nil {
		return err
	}

	// Write the model to disk.
	out, err := proto.Marshal(h)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(modelOutputFilename, out, 0644); err != nil {
		return err
	}

	return nil
}

func MainPredict(modelInputFilename string) error {
	in, err := ioutil.ReadFile(modelInputFilename)
	if err != nil {
		return err
	}
	h := NewHMM()
	if err := proto.Unmarshal(in, h); err != nil {
		return err
	}

	fmt.Println("--- EMISSION COMPARISONS ---")

	// Example 1: Evaluating a highly rigid token in a completely novel context
	// Does "Francisco" make sense as a generic VERB?
	fmt.Printf("\nToken: 'Francisco' | Evaluated State: 'VERB'\n")
	fmt.Printf("  No Smoothing:        %0.5f (Completely rejects it)\n", h.GetEmissionNoSmoothing("VERB", "Francisco"))
	fmt.Printf("  Kneser-Ney Sm.:      %0.5f (Keeps mass very low because it lacks versatility)\n", h.GetEmissionKneserNey("VERB", "Francisco"))

	// Example 2: Evaluating a versatile token in an unseen context
	// We've never explicitly trained "Apple" as a generic "DET" (determiner),
	// but it is a versatile word in the corpus.
	fmt.Printf("\nToken: 'Apple' | Evaluated State: 'DET'\n")
	fmt.Printf("  No Smoothing:        %0.5f (Strictly breaks the sequence sequence match)\n", h.GetEmissionNoSmoothing("DET", "Apple"))
	fmt.Printf("  Kneser-Ney Sm.:      %0.5f (Grants a higher fallback probability because 'Apple' is versatile)\n", h.GetEmissionKneserNey("DET", "Apple"))

	// Evaluation test sequence:
	// "Apple" is used here as a common NOUN in the object position, rather than the PROPN company name.
	testTokens := []string{"Eat", "an", "Apple"}

	fmt.Println("Tokens to predict:", testTokens)

	// Predict with NO smoothing on Emissions
	noSmoothPath := h.Viterbi(testTokens, false)
	fmt.Println("Prediction (No Smoothing): ", noSmoothPath)

	// Predict WITH Kneser-Ney smoothing on Emissions
	knSmoothPath := h.Viterbi(testTokens, true)
	fmt.Println("Prediction (Kneser-Ney):  ", knSmoothPath)

	return nil
}
