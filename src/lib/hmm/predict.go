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
	var model HMM
	data, err := os.ReadFile(inFilename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &model)
	if err != nil {
		return err
	}

	// Read in each sentence and predict its pitch.
	var sentences [][]SangoSyllable
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		words := strings.Fields(scanner.Text())
		tokens := make([]Syllable, len(words))
		for k, word := range words {
			tokens[k] = S(word)
		}
		tags, err := model.Predict(tokens)
		if err != nil {
			return err
		}
		var sentence []SangoSyllable
		for k, tag := range tags {
			sentence = append(sentence, SangoSyllable{
				Token:     tokens[k],
				TokenName: tokens[k].String(),
				TagIndex:  tag,
				TagName:   tag.String(),
			})
		}
		sentences = append(sentences, sentence)
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

// Predict uses the Viterbi algorithm operating in log space
func (h HMM) Predict(tokens []Syllable) ([]Tag, error) {
	if h.NumSentences <= 0 {
		return nil, fmt.Errorf("%s", "HMM.Predict called without having first called HMM.Generate")
	}
	if len(tokens) == 0 {
		return nil, fmt.Errorf("%s", "no tokens provided")
	}
	numTokens := len(tokens)
	viterbi := make([][NumTags]float64, numTokens)
	backpointer := make([][NumTags]Tag, numTokens)
	firstToken := tokens[0]
	for tag := range NumTags {
		emissionLog := h.Emission[tag][firstToken]
		viterbi[0][tag] = h.StartTags[tag] + emissionLog
	}
	for t, token := range tokens {
		if t == 0 {
			continue
		}
		for currTag := range NumTags {
			maxLogProb := math.Inf(-1)
			bestPrevTag := UnknownPitch
			currEmissionLog := h.Emission[currTag][token]
			for prevTag := range NumTags {
				logProb := viterbi[t-1][prevTag] + h.Transition[prevTag][currTag] + currEmissionLog
				if logProb > maxLogProb {
					maxLogProb = logProb
					bestPrevTag = prevTag
				}
			}
			viterbi[t][currTag] = maxLogProb
			backpointer[t][currTag] = bestPrevTag
		}
	}
	maxFinalLog := math.Inf(-1)
	bestFinalTag := UnknownPitch
	lastIndex := numTokens - 1
	for currTag := range NumTags {
		if viterbi[lastIndex][currTag] > maxFinalLog {
			maxFinalLog = viterbi[lastIndex][currTag]
			bestFinalTag = currTag
		}
	}
	resultTags := make([]Tag, numTokens)
	currTag := bestFinalTag
	for t := lastIndex; t >= 0; t-- {
		resultTags[t] = currTag
		currTag = backpointer[t][currTag]
	}
	return resultTags, nil
}
