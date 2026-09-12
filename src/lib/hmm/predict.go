// Uses Hidden Markov Model (HMM) with Kneser-Ney smoothing for diacritic restoration of Sango text.

package hmm

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strings"

	"github.com/zokwezo/sango/src/lib/sse"
	"google.golang.org/protobuf/proto"
)

// Global constants for Kneser-Ney
const Discount = 0.75

func MainPredict(inputTextFilename, modelInputFilename string) error {
	modelInputWireFormat, err := ioutil.ReadFile(modelInputFilename)
	if err != nil {
		return err
	}
	m := Model{}
	if err := proto.Unmarshal(modelInputWireFormat, &m); err != nil {
		return err
	}
	h := HMM{}
	h.FromModel(&m)

	inputText, err := os.ReadFile(inputTextFilename)
	if err != nil {
		return err
	}

	codes, tis := PrepareInputText(string(inputText))

	n := len(codes)
	index := 0
	outputIndex := 0
	for k := range tis {
		// Predict using Kneser-Ney smoothing on Emissions
		for i := range tis[k].State {
			tis[k].State[i] = ""
		}

		tis[k].State = h.Viterbi(tis[k].Token)
		if len(tis[k].State) != len(tis[k].Token) {
			panic("len(tis[k].State) != len(tis[k].Token)")
		}

		// Replace the source token with the predicted state if the toneless version matches.
		for i, modelToken := range tis[k].Token {
			outputTextBuilder := strings.Builder{}
			modelState := tis[k].State[i]
			index = tis[k].Index[i] - 1 // Index is one past the real index
			if index < 0 || index >= n {
				panic("index is out of bounds")
			}
			inputToken := sse.BuilderToString(codes[index].WriteAsTonelessTo)
			if inputToken != modelToken {
				panic("inputToken != modelToken")
			}
			newCodes, err := sse.Utf8ToSSEs(modelState, sse.FromLemma)
			if err != nil {
				panic(err)
			}
			if len(newCodes) != 1 {
				panic("len(newCodes) != 1")
			}
			if inputToken == sse.BuilderToString(newCodes[0].WriteAsTonelessTo) {
				const ignoreSpaceAndCase sse.SSE = 0x8FFF_FFFF_FFFF_FFFF
				if newCodes[0]&ignoreSpaceAndCase != codes[index]&ignoreSpaceAndCase {
					const pitchMask sse.SSE = 0x0_003_003_003_003_003
					codes[index] &= ^pitchMask
					codes[index] |= pitchMask & newCodes[0]
				}
			}
			for ; outputIndex <= index; outputIndex++ {
				codes[outputIndex].WriteAsLemmaTo(&outputTextBuilder)
			}
			fmt.Printf("%v", outputTextBuilder.String())
		}
	}

	return nil
}

// GetEmissionKneserNey calculates the smoothed emission probability.
// It uses absolute discounting and backs off to the token's structural versatility.
func (h *HMM) GetEmissionKneserNey(state, token string) float64 {
	totalStateCount := float64(h.StateCounts[state])
	if totalStateCount == 0.0 {
		return 1.0 / float64(len(h.Tokens)) // Uniform fallback if state is totally unseen
	}

	rawCount := h.Emissions[state][token]
	uniqueTokensInState := len(h.Emissions[state])

	// 1. Calculate Left Term: Discounted MLE
	leftTerm := math.Max(float64(rawCount)-Discount, 0.0) / totalStateCount

	// 2. Calculate Lambda: Back-off weight normalization
	lambda := (Discount / totalStateCount) * float64(uniqueTokensInState)

	// 3. Calculate Continuation Probability: Versatility of the token
	// How many unique states have emitted this specific token?
	statesEmittingToken := 0
	for _, tokens := range h.Emissions {
		if _, exists := tokens[token]; exists {
			statesEmittingToken++
		}
	}

	pContinuation := float64(statesEmittingToken) / float64(h.UniquePairs)

	// Final Kneser-Ney formula blending
	return leftTerm + (lambda * pContinuation)
}

// GetTransitionProbability calculates the transition probability P(toState | fromState)
// with a very simple Laplace add-alpha smoothing step to ensure no zero-probabilities exist.
func (h *HMM) GetTransitionProbability(fromState, toState string) float64 {
	const alpha = 0.1
	rawCount := h.Transitions[fromState][toState]
	totalTransitionsFromState := int64(0)
	for _, count := range h.Transitions[fromState] {
		totalTransitionsFromState += count
	}

	// Apply basic smoothing to transitions so they don't break Viterbi paths with 0.0
	return (float64(rawCount) + alpha) / (float64(totalTransitionsFromState) + (alpha * float64(len(h.States))))
}

// Viterbi predicts the most likely hidden state sequence for a slice of tokens.
func (h *HMM) Viterbi(tokens []string) []string {
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
		emissionProb = h.GetEmissionKneserNey(state, tokens[0])

		// Handle complete Out-of-Vocabulary tokens cleanly
		if emissionProb == 0 && !h.Tokens[tokens[0]] {
			emissionProb = 1.0 / float64(len(h.Tokens))
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
			emissionProb = h.GetEmissionKneserNey(currState, token)

			// Handle Out-of-Vocabulary tokens at current step
			if emissionProb == 0 && !h.Tokens[token] {
				emissionProb = 1.0 / float64(len(h.Tokens))
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
