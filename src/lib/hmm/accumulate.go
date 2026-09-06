package hmm

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Read and process sentences from stdin and output to a file.
func MainAccumulate(outFilename string) error {
	// Make sure we can actually write to the file.
	if err := os.WriteFile(outFilename, []byte("TEST"), 0644); err != nil {
		return err
	}
	var model HMM
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		s := scanner.Text()
		words := strings.Fields(s)
		// TODO: This is a toy. Change to reading real Sango text and splitting off the diacritics.
		if len(words)%2 != 0 {
			return fmt.Errorf("odd number (%v) of words in line %q", len(words), s)
		}
		numWords := len(words) / 2
		trainingSentence := make([]Syllable, numWords)
		trainingTags := make([]Tag, numWords)
		for k := range numWords {
			trainingSentence[k] = S(words[2*k])
			trainingTags[k] = T(words[2*k+1])
		}
		err := model.Accumulate(trainingSentence, trainingTags)
		if err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	jsonBytes, err := json.Marshal(model)
	if err != nil {
		return err
	}
	return os.WriteFile(outFilename, jsonBytes, 0644)
}

// Accumulates counts of tag transmissions and emissions
// Separate models can Accumulate concurrently and be merged at the end.
// NOTE: sentences must contain every known syllable at least once.
func (h *HMM) Accumulate(sentence []Syllable, tags []Tag) error {
	if h.NumSentences > 0 {
		return fmt.Errorf("%s", "Accumulate called after Generate has already been called on an HMM model")
	}
	h.NumSentences--
	prevTag := UnknownPitch
	for j := range sentence {
		syllable := sentence[j]
		currTag := tags[j]
		h.TagCounts[currTag]++
		if j == 0 {
			h.StartTags[currTag]++
		} else {
			h.Transition[prevTag][currTag]++
		}
		h.Emission[currTag][syllable]++
		prevTag = currTag
	}
	return nil
}
