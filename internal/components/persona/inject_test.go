package persona

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Thrasno/conpas-ai/internal/agents"
	"github.com/Thrasno/conpas-ai/internal/agents/claude"
	"github.com/Thrasno/conpas-ai/internal/agents/opencode"
	"github.com/Thrasno/conpas-ai/internal/assets"
	"github.com/Thrasno/conpas-ai/internal/model"
)

func claudeAdapter() agents.Adapter   { return claude.NewAdapter() }
func opencodeAdapter() agents.Adapter { return opencode.NewAdapter() }

func TestInjectClaudeArgentinoWritesSectionWithRealContent(t *testing.T) {
	home := t.TempDir()

	result, err := Inject(home, claudeAdapter(), model.PersonaArgentino)
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}
	if !result.Changed {
		t.Fatalf("Inject() changed = false")
	}

	path := filepath.Join(home, ".claude", "CLAUDE.md")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	text := string(content)
	if !strings.Contains(text, "<!-- gentle-ai:persona -->") {
		t.Fatal("CLAUDE.md missing open marker for persona")
	}
	if !strings.Contains(text, "<!-- /gentle-ai:persona -->") {
		t.Fatal("CLAUDE.md missing close marker for persona")
	}
	// Real content check — the embedded persona has these patterns.
	if !strings.Contains(text, "Senior Architect") {
		t.Fatal("CLAUDE.md missing real persona content (expected 'Senior Architect')")
	}
}

func TestInjectClaudeArgentinoWritesOutputStyleFile(t *testing.T) {
	home := t.TempDir()

	_, err := Inject(home, claudeAdapter(), model.PersonaArgentino)
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}

	// Verify output-style file was written.
	stylePath := filepath.Join(home, ".claude", "output-styles", "gentleman.md")
	content, err := os.ReadFile(stylePath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", stylePath, err)
	}

	text := string(content)
	if !strings.Contains(text, "name: Gentleman") {
		t.Fatal("Output style file missing YAML frontmatter 'name: Gentleman'")
	}
	if !strings.Contains(text, "keep-coding-instructions: true") {
		t.Fatal("Output style file missing 'keep-coding-instructions: true'")
	}
	if !strings.Contains(text, "Gentleman Output Style") {
		t.Fatal("Output style file missing 'Gentleman Output Style' heading")
	}
}

func TestInjectClaudeArgentinoMergesOutputStyleIntoSettings(t *testing.T) {
	home := t.TempDir()

	// Pre-create a settings.json with some existing content.
	settingsDir := filepath.Join(home, ".claude")
	if err := os.MkdirAll(settingsDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	existingSettings := `{"permissions": {"allow": ["Read"]}, "syntaxHighlightingDisabled": true}`
	if err := os.WriteFile(filepath.Join(settingsDir, "settings.json"), []byte(existingSettings), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	_, err := Inject(home, claudeAdapter(), model.PersonaArgentino)
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}

	// Verify settings.json has outputStyle merged in.
	settingsPath := filepath.Join(home, ".claude", "settings.json")
	settingsContent, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", settingsPath, err)
	}

	var settings map[string]any
	if err := json.Unmarshal(settingsContent, &settings); err != nil {
		t.Fatalf("Unmarshal settings.json error = %v", err)
	}

	outputStyle, ok := settings["outputStyle"]
	if !ok {
		t.Fatal("settings.json missing 'outputStyle' key")
	}
	if outputStyle != "Gentleman" {
		t.Fatalf("settings.json outputStyle = %q, want %q", outputStyle, "Gentleman")
	}

	// Verify existing keys were preserved.
	if _, ok := settings["permissions"]; !ok {
		t.Fatal("settings.json lost 'permissions' key during merge")
	}
	if _, ok := settings["syntaxHighlightingDisabled"]; !ok {
		t.Fatal("settings.json lost 'syntaxHighlightingDisabled' key during merge")
	}
}

func TestInjectClaudeArgentinoReturnsAllFiles(t *testing.T) {
	home := t.TempDir()

	result, err := Inject(home, claudeAdapter(), model.PersonaArgentino)
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}

	// Should return 3 files: CLAUDE.md, output-style, settings.json.
	if len(result.Files) != 3 {
		t.Fatalf("Inject() returned %d files, want 3: %v", len(result.Files), result.Files)
	}

	wantSuffixes := []string{"CLAUDE.md", "gentleman.md", "settings.json"}
	for _, suffix := range wantSuffixes {
		found := false
		for _, f := range result.Files {
			if strings.HasSuffix(f, suffix) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Inject() missing file with suffix %q in %v", suffix, result.Files)
		}
	}
}

func TestInjectClaudeNeutralWritesFullPersonaWithoutRegionalLanguage(t *testing.T) {
	home := t.TempDir()

	result, err := Inject(home, claudeAdapter(), model.PersonaNeutral)
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}
	if !result.Changed {
		t.Fatalf("Inject() changed = false")
	}

	path := filepath.Join(home, ".claude", "CLAUDE.md")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	text := string(content)
	// Neutral persona is the same teacher — should have Senior Architect.
	if !strings.Contains(text, "Senior Architect") {
		t.Fatal("Neutral persona should contain 'Senior Architect'")
	}
	// Should NOT have gentleman-specific regional language.
	if strings.Contains(text, "Rioplatense") {
		t.Fatal("Neutral persona should not contain Rioplatense language")
	}
}

func TestInjectClaudeNeutralDoesNotWriteOutputStyle(t *testing.T) {
	home := t.TempDir()

	result, err := Inject(home, claudeAdapter(), model.PersonaNeutral)
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}

	// Should only return CLAUDE.md, no output-style file.
	if len(result.Files) != 1 {
		t.Fatalf("Neutral persona returned %d files, want 1: %v", len(result.Files), result.Files)
	}

	// Output-style file should NOT exist.
	stylePath := filepath.Join(home, ".claude", "output-styles", "gentleman.md")
	if _, err := os.Stat(stylePath); !os.IsNotExist(err) {
		t.Fatal("Neutral persona should NOT write output-style file")
	}
}

func TestInjectCustomClaudeDoesNothing(t *testing.T) {
	home := t.TempDir()

	result, err := Inject(home, claudeAdapter(), model.PersonaCustom)
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}
	if result.Changed {
		t.Fatal("Custom persona should NOT change anything")
	}
	if len(result.Files) != 0 {
		t.Fatalf("Custom persona should return no files, got %v", result.Files)
	}

	// CLAUDE.md should NOT be created.
	claudeMD := filepath.Join(home, ".claude", "CLAUDE.md")
	if _, err := os.Stat(claudeMD); !os.IsNotExist(err) {
		t.Fatal("Custom persona should NOT create CLAUDE.md")
	}
}

func TestInjectCustomOpenCodeDoesNothing(t *testing.T) {
	home := t.TempDir()

	result, err := Inject(home, opencodeAdapter(), model.PersonaCustom)
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}
	if result.Changed {
		t.Fatal("Custom persona (OpenCode) should NOT change anything")
	}
	if len(result.Files) != 0 {
		t.Fatalf("Custom persona (OpenCode) should return no files, got %v", result.Files)
	}

	// AGENTS.md should NOT be created.
	agentsMD := filepath.Join(home, ".config", "opencode", "AGENTS.md")
	if _, err := os.Stat(agentsMD); !os.IsNotExist(err) {
		t.Fatal("Custom persona (OpenCode) should NOT create AGENTS.md")
	}
}

func TestInjectOpenCodeArgentinoWritesAgentsFile(t *testing.T) {
	home := t.TempDir()

	result, err := Inject(home, opencodeAdapter(), model.PersonaArgentino)
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}
	if !result.Changed {
		t.Fatalf("Inject() changed = false")
	}

	path := filepath.Join(home, ".config", "opencode", "AGENTS.md")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	text := string(content)
	if !strings.Contains(text, "Senior Architect") {
		t.Fatal("AGENTS.md missing real persona content")
	}
}

func TestInjectOpenCodeNeutralReplacesPersonaContent(t *testing.T) {
	home := t.TempDir()

	// First install argentino persona + simulate SDD/engram sections
	_, err := Inject(home, opencodeAdapter(), model.PersonaArgentino)
	if err != nil {
		t.Fatalf("Inject(gentleman) error = %v", err)
	}

	path := filepath.Join(home, ".config", "opencode", "AGENTS.md")

	// Simulate SDD and engram sections appended by sdd.Inject and engram.Inject
	existing, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	withSections := string(existing) + "\n\n<!-- gentle-ai:sdd-orchestrator -->\nSDD orchestrator content here\n<!-- /gentle-ai:sdd-orchestrator -->\n\n<!-- gentle-ai:engram-protocol -->\nEngram protocol content here\n<!-- /gentle-ai:engram-protocol -->\n"
	if err := os.WriteFile(path, []byte(withSections), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	// Now switch to neutral persona
	result, err := Inject(home, opencodeAdapter(), model.PersonaNeutral)
	if err != nil {
		t.Fatalf("Inject(neutral) error = %v", err)
	}
	if !result.Changed {
		t.Fatal("Inject(neutral) should report changed")
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() after neutral error = %v", err)
	}
	text := string(content)

	// Neutral content should be present
	if !strings.Contains(text, "Senior Architect") {
		t.Fatal("AGENTS.md missing neutral persona content")
	}
	// Neutral persona should NOT have Rioplatense language
	if strings.Contains(text, "Rioplatense") {
		t.Fatal("AGENTS.md has Rioplatense language in neutral persona — should be neutral tone")
	}

	// All standard personas do full-file replacement — SDD/engram sections
	// are re-injected by their own inject steps, not preserved by persona inject.
	// (Persona inject is always the first step in the pipeline.)
}

func TestInjectVSCodeNeutralReplacesPersonaContent(t *testing.T) {
	home := t.TempDir()

	vscodeAdapter, err := agents.NewAdapter("vscode-copilot")
	if err != nil {
		t.Fatalf("NewAdapter(vscode-copilot) error = %v", err)
	}

	_, err = Inject(home, vscodeAdapter, model.PersonaArgentino)
	if err != nil {
		t.Fatalf("Inject(gentleman) error = %v", err)
	}

	path := vscodeAdapter.SystemPromptFile(home)

	existing, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	withSections := string(existing) + "\n\n<!-- gentle-ai:sdd-orchestrator -->\nSDD content\n<!-- /gentle-ai:sdd-orchestrator -->\n"
	if err := os.WriteFile(path, []byte(withSections), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	_, err = Inject(home, vscodeAdapter, model.PersonaNeutral)
	if err != nil {
		t.Fatalf("Inject(neutral) error = %v", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() after neutral error = %v", err)
	}
	text := string(content)

	// Neutral persona content must be present.
	if !strings.Contains(text, "Senior Architect") {
		t.Fatal("instructions file missing neutral persona content")
	}
	if strings.Contains(text, "Rioplatense") {
		t.Fatal("instructions file has Rioplatense language in neutral persona")
	}
	// YAML frontmatter must be present (VSCode wraps with frontmatter).
	if !strings.Contains(text, "---\nname:") {
		t.Fatal("instructions file missing YAML frontmatter")
	}
	// All standard personas do full-file replacement — SDD/engram sections
	// are re-injected by their own inject steps, not preserved by persona inject.
}

func TestInjectNeutralOverwritesFileContent(t *testing.T) {
	home := t.TempDir()

	opencodeAdapter, err := agents.NewAdapter("opencode")
	if err != nil {
		t.Fatalf("NewAdapter(opencode) error = %v", err)
	}

	promptPath := opencodeAdapter.SystemPromptFile(home)
	if err := os.MkdirAll(filepath.Dir(promptPath), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	// File contains only a managed marker — persona inject replaces it with
	// fresh persona content (SDD sections are re-injected by their own step).
	markerOnly := "<!-- gentle-ai:sdd-orchestrator -->\nSDD content\n<!-- /gentle-ai:sdd-orchestrator -->\n"
	if err := os.WriteFile(promptPath, []byte(markerOnly), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	_, err = Inject(home, opencodeAdapter, model.PersonaNeutral)
	if err != nil {
		t.Fatalf("Inject(neutral) error = %v", err)
	}

	content, err := os.ReadFile(promptPath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	text := string(content)

	// Persona content must be written.
	if !strings.Contains(text, "Senior Architect") {
		t.Fatal("missing neutral persona content after inject")
	}
}

func TestInjectNeutralIsIdempotent(t *testing.T) {
	home := t.TempDir()

	opencodeAdapter, err := agents.NewAdapter("opencode")
	if err != nil {
		t.Fatalf("NewAdapter(opencode) error = %v", err)
	}

	// First neutral inject
	first, err := Inject(home, opencodeAdapter, model.PersonaNeutral)
	if err != nil {
		t.Fatalf("Inject(neutral) first error = %v", err)
	}
	if !first.Changed {
		t.Fatal("Inject(neutral) first should report changed (new file)")
	}

	// Second neutral inject — content unchanged, should be idempotent
	second, err := Inject(home, opencodeAdapter, model.PersonaNeutral)
	if err != nil {
		t.Fatalf("Inject(neutral) second error = %v", err)
	}
	if second.Changed {
		t.Fatal("Inject(neutral) second should NOT report changed (already up to date)")
	}

	promptPath := opencodeAdapter.SystemPromptFile(home)
	content, err := os.ReadFile(promptPath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	text := string(content)

	// Persona should appear exactly once — no duplication from idempotent inject.
	if strings.Count(text, "## Rules") != 1 {
		t.Fatal("neutral persona section duplicated after idempotent inject")
	}
}

func TestInjectClaudeIsIdempotent(t *testing.T) {
	home := t.TempDir()

	first, err := Inject(home, claudeAdapter(), model.PersonaArgentino)
	if err != nil {
		t.Fatalf("Inject() first error = %v", err)
	}
	if !first.Changed {
		t.Fatalf("Inject() first changed = false")
	}

	second, err := Inject(home, claudeAdapter(), model.PersonaArgentino)
	if err != nil {
		t.Fatalf("Inject() second error = %v", err)
	}
	if second.Changed {
		t.Fatalf("Inject() second changed = true")
	}
}

func TestInjectOpenCodeIsIdempotent(t *testing.T) {
	home := t.TempDir()

	first, err := Inject(home, opencodeAdapter(), model.PersonaArgentino)
	if err != nil {
		t.Fatalf("Inject() first error = %v", err)
	}
	if !first.Changed {
		t.Fatalf("Inject() first changed = false")
	}

	second, err := Inject(home, opencodeAdapter(), model.PersonaArgentino)
	if err != nil {
		t.Fatalf("Inject() second error = %v", err)
	}
	if second.Changed {
		t.Fatalf("Inject() second changed = true")
	}
}

func TestInjectWindsurfIsIdempotent(t *testing.T) {
	home := t.TempDir()

	windsurfAdapter, err := agents.NewAdapter("windsurf")
	if err != nil {
		t.Fatalf("NewAdapter(windsurf) error = %v", err)
	}

	first, err := Inject(home, windsurfAdapter, model.PersonaArgentino)
	if err != nil {
		t.Fatalf("Inject() first error = %v", err)
	}
	if !first.Changed {
		t.Fatalf("Inject() first changed = false")
	}

	promptPath := windsurfAdapter.SystemPromptFile(home)
	contentAfterFirst, err := os.ReadFile(promptPath)
	if err != nil {
		t.Fatalf("ReadFile() after first inject error = %v", err)
	}

	second, err := Inject(home, windsurfAdapter, model.PersonaArgentino)
	if err != nil {
		t.Fatalf("Inject() second error = %v", err)
	}
	if second.Changed {
		t.Fatalf("Inject() second changed = true — persona was duplicated in global_rules.md")
	}

	contentAfterSecond, err := os.ReadFile(promptPath)
	if err != nil {
		t.Fatalf("ReadFile() after second inject error = %v", err)
	}

	if string(contentAfterFirst) != string(contentAfterSecond) {
		t.Fatal("global_rules.md content changed on second inject — persona was duplicated")
	}
}

func TestInjectCursorArgentinoWritesRulesFileWithRealContent(t *testing.T) {
	home := t.TempDir()

	cursorAdapter, err := agents.NewAdapter("cursor")
	if err != nil {
		t.Fatalf("NewAdapter(cursor) error = %v", err)
	}

	result, injectErr := Inject(home, cursorAdapter, model.PersonaArgentino)
	if injectErr != nil {
		t.Fatalf("Inject(cursor) error = %v", injectErr)
	}

	if !result.Changed {
		t.Fatalf("Inject(cursor, gentleman) changed = false")
	}

	// Verify the generic persona content was used — not just neutral one-liner.
	path := filepath.Join(home, ".cursor", "rules", "gentle-ai.mdc")
	content, readErr := os.ReadFile(path)
	if readErr != nil {
		t.Fatalf("ReadFile(%q) error = %v", path, readErr)
	}

	text := string(content)
	if !strings.Contains(text, "Senior Architect") {
		t.Fatal("Cursor persona missing 'Senior Architect' — got neutral fallback instead of generic persona")
	}
	if !strings.Contains(text, "Skills") {
		t.Fatal("Cursor persona missing skills section")
	}
}

func TestInjectGeminiArgentinoWritesSystemPromptWithRealContent(t *testing.T) {
	home := t.TempDir()

	geminiAdapter, err := agents.NewAdapter("gemini-cli")
	if err != nil {
		t.Fatalf("NewAdapter(gemini-cli) error = %v", err)
	}

	result, injectErr := Inject(home, geminiAdapter, model.PersonaArgentino)
	if injectErr != nil {
		t.Fatalf("Inject(gemini) error = %v", injectErr)
	}

	if !result.Changed {
		t.Fatal("Inject(gemini, gentleman) changed = false")
	}

	path := filepath.Join(home, ".gemini", "GEMINI.md")
	content, readErr := os.ReadFile(path)
	if readErr != nil {
		t.Fatalf("ReadFile(%q) error = %v", path, readErr)
	}

	text := string(content)
	if !strings.Contains(text, "Senior Architect") {
		t.Fatal("Gemini persona missing 'Senior Architect'")
	}
}

func TestInjectVSCodeArgentinoWritesInstructionsFile(t *testing.T) {
	home := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))

	vscodeAdapter, err := agents.NewAdapter("vscode-copilot")
	if err != nil {
		t.Fatalf("NewAdapter(vscode-copilot) error = %v", err)
	}

	result, injectErr := Inject(home, vscodeAdapter, model.PersonaArgentino)
	if injectErr != nil {
		t.Fatalf("Inject(vscode) error = %v", injectErr)
	}

	if !result.Changed {
		t.Fatal("Inject(vscode, gentleman) changed = false")
	}

	path := vscodeAdapter.SystemPromptFile(home)
	content, readErr := os.ReadFile(path)
	if readErr != nil {
		t.Fatalf("ReadFile(%q) error = %v", path, readErr)
	}

	text := string(content)
	if !strings.Contains(text, "applyTo: \"**\"") {
		t.Fatal("VS Code instructions file missing YAML frontmatter applyTo pattern")
	}
	if !strings.Contains(text, "Senior Architect") {
		t.Fatal("VS Code persona missing 'Senior Architect'")
	}
}

// --- Auto-heal tests: Claude Code stale free-text persona ---

// legacyClaudePersonaBlock simulates a Gentleman persona block that was written
// directly (without markers) by an old installer or manually by the user.
const legacyClaudePersonaBlock = `## Rules

- NEVER add "Co-Authored-By" or any AI attribution to commits. Use conventional commits format only.
- Never build after changes.

## Personality

Senior Architect, 15+ years experience, GDE & MVP.

## Language

- Spanish input → Rioplatense Spanish.

## Behavior

- Push back when user asks for code without context.

`

func TestInjectClaudeAutoHealsStaleFreeTextPersona(t *testing.T) {
	home := t.TempDir()

	// Pre-populate CLAUDE.md with legacy persona content (no markers) followed
	// by a properly-marked section from a previous installer run.
	claudeMD := filepath.Join(home, ".claude", "CLAUDE.md")
	if err := os.MkdirAll(filepath.Dir(claudeMD), 0o755); err != nil {
		t.Fatalf("MkdirAll error = %v", err)
	}

	// Simulate a stale install: free-text persona block at top, then a different
	// marked section below (e.g., from a previous SDD install).
	stalePreamble := legacyClaudePersonaBlock + "\n<!-- gentle-ai:sdd -->\nOld SDD content.\n<!-- /gentle-ai:sdd -->\n"
	if err := os.WriteFile(claudeMD, []byte(stalePreamble), 0o644); err != nil {
		t.Fatalf("WriteFile error = %v", err)
	}

	result, err := Inject(home, claudeAdapter(), model.PersonaArgentino)
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}
	if !result.Changed {
		t.Fatal("Inject() should have changed the file to remove the legacy block")
	}

	content, err := os.ReadFile(claudeMD)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	text := string(content)

	// The file should now have the persona inside markers, not as free text.
	if !strings.Contains(text, "<!-- gentle-ai:persona -->") {
		t.Fatal("CLAUDE.md missing persona marker after heal")
	}
	if !strings.Contains(text, "<!-- /gentle-ai:persona -->") {
		t.Fatal("CLAUDE.md missing persona close marker after heal")
	}

	// The existing SDD section must be preserved.
	if !strings.Contains(text, "<!-- gentle-ai:sdd -->") {
		t.Fatal("CLAUDE.md lost the sdd section during heal")
	}
	if !strings.Contains(text, "Old SDD content.") {
		t.Fatal("CLAUDE.md lost the sdd section content during heal")
	}

	// The persona content must NOT appear twice (no duplicate blocks).
	firstPersonaIdx := strings.Index(text, "Senior Architect")
	if firstPersonaIdx < 0 {
		t.Fatal("CLAUDE.md missing 'Senior Architect' persona content")
	}
	// Verify there's no second occurrence outside the markers.
	lastPersonaIdx := strings.LastIndex(text, "Senior Architect")
	if firstPersonaIdx != lastPersonaIdx {
		// It's OK if the same string appears inside the single persona marker block
		// multiple times (e.g., content + newlines), but there must not be a
		// separate free-text block also containing it.
		// Check: everything before the open marker should NOT contain "Senior Architect".
		openMarkerIdx := strings.Index(text, "<!-- gentle-ai:persona -->")
		if openMarkerIdx >= 0 && strings.Contains(text[:openMarkerIdx], "Senior Architect") {
			t.Fatal("CLAUDE.md still has 'Senior Architect' before the persona marker — legacy block not fully stripped")
		}
	}
}

func TestInjectClaudeAutoHealStalePersonaOnlyFile(t *testing.T) {
	home := t.TempDir()

	// CLAUDE.md contains ONLY the legacy persona block (no markers at all).
	claudeMD := filepath.Join(home, ".claude", "CLAUDE.md")
	if err := os.MkdirAll(filepath.Dir(claudeMD), 0o755); err != nil {
		t.Fatalf("MkdirAll error = %v", err)
	}
	if err := os.WriteFile(claudeMD, []byte(legacyClaudePersonaBlock), 0o644); err != nil {
		t.Fatalf("WriteFile error = %v", err)
	}

	result, err := Inject(home, claudeAdapter(), model.PersonaArgentino)
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}
	if !result.Changed {
		t.Fatal("Inject() should have changed the file")
	}

	content, err := os.ReadFile(claudeMD)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	text := string(content)

	// Must have markers now.
	if !strings.Contains(text, "<!-- gentle-ai:persona -->") {
		t.Fatal("CLAUDE.md missing persona marker")
	}

	// Must NOT have the legacy free-text block before markers.
	openMarkerIdx := strings.Index(text, "<!-- gentle-ai:persona -->")
	if openMarkerIdx >= 0 {
		before := text[:openMarkerIdx]
		if strings.Contains(before, "## Rules") {
			t.Fatal("legacy '## Rules' block still present before persona marker")
		}
	}
}

func TestInjectClaudeHealDoesNotTouchNonPersonaContent(t *testing.T) {
	home := t.TempDir()

	// CLAUDE.md has user content that does NOT match persona fingerprints.
	claudeMD := filepath.Join(home, ".claude", "CLAUDE.md")
	if err := os.MkdirAll(filepath.Dir(claudeMD), 0o755); err != nil {
		t.Fatalf("MkdirAll error = %v", err)
	}
	userContent := "# My custom config\n\nI like turtles.\n"
	if err := os.WriteFile(claudeMD, []byte(userContent), 0o644); err != nil {
		t.Fatalf("WriteFile error = %v", err)
	}

	result, err := Inject(home, claudeAdapter(), model.PersonaArgentino)
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}
	if !result.Changed {
		t.Fatal("Inject() should write persona section")
	}

	content, err := os.ReadFile(claudeMD)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	text := string(content)

	// User content must be preserved.
	if !strings.Contains(text, "I like turtles.") {
		t.Fatal("user content was erased — heal was too aggressive")
	}
	// Persona section must be appended.
	if !strings.Contains(text, "<!-- gentle-ai:persona -->") {
		t.Fatal("persona section not appended")
	}
}

// --- Auto-heal tests: VSCode stale legacy path cleanup ---

func TestInjectVSCodeCleansLegacyGitHubPersonaFile(t *testing.T) {
	home := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))

	// Plant an old-style Gentleman persona file at the legacy path.
	legacyPath := filepath.Join(home, ".github", "copilot-instructions.md")
	if err := os.MkdirAll(filepath.Dir(legacyPath), 0o755); err != nil {
		t.Fatalf("MkdirAll error = %v", err)
	}
	// Old installer wrote raw persona content without YAML frontmatter.
	oldContent := "## Personality\n\nSenior Architect, 15+ years experience.\n"
	if err := os.WriteFile(legacyPath, []byte(oldContent), 0o644); err != nil {
		t.Fatalf("WriteFile error = %v", err)
	}

	vscodeAdapter, err := agents.NewAdapter("vscode-copilot")
	if err != nil {
		t.Fatalf("NewAdapter(vscode-copilot) error = %v", err)
	}

	result, injectErr := Inject(home, vscodeAdapter, model.PersonaArgentino)
	if injectErr != nil {
		t.Fatalf("Inject(vscode) error = %v", injectErr)
	}
	if !result.Changed {
		t.Fatal("Inject(vscode) should report changed (legacy cleanup + new file write)")
	}

	// Legacy file must be gone.
	if _, statErr := os.Stat(legacyPath); !os.IsNotExist(statErr) {
		t.Fatal("legacy ~/.github/copilot-instructions.md was NOT removed by auto-heal")
	}

	// New file must exist at the current path.
	newPath := vscodeAdapter.SystemPromptFile(home)
	content, readErr := os.ReadFile(newPath)
	if readErr != nil {
		t.Fatalf("ReadFile new path %q error = %v", newPath, readErr)
	}
	if !strings.Contains(string(content), "applyTo: \"**\"") {
		t.Fatal("new VSCode instructions file missing YAML frontmatter")
	}
}

func TestInjectVSCodePreservesNonPersonaGitHubFile(t *testing.T) {
	home := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))

	// Plant a .github/copilot-instructions.md that has user content (not a
	// Gentleman persona) — it must NOT be deleted.
	legacyPath := filepath.Join(home, ".github", "copilot-instructions.md")
	if err := os.MkdirAll(filepath.Dir(legacyPath), 0o755); err != nil {
		t.Fatalf("MkdirAll error = %v", err)
	}
	userContent := "# My custom Copilot instructions\n\nAlways be concise.\n"
	if err := os.WriteFile(legacyPath, []byte(userContent), 0o644); err != nil {
		t.Fatalf("WriteFile error = %v", err)
	}

	vscodeAdapter, err := agents.NewAdapter("vscode-copilot")
	if err != nil {
		t.Fatalf("NewAdapter(vscode-copilot) error = %v", err)
	}

	_, injectErr := Inject(home, vscodeAdapter, model.PersonaArgentino)
	if injectErr != nil {
		t.Fatalf("Inject(vscode) error = %v", injectErr)
	}

	// User's file must still exist.
	remaining, readErr := os.ReadFile(legacyPath)
	if readErr != nil {
		t.Fatalf("legacy user file was deleted: ReadFile error = %v", readErr)
	}
	if string(remaining) != userContent {
		t.Fatalf("user file content was modified: got %q", string(remaining))
	}
}

func TestNeutralAndGentlemanToneSectionsMatch(t *testing.T) {
	// With the base+variant architecture, all personas share the same ## Tone section
	// from persona-base.md. Verify the base file contains the Tone section.
	base := assets.MustRead("generic/persona-base.md")

	if !strings.Contains(base, "## Tone") {
		t.Fatal("generic/persona-base.md missing ## Tone section — all personas inherit this")
	}
}

// TestPersonaContentConcatenation verifies that personaContent() (base + variant)
// produces valid markdown that includes all 7 required sections for each of the 6 personas.
func TestPersonaContentConcatenation(t *testing.T) {
	requiredSections := []string{
		"## Rules",
		"## Personality",
		"## Tone",
		"## Philosophy",
		"## Expertise",
		"## Behavior",
		"## Skills",
	}

	personas := []model.PersonaID{
		model.PersonaArgentino,
		model.PersonaNeutral,
		model.PersonaGalleguinho,
		model.PersonaAsturianu,
		model.PersonaSargentoDeHierro,
		model.PersonaStark,
	}

	for _, persona := range personas {
		persona := persona
		t.Run(string(persona), func(t *testing.T) {
			// Use generic agent (non-Claude) to test generic base.
			content := personaContent(model.AgentOpenCode, persona)
			if content == "" {
				t.Fatalf("personaContent(%q) returned empty string", persona)
			}
			for _, section := range requiredSections {
				if !strings.Contains(content, section) {
					t.Errorf("persona %q missing section %q in concatenated content", persona, section)
				}
			}
		})
	}
}

func TestInjectVSCodeIdempotentAfterHeal(t *testing.T) {
	home := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))

	// Plant legacy file and run inject twice — second run should be idempotent.
	legacyPath := filepath.Join(home, ".github", "copilot-instructions.md")
	if err := os.MkdirAll(filepath.Dir(legacyPath), 0o755); err != nil {
		t.Fatalf("MkdirAll error = %v", err)
	}
	if err := os.WriteFile(legacyPath, []byte("## Personality\n\nSenior Architect.\n"), 0o644); err != nil {
		t.Fatalf("WriteFile error = %v", err)
	}

	vscodeAdapter, err := agents.NewAdapter("vscode-copilot")
	if err != nil {
		t.Fatalf("NewAdapter(vscode-copilot) error = %v", err)
	}

	first, err := Inject(home, vscodeAdapter, model.PersonaArgentino)
	if err != nil {
		t.Fatalf("Inject() first error = %v", err)
	}
	if !first.Changed {
		t.Fatal("first inject should have changed")
	}

	second, err := Inject(home, vscodeAdapter, model.PersonaArgentino)
	if err != nil {
		t.Fatalf("Inject() second error = %v", err)
	}
	if second.Changed {
		t.Fatalf("second inject should be idempotent (changed = false), but changed = true")
	}
}

// TestInjectAllPersonasProduceValidContent verifies that Inject() works correctly
// for each of the 7 standard personas: writes files and includes correct sections.
func TestInjectAllPersonasProduceValidContent(t *testing.T) {
	cases := []struct {
		persona          model.PersonaID
		wantLanguageHint string // distinctive string in that persona's ## Language section
		wantNoLanguage   string // string that must NOT appear (from other personas)
	}{
		{
			persona:          model.PersonaArgentino,
			wantLanguageHint: "Rioplatense",
		},
		{
			persona:          model.PersonaNeutral,
			wantLanguageHint: "## Language",
			wantNoLanguage:   "Rioplatense",
		},
		{
			persona:          model.PersonaGalleguinho,
			wantLanguageHint: "Galician",
		},
		{
			persona:          model.PersonaAsturianu,
			wantLanguageHint: "Asturian",
		},
		{
			persona:          model.PersonaSargentoDeHierro,
			wantLanguageHint: "## Language",
			wantNoLanguage:   "Rioplatense",
		},
		{
			persona:          model.PersonaStark,
			wantLanguageHint: "## Language",
		},
		{
			persona:          model.PersonaLittleYoda,
			wantLanguageHint: "Yoda flavor",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(string(tc.persona), func(t *testing.T) {
			home := t.TempDir()

			result, err := Inject(home, opencodeAdapter(), tc.persona)
			if err != nil {
				t.Fatalf("Inject(%q) error = %v", tc.persona, err)
			}
			if !result.Changed {
				t.Fatalf("Inject(%q) changed = false", tc.persona)
			}

			path := filepath.Join(home, ".config", "opencode", "AGENTS.md")
			content, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("ReadFile() error = %v", err)
			}
			text := string(content)

			// All standard personas must have base content.
			if !strings.Contains(text, "Senior Architect") {
				t.Errorf("persona %q missing 'Senior Architect' in base content", tc.persona)
			}
			// All standard personas must have ## Language section.
			if !strings.Contains(text, "## Language") {
				t.Errorf("persona %q missing ## Language section", tc.persona)
			}
			// Persona-specific language hint.
			if tc.wantLanguageHint != "" && tc.wantLanguageHint != "## Language" {
				if !strings.Contains(text, tc.wantLanguageHint) {
					t.Errorf("persona %q missing language hint %q", tc.persona, tc.wantLanguageHint)
				}
			}
			// Ensure unwanted language markers are absent.
			if tc.wantNoLanguage != "" && strings.Contains(text, tc.wantNoLanguage) {
				t.Errorf("persona %q should NOT contain %q", tc.persona, tc.wantNoLanguage)
			}
		})
	}
}
