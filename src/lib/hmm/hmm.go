// Hidden Markov Model (HMM) with Kneser-Ney smoothing

// Used for diacritic restoration of Sango text.

package hmm

import (
	"fmt"
	"io/ioutil"
	"math"

	"google.golang.org/protobuf/proto"
)

// Global constants for Kneser-Ney
const Discount = 0.75

// TrainingInstance represents a sequence of tokens and their true state/label sequence.
// In this context, a Token is a Sango word without diacritics and a State is the
// same word with diacritics.
type TrainingInstance struct {
	Tokens []string
	States []string
}

func NewHMM() *HMM {
	return &HMM{
		States:              make(map[string]bool),
		Vocab:               make(map[string]bool),
		TransitionFromState: make(map[string]*Then),
		EmissionFromState:   make(map[string]*Then),
		StateCount:          make(map[string]float64),
	}
}

// Populates the raw count matrices from our tiny low-resource dataset.
// This is a two-pass algorithm
func (h *HMM) Train(data []TrainingInstance) error {
	for _, inst := range data {
		for i := 0; i < len(inst.Tokens); i++ {
			state := inst.States[i]
			word := inst.Tokens[i]

			h.States[state] = true
			h.Vocab[word] = true
			h.StateCount[state]++

			// Count Emissions
			if _, exists := h.EmissionFromState[state]; !exists {
				h.EmissionFromState[state] = &Then{To: make(map[string]float64)}
			}
			h.EmissionFromState[state].To[word]++

			// Count Transitions (for simplicity, ignoring structural boundary tokens here)
			if i > 0 {
				prevState := inst.States[i-1]
				if _, exists := h.TransitionFromState[prevState]; !exists {
					h.TransitionFromState[prevState] = &Then{To: make(map[string]float64)}
				}
				h.TransitionFromState[prevState].To[state]++
			}
		}
	}

	// Count unique (state, word) type configurations for Kneser-Ney denominator
	for _, emissionFromState := range h.EmissionFromState {
		h.UniquePairs += float64(len(emissionFromState.To))
	}

	return nil
}

// GetEmissionNoSmoothing calculates pure MLE probabilities.
// Returns 0.0 for any token not explicitly bound to that state in training.
func (h *HMM) GetEmissionNoSmoothing(state, word string) float64 {
	totalStateCount := h.StateCount[state]
	if totalStateCount == 0 {
		return 0.0
	}
	if _, exists := h.EmissionFromState[state]; !exists {
		h.EmissionFromState[state] = &Then{To: make(map[string]float64)}
	}
	return h.EmissionFromState[state].To[word] / totalStateCount
}

// GetEmissionKneserNey calculates the smoothed emission probability.
// It uses absolute discounting and backs off to the token's structural versatility.
func (h *HMM) GetEmissionKneserNey(state, word string) float64 {
	totalStateCount := h.StateCount[state]
	if totalStateCount == 0 {
		return 1.0 / float64(len(h.Vocab)) // Uniform fallback if state is totally unseen
	}

	if _, exists := h.EmissionFromState[state]; !exists {
		h.EmissionFromState[state] = &Then{To: make(map[string]float64)}
	}
	rawCount := h.EmissionFromState[state].To[word]
	uniqueWordsInState := float64(len(h.EmissionFromState[state].To))

	// 1. Calculate Left Term: Discounted MLE
	leftTerm := math.Max(rawCount-Discount, 0.0) / totalStateCount

	// 2. Calculate Lambda: Back-off weight normalization
	lambda := (Discount / totalStateCount) * uniqueWordsInState

	// 3. Calculate Continuation Probability: Versatility of the word
	// How many unique states have emitted this specific word?
	statesEmittingWord := 0.0
	for _, emissionFromState := range h.EmissionFromState {
		if _, exists := emissionFromState.To[word]; exists {
			statesEmittingWord++
		}
	}

	pContinuation := statesEmittingWord / h.UniquePairs

	// Final Kneser-Ney formula blending
	return leftTerm + (lambda * pContinuation)
}

//////////////////////////////////////////////////////////////////////////////

// GetTransitionProbability calculates the transition probability P(toState | fromState)
// with a very simple Laplace add-alpha smoothing step to ensure no zero-probabilities exist.
func (h *HMM) GetTransitionProbability(fromState, toState string) float64 {
	const alpha = 0.1
	if _, exists := h.TransitionFromState[fromState]; !exists {
		h.TransitionFromState[fromState] = &Then{To: make(map[string]float64)}
	}
	rawCount := h.TransitionFromState[fromState].To[toState]
	totalTransitionsFromState := 0.0
	for _, count := range h.TransitionFromState[fromState].To {
		totalTransitionsFromState += count
	}

	// Apply basic smoothing to transitions so they don't break Viterbi paths with 0.0
	return (rawCount + alpha) / (totalTransitionsFromState + (alpha * float64(len(h.States))))
}

// Viterbi predicts the most likely hidden state sequence for a slice of tokens.
// Set 'useSmoothing' to true to see Kneser-Ney in action, or false for MLE.
func (h *HMM) Viterbi(tokens []string, useSmoothing bool) []string {
	n := len(tokens)
	if n == 0 {
		return nil
	}

	// Create list of unique states to index consistently
	stateList := make([]string, 0, len(h.States))
	for state := range h.States {
		stateList = append(stateList, state)
	}
	numStates := len(stateList)

	// viterbi[t][s] stores the min negative log probability at time t for state s
	viterbi := make([]map[string]float64, n)
	// backpointer[t][s] stores the name of the state at t-1 that led to state s at t
	backpointer := make([]map[string]string, n)

	for i := range viterbi {
		viterbi[i] = make(map[string]float64)
		backpointer[i] = make(map[string]string)
	}

	// 1. Initialization (t = 0)
	// For simplicity, we assume a uniform starting distribution across all states
	startProb := 1.0 / float64(numStates)
	for _, state := range stateList {
		var emissionProb float64
		if useSmoothing {
			emissionProb = h.GetEmissionKneserNey(state, tokens[0])
		} else {
			emissionProb = h.GetEmissionNoSmoothing(state, tokens[0])
		}

		// Handle complete Out-of-Vocabulary tokens cleanly
		if emissionProb == 0 && !h.Vocab[tokens[0]] {
			emissionProb = 1.0 / float64(len(h.Vocab))
		}

		// We work in negative log space: lower cost = higher probability
		if emissionProb > 0 {
			viterbi[0][state] = -math.Log(startProb) - math.Log(emissionProb)
		} else {
			viterbi[0][state] = math.Inf(1) // Absolute zero probability path
		}
	}

	// 2. Recursion (t = 1 to n-1)
	for t := 1; t < n; t++ {
		token := tokens[t]
		for _, currState := range stateList {

			var emissionProb float64
			if useSmoothing {
				emissionProb = h.GetEmissionKneserNey(currState, token)
			} else {
				emissionProb = h.GetEmissionNoSmoothing(currState, token)
			}

			// Handle Out-of-Vocabulary tokens at current step
			if emissionProb == 0 && !h.Vocab[token] {
				emissionProb = 1.0 / float64(len(h.Vocab))
			}

			minCost := math.Inf(1)
			bestPrevState := ""

			// Find the previous state that minimizes the total sequence cost
			for _, prevState := range stateList {
				transProb := h.GetTransitionProbability(prevState, currState)

				// Total cost up to this point
				cost := viterbi[t-1][prevState] - math.Log(transProb)

				if cost < minCost {
					minCost = cost
					bestPrevState = prevState
				}
			}

			viterbi[t][currState] = minCost - math.Log(emissionProb)
			backpointer[t][currState] = bestPrevState
		}
	}

	// 3. Termination & Backtracking
	minFinalCost := math.Inf(1)
	bestFinalState := ""
	for _, state := range stateList {
		if viterbi[n-1][state] < minFinalCost {
			minFinalCost = viterbi[n-1][state]
			bestFinalState = state
		}
	}

	// Trace back pointers to extract the optimal structural labels
	path := make([]string, n)
	if bestFinalState == "" {
		// Fallback if everything zeroed out
		bestFinalState = stateList[0]
	}

	path[n-1] = bestFinalState
	for t := n - 1; t > 0; t-- {
		path[t-1] = backpointer[t][path[t]]
	}

	return path
}

//////////////////////////////////////////////////////////////////////////////

func MainTrain(modelOutputFilename string) error {
	// Tiny low-resource corpus
	// Notice how "Francisco" ALWAYS appears with "NOUN", making it highly rigid.
	// "Apple" appears with both "NOUN" and "PROPN" (the company), making it versatile.
	trainingData := []TrainingInstance{
		{Tokens: []string{"San", "Francisco"}, States: []string{"PROPN", "PROPN"}},
		{Tokens: []string{"Eat", "an", "apple"}, States: []string{"VERB", "DET", "NOUN"}},
		{Tokens: []string{"Apple", "announces", "iPhone"}, States: []string{"PROPN", "VERB", "PROPN"}},
	}

	hmm := NewHMM()
	err := hmm.Train(trainingData)
	if err != nil {
		return err
	}

	// Write the model to disk.
	out, err := proto.Marshal(hmm)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(modelOutputFilename, out, 0644); err != nil {
		return err
	}

	return nil
}

func MainPredict(modelInputFilename string) error {
	in, err := ioutil.ReadFile(modelInputFilename)
	if err != nil {
		return err
	}
	hmm := NewHMM()
	if err := proto.Unmarshal(in, hmm); err != nil {
		return err
	}

	fmt.Println("--- EMISSION COMPARISONS ---")

	// Example 1: Evaluating a highly rigid token in a completely novel context
	// Does "Francisco" make sense as a generic VERB?
	fmt.Printf("\nToken: 'Francisco' | Evaluated State: 'VERB'\n")
	fmt.Printf("  No Smoothing:        %0.5f (Completely rejects it)\n", hmm.GetEmissionNoSmoothing("VERB", "Francisco"))
	fmt.Printf("  Kneser-Ney Sm.:      %0.5f (Keeps mass very low because it lacks versatility)\n", hmm.GetEmissionKneserNey("VERB", "Francisco"))

	// Example 2: Evaluating a versatile token in an unseen context
	// We've never explicitly trained "Apple" as a generic "DET" (determiner),
	// but it is a versatile word in the corpus.
	fmt.Printf("\nToken: 'Apple' | Evaluated State: 'DET'\n")
	fmt.Printf("  No Smoothing:        %0.5f (Strictly breaks the sequence sequence match)\n", hmm.GetEmissionNoSmoothing("DET", "Apple"))
	fmt.Printf("  Kneser-Ney Sm.:      %0.5f (Grants a higher fallback probability because 'Apple' is versatile)\n", hmm.GetEmissionKneserNey("DET", "Apple"))

	// Evaluation test sequence:
	// "Apple" is used here as a common NOUN in the object position, rather than the PROPN company name.
	testTokens := []string{"Eat", "an", "Apple"}

	fmt.Println("Tokens to predict:", testTokens)

	// Predict with NO smoothing on Emissions
	noSmoothPath := hmm.Viterbi(testTokens, false)
	fmt.Println("Prediction (No Smoothing): ", noSmoothPath)

	// Predict WITH Kneser-Ney smoothing on Emissions
	knSmoothPath := hmm.Viterbi(testTokens, true)
	fmt.Println("Prediction (Kneser-Ney):  ", knSmoothPath)

	return nil
}
