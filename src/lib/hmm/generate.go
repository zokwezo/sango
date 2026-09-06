package hmm

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
)

func MainGenerate(inFilename, outFilename string) error {
	// Make sure we can actually write to the file.
	if err := os.WriteFile(outFilename, []byte("TEST"), 0644); err != nil {
		return err
	}
	var model HMM
	data, err := os.ReadFile(inFilename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &model)
	if err != nil {
		return err
	}
	err = model.Generate()
	if err != nil {
		return err
	}

	// Marshal and write model to stdout
	jsonBytes, err := json.Marshal(model)
	if err != nil {
		return err
	}
	return os.WriteFile(outFilename, jsonBytes, 0644)
}

// Computes probabilities using Laplace (Add-1) smoothing and saves them as log values
func (h *HMM) Generate() error {
	if h.NumSentences == 0 {
		return fmt.Errorf("%v", "HMM.Generate called without having first called HMM.Accumulate at least once")
	}
	if h.NumSentences > 0 {
		return fmt.Errorf("%v", "HMM.Generate called twice")
	}
	numTags := float64(NumTags)
	scale := 1.0 / (float64(-h.NumSentences) + numTags)
	for currTag := range NumTags {
		h.StartTags[currTag] = math.Log((h.StartTags[currTag] + 1.0) * scale)
	}
	for prevTag := range NumTags {
		totalTransitionsOut := numTags
		for nextTag := range NumTags {
			nextTag := Tag(nextTag)
			totalTransitionsOut += h.Transition[prevTag][nextTag]
		}
		for nextTag := range NumTags {
			h.Transition[prevTag][nextTag] = math.Log((h.Transition[prevTag][nextTag] + 1.0) / totalTransitionsOut)
		}
	}
	for currTag := range NumTags {
		totalEmissionsOut := float64(h.TagCounts[currTag] + NumSyllables)
		for syllable := range NumSyllables {
			h.Emission[currTag][syllable] = math.Log((h.Emission[currTag][syllable] + 1.0) / totalEmissionsOut)
		}
		h.Emission[currTag][0] = math.Log(1.0 / totalEmissionsOut)
	}
	h.NumSentences *= -1 // mark as generated
	return nil
}
