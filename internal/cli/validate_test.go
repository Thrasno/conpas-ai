package cli

import (
	"testing"

	"github.com/Thrasno/conpas-ai/internal/model"
)

// TestNormalizePersonaBackwardCompat verifies that "gentleman" input is silently
// mapped to "argentino" so existing scripts and config files continue to work.
func TestNormalizePersonaBackwardCompat(t *testing.T) {
	got, err := normalizePersona("gentleman")
	if err != nil {
		t.Fatalf("normalizePersona(%q) unexpected error = %v", "gentleman", err)
	}
	if got != model.PersonaArgentino {
		t.Errorf("normalizePersona(%q) = %q, want %q", "gentleman", got, model.PersonaArgentino)
	}
}

// TestNormalizePersonaAllVariants verifies that all 7 standard personas and custom pass through.
func TestNormalizePersonaAllVariants(t *testing.T) {
	cases := []struct {
		input string
		want  model.PersonaID
	}{
		{"argentino", model.PersonaArgentino},
		{"neutral", model.PersonaNeutral},
		{"galleguinho", model.PersonaGalleguinho},
		{"asturianu", model.PersonaAsturianu},
		{"sargentoDeHierro", model.PersonaSargentoDeHierro},
		{"stark", model.PersonaStark},
		{"littleYoda", model.PersonaLittleYoda},
		{"custom", model.PersonaCustom},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			got, err := normalizePersona(tc.input)
			if err != nil {
				t.Fatalf("normalizePersona(%q) unexpected error = %v", tc.input, err)
			}
			if got != tc.want {
				t.Errorf("normalizePersona(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

// TestNormalizePersonaEmptyDefaultsToArgentino verifies that empty string defaults to argentino.
func TestNormalizePersonaEmptyDefaultsToArgentino(t *testing.T) {
	got, err := normalizePersona("")
	if err != nil {
		t.Fatalf("normalizePersona(%q) unexpected error = %v", "", err)
	}
	if got != model.PersonaArgentino {
		t.Errorf("normalizePersona(%q) = %q, want %q", "", got, model.PersonaArgentino)
	}
}

// TestNormalizePersonaUnknownReturnsError verifies that invalid personas return an error.
func TestNormalizePersonaUnknownReturnsError(t *testing.T) {
	_, err := normalizePersona("pirate")
	if err == nil {
		t.Fatalf("normalizePersona(%q) expected error, got nil", "pirate")
	}
}

// TestNormalizePresetBackwardCompat verifies that "full-gentleman" maps to "full".
func TestNormalizePresetBackwardCompat(t *testing.T) {
	got, err := normalizePreset("full-gentleman")
	if err != nil {
		t.Fatalf("normalizePreset(%q) unexpected error = %v", "full-gentleman", err)
	}
	if got != model.PresetFull {
		t.Errorf("normalizePreset(%q) = %q, want %q", "full-gentleman", got, model.PresetFull)
	}
}

// TestNormalizePresetFull verifies that "full" works directly.
func TestNormalizePresetFull(t *testing.T) {
	got, err := normalizePreset("full")
	if err != nil {
		t.Fatalf("normalizePreset(%q) unexpected error = %v", "full", err)
	}
	if got != model.PresetFull {
		t.Errorf("normalizePreset(%q) = %q, want %q", "full", got, model.PresetFull)
	}
}
