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
	Transition   [NumTags][NumTags]float64      // [prevTag][nextTag] = log P(nextTag  | prevTag)
	Emission     [NumTags][NumSyllables]float64 // [tag][syllable]    = log P(syllable | tag)
	StartTags    [NumTags]float64               // [tag]              = log P(tag)
	TagCounts    [NumTags]int                   // [tag] = # occurances of tag in training corpus
	NumSentences int                            // < 0 during Accumulate, > 0 after Generate
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

// Accumulates counts of tag transmissions and emissions
// Separate models can Accumulate concurrently and be merged at the end.
// NOTE: sentences must contain every known syllable at least once.
func (h *HMM) Accumulate(sentence []Syllable, tags []Tag) {
	if h.NumSentences > 0 {
		panic("HMM.Accumulate called after HMM.Generate has already been called")
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
}

func (h *HMM) Merge(rhs HMM) {
	if h.NumSentences > 0 {
		panic("HMM.Merge called after HMM.Generate has already been called")
	}
	if rhs.NumSentences > 0 {
		panic("HMM.Merge called with argument on which HMM.Generate has already been called")
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
}

// Computes probabilities using Laplace (Add-1) smoothing and saves them as log values
func (h *HMM) Generate() {
	if h.NumSentences == 0 {
		panic("HMM.Generate called without having first called HMM.Accumulate at least once")
	}
	if h.NumSentences > 0 {
		panic("HMM.Generate called twice")
	}
	numTags := float64(NumTags)
	scale := 1.0 / (float64(-h.NumSentences) + numTags)
	for currTag := range NumTags {
		h.StartTags[currTag] = math.Log((h.StartTags[currTag] + 1.0) * scale)
	}
	for prevTag := range NumTags {
		totalTransitionsOut := numTags
		for nextTag := range NumTags {
			nextTag := Tag(nextTag)
			totalTransitionsOut += h.Transition[prevTag][nextTag]
		}
		for nextTag := range NumTags {
			h.Transition[prevTag][nextTag] = math.Log((h.Transition[prevTag][nextTag] + 1.0) / totalTransitionsOut)
		}
	}
	for currTag := range NumTags {
		totalEmissionsOut := float64(h.TagCounts[currTag] + NumSyllables)
		for syllable := range NumSyllables {
			h.Emission[currTag][syllable] = math.Log((h.Emission[currTag][syllable] + 1.0) / totalEmissionsOut)
		}
		h.Emission[currTag][0] = math.Log(1.0 / totalEmissionsOut)
	}
	h.NumSentences *= -1 // mark as generated
}

// Predict uses the Viterbi algorithm operating in log space
func (h HMM) Predict(tokens []Syllable) []Tag {
	if h.NumSentences <= 0 {
		panic("HMM.Predict called without having first called HMM.Generate")
	}
	if len(tokens) == 0 {
		return nil
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
	return resultTags
}

// Evaluate runs predictions on test data and prints Precision, Recall, and F1 per tag
func Evaluate(predictedTags, expectedTags [][]Tag) [NumTags]Metrics {
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
	// TODO: Accumulate (unlike Predict) can be streamed and processed sentence by sentence.
	// TODO: Models can be merged for concurrent processing.
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
		// TODO: Process sentences concurrently in separate models then merge models with HMM.Merge
		for i := range trainingSentences {
			model.Accumulate(trainingSentences[i], trainingTags[i])
		}
		model.Generate()
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
	metrics := Evaluate(predictedTags, expectedTags)

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
