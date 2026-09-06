package hmm

import (
	"encoding/json"
	"fmt"
	"os"
)

func MainMerge(inFilenames []string, outFilename string) error {
	// Make sure we can actually write to the file.
	if err := os.WriteFile(outFilename, []byte("TEST"), 0644); err != nil {
		return err
	}
	var model HMM
	for _, inFilename := range inFilenames {
		data, err := os.ReadFile(inFilename)
		if err != nil {
			return err
		}

		// Convert the byte slice into a string variable
		var h HMM
		err = json.Unmarshal(data, &h)
		if err != nil {
			return err
		}
		err = model.Merge(h)
		if err != nil {
			return err
		}
	}

	// Marshal and write model to stdout
	jsonBytes, err := json.Marshal(model)
	if err != nil {
		return err
	}
	return os.WriteFile(outFilename, jsonBytes, 0644)
}

func (h *HMM) Merge(rhs HMM) error {
	if h.NumSentences > 0 {
		return fmt.Errorf("%s", "HMM.Merge called after HMM.Generate has already been called")
	}
	if rhs.NumSentences > 0 {
		return fmt.Errorf("%s", "HMM.Merge called with argument on which HMM.Generate has already been called")
	}
	h.NumSentences += rhs.NumSentences
	for i := range NumTags {
		h.TagCounts[i] += rhs.TagCounts[i]
		h.StartTags[i] += rhs.StartTags[i]
		for j := range NumTags {
			h.Transition[i][j] += rhs.Transition[i][j]
		}
		for j := range NumSyllables {
			h.Emission[i][j] += rhs.Emission[i][j]
		}
	}
	return nil
}
