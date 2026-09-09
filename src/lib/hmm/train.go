// Trains Hidden Markov Model (HMM) with Kneser-Ney smoothing for diacritic restoration of Sango text.

package hmm

import (
	"io/ioutil"
	"os"

	"google.golang.org/protobuf/proto"
)

func MainTrain(trainingTextFilename, modelOutputFilename string) error {
	trainingText, err := os.ReadFile(trainingTextFilename)
	if err != nil {
		return err
	}
	trainingData := PrepareInputText(string(trainingText))

	h := HMM{
		States:      make(map[string]bool),
		Tokens:      make(map[string]bool),
		Transitions: make(map[string]map[string]int64),
		Emissions:   make(map[string]map[string]int64),
		StateCounts: make(map[string]int64),
	}
	h.Train(trainingData)
	m := h.ToModel()

	// Write the model to disk.
	out, err := proto.Marshal(&m)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(modelOutputFilename, out, 0644); err != nil {
		return err
	}

	return nil
}

// Populates the raw count matrices from our tiny low-resource dataset.
// This is a two-pass algorithm
func (h *HMM) Train(data []TrainingInstance) {
	for _, inst := range data {
		for i := 0; i < len(inst.Tokens); i++ {
			state := inst.States[i]
			word := inst.Tokens[i]

			h.States[state] = true
			h.Tokens[word] = true
			h.StateCounts[state]++

			// Count Emissions
			if _, exists := h.Emissions[state]; !exists {
				h.Emissions[state] = make(map[string]int64)
			}
			h.Emissions[state][word]++

			// Count Transitions (for simplicity, ignoring structural boundary tokens here)
			if i > 0 {
				prevState := inst.States[i-1]
				if _, exists := h.Transitions[prevState]; !exists {
					h.Transitions[prevState] = make(map[string]int64)
				}
				h.Transitions[prevState][state]++
			}
		}
	}

	// Count unique (state, word) type configurations for Kneser-Ney denominator
	for _, words := range h.Emissions {
		h.UniquePairs += int64(len(words))
	}
}
