// Hidden Markov Model (HMM) with Kneser-Ney smoothing

// Used for diacritic restoration of Sango text.

package hmm

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/unicode/norm"
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

var (
	reFinalPunctuation      = regexp.MustCompile(`[.!?][ ]`)
	reNonTextWithDiacritics = regexp.MustCompile(`[^a-z .?!\x{302}\x{308}\x{323}]`)
	reNonText               = regexp.MustCompile(`[^a-z .?!]`)
	reCompressSpaces        = regexp.MustCompile(`[ ]{2,}`)
)

func PrepareInputText(s string) []TrainingInstance {
	log.Printf("S: <%s>\n", s)
	s = norm.NFD.String(cases.Lower(language.English).String(s))
	log.Printf("S: <%s>\n", s)
	s = strings.ReplaceAll(s, "\n", " ")
	log.Printf("S: <%s>\n", s)
	s = strings.Trim(s, " ")
	log.Printf("S: <%s>\n", s)
	s = reNonTextWithDiacritics.ReplaceAllLiteralString(s, " ")
	log.Printf("S: <%s>\n", s)
	s = reCompressSpaces.ReplaceAllLiteralString(s, " ")
	log.Printf("S: <%s>\n", s)
	s = reFinalPunctuation.ReplaceAllLiteralString(s, "\n")
	log.Printf("S: <%s>\n", s)
	tis := []TrainingInstance{}
	for i, sentence := range strings.Split(strings.Trim(s, " "), "\n") {
	  log.Printf("S[%v]: <%s>\n", i, sentence)
	  sentence = strings.Trim(sentence, " ")
	  log.Printf("S[%v]: <%s>\n", i, sentence)
		if sentence != "" {
	  log.Printf("S[%v]: <%s>\n", i, sentence)
		ti := TrainingInstance{}
		ti.States = strings.Split(sentence, " ")
		n := len(ti.States)
		log.Printf("n[%v] = %v\n", i, n)
		ti.Tokens = make([]string, n)
		for k := range ti.States {
	    log.Printf("S[%v][%v]: <%s>\n", i, k, ti.States[k])
			ti.Tokens[k] = reNonText.ReplaceAllLiteralString(ti.States[k], "")
	    log.Printf("T[%v][%v]: <%s>\n", i, k, ti.Tokens[k])
			ti.States[k] = norm.NFC.String(ti.States[k])
	    log.Printf("S[%v][%v]: <%s>\n", i, k, ti.States[k])
		}
		tis = append(tis, ti)
		}
	}
	log.Printf("n = %v\n", len(tis))
	log.Printf("%#v\n", tis)
	return tis
}

func MainTrain(trainingTextFilename, modelOutputFilename string) error {
	trainingText, err := os.ReadFile(trainingTextFilename)
	if err != nil {
		return err
	}
	trainingData := PrepareInputText(string(trainingText))
	log.Printf("trainingData = %#v\n", trainingData)

	h := NewHMM()
	err = h.Train(trainingData)
	if err != nil {
		return err
	}
	log.Printf("h = %#v\n", *h)

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
