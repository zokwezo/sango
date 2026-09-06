package hmm

import (
	"encoding/json"
	"fmt"
	"os"
)

func MainEvaluate(actual, expect string) error {
	inFilenames := [2]string{actual, expect}
	var sentences [2][][]SangoSyllable // {actual, expect}
	for k, inFilename := range inFilenames {
		data, err := os.ReadFile(inFilename)
		if err != nil {
			return err
		}
		err = json.Unmarshal(data, &sentences[k])
		if err != nil {
			return err
		}
	}

	metrics, err := Evaluate(sentences)
	if err != nil {
		return err
	}

	// 5. Output metrics
	fmt.Println("")
	fmt.Println("PITCH | PRECISION |  RECALL   | F1-SCORE   ")
	fmt.Println("------+-----------+-----------+------------")
	for t, m := range metrics {
		tag := Tag(t)
		if m.TP+m.FP > 0 {
			m.Precision = float64(m.TP) / float64(m.TP+m.FP)
		}
		if m.TP+m.FN > 0 {
			m.Recall = float64(m.TP) / float64(m.TP+m.FN)
		}
		if m.Precision+m.Recall > 0 {
			m.F1Score = 2 * (m.Precision * m.Recall) / (m.Precision + m.Recall)
		}
		fmt.Printf("  %s   | %6.2f %%  | %6.2f %%  | %6.2f %%\n",
			tag, m.Precision*100, m.Recall*100, m.F1Score*100)
	}
	return nil
}

// Evaluate runs predictions on test data and prints Precision, Recall, and F1 per tag
func Evaluate(sentences [2][][]SangoSyllable) ([NumTags]Metrics, error) {
	metrics := [NumTags]Metrics{}
	n := len(sentences[0])
	if len(sentences[1]) != n {
		return metrics, fmt.Errorf("bad number of sentences passed to Evaluate")
	}
	for i := range n {
		m := len(sentences[0][i])
		if len(sentences[1][i]) != m {
			return metrics, fmt.Errorf("bad number of syllables in sentence %v passed to Evaluate", i)
		}
		for j := range m {
			actual := sentences[0][i][j].TagIndex
			expect := sentences[1][i][j].TagIndex
			if actual == expect {
				metrics[expect].TP++
			} else {
				metrics[actual].FP++
				metrics[expect].FN++
			}
		}
	}
	return metrics, nil
}
