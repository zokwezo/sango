// Hidden Markov Model (HMM) with Kneser-Ney smoothing

// Used for diacritic restoration of Sango text.

package hmm

import (
	"math"
)

// Global constants for Kneser-Ney
const Discount = 0.75

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
