package restore

import (
	"testing"
)

func TestAlreadyCorrectSangoVowels(t *testing.T) {
	original := "Na pekö tî sô lo tambûla ngbii piî na ndäpêrêrê asï na lâ-kûî, na löndöngɔ̈ tî lâ asï na sïgïngɔ̈ tî nzɛ; awɛ so, lo sï na bariëre sô azîâ tî kânga na yângâ tî kɔ̈dɔ̈rɔ̈ tî Ngiba sô. mbɛ̂nî turûgu tî bätängɔ̈ gbïä tî kɔ̈dɔ̈rɔ̈ aîri lo."
	expected := original
	actually := RestoreSangoVowels(original)
	if actually != expected {
		t.Errorf("ORIGINAL = %s\n", original)
		t.Errorf("ACTUALLY = %s\n", actually)
		t.Errorf("EXPECTED = %s\n", expected)
	}
}

func TestCorrectSangoVowelHeight(t *testing.T) {
	original := "Na peko tî sô lo tambûla ngbii piî na ndäpêrêrê asï na lâkûi, na löndöngö tî lâ asï na sigingo tî nze; awe so, lo si na bariëre sô azîa tî kânga na yângâ tî ködörö tî Ngiba sô. Mbênî turûgu tî bätängö gbïä ti ködörö aîri lo."
	expected := "Na pekö tî sô lo tambûla ngbii piî na ndäpêrêrê asï na lâ-kûî, na löndöngɔ̈ tî lâ asï na sïgïngɔ̈ tî nzɛ; awɛ so, lo sî|sï na bariëre sô azîâ tî kânga na yângâ tî kɔ̈dɔ̈rɔ̈ tî Ngiba sô. mbɛ̂nî turûgu tî bâ-tângo|bätängɔ̈ gbïä tî|tï kɔ̈dɔ̈rɔ̈ aîri lo."
	actually := RestoreSangoVowels(original)
	if actually != expected {
		t.Errorf("ORIGINAL = %s\n", original)
		t.Errorf("ACTUALLY = %s\n", actually)
		t.Errorf("EXPECTED = %s\n", expected)
	}
}

func TestRestoreSangoVowelHeightAndPitch(t *testing.T) {
	original := "Na peko ti so lo tambula ngbii pii na ndaperere asi na lakui, na londongo ti la asi na sigingo ti nze; awe so, lo si na bariere so azia ti kanga na yanga ti kodoro ti Ngiba so. Mbeni turugu ti batango gbia ti kodoro airi lo."
	expected := "Na pekö tî|tï so lo tambûla ngbii pii na ndäpêrêrê asî|asï na lâ-kûî, na löndöngɔ̈ tî|tï lâ asî|asï na sïgïngɔ̈ tî|tï nzɛ; awɛ so, lo sî|sï na bariere so azîâ tî|tï kânga|kângâ|känga na yângâ tî|tï kɔ̈dɔ̈rɔ̈ tî|tï Ngiba so. mbɛ̂nî turûgu tî|tï bâ-tângo|bätängɔ̈ gbïä tî|tï kɔ̈dɔ̈rɔ̈ aîri|âïrï lo."
	actually := RestoreSangoVowels(original)
	if actually != expected {
		t.Errorf("ORIGINAL = %s\n", original)
		t.Errorf("ACTUALLY = %s\n", actually)
		t.Errorf("EXPECTED = %s\n", expected)
	}
}
