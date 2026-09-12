// Uses Hidden Markov Model (HMM) with Kneser-Ney smoothing for diacritic restoration of Sango text.

package hmm

import (
	"cmp"
	"slices"

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
// same word with diacritics. The index is one past the source code to distinguish it
// from the default value of zero that implies a missing value.
type TrainingInstance struct {
	Index []int // index + 1, where index points back to input SSE to backreplace after prediction
	Token []string
	State []string
}

// These three sentence final code points create a new training instance.
const (
	newline rune = 0x0A
	dot     rune = 0x2E
	huh     rune = 0x3F
	bang    rune = 0x21
)

func PrepareInputText(s string) (sse.SSEs, []TrainingInstance) {
	tis := []TrainingInstance{
		TrainingInstance{
			Index: []int{},
			Token: []string{},
			State: []string{},
		},
	}

	codes, err := sse.Utf8ToSSEs(norm.NFC.String(s), sse.FromLemma)
	if err != nil {
		panic(err)
	}
	ti := &tis[0]
	numConsecutiveNewlines := 0
	flush := false
	for index, code := range codes {
		switch code.IsSango() {
		case false:
			// Skip over any Unicode symbol, except that if it is sentence final,
			// then create a new sentence.
			for c := uint64(code); c != 0; c >>= 16 {
				r := rune(c & 0xFFFF)
				if r == '\x00' {
					continue
				}
				switch r {
				case dot:
					fallthrough
				case huh:
					fallthrough
				case bang:
					flush = true
					break
				case newline:
					numConsecutiveNewlines++
					if numConsecutiveNewlines > 1 {
						numConsecutiveNewlines = 0
						flush = true
						break
					}
				default:
					if numConsecutiveNewlines != 0 {
						numConsecutiveNewlines = 0
					}
				}
			}
		case true:
			if flush {
				k := len(tis)
				tis = append(tis, TrainingInstance{
					Index: []int{},
					Token: []string{},
					State: []string{},
				})
				ti = &tis[k]
				flush = false
			}
			code &= 0x0FFF_FFFF_FFFF_FFFF // clear 4 MSB
			code |= 0x9000_0000_0000_0000 // force Sango, no-space, lowercase
			// Remember the place of this code so that after diacritic restoration
			// we know which token to update using the predicted state.
			state := sse.BuilderToString(code.WriteAsLemmaTo)
			token := sse.BuilderToString(code.WriteAsTonelessTo)
			if a, b, c := len(ti.Index), len(ti.Token), len(ti.State); a != b || a != c {
				panic("unbalanced TrainingInstance slices")
			}
			ti.Index = append(ti.Index, index+1) // one past the source code index
			ti.Token = append(ti.Token, token)
			ti.State = append(ti.State, state)
		}
	}
	return codes, tis
}
