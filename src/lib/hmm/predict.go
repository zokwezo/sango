package hmm

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
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
			emissionLog, exists := hmmPerTag.Emission[firstToken]
			if !exists {
				emissionLog = -10.0
			}
			viterbi[0][tag] = hmmPerTag.StartTag + emissionLog
		}
	}

	for t := 1; t < numTokens; t++ {
		viterbi[t] = make(map[string]float64)
		backpointer[t] = make(map[string]string)
		token := strings.ToLower(sangoTokens[t].Token)

		for tag, hmmPerTag := range model.HmmPerTag {
			maxLogProb := math.Inf(-1)
			bestPrevTag := UnknownTag

			currEmissionLog, exists := hmmPerTag.Emission[token]
			if !exists {
				currEmissionLog = math.Inf(-1)
			}

			for prevTag, hmmPerTag := range model.HmmPerTag {
				currTransitionLog, exists := hmmPerTag.Transition[tag]
				if !exists {
					currTransitionLog = math.Inf(-1)
				}
				logProb := viterbi[t-1][prevTag] + currTransitionLog + currEmissionLog
				if logProb > maxLogProb {
					maxLogProb = logProb
					bestPrevTag = prevTag
				}
			}
			viterbi[t][tag] = maxLogProb
			backpointer[t][tag] = bestPrevTag
		}
	}

	maxFinalLog := math.Inf(-1)
	bestFinalTag := UnknownTag
	lastIdx := numTokens - 1

	for tag := range model.HmmPerTag {
		if viterbi[lastIdx][tag] > maxFinalLog {
			maxFinalLog = viterbi[lastIdx][tag]
			bestFinalTag = tag
		}
	}

	tag := bestFinalTag
	for t := lastIdx; t >= 0; t-- {
		sangoTokens[t].Tag = tag
		tag = backpointer[t][tag]
	}

	return nil
}
