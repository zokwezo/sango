package hmm

const UnknownToken string = "<UNKNOWN>"
const UnknownTag string = "?"

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

func UnknownTagForToken(token string) string {
	// TODO: Replace with "???...", one ? for each vowel.
	return "?"
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
