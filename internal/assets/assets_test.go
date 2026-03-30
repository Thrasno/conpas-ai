package assets

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestAllEmbeddedAssetsAreReadable verifies that every expected embedded file
// can be loaded via Read() without error. This catches missing/misnamed files
// at test time rather than at runtime.
func TestAllEmbeddedAssetsAreReadable(t *testing.T) {
	expectedFiles := []string{
		// Claude agent files
		"claude/engram-protocol.md",
		"claude/persona-base.md",
		"claude/sdd-orchestrator.md",

		// Generic persona files
		"generic/persona-base.md",
		"generic/persona-argentino.md",
		"generic/persona-neutral.md",
		"generic/persona-galleguinho.md",
		"generic/persona-asturianu.md",
		"generic/persona-sargentoDeHierro.md",
		"generic/persona-stark.md",

		// OpenCode agent files
		"opencode/sdd-overlay-single.json",
		"opencode/sdd-overlay-multi.json",
		"opencode/commands/sdd-apply.md",
		"opencode/commands/sdd-archive.md",
		"opencode/commands/sdd-continue.md",
		"opencode/commands/sdd-explore.md",
		"opencode/commands/sdd-ff.md",
		"opencode/commands/sdd-init.md",
		"opencode/commands/sdd-new.md",
		"opencode/commands/sdd-verify.md",
		"opencode/plugins/background-agents.ts",

		// Gemini agent files
		"gemini/sdd-orchestrator.md",

		// Codex agent files
		"codex/sdd-orchestrator.md",

		// Cursor agent files
		"cursor/sdd-orchestrator.md",
		"cursor/agents/sdd-init.md",
		"cursor/agents/sdd-explore.md",
		"cursor/agents/sdd-propose.md",
		"cursor/agents/sdd-spec.md",
		"cursor/agents/sdd-design.md",
		"cursor/agents/sdd-tasks.md",
		"cursor/agents/sdd-apply.md",
		"cursor/agents/sdd-verify.md",
		"cursor/agents/sdd-archive.md",

		// SDD skills
		"skills/sdd-init/SKILL.md",
		"skills/sdd-apply/SKILL.md",
		"skills/sdd-archive/SKILL.md",
		"skills/sdd-design/SKILL.md",
		"skills/sdd-explore/SKILL.md",
		"skills/sdd-propose/SKILL.md",
		"skills/sdd-spec/SKILL.md",
		"skills/sdd-tasks/SKILL.md",
		"skills/sdd-verify/SKILL.md",
		"skills/skill-registry/SKILL.md",
		"skills/_shared/persistence-contract.md",
		"skills/_shared/engram-convention.md",
		"skills/_shared/openspec-convention.md",
		"skills/_shared/sdd-phase-common.md",

		// Foundation skills
		"skills/go-testing/SKILL.md",
		"skills/skill-creator/SKILL.md",
		"skills/zoho-deluge/SKILL.md",
	}

	for _, path := range expectedFiles {
		t.Run(path, func(t *testing.T) {
			content, err := Read(path)
			if err != nil {
				t.Fatalf("Read(%q) error = %v", path, err)
			}

			if len(strings.TrimSpace(content)) == 0 {
				t.Fatalf("Read(%q) returned empty content", path)
			}

			// Real content should be substantial, not a one-line stub.
			if len(content) < 50 {
				t.Fatalf("Read(%q) content is suspiciously short (%d bytes) — possible stub", path, len(content))
			}
		})
	}
}

func TestOpenCodeEmbeddedAssetLayout(t *testing.T) {
	entries, err := FS.ReadDir("opencode")
	if err != nil {
		t.Fatalf("ReadDir(opencode) error = %v", err)
	}

	seen := map[string]bool{}
	for _, entry := range entries {
		seen[entry.Name()] = true
	}

	for _, name := range []string{"commands", "plugins", "sdd-overlay-single.json", "sdd-overlay-multi.json"} {
		if !seen[name] {
			t.Fatalf("opencode embedded assets missing %q", name)
		}
	}

	commandEntries, err := FS.ReadDir("opencode/commands")
	if err != nil {
		t.Fatalf("ReadDir(opencode/commands) error = %v", err)
	}
	if len(commandEntries) != 8 {
		t.Fatalf("opencode commands count = %d, want 8", len(commandEntries))
	}

	pluginEntries, err := FS.ReadDir("opencode/plugins")
	if err != nil {
		t.Fatalf("ReadDir(opencode/plugins) error = %v", err)
	}
	if len(pluginEntries) != 1 {
		t.Fatalf("opencode plugins count = %d, want 1", len(pluginEntries))
	}
	if pluginEntries[0].Name() != "background-agents.ts" {
		t.Fatalf("plugin entry = %q, want background-agents.ts", pluginEntries[0].Name())
	}
}

// TestMustReadPanicsOnMissingFile verifies that MustRead panics for a
// nonexistent file, confirming the safety mechanism works.
func TestMustReadPanicsOnMissingFile(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("MustRead() did not panic for missing file")
		}
	}()

	MustRead("nonexistent/file.md")
}

// TestEmbeddedAssetCount verifies we have the expected number of embedded files.
// This catches accidental deletions of asset files.
func TestEmbeddedAssetCount(t *testing.T) {
	// Count skill files.
	entries, err := FS.ReadDir("skills")
	if err != nil {
		t.Fatalf("ReadDir(skills) error = %v", err)
	}

	skillDirs := 0
	for _, entry := range entries {
		if entry.IsDir() {
			skillDirs++
		}
	}

	// We expect 17 skill directories (9 SDD + judgment-day + 5 foundation + zoho-deluge + _shared).
	if skillDirs != 17 {
		t.Fatalf("expected 17 skill directories, got %d", skillDirs)
	}

	// Verify each skill directory has a SKILL.md.
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if entry.Name() == "_shared" {
			for _, sharedFile := range []string{"persistence-contract.md", "engram-convention.md", "openspec-convention.md", "sdd-phase-common.md", "skill-resolver.md"} {
				sharedPath := "skills/_shared/" + sharedFile
				if _, err := Read(sharedPath); err != nil {
					t.Fatalf("shared directory missing %q: %v", sharedFile, err)
				}
			}
			continue
		}
		skillPath := "skills/" + entry.Name() + "/SKILL.md"
		if _, err := Read(skillPath); err != nil {
			t.Fatalf("skill directory %q missing SKILL.md: %v", entry.Name(), err)
		}
	}
}

// TestPersonaFilesHaveZohoDelugeAutoLoad verifies that the persona base files
// contain an auto-load row for the zoho-deluge skill.
func TestPersonaFilesHaveZohoDelugeAutoLoad(t *testing.T) {
	personaFiles := []string{
		"claude/persona-base.md",
		"generic/persona-base.md",
	}

	for _, path := range personaFiles {
		t.Run(path, func(t *testing.T) {
			content := MustRead(path)
			if !strings.Contains(content, "zoho-deluge") {
				t.Fatalf("%q missing zoho-deluge auto-load row in Skills table", path)
			}
		})
	}
}

// TestZohoDelugeSkillHasFrontmatter verifies that the zoho-deluge SKILL.md
// contains the required YAML frontmatter block for skill-registry compatibility.
func TestZohoDelugeSkillHasFrontmatter(t *testing.T) {
	content := MustRead("skills/zoho-deluge/SKILL.md")

	for _, want := range []string{
		"name: zoho-deluge",
		"license: Apache-2.0",
		"author: gentleman-programming",
	} {
		if !strings.Contains(content, want) {
			t.Fatalf("skills/zoho-deluge/SKILL.md missing frontmatter field %q", want)
		}
	}
}

func TestSDDPhaseCommonEnforcesExecutorBoundary(t *testing.T) {
	content := MustRead("skills/_shared/sdd-phase-common.md")

	// Must enforce executor boundary — no delegation allowed.
	for _, want := range []string{
		"EXECUTOR, not an orchestrator",
		"Do NOT launch sub-agents",
		"do NOT call `delegate`/`task`",
	} {
		if !strings.Contains(content, want) {
			t.Fatalf("sdd-phase-common missing executor boundary rule %q", want)
		}
	}

	// Must instruct phase agents to search the skill registry themselves
	// when no explicit skill path was provided — this is skill LOADING, not delegation.
	if !strings.Contains(content, `mem_search(query: "skill-registry"`) {
		t.Fatal("sdd-phase-common must instruct phase agents to search skill-registry themselves for skill loading")
	}

	// Must NOT tell agents to launch sub-agents or delegate tasks.
	for _, forbidden := range []string{
		"launch a sub-agent",
		"delegate this to",
	} {
		if strings.Contains(content, forbidden) {
			t.Fatalf("sdd-phase-common should not contain delegation instruction %q", forbidden)
		}
	}
}

func TestOpenCodeSDDOverlaySubagentsAreExplicitExecutors(t *testing.T) {
	for _, assetPath := range []string{"opencode/sdd-overlay-single.json", "opencode/sdd-overlay-multi.json"} {
		t.Run(assetPath, func(t *testing.T) {
			var root map[string]any
			if err := json.Unmarshal([]byte(MustRead(assetPath)), &root); err != nil {
				t.Fatalf("Unmarshal(%q) error = %v", assetPath, err)
			}

			agents, ok := root["agent"].(map[string]any)
			if !ok {
				t.Fatalf("%q missing agent map", assetPath)
			}

			for _, phase := range []string{"sdd-init", "sdd-explore", "sdd-propose", "sdd-spec", "sdd-design", "sdd-tasks", "sdd-apply", "sdd-verify", "sdd-archive"} {
				agentDef, ok := agents[phase].(map[string]any)
				if !ok {
					t.Fatalf("%q missing %s agent", assetPath, phase)
				}
				prompt, _ := agentDef["prompt"].(string)
				for _, want := range []string{"not the orchestrator", "Do NOT delegate", "Do NOT call task/delegate", "Do NOT launch sub-agents"} {
					if !strings.Contains(prompt, want) {
						t.Fatalf("%q phase %s prompt missing %q", assetPath, phase, want)
					}
				}
			}
		})
	}
}

// TestPersonaBaseHasAllSections verifies that both persona-base.md files (generic
// and claude) contain the 7 required common sections.
func TestPersonaBaseHasAllSections(t *testing.T) {
	requiredSections := []string{
		"## Rules",
		"## Personality",
		"## Tone",
		"## Philosophy",
		"## Expertise",
		"## Behavior",
		"## Skills",
	}

	baseFiles := []string{
		"generic/persona-base.md",
		"claude/persona-base.md",
	}

	for _, path := range baseFiles {
		t.Run(path, func(t *testing.T) {
			content := MustRead(path)
			for _, section := range requiredSections {
				if !strings.Contains(content, section) {
					t.Errorf("%q missing required section %q", path, section)
				}
			}
		})
	}
}

// TestPersonaVariantsHaveLanguage verifies that all 6 variant files contain
// the ## Language section.
func TestPersonaVariantsHaveLanguage(t *testing.T) {
	variantFiles := []string{
		"generic/persona-argentino.md",
		"generic/persona-neutral.md",
		"generic/persona-galleguinho.md",
		"generic/persona-asturianu.md",
		"generic/persona-sargentoDeHierro.md",
		"generic/persona-stark.md",
	}

	for _, path := range variantFiles {
		t.Run(path, func(t *testing.T) {
			content := MustRead(path)
			if !strings.Contains(content, "## Language") {
				t.Errorf("%q missing required ## Language section", path)
			}
		})
	}
}

// TestPersonaVariantsNoDuplication verifies that variant files do NOT contain
// sections that belong exclusively to the base file (Rules, Tone, Philosophy).
func TestPersonaVariantsNoDuplication(t *testing.T) {
	variantFiles := []string{
		"generic/persona-argentino.md",
		"generic/persona-neutral.md",
		"generic/persona-galleguinho.md",
		"generic/persona-asturianu.md",
		"generic/persona-sargentoDeHierro.md",
	}

	forbiddenSections := []string{
		"## Rules",
		"## Tone",
		"## Philosophy",
		"## Expertise",
		"## Behavior",
	}

	for _, path := range variantFiles {
		t.Run(path, func(t *testing.T) {
			content := MustRead(path)
			for _, section := range forbiddenSections {
				if strings.Contains(content, section) {
					t.Errorf("%q must not contain base section %q (move it to persona-base.md)", path, section)
				}
			}
		})
	}
}

func TestSDDOrchestratorAssetsScopedToDedicatedAgent(t *testing.T) {
	for _, assetPath := range []string{
		"generic/sdd-orchestrator.md",
		"claude/sdd-orchestrator.md",
		"gemini/sdd-orchestrator.md",
		"codex/sdd-orchestrator.md",
		"cursor/sdd-orchestrator.md",
	} {
		t.Run(assetPath, func(t *testing.T) {
			content := MustRead(assetPath)
			if !strings.Contains(content, "dedicated `sdd-orchestrator` agent or rule only") {
				t.Fatalf("%q missing dedicated-agent scoping note", assetPath)
			}
			if !strings.Contains(content, "Do NOT apply it to executor phase agents") {
				t.Fatalf("%q missing executor exclusion note", assetPath)
			}
		})
	}
}
