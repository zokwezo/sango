package hmm

import (
	"cmp"
	"fmt"
	"io/ioutil"
	"math"
	"slices"
	"strings"

	"github.com/zokwezo/sango/src/lib/sse"
	"google.golang.org/protobuf/proto"
)

// Metric holds the classification performance statistics
type Metric struct {
	Precision float64 // fraction of predicted that are correct
	Recall    float64 // fraction of correct that are predicted
	F1Score   float64 // harmonic mean of Precision and Recall
	TP        int     // # true  positives
	FP        int     // # false positives
	FN        int     // # false negatives
}
type MetricMap map[sse.SSE]*Metric

func (mm MetricMap) forState(state sse.SSE) *Metric {
	if mm[state] == nil {
		mm[state] = &Metric{}
	}
	return mm[state]
}

type CodeLemmaMetricPair = struct {
	code   sse.SSE
	lemma  string
	metric Metric
}

func MainEvaluate(actualFilename, expectFilename string) error {
	bothSSECodes := [2]SSECodes{} // {actual, expect}
	{
		inFilenames := [2]string{actualFilename, expectFilename}
		for k, inFilename := range inFilenames {
			data, err := ioutil.ReadFile(inFilename)
			if err != nil {
				return err
			}
			if err := proto.Unmarshal(data, &bothSSECodes[k]); err != nil {
				return err
			}
		}
	}

	// Verify that actual and expect have the same topology.
	a := &bothSSECodes[0].ShortCodes
	e := &bothSSECodes[1].ShortCodes
	nsa := len(*a)
	nse := len(*e)
	fmt.Printf("Corpus has %v actual and %v expected codes\n", nsa, nse)
	if nsa != nse {
		return fmt.Errorf("Corpus has %v actual but %v expected codes\n", nsa, nse)
	}

	metricMap, err := Evaluate(&bothSSECodes)
	if err != nil {
		return err
	}

	codeLemmaMetricPairs := make([]CodeLemmaMetricPair, 0, len(metricMap))
	maxLemmaLen := 7
	for code, metric := range metricMap {
		lemma := sse.BuilderToString(code.WriteAsLemmaTo)
		codeLemmaMetricPairs = append(codeLemmaMetricPairs, CodeLemmaMetricPair{code: code, lemma: lemma, metric: *metric})
		maxLemmaLen = max(maxLemmaLen, len(lemma))
	}
	slices.SortStableFunc(codeLemmaMetricPairs, func(lhs, rhs CodeLemmaMetricPair) int {
		if c := lhs.code.Compare(rhs.code); c != 0 {
			return c
		}
		return cmp.Compare(lhs.lemma, rhs.lemma)
	})

	maxRankLen := max(2, int(math.Ceil(math.Log10(float64(len(codeLemmaMetricPairs)+1)))))
	fmt.Println("")
	fmt.Printf("# %s |       CODE       | LEMMA %s | PRECISION |  RECALL   | F1-SCORE   \n",
		strings.Repeat(" ", maxRankLen-2), strings.Repeat(" ", maxLemmaLen-7))
	fmt.Printf("--%s-+------------------+-------%s-+-----------+-----------+------------\n",
		strings.Repeat("-", maxRankLen-2), strings.Repeat("-", maxLemmaLen-7))
	for k, e := range codeLemmaMetricPairs {
		if e.metric.TP+e.metric.FP > 0 {
			e.metric.Precision = float64(e.metric.TP) / float64(e.metric.TP+e.metric.FP)
		}
		if e.metric.TP+e.metric.FN > 0 {
			e.metric.Recall = float64(e.metric.TP) / float64(e.metric.TP+e.metric.FN)
		}
		if e.metric.Precision+e.metric.Recall > 0 {
			e.metric.F1Score = 2 * (e.metric.Precision * e.metric.Recall) / (e.metric.Precision + e.metric.Recall)
		}
		fmt.Printf("%-*v | %016X | %-*q | %6.2f %%  | %6.2f %%  | %6.2f %%\n",
			maxRankLen, k, uint64(e.code), maxLemmaLen, e.lemma, e.metric.Precision*100, e.metric.Recall*100, e.metric.F1Score*100)
	}

	return nil
}

// Evaluate runs predictions on test data and prints Precision, Recall, and F1 per state
func Evaluate(bothSSECodes *[2]SSECodes) (MetricMap, error) {
	const (
		noSpacePrefixAnd sse.SSE = 0x8FFF_FFFF_FFFF_FFFF
		noSpacePrefixOr  sse.SSE = 0x9000_0000_0000_0000
	)
	m := MetricMap{}
	a := &bothSSECodes[0].ShortCodes
	e := &bothSSECodes[1].ShortCodes
	n := len(*a)
	if len(*e) != n {
		return m, fmt.Errorf("bad number of bothSSECodes passed to Evaluate")
	}
	for k := range n {
		actual := sse.FromShortCode((*a)[k])
		expect := sse.FromShortCode((*e)[k])
		if actual.IsSango() && expect.IsSango() {
			aa := uint64(actual)
			ee := uint64(expect)
			for x := aa; x != 0; x >>= 12 {
				if x&0x0_000_000_000_000_FFF != 0 &&
					x&0x0_000_000_000_000_003 == 0 {
					fmt.Printf("actual[%v] = %016X -> %016X\n", k, (*a)[k], aa)
					panic("Unknown pitch")
				}
			}
			for x := ee; x != 0; x >>= 12 {
				if x&0x0_000_000_000_000_FFF != 0 &&
					x&0x0_000_000_000_000_003 == 0 {
					fmt.Printf("expect[%v] = %016X -> %016X\n", k, (*e)[k], ee)
					panic("Unknown pitch")
				}
			}
			actual &= noSpacePrefixAnd // set to no space prefix
			expect &= noSpacePrefixAnd // set to no space prefix
			actual |= noSpacePrefixOr  // set to lower case
			expect |= noSpacePrefixOr  // set to lower case
			if actual == expect {
				m.forState(expect).TP++
			} else {
				m.forState(actual).FP++
				m.forState(expect).FN++
			}
		}
	}
	return m, nil
}
