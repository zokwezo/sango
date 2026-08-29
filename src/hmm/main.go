package main

import (
	"fmt"
	"math"
	"strings"
)

type HMM struct {
	Transition map[string]map[string]float64
	Emission   map[string]map[string]float64
	StartTags  map[string]float64
	TagCounts  map[string]int
	Vocab      map[string]bool
	Tags       []string
}

// Metrics holds the classification performance statistics
type Metrics struct {
	Precision float64
	Recall    float64
	F1Score   float64
	TP        int // True Positives
	FP        int // False Positives
	FN        int // False Negatives
}

func NewHMM() *HMM {
	return &HMM{
		Transition: make(map[string]map[string]float64),
		Emission:   make(map[string]map[string]float64),
		StartTags:  make(map[string]float64),
		TagCounts:  make(map[string]int),
		Vocab:      make(map[string]bool),
		Tags:       []string{},
	}
}

// Train computes probabilities using Laplace (Add-1) smoothing and saves them as log values
func (h *HMM) Train(sentences [][]string, tags [][]string) {
	totalSentences := float64(len(sentences))
	rawStart := make(map[string]float64)
	rawTrans := make(map[string]map[string]float64)
	rawEmis := make(map[string]map[string]float64)

	tagMap := make(map[string]bool)
	for i := range sentences {
		prevTag := ""
		for j := range sentences[i] {
			word := strings.ToLower(sentences[i][j])
			tag := tags[i][j]

			h.TagCounts[tag]++
			h.Vocab[word] = true
			tagMap[tag] = true

			if j == 0 {
				rawStart[tag]++
			} else {
				if rawTrans[prevTag] == nil {
					rawTrans[prevTag] = make(map[string]float64)
				}
				rawTrans[prevTag][tag]++
			}

			if rawEmis[tag] == nil {
				rawEmis[tag] = make(map[string]float64)
			}
			rawEmis[tag][word]++
			prevTag = tag
		}
	}

	for tag := range tagMap {
		h.Tags = append(h.Tags, tag)
	}

	numTags := float64(len(h.Tags))
	numWords := float64(len(h.Vocab))

	for _, tag := range h.Tags {
		count := rawStart[tag]
		h.StartTags[tag] = math.Log((count + 1.0) / (totalSentences + numTags))
	}

	for _, prevTag := range h.Tags {
		h.Transition[prevTag] = make(map[string]float64)
		totalTransitionsOut := 0.0
		for _, nextTag := range h.Tags {
			totalTransitionsOut += rawTrans[prevTag][nextTag]
		}
		for _, nextTag := range h.Tags {
			count := rawTrans[prevTag][nextTag]
			prob := (count + 1.0) / (totalTransitionsOut + numTags)
			h.Transition[prevTag][nextTag] = math.Log(prob)
		}
	}

	for _, tag := range h.Tags {
		h.Emission[tag] = make(map[string]float64)
		totalEmissionsOut := float64(h.TagCounts[tag])
		for word := range h.Vocab {
			count := rawEmis[tag][word]
			prob := (count + 1.0) / (totalEmissionsOut + numWords)
			h.Emission[tag][word] = math.Log(prob)
		}
		h.Emission[tag]["<UNKNOWN>"] = math.Log(1.0 / (totalEmissionsOut + numWords))
	}
}

// Predict uses the Viterbi algorithm operating in log space
func (h *HMM) Predict(tokens []string) []string {
	if len(tokens) == 0 {
		return nil
	}

	numTokens := len(tokens)
	viterbi := make([]map[string]float64, numTokens)
	backpointer := make([]map[string]string, numTokens)

	viterbi[0] = make(map[string]float64)
	backpointer[0] = make(map[string]string)
	firstWord := strings.ToLower(tokens[0])

	for _, tag := range h.Tags {
		emissionLog := h.Emission[tag][firstWord]
		if !h.Vocab[firstWord] {
			emissionLog = h.Emission[tag]["<UNKNOWN>"]
		}
		viterbi[0][tag] = h.StartTags[tag] + emissionLog
	}

	for t := 1; t < numTokens; t++ {
		viterbi[t] = make(map[string]float64)
		backpointer[t] = make(map[string]string)
		word := strings.ToLower(tokens[t])

		for _, currTag := range h.Tags {
			maxLogProb := math.Inf(-1)
			bestPrevTag := h.Tags[0]

			currEmissionLog := h.Emission[currTag][word]
			if !h.Vocab[word] {
				currEmissionLog = h.Emission[currTag]["<UNKNOWN>"]
			}

			for _, prevTag := range h.Tags {
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
	bestFinalTag := h.Tags[0]
	lastIdx := numTokens - 1

	for _, tag := range h.Tags {
		if viterbi[lastIdx][tag] > maxFinalLog {
			maxFinalLog = viterbi[lastIdx][tag]
			bestFinalTag = tag
		}
	}

	resultTags := make([]string, numTokens)
	currTag := bestFinalTag
	for t := lastIdx; t >= 0; t-- {
		resultTags[t] = currTag
		currTag = backpointer[t][currTag]
	}

	return resultTags
}

// Evaluate runs predictions on test data and prints Precision, Recall, and F1 per tag
func (h *HMM) Evaluate(testSentences [][]string, testTags [][]string) map[string]*Metrics {
	metricsMap := make(map[string]*Metrics)
	for _, tag := range h.Tags {
		metricsMap[tag] = &Metrics{}
	}

	// 1. Accumulate True Positives, False Positives, and False Negatives
	for i := range testSentences {
		predicted := h.Predict(testSentences[i])
		expected := testTags[i]

		for j := range expected {
			trueTag := expected[j]
			predTag := predicted[j]

			// Handle potential unseen tags in testing data safely
			if metricsMap[trueTag] == nil {
				metricsMap[trueTag] = &Metrics{}
			}
			if metricsMap[predTag] == nil {
				metricsMap[predTag] = &Metrics{}
			}

			if trueTag == predTag {
				metricsMap[trueTag].TP++
			} else {
				metricsMap[predTag].FP++
				metricsMap[trueTag].FN++
			}
		}
	}

	// 2. Compute finalized rates
	fmt.Printf("\n%-12s | %-10s | %-10s | %-10s\n", "TAG", "PRECISION", "RECALL", "F1-SCORE")
	fmt.Println(strings.Repeat("-", 53))

	for tag, m := range metricsMap {
		if m.TP+m.FP > 0 {
			m.Precision = float64(m.TP) / float64(m.TP+m.FP)
		}
		if m.TP+m.FN > 0 {
			m.Recall = float64(m.TP) / float64(m.TP+m.FN)
		}
		if m.Precision+m.Recall > 0 {
			m.F1Score = 2 * (m.Precision * m.Recall) / (m.Precision + m.Recall)
		}

		fmt.Printf("%-12s | %-10.2f%% | %-10.2f%% | %-10.2f%%\n",
			tag, m.Precision*100, m.Recall*100, m.F1Score*100)
	}

	return metricsMap
}

func main() {
	// 1. Train Data Setup
	trainingSentences := [][]string{
		{"Steve", "Jobs", "visited", "Cupertino"},
		{"Alice", "saw", "Paris"},
		{"Bob", "went", "to", "London"},
	}
	trainingTags := [][]string{
		{"B-PER", "I-PER", "O", "B-LOC"},
		{"B-PER", "O", "B-LOC"},
		{"B-PER", "O", "O", "B-LOC"},
	}

	model := NewHMM()
	model.Train(trainingSentences, trainingTags)

	// 2. Separate Validation/Test Dataset Setup
	validationSentences := [][]string{
		{"Alice", "visited", "Cupertino"}, // Mix of known names & places
		{"Charlie", "saw", "London"},      // "Charlie" is a completely new unknown word token
	}
	validationTags := [][]string{
		{"B-PER", "O", "B-LOC"},
		{"B-PER", "O", "B-LOC"},
	}

	// 3. Execute Matrix Evaluation Pipeline
	fmt.Println("Evaluating Model Metrics against Validation Set...")
	model.Evaluate(validationSentences, validationTags)
}
