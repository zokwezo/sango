package hmm

const (
	UnknownToken string = ""
	UnknownTag   string = ""
	// Transition dampening better handle unknown tokens at the cost of tag sequence matching.
	// Lower values deweight the effect of tag 1-grams and 2-grams relative to (token,tag) associations.
	// NOTE: The hyperparameters below are brittle for small training corpora.
	// TODO: Replace Laplace Smoothing with Kneser-Ney Smoothing (https://share.google/aimode/yxkByUFOWe6YpzZ4M)
	StartTagDampening     float64 = 0.01 // must be in [0,1]
	TransitionDampening   float64 = 0.01 // must be in [0,1]
	EmissionDampening     float64 = 1.0  // must be in [0,1]
	UnknownTokenDampening float64 = 0.1  // must be in (0,1]
)

type HmmMap map[string]float64
type HMMPerTag struct {
	StartTag   float64 //           = log P(tag)
	Transition HmmMap  // [nextTag] = log P(nextTag | tag)
	Emission   HmmMap  // [token]   = log P(token   | tag)
}

type HMM struct {
	HmmPerTag    map[string]*HMMPerTag // [tag]
	NumSentences int                   // < 0 during Accumulate, > 0 after Generate
}

func (h *HMM) ensureDefined() {
	if h.HmmPerTag == nil {
		h.HmmPerTag = map[string]*HMMPerTag{}
	}
	if h.HmmPerTag[UnknownTag] == nil {
		h.HmmPerTag[UnknownTag] = &HMMPerTag{Transition: HmmMap{}, Emission: HmmMap{}}
	}
}
func (h *HMM) ensureDefinedForTag(tag string) {
	h.ensureDefined()
	if h.HmmPerTag[tag] == nil {
		h.HmmPerTag[tag] = &HMMPerTag{Transition: HmmMap{}, Emission: HmmMap{}}
	}
}
func (h *HMM) forTag(tag string) *HMMPerTag {
	h.ensureDefinedForTag(tag)
	return h.HmmPerTag[tag]
}

type SangoToken struct {
	Token string
	Tag   string
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
type MetricsMap map[string]*Metrics

func (mm MetricsMap) forTag(tag string) *Metrics {
	if mm[tag] == nil {
		mm[tag] = &Metrics{}
	}
	return mm[tag]
}

var TagOrderMap = func() map[rune]int {
	orderMap := make(map[rune]int)
	for i, r := range "_:^" {
		orderMap[r] = i
	}
	return orderMap
}()

func TagCompare(lhs, rhs string) int {
	lhsRunes := []rune(lhs)
	rhsRunes := []rune(rhs)
	minLength := len(lhsRunes)
	if len(rhsRunes) < minLength {
		minLength = len(rhsRunes)
	}
	for i := 0; i < minLength; i++ {
		rankLhs := TagOrderMap[lhsRunes[i]]
		rankRhs := TagOrderMap[rhsRunes[i]]
		if rankLhs != rankRhs {
			return rankLhs - rankRhs
		}
	}
	return len(lhsRunes) - len(rhsRunes)
}
