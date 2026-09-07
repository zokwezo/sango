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
	h := HMM{}
	for _, inFilename := range inFilenames {
		data, err := os.ReadFile(inFilename)
		if err != nil {
			return err
		}

		// Convert the byte slice into a string variable
		rhs := HMM{}
		err = json.Unmarshal(data, &rhs)
		if err != nil {
			return err
		}
		if len(rhs.HmmPerTag) == 0 {
			return fmt.Errorf("parse error reading model from %v", inFilename)
		}
		if err = h.Merge(rhs); err != nil {
			return err
		}
	}

	// Marshal and write model to stdout
	jsonBytes, err := json.Marshal(h)
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
	rhs.ensureDefined()
	for tag := range rhs.HmmPerTag {
		src := rhs.forTag(tag)
		tgt := h.forTag(tag)
		tgt.StartTag += src.StartTag
		for nextTag, count := range src.Transition {
			tgt.Transition[nextTag] += count
		}
		for token, count := range src.Emission {
			tgt.Emission[token] += count
		}
	}
	return nil
}
