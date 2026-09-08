// Hidden Markov Model (HMM) with Kneser-Ney smoothing

// Used for diacritic restoration of Sango text.

package hmm

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
