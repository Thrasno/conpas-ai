# Usage

← [Back to README](../README.md)

---

## Persona Modes

| Persona | ID | Description |
|---------|-----|-------------|
| Argentino | `argentino` | Teaching-oriented mentor with Rioplatense Spanish flavor — passionate, direct, uses voseo |
| Gentleman | `gentleman` | Alias for `argentino` — preserved for backward compatibility |
| Neutral | `neutral` | Same teacher, same philosophy, no regional language — warm and professional |
| Galleguinho | `galleguinho` | Galician-Spanish flavor — friendly warmth with Galician expressions |
| Asturianu | `asturianu` | Asturian dialect flavor — peer-level, fun expressions like "guaje" and "ta de mieu" |
| Sargento de Hierro | `sargentoDeHierro` | Drill-sergeant mode — blunt, demanding, zero tolerance for excuses |
| Stark | `stark` | Tony Stark mode — brilliant, confident, tech-billionaire wit |
| Little Yoda | `littleYoda` | Cryptic Jedi Master — Yoda syntax, ERP wisdom, minimal verbosity |
| Custom | `custom` | Bring your own persona instructions |

---

## Interactive TUI

Just run it — the Bubbletea TUI guides you through agent selection, components, skills, and presets:

```bash
gentle-ai
```

---

## CLI Commands

### install

First-time setup — detects your tools, configures agents, injects all components:

```bash
# Full ecosystem for multiple agents
gentle-ai install \
  --agent claude-code,opencode,gemini-cli \
  --preset full

# Minimal setup for Cursor
gentle-ai install \
  --agent cursor \
  --preset minimal

# Pick specific components and skills
gentle-ai install \
  --agent claude-code \
  --component engram,sdd,skills,context7,persona,permissions \
  --skill go-testing,skill-creator,branch-pr,issue-creation \
  --persona argentino

# Dry-run first (preview plan without applying changes)
gentle-ai install --dry-run \
  --agent claude-code,opencode \
  --preset full
```

### sync

Refresh managed assets to the current version. Use after `brew upgrade gentle-ai` or when you want your local configs aligned with the latest release. Does NOT reinstall binaries (engram, GGA) — only updates prompt content, skills, MCP configs, and SDD orchestrators.

```bash
# Sync all installed agents
gentle-ai sync

# Sync specific agents only
gentle-ai sync --agent cursor --agent windsurf

# Sync a specific component
gentle-ai sync --component sdd
gentle-ai sync --component skills
gentle-ai sync --component engram
```

Sync is safe and idempotent — running it twice produces no changes the second time.

### update / upgrade

Check for and install new versions of `gentle-ai` itself:

```bash
# Check if a newer version is available
gentle-ai update

# Upgrade to the latest release (downloads new binary, replaces current)
gentle-ai upgrade
```

After upgrading, run `gentle-ai sync` to refresh all managed assets to the new version's content.

### version

```bash
gentle-ai version
gentle-ai --version
gentle-ai -v
```

---

## CLI Flags (install)

| Flag | Description |
|------|-------------|
| `--agent`, `--agents` | Agents to configure (comma-separated) |
| `--component`, `--components` | Components to install (comma-separated) |
| `--skill`, `--skills` | Skills to install (comma-separated) |
| `--persona` | Persona mode: `argentino`, `gentleman` (alias), `neutral`, `galleguinho`, `asturianu`, `sargentoDeHierro`, `stark`, `littleYoda`, `custom` |
| `--preset` | Preset: `full`, `full-gentleman` (alias), `ecosystem-only`, `minimal`, `custom` |
| `--dry-run` | Preview the install plan without applying changes |

## CLI Flags (sync)

| Flag | Description |
|------|-------------|
| `--agent`, `--agents` | Agents to sync (defaults to all installed agents) |
| `--component` | Sync a specific component only: `sdd`, `engram`, `context7`, `skills`, `gga`, `permissions`, `theme` |
| `--include-permissions` | Include permissions sync (opt-in) |
| `--include-theme` | Include theme sync (opt-in) |

---

## Typical Workflow

```bash
# First time: install everything
brew install gentleman-programming/tap/gentle-ai
gentle-ai install --agent claude-code,cursor --preset full

# After a new release: upgrade + sync
brew upgrade gentle-ai
gentle-ai sync

# Adding a new agent later
gentle-ai install --agent windsurf --preset full
```

---

## Dependency Management

`gentle-ai` auto-detects prerequisites before installation and provides platform-specific guidance:

- **Detected tools**: git, curl, node, npm, brew, go
- **Version checks**: validates minimum versions where applicable
- **Platform-aware hints**: suggests `brew install`, `apt install`, `pacman -S`, `dnf install`, or `winget install` depending on your OS
- **Node LTS alignment**: on apt/dnf systems, Node.js hints use NodeSource LTS bootstrap before package install
- **Dependency-first approach**: detects what's installed, calculates what's needed, shows the full dependency tree before installing anything, then verifies each dependency after installation
