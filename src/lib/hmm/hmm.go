// Uses Hidden Markov Model (HMM) with Kneser-Ney smoothing for diacritic restoration of Sango text.

package hmm

import (
	"cmp"
	"regexp"
	"slices"
	"strings"

	"github.com/zokwezo/sango/src/lib/sse"
	"golang.org/x/text/unicode/norm"
)

// TODO: Switch from string to sse.SSE for robustness and ease of coding.
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
	Index []int // points back to input SSE to backreplace after prediction
	Token []string
	State []string
}

var (
	reFinalPunctuation      = regexp.MustCompile(`[.!?]`)
	reNonTextWithDiacritics = regexp.MustCompile(`[^a-z .?!\x{302}\x{308}\x{323}]`)
	reNonText               = regexp.MustCompile(`[^a-z .?!]`)
	reCompressSpaces        = regexp.MustCompile(`[ ]{2,}`)
)

// These three sentence final code points create a new training instance.
const (
	dot  rune = 0x2E
	huh  rune = 0x3F
	bang rune = 0x21
)

func PrepareInputText(s string) []TrainingInstance {
	tis := []TrainingInstance{
		TrainingInstance{
			Index: []int{},
			Token: []string{},
			State: []string{},
		},
	}

	// TODO: Return sses from this function and in Predict use tis[k].State[j] to
	// update the diacritics of sses[tis[k].Index[j]].
	sses, err := sse.Utf8ToSSEs(norm.NFC.String(s), sse.FromLemma)
	if err != nil {
		panic(err)
	}
	ti := &tis[0]
	for index, code := range sses {
		switch code.IsSango() {
		case false:
			// Skip over any Unicode symbol, except that if it is sentence final,
			// then create a new sentence.
			flush := false
			c := uint64(code)
			for _ = range 4 {
				r := rune(c & 0xFFFF)
				switch r {
				case dot:
					fallthrough
				case huh:
					fallthrough
				case bang:
					flush = true
					break
				}
				c >>= 16
			}
			if flush {
				k := len(tis)
				tis = append(tis, TrainingInstance{
					Index: []int{},
					Token: []string{},
					State: []string{},
				})
				ti = &tis[k]
			}
		case true:
			// Remember the place of this code so that after diacritic restoration
			// we know which token to update using the predicted state.
			state := strings.ToLower(strings.Trim(sse.BuilderToString(code.WriteAsLemmaTo), " "))
			token := strings.ToLower(strings.Trim(sse.BuilderToString(code.WriteAsTonelessTo), " "))
			if a, b, c := len(ti.Index), len(ti.Token), len(ti.State); a != b || a != c {
				panic("unbalanced TrainingInstance slices")
			}
			ti.Index = append(ti.Index, index)
			ti.Token = append(ti.Token, token)
			ti.State = append(ti.State, state)
		}
	}
	return tis
}
