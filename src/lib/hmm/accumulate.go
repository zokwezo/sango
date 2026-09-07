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
	h := HMM{}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		s := scanner.Text()
		words := strings.Fields(s)
		// TODO: This is a toy. Change to reading real Sango text and splitting off the diacritics.
		if len(words)%2 != 0 {
			return fmt.Errorf("odd number (%v) of words in line %q", len(words), s)
		}
		numWords := len(words) / 2
		sentence := make([]SangoToken, numWords)
		for k := range numWords {
			sentence[k] = SangoToken{Token: words[2*k], Tag: words[2*k+1]}
		}
		err := h.Accumulate(sentence)
		if err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	jsonBytes, err := json.Marshal(h)
	if err != nil {
		return err
	}
	return os.WriteFile(outFilename, jsonBytes, 0644)
}

// Accumulates counts of tag transmissions and emissions
// Separate models can Accumulate concurrently and be merged at the end.
// NOTE: sentences must contain every known token at least once.
func (h *HMM) Accumulate(sentence []SangoToken) error {
	if h.NumSentences > 0 {
		return fmt.Errorf("%s", "HMM.Accumulate called after Generate has already been called on an HMM model")
	}
	if len(sentence) == 0 {
		return nil
	}
	h.NumSentences--
	prevTag := UnknownTagForToken(sentence[0].Token)
	for j := range sentence {
		sangoToken := sentence[j]
		token := sangoToken.Token
		tag := sangoToken.Tag
		if j == 0 {
			h.forTag(tag).StartTag++
		} else {
			h.forTag(prevTag).Transition[tag]++
		}
		h.forTag(tag).Emission[token]++
		prevTag = tag
	}
	return nil
}
