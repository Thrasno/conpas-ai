package model

import "testing"

// TestPersonaConstants verifies all 7 PersonaID constants exist with the correct values.
func TestPersonaConstants(t *testing.T) {
	wantConstants := map[PersonaID]string{
		PersonaGentleman:        "gentleman",
		PersonaNeutral:          "neutral",
		PersonaCustom:           "custom",
		PersonaArgentino:        "argentino",
		PersonaGalleguinho:      "galleguinho",
		PersonaAsturianu:        "asturianu",
		PersonaSargentoDeHierro: "sargentoDeHierro",
		PersonaStark:            "stark",
	}

	// 7 named constants + PersonaCustom = 8 total; ensure all expected values are correct.
	for constant, wantValue := range wantConstants {
		if string(constant) != wantValue {
			t.Errorf("PersonaID constant %q has value %q, want %q", constant, string(constant), wantValue)
		}
	}

	// Verify uniqueness — no two constants may share the same string value.
	seen := map[string]bool{}
	for constant := range wantConstants {
		val := string(constant)
		if seen[val] {
			t.Errorf("duplicate PersonaID value %q", val)
		}
		seen[val] = true
	}

	// Verify PersonaGentleman is unchanged (backward compat).
	if PersonaGentleman != "gentleman" {
		t.Errorf("PersonaGentleman = %q, want %q", PersonaGentleman, "gentleman")
	}
}
