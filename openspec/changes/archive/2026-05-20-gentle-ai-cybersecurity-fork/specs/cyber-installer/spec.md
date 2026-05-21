# Cybersecurity Installer Specification

## Purpose

Defines the `PresetCyber` preset in the TUI installer, the `--cyber` CLI flag, and the cybersecurity edition's behavior during installation — including skill selection, MCP configuration, and agent augmentation.

## Requirements

### Requirement: PresetCyber Constant

The system MUST define a `PresetCyber` constant of type `PresetID` with value `"cyber"` in `internal/model/types.go`. It SHALL be added alongside existing presets without removing or altering them.

#### Scenario: PresetCyber is defined

- GIVEN the codebase is compiled
- WHEN `model.PresetCyber` is referenced
- THEN its value is the string `"cyber"`

#### Scenario: Existing presets are unchanged

- GIVEN the fork is built
- WHEN the `PresetOptions()` function returns available presets
- THEN `full-gentleman`, `ecosystem-only`, `minimal`, and `custom` are all present
- AND `cyber` is also present

---

### Requirement: TUI Preset Selection

The TUI preset selection screen MUST display `PresetCyber` as a selectable option. The description SHALL clearly indicate it is the cybersecurity edition and SHALL mention that destructive tools require manual confirmation.

#### Scenario: Cyber preset appears in TUI preset list

- GIVEN the TUI is launched
- WHEN the preset selection screen is displayed
- THEN "cyber" is listed as one of the selectable presets
- AND the description includes the text "Cybersecurity Edition"
- AND the description includes a warning about destructive tool confirmation

#### Scenario: Cyber preset description is distinct from other presets

- GIVEN the TUI preset screen is rendered
- WHEN comparing the cyber preset description to other presets
- THEN the cyber description is unique (not identical to full-gentleman, ecosystem-only, minimal, or custom)

#### Scenario: User can navigate to and select cyber preset

- GIVEN the TUI is on the preset selection screen
- WHEN the user navigates to the cyber preset and presses enter/select
- THEN `m.Selection.Preset` is set to `model.PresetCyber`
- AND the installer proceeds to the next screen

---

### Requirement: Cyber Preset Composition

Selecting `PresetCyber` MUST configure the following components compared to `PresetFullGentleman`:
- All full-gentleman components (Engram, SDD, persona, existing skills, docs)
- PLUS all 11 cybersecurity skills
- PLUS kali-mcp, shodan-mcp, and virustotal-mcp (with kali-mcp requiring explicit opt-in)
- PLUS gentleman-soc agent augmentation
- PLUS destructive-tool warning blocks in agent prompts

#### Scenario: Cyber preset includes full-gentleman base

- GIVEN `PresetCyber` is selected
- WHEN the installer resolves component dependencies
- THEN all components from `PresetFullGentleman` are included
- AND no full-gentleman component is removed

#### Scenario: Cyber preset adds 11 security skills

- GIVEN `PresetCyber` is selected
- WHEN the skill selection step computes the skill list
- THEN all 11 cybersecurity skills are present in the plan
- AND they are grouped under their respective categories (red-team, blue-team, soc)

#### Scenario: Cyber preset adds 3 MCP servers

- GIVEN `PresetCyber` is selected
- WHEN the MCP configuration step runs
- THEN shodan-mcp and virustotal-mcp are enabled by default
- AND kali-mcp is offered with a destructive capability warning and requires explicit opt-in

---

### Requirement: --cyber CLI Flag

The binary MUST accept a `--cyber` CLI flag that sets the active preset to `PresetCyber` at startup, equivalent to `--preset=cyber`.

#### Scenario: --cyber flag selects cyber preset

- GIVEN the binary is invoked with `gentle-ai-cyber --cyber`
- WHEN the application initializes
- THEN `Selection.Preset` is set to `model.PresetCyber`
- AND the TUI starts at the component review screen (skipping preset selection)

#### Scenario: --cyber flag is mutually exclusive with --preset

- GIVEN the binary is invoked with `gentle-ai-cyber --cyber --preset=minimal`
- WHEN the application initializes
- THEN an error is returned indicating the flags conflict
- AND the application exits with a non-zero status code

#### Scenario: --cyber flag is absent in base gentle-ai builds

- GIVEN the base gentle-ai binary (without cyber additions)
- WHEN `gentle-ai --help` is invoked
- THEN `--cyber` is NOT listed in the help output

---

### Requirement: Installer Does Not Break Existing Flows

The addition of the cyber preset MUST NOT alter the behavior of existing presets. Selecting `full-gentleman`, `ecosystem-only`, `minimal`, or `custom` MUST produce identical results to upstream gentle-ai.

#### Scenario: Full-gentleman preset unchanged

- GIVEN `PresetFullGentleman` is selected in the forked binary
- WHEN installation completes
- THEN the installed components match exactly what upstream gentle-ai would install for `full-gentleman`
- AND no cybersecurity skills, MCPs, or prompt warnings are installed

#### Scenario: Custom preset retains full manual selection capability

- GIVEN `PresetCustom` is selected in the forked binary
- WHEN the user manually selects components
- THEN cybersecurity skills and MCPs are available for manual selection
- AND destructive tool warnings are injected only if kali-mcp is manually enabled

---

### Requirement: Binary Identity

The forked binary SHALL be named `gentle-ai-cyber`. The Go module path SHALL be `github.com/rortizs/gentle-ai-cyber`. The binary description SHALL reference "cybersecurity edition".

#### Scenario: Binary compiles with correct name

- GIVEN the fork repository is checked out
- WHEN `go build ./cmd/gentle-ai-cyber` is executed
- THEN compilation succeeds without errors
- AND the output binary is named `gentle-ai-cyber`

#### Scenario: Go module path is set to fork identity

- GIVEN `go.mod` is read
- WHEN the module declaration is inspected
- THEN the module path is `github.com/rortizs/gentle-ai-cyber`
- AND it is different from `github.com/gentleman-programming/gentle-ai`
