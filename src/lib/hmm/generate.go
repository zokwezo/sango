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
	model := HMM{}
	data, err := os.ReadFile(inFilename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &model)
	if err != nil {
		return err
	}
	if len(model.HmmPerTag) == 0 {
		return fmt.Errorf("parse error reading model from %v", inFilename)
	}
	if err = model.Generate(); err != nil {
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
	h.ensureDefined()
	// All the + 1.0 are Laplace smoothing to account for hapax legomena.
	totalStartTags := 0.0
	for tag := range h.HmmPerTag {
		hh := h.forTag(tag)
		totalStartTags += hh.StartTag + 1.0
	}
	for tag := range h.HmmPerTag {
		hh := h.forTag(tag)
		hh.StartTag = math.Log((hh.StartTag + 1.0) / totalStartTags)
		totalTransitions := 0.0
		for _, transition := range hh.Transition {
			totalTransitions += transition + 1.0
		}
		for nextTag, transition := range hh.Transition {
			hh.Transition[nextTag] = math.Log((transition + 1.0) / totalTransitions)
		}
		totalEmissions := 0.0
		for _, emission := range hh.Emission {
			totalEmissions += emission + 1.0
		}
		for token, emission := range hh.Emission {
			hh.Emission[token] = math.Log((emission + UnknownTokenDampening) / totalEmissions)
		}
		totalEmissions++
		hh.Emission[UnknownToken] = math.Log(UnknownTokenDampening / totalEmissions)
	}
	h.NumSentences *= -1 // mark as generated
	return nil
}
