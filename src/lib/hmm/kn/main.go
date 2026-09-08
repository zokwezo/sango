package main

import (
	"fmt"
	"math"
)

// Global constants for Kneser-Ney
const Discount = 0.75

// TrainingInstance represents a sequence of tokens and their true state/label sequence.
type TrainingInstance struct {
	Tokens []string
	States []string
}

type HMM struct {
	States      map[string]bool
	Vocab       map[string]bool
	Transitions map[string]map[string]float64 // count(s_t | s_{t-1})
	Emissions   map[string]map[string]float64 // count(w_t | s_t)
	StateCounts map[string]float64            // total occurrences of each state
	UniquePairs float64                       // total unique (state, word) types in corpus
}

func NewHMM() *HMM {
	return &HMM{
		States:      make(map[string]bool),
		Vocab:       make(map[string]bool),
		Transitions: make(map[string]map[string]float64),
		Emissions:   make(map[string]map[string]float64),
		StateCounts: make(map[string]float64),
	}
}

// Train populates the raw count matrices from our tiny low-resource dataset.
func (h *HMM) Train(data []TrainingInstance) {
	for _, inst := range data {
		for i := 0; i < len(inst.Tokens); i++ {
			state := inst.States[i]
			word := inst.Tokens[i]

			h.States[state] = true
			h.Vocab[word] = true
			h.StateCounts[state]++

			// Count Emissions
			if _, exists := h.Emissions[state]; !exists {
				h.Emissions[state] = make(map[string]float64)
			}
			h.Emissions[state][word]++

			// Count Transitions (for simplicity, ignoring structural boundary tokens here)
			if i > 0 {
				prevState := inst.States[i-1]
				if _, exists := h.Transitions[prevState]; !exists {
					h.Transitions[prevState] = make(map[string]float64)
				}
				h.Transitions[prevState][state]++
			}
		}
	}

	// Count unique (state, word) type configurations for Kneser-Ney denominator
	for _, words := range h.Emissions {
		h.UniquePairs += float64(len(words))
	}
}

// GetEmissionNoSmoothing calculates pure MLE probabilities.
// Returns 0.0 for any token not explicitly bound to that state in training.
func (h *HMM) GetEmissionNoSmoothing(state, word string) float64 {
	totalStateCount := h.StateCounts[state]
	if totalStateCount == 0 {
		return 0.0
	}
	return h.Emissions[state][word] / totalStateCount
}

// GetEmissionKneserNey calculates the smoothed emission probability.
// It uses absolute discounting and backs off to the token's structural versatility.
func (h *HMM) GetEmissionKneserNey(state, word string) float64 {
	totalStateCount := h.StateCounts[state]
	if totalStateCount == 0 {
		return 1.0 / float64(len(h.Vocab)) // Uniform fallback if state is totally unseen
	}

	rawCount := h.Emissions[state][word]
	uniqueWordsInState := float64(len(h.Emissions[state]))

	// 1. Calculate Left Term: Discounted MLE
	leftTerm := math.Max(rawCount-Discount, 0.0) / totalStateCount

	// 2. Calculate Lambda: Back-off weight normalization
	lambda := (Discount / totalStateCount) * uniqueWordsInState

	// 3. Calculate Continuation Probability: Versatility of the word
	// How many unique states have emitted this specific word?
	statesEmittingWord := 0.0
	for _, words := range h.Emissions {
		if _, exists := words[word]; exists {
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
	rawCount := h.Transitions[fromState][toState]
	totalTransitionsFromState := 0.0
	for _, count := range h.Transitions[fromState] {
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

func main() {
	// Tiny low-resource corpus
	// Notice how "Francisco" ALWAYS appears with "NOUN", making it highly rigid.
	// "Apple" appears with both "NOUN" and "PROPN" (the company), making it versatile.
	trainingData := []TrainingInstance{
		{Tokens: []string{"San", "Francisco"}, States: []string{"PROPN", "PROPN"}},
		{Tokens: []string{"Eat", "an", "apple"}, States: []string{"VERB", "DET", "NOUN"}},
		{Tokens: []string{"Apple", "announces", "iPhone"}, States: []string{"PROPN", "VERB", "PROPN"}},
	}

	hmm := NewHMM()
	hmm.Train(trainingData)

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
}
