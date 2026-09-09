// Uses Hidden Markov Model (HMM) with Kneser-Ney smoothing for diacritic restoration of Sango text.

package hmm

import (
	"cmp"
	"regexp"
	"slices"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/unicode/norm"
)

type HMM struct {
	States      map[string]bool
	Tokens      map[string]bool
	Transitions map[string]map[string]int64 // count(s_t | s_{t-1})
	Emissions   map[string]map[string]int64 // count(w_t | s_t)
	StateCounts map[string]int64            // total occurrences of each state
	UniquePairs int64                       // total unique (state, word) types in corpus
}

func FromToCountCompare(lhs, rhs *FromToCount) int {
	if c := cmp.Compare(lhs.From, rhs.From); c != 0 {
		return c
	}
	if c := cmp.Compare(lhs.To, rhs.To); c != 0 {
		return c
	}
	if c := cmp.Compare(lhs.Count, rhs.Count); c != 0 {
		return c
	}
	return 0
}

func (h *HMM) ToModel() Model {
	m := Model{
		States:      []string{},
		Tokens:      []string{},
		Transitions: []*FromToCount{},
		Emissions:   []*FromToCount{},
		StateCounts: []*FromToCount{},
	}
	for s := range h.States {
		m.States = append(m.States, s)
	}
	slices.Sort(m.States)
	mStates := make(map[string]int32, len(m.States))
	for k, s := range m.States {
		mStates[s] = int32(k)
	}
	for s := range h.Tokens {
		m.Tokens = append(m.Tokens, s)
	}
	slices.Sort(m.Tokens)
	mTokens := make(map[string]int32, len(m.Tokens))
	for k, s := range m.Tokens {
		mTokens[s] = int32(k)
	}
	for sFrom, p := range h.Transitions {
		if kFrom, found := mStates[sFrom]; found {
			for sTo, v := range p {
				if kTo, found := mStates[sTo]; found {
					m.Transitions = append(m.Transitions, &FromToCount{From: kFrom, To: kTo, Count: int64(v)})
				} else {
					panic("Transitions key missing from States")
				}
			}
		} else {
			panic("Transitions key missing from States")
		}
	}
	slices.SortFunc(m.Transitions, FromToCountCompare)
	for sFrom, p := range h.Emissions {
		if kFrom, found := mStates[sFrom]; found {
			for sTo, v := range p {
				if kTo, found := mTokens[sTo]; found {
					m.Emissions = append(m.Emissions, &FromToCount{From: kFrom, To: kTo, Count: int64(v)})
				} else {
					panic("Emissions key missing from Tokens")
				}
			}
		} else {
			panic("Emissions key missing from States")
		}
	}
	slices.SortFunc(m.Emissions, FromToCountCompare)
	for s, c := range h.StateCounts {
		if k, found := mStates[s]; found {
			m.StateCounts = append(m.StateCounts, &FromToCount{From: k, Count: int64(c)})
		} else {
			panic("StateCounts key missing from States")
		}
	}
	slices.SortFunc(m.StateCounts, FromToCountCompare)
	m.UniquePairs = int64(h.UniquePairs)
	return m
}

func (h *HMM) FromModel(m *Model) {
	if h == nil || m == nil {
		return
	}
	if m.States == nil {
		m.States = []string{}
	}
	if m.Tokens == nil {
		m.Tokens = []string{}
	}
	if m.Transitions == nil {
		m.Transitions = []*FromToCount{}
	}
	if m.Emissions == nil {
		m.Emissions = []*FromToCount{}
	}
	if m.StateCounts == nil {
		m.StateCounts = []*FromToCount{}
	}
	*h = HMM{
		States:      make(map[string]bool),
		Tokens:      make(map[string]bool),
		Transitions: make(map[string]map[string]int64),
		Emissions:   make(map[string]map[string]int64),
		StateCounts: make(map[string]int64),
	}

	numStates := int32(len(m.States))
	numTokens := int32(len(m.Tokens))
	for _, state := range m.States {
		h.States[state] = true
	}
	for _, word := range m.Tokens {
		h.Tokens[word] = true
	}
	for _, ftc := range m.Transitions {
		if ftc != nil {
			if kFrom, kTo := ftc.From, ftc.To; kFrom >= 0 && kTo >= 0 && kFrom < numStates && kTo < numStates {
				if h.Transitions[m.States[kFrom]] == nil {
					h.Transitions[m.States[kFrom]] = make(map[string]int64)
				}
				h.Transitions[m.States[kFrom]][m.States[kTo]] = ftc.Count
			}
		}
	}
	for _, ftc := range m.Emissions {
		if ftc != nil {
			if kFrom, kTo := ftc.From, ftc.To; kFrom >= 0 && kTo >= 0 && kFrom < numStates && kTo < numTokens {
				if h.Emissions[m.States[kFrom]] == nil {
					h.Emissions[m.States[kFrom]] = make(map[string]int64)
				}
				h.Emissions[m.States[kFrom]][m.Tokens[kTo]] = ftc.Count
			}
		}
	}
	for _, ftc := range m.StateCounts {
		if ftc != nil {
			if kFrom := ftc.From; kFrom >= 0 && kFrom < numStates {
				h.StateCounts[m.States[kFrom]] = ftc.Count
			}
		}
	}
	h.UniquePairs = m.UniquePairs
}

// TrainingInstance represents a sequence of tokens and their true state/label sequence.
// In this context, a Token is a Sango word without diacritics and a State is the
// same word with diacritics.
type TrainingInstance struct {
	Tokens []string
	States []string
}

var (
	reFinalPunctuation      = regexp.MustCompile(`[.!?]`)
	reNonTextWithDiacritics = regexp.MustCompile(`[^a-z .?!\x{302}\x{308}\x{323}]`)
	reNonText               = regexp.MustCompile(`[^a-z .?!]`)
	reCompressSpaces        = regexp.MustCompile(`[ ]{2,}`)
)

func PrepareInputText(s string) []TrainingInstance {
	s = norm.NFD.String(cases.Lower(language.English).String(s))
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.Trim(s, " ")
	s = reNonTextWithDiacritics.ReplaceAllLiteralString(s, " ")
	s = reCompressSpaces.ReplaceAllLiteralString(s, " ")
	s = reFinalPunctuation.ReplaceAllLiteralString(s, "\n")
	tis := []TrainingInstance{}
	for _, sentence := range strings.Split(strings.Trim(s, " "), "\n") {
		sentence = strings.Trim(sentence, " ")
		if sentence != "" {
			ti := TrainingInstance{}
			ti.States = strings.Split(sentence, " ")
			n := len(ti.States)
			ti.Tokens = make([]string, n)
			for k := range ti.States {
				ti.Tokens[k] = reNonText.ReplaceAllLiteralString(ti.States[k], "")
				ti.States[k] = norm.NFC.String(ti.States[k])
			}
			tis = append(tis, ti)
		}
	}
	return tis
}
