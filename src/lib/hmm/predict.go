package hmm

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func MainPredict(inFilename string) error {
	// Read in the model.
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

	// Read in each sentence and predict its pitch.
	sentences := [][]SangoToken{}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		tokens := strings.Fields(scanner.Text())
		sangoTokens := make([]SangoToken, len(tokens))
		for k, token := range tokens {
			sangoTokens[k].Token = token
		}
		if err := model.Predict(sangoTokens); err != nil {
			return err
		}
		sentences = append(sentences, sangoTokens)
	}

	// Check if the loop terminated due to a reading error
	if err := scanner.Err(); err != nil {
		return err
	}
	jsonBytes, err := json.Marshal(sentences)
	if err != nil {
		return err
	}
	_, err = os.Stdout.Write(jsonBytes)
	return err
}

type VBPerTokenPerTag = struct {
	viberbi     float64
	backpointer string
}
type VBPerToken = map[string]VBPerTokenPerTag
type VB map[string]*VBPerToken

func (vb VB) forToken(token string) *VBPerToken {
	if vb[token] == nil {
		vb[token] = &VBPerToken{}
	}
	return vb[token]
}

// Predict uses the Viterbi algorithm operating in log space
func (model HMM) Predict(sangoTokens []SangoToken) error {
	if model.NumSentences <= 0 {
		return fmt.Errorf("%s", "HMM.Predict called without having first called HMM.Generate")
	}
	if len(model.HmmPerTag) == 0 {
		return fmt.Errorf("bad model passed to HMM.Predict")
	}
	numTokens := len(sangoTokens)
	if numTokens == 0 {
		return fmt.Errorf("%s", "no tokens provided")
	}

	viterbi := make([]map[string]float64, numTokens)
	backpointer := make([]map[string]string, numTokens)

	viterbi[0] = make(map[string]float64)
	backpointer[0] = make(map[string]string)
	firstToken := strings.ToLower(sangoTokens[0].Token)

	for tag, hmmPerTag := range model.HmmPerTag {
		if hmmPerTag != nil {
			logStartTag := hmmPerTag.StartTag
			if logStartTag == 0.0 {
				logStartTag = -50.0
			}
			logEmission := hmmPerTag.Emission[firstToken]
			if logEmission == 0.0 {
				logEmission = max(-50.0, hmmPerTag.Emission[UnknownToken])
			}
			viterbi[0][tag] = StartTagDampening*logStartTag + EmissionDampening*logEmission
		}
	}

	for t := 1; t < numTokens; t++ {
		viterbi[t] = make(map[string]float64)
		backpointer[t] = make(map[string]string)
		token := strings.ToLower(sangoTokens[t].Token)

		for tag, hmmPerTag := range model.HmmPerTag {
			maxLogProb := -50.0
			bestPrevTag := UnknownTag

			logCurrEmission := hmmPerTag.Emission[token]
			if logCurrEmission == 0.0 {
				logCurrEmission = max(-50.0, hmmPerTag.Emission[UnknownToken])
			}

			for prevTag, hmmPerTag := range model.HmmPerTag {
				logCurrTransition := hmmPerTag.Transition[tag]
				if logCurrTransition == 0.0 {
					logCurrTransition = -50.0
				}
				logProb := viterbi[t-1][prevTag] + TransitionDampening*logCurrTransition + EmissionDampening*logCurrEmission
				if logProb > maxLogProb {
					maxLogProb = logProb
					bestPrevTag = prevTag
				}
			}
			viterbi[t][tag] = maxLogProb
			backpointer[t][tag] = bestPrevTag
		}
	}

	maxLogProbFinal := -50.0
	bestTagFinal := UnknownTag
	lastIdx := numTokens - 1

	for tag := range model.HmmPerTag {
		if viterbi[lastIdx][tag] > maxLogProbFinal {
			maxLogProbFinal = viterbi[lastIdx][tag]
			bestTagFinal = tag
		}
	}

	tag := bestTagFinal
	for t := lastIdx; t >= 0; t-- {
		sangoTokens[t].Tag = tag
		tag = backpointer[t][tag]
	}

	return nil
}
