package main

import (
	"encoding/json"
	"fmt"
	"math"
)

type (
	Syllable int8
	Tag      int8
)

const NumSyllables int = 13

func (syllable Syllable) String() string {
	return [NumSyllables]string{
		"UNKNOWN",
		"alice",
		"bob",
		"charlie",
		"cupertino",
		"jobs",
		"london",
		"paris",
		"saw",
		"steve",
		"to",
		"visited",
		"went",
	}[syllable]
}

var syllableFromNameMap_ = func() map[string]Syllable {
	m := make(map[string]Syllable, NumSyllables)
	for i := range NumSyllables {
		s := Syllable(i)
		m[s.String()] = s
	}
	return m
}()

func S(syllableName string) Syllable {
	return syllableFromNameMap_[syllableName]
}

// TAGS
const (
	UnknownPitch Tag = iota
	LowPitch
	MedPitch
	HighPitch
	NumTags
)

func (tag Tag) String() string {
	return [NumTags]string{
		"?",
		"_",
		":",
		"^",
	}[tag]
}

type HMM struct {
	Transition [NumTags][NumTags]float64      `json:"transition"` // [prevTag][nextTag] = log P(nextTag  | prevTag)
	Emission   [NumTags][NumSyllables]float64 `json:"emission"`   // [tag][syllable]    = log P(syllable | tag)
	StartTags  [NumTags]float64               `json:"startTags"`  // [tag]              = log P(tag)
	TagCounts  [NumTags]int                   `json:"tagCounts"`  // [tag] = # occurances of tag in training corpus
}

// Metrics holds the classification performance statistics
type Metrics struct {
	Precision float64 // fraction of predicted that are correct
	Recall    float64 // fraction of correct that are predicted
	F1Score   float64 // harmonic mean of Precision and Recall
	TP        int     // # true  positives
	FP        int     // # false positives
	FN        int     // # false negatives
}

// Train computes probabilities using Laplace (Add-1) smoothing and saves them as log values
// NOTE: Training corpus must contain every syllable at least once.
func (h *HMM) Train(sentences [][]Syllable, tags [][]Tag) {
	totalSentences := float64(len(sentences))
	rawStart := [NumTags]float64{}
	rawTrans := [NumTags][NumTags]float64{}
	rawEmis := [NumTags][NumSyllables]float64{}
	for i := range sentences {
		prevTag := UnknownPitch
		for j := range sentences[i] {
			syllable := sentences[i][j]
			tag := tags[i][j]
			h.TagCounts[tag]++

			if j == 0 {
				rawStart[tag]++
			} else {
				rawTrans[prevTag][tag]++
			}

			rawEmis[tag][syllable]++
			prevTag = tag
		}
	}
	numTags := float64(NumTags)
	numSyllables := float64(NumSyllables)
	for t := range NumTags {
		tag := Tag(t)
		count := rawStart[tag]
		h.StartTags[tag] = math.Log((count + 1.0) / (totalSentences + numTags))
	}
	for p := range NumTags {
		prevTag := Tag(p)
		h.Transition[prevTag] = [NumTags]float64{}
		totalTransitionsOut := 0.0
		for n := range NumTags {
			nextTag := Tag(n)
			totalTransitionsOut += rawTrans[prevTag][nextTag]
		}
		for n := range NumTags {
			nextTag := Tag(n)
			count := rawTrans[prevTag][nextTag]
			prob := (count + 1.0) / (totalTransitionsOut + numTags)
			h.Transition[prevTag][nextTag] = math.Log(prob)
		}
	}
	for t := range NumTags {
		tag := Tag(t)
		h.Emission[tag] = [NumSyllables]float64{}
		totalEmissionsOut := float64(h.TagCounts[tag])
		for syllable := range NumSyllables {
			count := rawEmis[tag][syllable]
			prob := (count + 1.0) / (totalEmissionsOut + numSyllables)
			h.Emission[tag][syllable] = math.Log(prob)
		}
		h.Emission[tag][0] = math.Log(1.0 / (totalEmissionsOut + numSyllables))
	}
}

// Predict uses the Viterbi algorithm operating in log space
func (h HMM) Predict(tokens []Syllable) []Tag {
	if len(tokens) == 0 {
		return nil
	}
	numTokens := len(tokens)
	viterbi := make([][NumTags]float64, numTokens)
	backpointer := make([][NumTags]Tag, numTokens)
	viterbi[0] = [NumTags]float64{}
	backpointer[0] = [NumTags]Tag{}
	firstSyllable := tokens[0]
	for t := range NumTags {
		tag := Tag(t)
		emissionLog := h.Emission[tag][firstSyllable]
		viterbi[0][tag] = h.StartTags[tag] + emissionLog
	}
	for t := 1; t < numTokens; t++ {
		viterbi[t] = [NumTags]float64{}
		backpointer[t] = [NumTags]Tag{}
		syllable := tokens[t]
		for c := range NumTags {
			currTag := Tag(c)
			maxLogProb := math.Inf(-1)
			bestPrevTag := UnknownPitch
			currEmissionLog := h.Emission[currTag][syllable]
			for p := range NumTags {
				prevTag := Tag(p)
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
	lastIdx := numTokens - 1
	for t := range NumTags {
		tag := Tag(t)
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
func (h HMM) Evaluate(predictedTags, expectedTags [][]Tag) [NumTags]Metrics {
	metrics := [NumTags]Metrics{}
	if len(predictedTags) != len(expectedTags) {
		panic("bad number of expected tags passed to Evaluate")
	}
	for i := range predictedTags {
		predicted := predictedTags[i]
		expected := expectedTags[i]
		for j := range expected {
			trueTag := expected[j]
			predTag := predicted[j]
			if trueTag == predTag {
				metrics[trueTag].TP++
			} else {
				metrics[predTag].FP++
				metrics[trueTag].FN++
			}
		}
	}
	return metrics
}

func main() {
	// 1. Train model and persist to disk in JSON format
	var serializedModel string
	{
		trainingSentences := [][]Syllable{
			{S("steve"), S("visited"), S("cupertino")},
			{S("alice"), S("saw"), S("paris")},
			{S("bob"), S("went"), S("to"), S("london")},
		}
		trainingTags := [][]Tag{
			{MedPitch, HighPitch, LowPitch},
			{MedPitch, HighPitch, LowPitch},
			{MedPitch, HighPitch, HighPitch, LowPitch},
		}
		model := HMM{}
		model.Train(trainingSentences, trainingTags)
		jsonBytes, err := json.Marshal(model)
		if err != nil {
			panic(err)
		}
		serializedModel = string(jsonBytes)
	}

	// 2. Read in JSON-serialized model from disk
	fmt.Println("Marshaled JSON string:")
	fmt.Println(serializedModel)
	var model HMM
	err := json.Unmarshal([]byte(serializedModel), &model)
	if err != nil {
		panic(err)
	}
	fmt.Println("Unmarshaled model:")
	fmt.Printf("%#v\n", model)

	// 3. Predict tags on test data
	testSentences := [][]Syllable{
		{S("alice"), S("visited"), S("cupertino")}, // Mix of known names & places
		{S("charlie"), S("saw"), S("london")},      // "charlie" is a completely new unknown syllable
	}
	var predictedTags [][]Tag
	for i := range testSentences {
		predictedTags = append(predictedTags, model.Predict(testSentences[i]))
	}

	// 4. Evaluate against expected tags and generate metrics
	expectedTags := [][]Tag{
		{LowPitch, HighPitch, LowPitch},
		{MedPitch, HighPitch, LowPitch},
	}
	fmt.Println("Evaluating Model Metrics against Test Set...")
	metrics := model.Evaluate(predictedTags, expectedTags)

	// 5. Output metrics
	fmt.Println("")
	fmt.Println("PITCH | PRECISION |  RECALL   | F1-SCORE   ")
	fmt.Println("------+-----------+-----------+------------")
	for t, m := range metrics {
		tag := Tag(t)
		if m.TP+m.FP > 0 {
			m.Precision = float64(m.TP) / float64(m.TP+m.FP)
		}
		if m.TP+m.FN > 0 {
			m.Recall = float64(m.TP) / float64(m.TP+m.FN)
		}
		if m.Precision+m.Recall > 0 {
			m.F1Score = 2 * (m.Precision * m.Recall) / (m.Precision + m.Recall)
		}
		fmt.Printf("  %s   | %6.2f %%  | %6.2f %%  | %6.2f %%\n",
			tag, m.Precision*100, m.Recall*100, m.F1Score*100)
	}
}
