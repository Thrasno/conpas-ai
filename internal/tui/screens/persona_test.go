package screens

import (
	"testing"

	"github.com/Thrasno/conpas-ai/internal/model"
)

// TestPersonaOptionsReturns8Items verifies that PersonaOptions returns all 8 personas.
func TestPersonaOptionsReturns8Items(t *testing.T) {
	options := PersonaOptions()
	if len(options) != 8 {
		t.Fatalf("PersonaOptions() returned %d items, want 8: %v", len(options), options)
	}

	wantPersonas := []model.PersonaID{
		model.PersonaArgentino,
		model.PersonaNeutral,
		model.PersonaGalleguinho,
		model.PersonaAsturianu,
		model.PersonaSargentoDeHierro,
		model.PersonaStark,
		model.PersonaLittleYoda,
		model.PersonaCustom,
	}

	for i, want := range wantPersonas {
		if options[i] != want {
			t.Errorf("PersonaOptions()[%d] = %q, want %q", i, options[i], want)
		}
	}
}

// TestPersonaLabelContainsDescription verifies that each persona has a non-empty
// descriptive label (not just the raw persona ID).
func TestPersonaLabelContainsDescription(t *testing.T) {
	personas := []model.PersonaID{
		model.PersonaArgentino,
		model.PersonaNeutral,
		model.PersonaGalleguinho,
		model.PersonaAsturianu,
		model.PersonaSargentoDeHierro,
		model.PersonaStark,
		model.PersonaLittleYoda,
		model.PersonaCustom,
	}

	for _, persona := range personas {
		persona := persona
		t.Run(string(persona), func(t *testing.T) {
			label := PersonaLabel(persona)
			if label == "" {
				t.Fatalf("PersonaLabel(%q) returned empty string", persona)
			}
			// Label should contain a dash separator (ID - description format)
			if len(label) <= len(string(persona)) {
				t.Errorf("PersonaLabel(%q) = %q — appears to be only the ID, want a descriptive label", persona, label)
			}
		})
	}
}
