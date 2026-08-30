package main

import (
	"fmt"
	"math"
)

type (
	Syllable byte
	Tag      byte
)

var S = map[string]Syllable{
	"alice":     Syllable(65),
	"bob":       Syllable(66),
	"charlie":   Syllable(67),
	"cupertino": Syllable(68),
	"jobs":      Syllable(69),
	"london":    Syllable(70),
	"paris":     Syllable(71),
	"saw":       Syllable(72),
	"steve":     Syllable(73),
	"to":        Syllable(74),
	"visited":   Syllable(75),
	"went":      Syllable(76),
}

var T = map[string]Tag{
	"B-LOC": Tag(49),
	"B-PER": Tag(50),
	"I-PER": Tag(51),
	"O":     Tag(52),
}

func (tag Tag) String() string {
	for s, t := range T {
		if t == tag {
			return s
		}
	}
	return "UNKNOWN_TAG"
}

const (
	UNKNOWN_SYLLABLE Syllable = Syllable(33)
	UNKNOWN_TAG      Tag      = Tag(48)
	LO_TAG           Tag      = Tag(49)
	ME_TAG           Tag      = Tag(50)
	HI_TAG           Tag      = Tag(51)
)

type HMM struct {
	Transition map[Tag]map[Tag]float64
	Emission   map[Tag]map[Syllable]float64
	StartTags  map[Tag]float64
	TagCounts  map[Tag]int
	Vocab      map[Syllable]bool
	Tags       []Tag
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
		Transition: make(map[Tag]map[Tag]float64),
		Emission:   make(map[Tag]map[Syllable]float64),
		StartTags:  make(map[Tag]float64),
		TagCounts:  make(map[Tag]int),
		Vocab:      make(map[Syllable]bool),
		Tags:       []Tag{},
	}
}

// Train computes probabilities using Laplace (Add-1) smoothing and saves them as log values
func (h *HMM) Train(sentences [][]Syllable, tags [][]Tag) {
	totalSentences := float64(len(sentences))
	rawStart := make(map[Tag]float64)
	rawTrans := make(map[Tag]map[Tag]float64)
	rawEmis := make(map[Tag]map[Syllable]float64)

	tagMap := make(map[Tag]bool)
	for i := range sentences {
		prevTag := UNKNOWN_TAG
		for j := range sentences[i] {
			syllable := sentences[i][j]
			tag := tags[i][j]

			h.TagCounts[tag]++
			h.Vocab[syllable] = true
			tagMap[tag] = true

			if j == 0 {
				rawStart[tag]++
			} else {
				if rawTrans[prevTag] == nil {
					rawTrans[prevTag] = make(map[Tag]float64)
				}
				rawTrans[prevTag][tag]++
			}

			if rawEmis[tag] == nil {
				rawEmis[tag] = make(map[Syllable]float64)
			}
			rawEmis[tag][syllable]++
			prevTag = tag
		}
	}

	for tag := range tagMap {
		h.Tags = append(h.Tags, tag)
	}

	numTags := float64(len(h.Tags))
	numSyllables := float64(len(h.Vocab))

	for _, tag := range h.Tags {
		count := rawStart[tag]
		h.StartTags[tag] = math.Log((count + 1.0) / (totalSentences + numTags))
	}

	for _, prevTag := range h.Tags {
		h.Transition[prevTag] = make(map[Tag]float64)
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
		h.Emission[tag] = make(map[Syllable]float64)
		totalEmissionsOut := float64(h.TagCounts[tag])
		for syllable := range h.Vocab {
			count := rawEmis[tag][syllable]
			prob := (count + 1.0) / (totalEmissionsOut + numSyllables)
			h.Emission[tag][syllable] = math.Log(prob)
		}
		h.Emission[tag][UNKNOWN_SYLLABLE] = math.Log(1.0 / (totalEmissionsOut + numSyllables))
	}
}

// Predict uses the Viterbi algorithm operating in log space
// NOTE: what is the difference between a token and a syllable
func (h *HMM) Predict(tokens []Syllable) []Tag {
	if len(tokens) == 0 {
		return nil
	}

	numTokens := len(tokens)
	viterbi := make([]map[Tag]float64, numTokens)
	backpointer := make([]map[Tag]Tag, numTokens)

	viterbi[0] = make(map[Tag]float64)
	backpointer[0] = make(map[Tag]Tag)
	firstSyllable := tokens[0]

	for _, tag := range h.Tags {
		emissionLog := h.Emission[tag][firstSyllable]
		if !h.Vocab[firstSyllable] {
			emissionLog = h.Emission[tag][UNKNOWN_SYLLABLE]
		}
		viterbi[0][tag] = h.StartTags[tag] + emissionLog
	}

	for t := 1; t < numTokens; t++ {
		viterbi[t] = make(map[Tag]float64)
		backpointer[t] = make(map[Tag]Tag)
		syllable := tokens[t]

		for _, currTag := range h.Tags {
			maxLogProb := math.Inf(-1)
			bestPrevTag := h.Tags[0]

			currEmissionLog := h.Emission[currTag][syllable]
			if !h.Vocab[syllable] {
				currEmissionLog = h.Emission[currTag][UNKNOWN_SYLLABLE]
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

	resultTags := make([]Tag, numTokens)
	currTag := bestFinalTag
	for t := lastIdx; t >= 0; t-- {
		resultTags[t] = currTag
		currTag = backpointer[t][currTag]
	}

	return resultTags
}

// Evaluate runs predictions on test data and prints Precision, Recall, and F1 per tag
func (h *HMM) Evaluate(testSentences [][]Syllable, testTags [][]Tag) map[Tag]*Metrics {
	metricsMap := make(map[Tag]*Metrics)
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
	fmt.Println("-----------------------------------------------------")

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
	trainingSentences := [][]Syllable{
		{S["steve"], S["jobs"], S["visited"], S["cupertino"]},
		{S["alice"], S["saw"], S["paris"]},
		{S["bob"], S["went"], S["to"], S["london"]},
	}
	trainingTags := [][]Tag{
		{T["B-PER"], T["I-PER"], T["O"], T["B-LOC"]},
		{T["B-PER"], T["O"], T["B-LOC"]},
		{T["B-PER"], T["O"], T["O"], T["B-LOC"]},
	}

	model := NewHMM()
	model.Train(trainingSentences, trainingTags)

	// 2. Separate Validation/Test Dataset Setup
	validationSentences := [][]Syllable{
		{S["alice"], S["visited"], S["cupertino"]}, // Mix of known names & places
		{S["charlie"], S["saw"], S["london"]},      // "charlie" is a completely new unknown syllable
	}
	validationTags := [][]Tag{
		{T["B-PER"], T["O"], T["B-LOC"]},
		{T["B-PER"], T["O"], T["B-LOC"]},
	}

	// 3. Execute Matrix Evaluation Pipeline
	fmt.Println("Evaluating Model Metrics against Validation Set...")
	model.Evaluate(validationSentences, validationTags)
}
