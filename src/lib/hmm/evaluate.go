package hmm

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

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

func (mm MetricsMap) forState(state string) *Metrics {
	if mm[state] == nil {
		mm[state] = &Metrics{}
	}
	return mm[state]
}

func MainEvaluate(actual, expect string) error {
	inFilenames := [2]string{actual, expect}
	var sentences [2][][]string // {actual, expect}
	for k, inFilename := range inFilenames {
		inputText, err := os.ReadFile(inFilename)
		if err != nil {
			return err
		}
		for _, sentence := range strings.Split(strings.Trim(string(inputText), " "), "\n") {
			sentence = strings.Trim(sentence, " ")
			if sentence != "" {
				sentences[k] = append(sentences[k], strings.Split(sentence, " "))
			}
		}
	}

	// Verify that actual and expect have the same topology.
	nsa := len(sentences[0])
	nse := len(sentences[1])
	if nsa != nse {
		return fmt.Errorf("actual has %v sentences but expect has %v sentences", nsa, nse)
	}
	nWords := 0
	for k := range nsa {
		nwa := len(sentences[0][k])
		nwe := len(sentences[1][k])
		if nwa != nwe {
			return fmt.Errorf("actual[%v] has %v words but expect[%v] has %v words", k, nwa, k, nwe)
		}
		nWords += nwa
	}
	fmt.Printf("Corpus has %v sentences and %v words\n", nsa, nWords)

	metricsMap, err := Evaluate(sentences)
	if err != nil {
		return err
	}

	sortedStates := make([]string, 0, len(metricsMap))
	for k := range metricsMap {
		sortedStates = append(sortedStates, k)
	}
	slices.Sort(sortedStates)

	fmt.Println("")
	fmt.Println("STATE | PRECISION |  RECALL   | F1-SCORE   ")
	fmt.Println("------+-----------+-----------+------------")
	for _, state := range sortedStates {
		m := metricsMap[state]
		if m.TP+m.FP > 0 {
			m.Precision = float64(m.TP) / float64(m.TP+m.FP)
		}
		if m.TP+m.FN > 0 {
			m.Recall = float64(m.TP) / float64(m.TP+m.FN)
		}
		if m.Precision+m.Recall > 0 {
			m.F1Score = 2 * (m.Precision * m.Recall) / (m.Precision + m.Recall)
		}
		if state == "" {
			state = "?"
		}
		fmt.Printf("  %s   | %6.2f %%  | %6.2f %%  | %6.2f %%\n",
			state, m.Precision*100, m.Recall*100, m.F1Score*100)
	}

	return nil
}

// Evaluate runs predictions on test data and prints Precision, Recall, and F1 per state
func Evaluate(sentences [2][][]string) (MetricsMap, error) {
	metricsMap := MetricsMap{}
	n := len(sentences[0])
	if len(sentences[1]) != n {
		return metricsMap, fmt.Errorf("bad number of sentences passed to Evaluate")
	}
	for i := range n {
		m := len(sentences[0][i])
		if len(sentences[1][i]) != m {
			return metricsMap, fmt.Errorf("bad number of tokens in sentence %v passed to Evaluate", i)
		}
		for j := range m {
			actual := sentences[0][i][j]
			expect := sentences[1][i][j]
			if actual == expect {
				metricsMap.forState(expect).TP++
			} else {
				metricsMap.forState(actual).FP++
				metricsMap.forState(expect).FN++
			}
		}
	}
	return metricsMap, nil
}
