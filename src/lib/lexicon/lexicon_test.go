package lexicon

import (
	"fmt"
	"os"
	"testing"

	cuckoo "github.com/panmari/cuckoofilter"
)

func TestGenerateEncodedCuckooFilterOfLemma(t *testing.T) {
	cf := cuckoo.NewFilter(uint(len(LexiconCols().LemmaUTF8)))
	keys := map[string]struct{}{}
	for _, sangoUTF8 := range LexiconCols().LemmaUTF8 {
		key := string(sangoUTF8)
		if _, found := keys[key]; !found {
			keys[key] = struct{}{}
			cf.Insert(sangoUTF8)
		}
	}
	err := os.WriteFile("/tmp/wordlist_sg.cf", cf.Encode(), 0664)
	if err != nil {
		t.Error(err)
	}
}

func TestGenerateEncodedCuckooFilterOfToneless(t *testing.T) {
	cf := cuckoo.NewFilter(uint(len(LexiconCols().Toneless)))
	keys := map[string]struct{}{}
	for _, toneless := range LexiconCols().Toneless {
		key := string(toneless)
		if _, found := keys[key]; !found {
			keys[key] = struct{}{}
			cf.Insert(toneless)
		}
	}
	err := os.WriteFile("/tmp/wordlist_sg_toneless.cf", cf.Encode(), 0664)
	if err != nil {
		t.Error(err)
	}
}

func TestGenerateEncodedCuckooFilterOfHeightless(t *testing.T) {
	cf := cuckoo.NewFilter(uint(len(LexiconCols().Heightless)))
	keys := map[string]struct{}{}
	for _, heightless := range LexiconCols().Heightless {
		key := string(heightless)
		if _, found := keys[key]; !found {
			keys[key] = struct{}{}
			cf.Insert(heightless)
		}
	}
	err := os.WriteFile("/tmp/wordlist_sg_heightless.cf", cf.Encode(), 0664)
	if err != nil {
		t.Error(err)
	}
}

func TestConsistencyBetweenRowAndColMajorOrder(t *testing.T) {
	name := "Lexicon"
	rows := LexiconRows()
	cols := LexiconCols()

	// Check that all array lengths match.
	nr := len(rows)
	if nr == 0 {
		t.Error(name, ": Rows is empty")
	}
	if nc := len(cols.Toneless); nc != nr {
		t.Error(name, ": Toneless cols has ", nc, " entries, but rows has ", nr, " entries")
	}
	if nc := len(cols.LemmaRunes); nc != nr {
		t.Error(name, ": Lemma cols has ", nc, " entries, but rows has ", nr, " entries")
	}
	if nc := len(cols.LemmaUTF8); nc != nr {
		t.Error(name, ": LemmaUTF8 cols has ", nc, " entries, but rows has ", nr, " entries")
	}
	if nc := len(cols.UDPos); nc != nr {
		t.Error(name, ": UDPos cols has ", nc, " entries, but rows has ", nr, " entries")
	}
	if nc := len(cols.UDFeature); nc != nr {
		t.Error(name, ": UDFeature cols has ", nc, " entries, but rows has ", nr, " entries")
	}
	if nc := len(cols.Category); nc != nr {
		t.Error(name, ": Category cols has ", nc, " entries, but rows has ", nr, " entries")
	}
	if nc := len(cols.Frequency); nc != nr {
		t.Error(name, ": Frequency cols has ", nc, " entries, but rows has ", nr, " entries")
	}
	if nc := len(cols.EnglishTranslation); nc != nr {
		t.Error(name, ": EnglishTranslation cols has ", nc, " entries, but rows has ", nr, " entries")
	}
	if nc := len(cols.EnglishDefinition); nc != nr {
		t.Error(name, ": EnglishDefinition cols has ", nc, " entries, but rows has ", nr, " entries")
	}

	// Check that backing arrays are not empty.
	if len(cols.Bytes) == 0 {
		t.Error(name, ": Bytes col is empty")
	}
	if len(cols.Runes) == 0 {
		t.Error(name, ": Runes col is empty")
	}

	for k, row := range rows {
		if len(row.Toneless) == 0 {
			if len(row.Lemma) == 0 {
				// Skip metadata header row
				continue
			}
			t.Error(name, ": Toneless row is empty")
		}
		if s := string(cols.Toneless[k]); s != row.Toneless {
			t.Error(name, ": Toneless[", k, "] col (", s, ") != row (", row.Toneless, ")")
		}

		if len(row.Lemma) == 0 {
			t.Error(name, ": Lemma row is empty")
		}
		if s := string(cols.LemmaRunes[k]); s != row.Lemma {
			t.Error(name, ": Lemma[", k, "] col (", s, ") != row (", row.Lemma, ")")
		}
		if s := string(cols.LemmaUTF8[k]); s != row.Lemma {
			t.Error(name, ": LemmaUTF8[", k, "] col (", s, ") != row (", row.Lemma, ")")
		}

		if len(row.UDPos) == 0 {
			t.Error(name, ": UDPos row is empty")
		}
		if s := string(cols.UDPos[k]); s != row.UDPos {
			t.Error(name, ": UDPos[", k, "] col (", s, ") != row (", row.UDPos, ")")
		}

		if s := string(cols.UDFeature[k]); s != row.UDFeature {
			t.Error(name, ": UDFeature[", k, "] col (", s, ") != row (", row.UDFeature, ")")
		}

		if s := string(cols.Category[k]); s != row.Category {
			t.Error(name, ": Category[", k, "] col (", s, ") != row (", row.Category, ")")
		}

		if row.Frequency < 1 || row.Frequency > 9 {
			t.Error(name, ": Invalid Frequency[", k, "] (", row.Frequency, ")")
		}
		if s := cols.Frequency[k]; s != row.Frequency {
			t.Error(name, ": Frequency[", k, "] col (", s, ") != row (", row.Frequency, ")")
		}

		if len(row.EnglishTranslation) == 0 {
			t.Error(name, ": EnglishTranslation row is empty")
		}
		if s := string(cols.EnglishTranslation[k]); s != row.EnglishTranslation {
			t.Error(name, ": EnglishTranslation[", k, "] col (", s, ") != row (", row.EnglishTranslation, ")")
		}

		if len(row.EnglishDefinition) == 0 {
			t.Error(name, ": EnglishDefinition row is empty")
		}
		if s := string(cols.EnglishDefinition[k]); s != row.EnglishDefinition {
			t.Error(name, ": EnglishDefinition[", k, "] col (", s, ") != row (", row.EnglishDefinition, ")")
		}
	}
}

func TestRowsMatchingLemma(t *testing.T) {
	actual := RowsMatchingLemma()["dɛ̈"]
	expect := DictRows{
		{"de", "dë", "dɛ̈", "VERB", "Subcat=Tran", "ACT", 2, "cut-or-grow", "cut, slice; grow, cultivate"},
		{"de", "dë", "dɛ̈", "VERB", "Subcat=Tran", "INTERACT", 3, "emit", "emit"},
	}
	actualStr := fmt.Sprintf("%v", actual)
	expectStr := fmt.Sprintf("%v", expect)
	if actualStr != expectStr {
		t.Error("actual: " + actualStr)
		t.Error("expect: " + expectStr)
	}
}

func TestRowsMatchingHeightless(t *testing.T) {
	actual := RowsMatchingHeightless()["dë"]
	expect := DictRows{
		{"de", "dë", "dɛ̈", "VERB", "Subcat=Tran", "ACT", 2, "cut-or-grow", "cut, slice; grow, cultivate"},
		{"de", "dë", "dɛ̈", "VERB", "Subcat=Tran", "INTERACT", 3, "emit", "emit"},
		{"de", "dë", "dë", "VERB", "Subcat=Intr", "HOW", 3, "be-cold", "be cold"},
	}
	actualStr := fmt.Sprintf("%v", actual)
	expectStr := fmt.Sprintf("%v", expect)
	if actualStr != expectStr {
		t.Error("actual: " + actualStr)
		t.Error("expect: " + expectStr)
	}
}

func TestRowsMatchingToneless(t *testing.T) {
	actual := RowsMatchingToneless()["de"]
	expect := DictRows{
		{"de", "de", "de", "VERB", "Subcat=Intr", "BODY", 3, "vomit", "vomit"},
		{"de", "de", "dɛ", "VERB", "", "STATE", 2, "remain", "remain"},
		{"de", "dê", "dê", "NOUN", "", "HOW", 3, "coldness", "coldness, shade"},
		{"de", "dë", "dɛ̈", "VERB", "Subcat=Tran", "ACT", 2, "cut-or-grow", "cut, slice; grow, cultivate"},
		{"de", "dë", "dɛ̈", "VERB", "Subcat=Tran", "INTERACT", 3, "emit", "emit"},
		{"de", "dë", "dë", "VERB", "Subcat=Intr", "HOW", 3, "be-cold", "be cold"},
	}
	actualStr := fmt.Sprintf("%v", actual)
	expectStr := fmt.Sprintf("%v", expect)
	if actualStr != expectStr {
		t.Error("actual: " + actualStr)
		t.Error("expect: " + expectStr)
	}
}
