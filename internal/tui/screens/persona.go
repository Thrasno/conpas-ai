package screens

import (
	"strings"

	"github.com/Thrasno/conpas-ai/internal/model"
	"github.com/Thrasno/conpas-ai/internal/tui/styles"
)

// personaLabels maps each persona ID to its human-readable TUI label.
var personaLabels = map[model.PersonaID]string{
	model.PersonaArgentino:        "Argentino - Rioplatense Spanish, passionate teacher",
	model.PersonaNeutral:          "Neutral - Professional, language-agnostic",
	model.PersonaGalleguinho:      "Galleguinho - Galician retranca (constructive irony)",
	model.PersonaAsturianu:        "Asturianu - Asturian expressions, friendly",
	model.PersonaSargentoDeHierro: "Sargento de Hierro - Minimal verbosity, hyper-technical",
	model.PersonaStark:            "Stark - Tony Stark personality (genius, witty)",
	model.PersonaLittleYoda:       "Little Yoda - Cryptic Jedi Master (ERP wisdom)",
	model.PersonaCustom:           "Custom - Use your own persona file",
}

// PersonaLabel returns the human-readable TUI label for a persona.
func PersonaLabel(persona model.PersonaID) string {
	if label, ok := personaLabels[persona]; ok {
		return label
	}
	return string(persona)
}

func PersonaOptions() []model.PersonaID {
	return []model.PersonaID{
		model.PersonaArgentino,
		model.PersonaNeutral,
		model.PersonaGalleguinho,
		model.PersonaAsturianu,
		model.PersonaSargentoDeHierro,
		model.PersonaStark,
		model.PersonaLittleYoda,
		model.PersonaCustom,
	}
}

func RenderPersona(selected model.PersonaID, cursor int) string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Choose your Persona"))
	b.WriteString("\n\n")
	b.WriteString(styles.SubtextStyle.Render("Your own Gentleman! teaches before it solves."))
	b.WriteString("\n\n")

	for idx, persona := range PersonaOptions() {
		isSelected := persona == selected
		focused := idx == cursor
		b.WriteString(renderRadio(PersonaLabel(persona), isSelected, focused))
	}

	b.WriteString("\n")
	b.WriteString(renderOptions([]string{"Back"}, cursor-len(PersonaOptions())))
	b.WriteString("\n")
	b.WriteString(styles.HelpStyle.Render("j/k: navigate • enter: select • esc: back"))

	return b.String()
}
