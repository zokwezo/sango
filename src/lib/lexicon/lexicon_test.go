package lexicon

import "testing"

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
	if nc := len(cols.Sango); nc != nr {
		t.Error(name, ": Sango cols has ", nc, " entries, but rows has ", nr, " entries")
	}
	if nc := len(cols.SangoUTF8); nc != nr {
		t.Error(name, ": SangoUTF8 cols has ", nc, " entries, but rows has ", nr, " entries")
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
			if len(row.Sango) == 0 {
				// Skip metadata header row
				continue
			}
			t.Error(name, ": Toneless row is empty")
		}
		if s := string(cols.Toneless[k]); s != row.Toneless {
			t.Error(name, ": Toneless[", k, "] col (", s, ") != row (", row.Toneless, ")")
		}

		if len(row.Sango) == 0 {
			t.Error(name, ": Sango row is empty")
		}
		if s := string(cols.Sango[k]); s != row.Sango {
			t.Error(name, ": Sango[", k, "] col (", s, ") != row (", row.Sango, ")")
		}
		if s := string(cols.SangoUTF8[k]); s != row.Sango {
			t.Error(name, ": SangoUTF8[", k, "] col (", s, ") != row (", row.Sango, ")")
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
