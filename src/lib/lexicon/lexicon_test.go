package lexicon

import "testing"

func TestConsistencyBetweenRowAndColMajorOrder(t *testing.T) {
	type RC struct {
		name string
		rows DictRows
		cols DictCols
	}
	var rcs = []RC{
		{"Affixes", AffixesRows(), AffixesCols()},
		{"Lexicon", LexiconRows(), LexiconCols()},
	}

	for _, rc := range rcs {
		// Check that all array lengths match.
		nr := len(rc.rows)
		if nr == 0 {
			t.Error(rc.name, ": Rows is empty")
		}
		if nc := len(rc.cols.Toneless); nc != nr {
			t.Error(rc.name, ": Toneless cols has ", nc, " entries, but rows has ", nr, " entries")
		}
		if nc := len(rc.cols.Sango); nc != nr {
			t.Error(rc.name, ": Sango cols has ", nc, " entries, but rows has ", nr, " entries")
		}
		if nc := len(rc.cols.SangoUTF8); nc != nr {
			t.Error(rc.name, ": SangoUTF8 cols has ", nc, " entries, but rows has ", nr, " entries")
		}
		if nc := len(rc.cols.LexPos); nc != nr {
			t.Error(rc.name, ": LexPos cols has ", nc, " entries, but rows has ", nr, " entries")
		}
		if nc := len(rc.cols.UDPos); nc != nr {
			t.Error(rc.name, ": UDPos cols has ", nc, " entries, but rows has ", nr, " entries")
		}
		if nc := len(rc.cols.UDFeature); nc != nr {
			t.Error(rc.name, ": UDFeature cols has ", nc, " entries, but rows has ", nr, " entries")
		}
		if nc := len(rc.cols.Category); nc != nr {
			t.Error(rc.name, ": Category cols has ", nc, " entries, but rows has ", nr, " entries")
		}
		if nc := len(rc.cols.Frequency); nc != nr {
			t.Error(rc.name, ": Frequency cols has ", nc, " entries, but rows has ", nr, " entries")
		}
		if nc := len(rc.cols.English); nc != nr {
			t.Error(rc.name, ": English cols has ", nc, " entries, but rows has ", nr, " entries")
		}

		// Check that backing arrays are not empty.
		if len(rc.cols.Bytes) == 0 {
			t.Error(rc.name, ": Bytes col is empty")
		}
		if len(rc.cols.Runes) == 0 {
			t.Error(rc.name, ": Runes col is empty")
		}

		for k, row := range rc.rows {
			if len(row.Toneless) == 0 {
				if len(row.Sango) == 0 {
					// Skip metadata header row
					continue
				}
				t.Error(rc.name, ": Toneless row is empty")
			}
			if s := string(rc.cols.Toneless[k]); s != row.Toneless {
				t.Error(rc.name, ": Toneless[", k, "] col (", s, ") != row (", row.Toneless, ")")
			}

			if len(row.Sango) == 0 {
				t.Error(rc.name, ": Sango row is empty")
			}
			if s := string(rc.cols.Sango[k]); s != row.Sango {
				t.Error(rc.name, ": Sango[", k, "] col (", s, ") != row (", row.Sango, ")")
			}
			if s := string(rc.cols.SangoUTF8[k]); s != row.Sango {
				t.Error(rc.name, ": SangoUTF8[", k, "] col (", s, ") != row (", row.Sango, ")")
			}

			if len(row.LexPos) == 0 {
				t.Error(rc.name, ": LexPos row is empty")
			}
			if s := string(rc.cols.LexPos[k]); s != row.LexPos {
				t.Error(rc.name, ": LexPos[", k, "] col (", s, ") != row (", row.LexPos, ")")
			}

			if len(row.UDPos) == 0 {
				t.Error(rc.name, ": UDPos row is empty")
			}
			if s := string(rc.cols.UDPos[k]); s != row.UDPos {
				t.Error(rc.name, ": UDPos[", k, "] col (", s, ") != row (", row.UDPos, ")")
			}

			if s := string(rc.cols.UDFeature[k]); s != row.UDFeature {
				t.Error(rc.name, ": UDFeature[", k, "] col (", s, ") != row (", row.UDFeature, ")")
			}

			if s := string(rc.cols.Category[k]); s != row.Category {
				t.Error(rc.name, ": Category[", k, "] col (", s, ") != row (", row.Category, ")")
			}

			if row.Frequency < 1 || row.Frequency > 9 {
				t.Error(rc.name, ": Invalid Frequency[", k, "] (", row.Frequency, ")")
			}
			if s := rc.cols.Frequency[k]; s != row.Frequency {
				t.Error(rc.name, ": Frequency[", k, "] col (", s, ") != row (", row.Frequency, ")")
			}

			if len(row.English) == 0 {
				t.Error(rc.name, ": English row is empty")
			}
			if s := string(rc.cols.English[k]); s != row.English {
				t.Error(rc.name, ": English[", k, "] col (", s, ") != row (", row.English, ")")
			}

			// Ensure that the rows are in strictly ascending order (ignoring Frequency and English fields).
			if k == 0 {
				continue
			}
			prevRow := &rc.rows[k-1]
			currRow := &rc.rows[k]
			if prevRow.Toneless < currRow.Toneless {
				continue
			}
			if prevRow.Toneless > currRow.Toneless {
				t.Error(rc.name, "[", k, "]: Bad Toneless order, prev={", prevRow, "}, curr={", currRow, "}")
			}
			if prevRow.Sango < currRow.Sango {
				continue
			}
			if prevRow.Sango > currRow.Sango {
				t.Error(rc.name, "[", k, "]: Bad Sango order, prev={", prevRow, "}, curr={", currRow, "}")
			}
			if prevRow.LexPos < currRow.LexPos {
				continue
			}
			if prevRow.LexPos > currRow.LexPos {
				t.Error(rc.name, "[", k, "]: Bad LexPos order, prev={", prevRow, "}, curr={", currRow, "}")
			}
			if prevRow.UDPos < currRow.UDPos {
				continue
			}
			if prevRow.UDPos > currRow.UDPos {
				t.Error(rc.name, "[", k, "]: Bad UDPos order, prev={", prevRow, "}, curr={", currRow, "}")
			}
			if prevRow.UDFeature < currRow.UDFeature {
				continue
			}
			if prevRow.UDFeature > currRow.UDFeature {
				t.Error(rc.name, "[", k, "]: Bad UDFeature order, prev={", prevRow, "}, curr={", currRow, "}")
			}
			if prevRow.Category < currRow.Category {
				continue
			}
			t.Error(rc.name, "[", k, "]: Bad Category order, prev={", prevRow, "}, curr={", currRow, "}")
		}
	}
}
