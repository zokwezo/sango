package hmm

import (
	"encoding/json"
	"fmt"
	"os"
)

func MainEvaluate(actual, expect string) error {
	inFilenames := [2]string{actual, expect}
	var sentences [2][][]SangoToken // {actual, expect}
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

	metricsMap, err := Evaluate(sentences)
	if err != nil {
		return err
	}

	fmt.Println("")
	fmt.Println("PITCH | PRECISION |  RECALL   | F1-SCORE   ")
	fmt.Println("------+-----------+-----------+------------")
	for tag, m := range metricsMap {
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
func Evaluate(sentences [2][][]SangoToken) (MetricsMap, error) {
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
			actual := sentences[0][i][j].Tag
			expect := sentences[1][i][j].Tag
			if actual == expect {
				metricsMap.forTag(expect).TP++
			} else {
				metricsMap.forTag(actual).FP++
				metricsMap.forTag(expect).FN++
			}
		}
	}
	return metricsMap, nil
}
