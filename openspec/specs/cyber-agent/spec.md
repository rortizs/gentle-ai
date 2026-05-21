# Cybersecurity Agent Definition Specification

## Purpose

Defines the `gentleman-soc.md` agent definition file: its content requirements, injection strategy, skill references, and relationship to existing agent identities.

## Requirements

### Requirement: Gentleman-Soc Agent Definition File

The system MUST include a valid Markdown file at `agents/gentleman-soc.md` that defines the SOC orchestrator agent. This file SHALL be embedded as a Go asset and installed into the agent's configuration directory when the cyber preset is selected.

#### Scenario: Gentleman-soc file exists and is valid markdown

- GIVEN the cybersecurity fork repository
- WHEN `agents/gentleman-soc.md` is read
- THEN the file exists
- AND it contains valid markdown content with at least one markdown heading

#### Scenario: Gentleman-soc file is embedded as a Go asset

- GIVEN the codebase is compiled
- WHEN the `internal/assets/` package is inspected
- THEN `gentleman-soc.md` is accessible as an embedded asset
- AND `go:embed` directive references it correctly

---

### Requirement: PICERL Pipeline Structure

The gentleman-soc definition MUST describe a structured incident response pipeline using a phase-based framework. The pipeline SHALL cover at minimum: Preparation, Identification, Containment, Eradication, Recovery, and Lessons Learned (PICERL or equivalent).

#### Scenario: Pipeline phases are documented

- GIVEN the `agents/gentleman-soc.md` content
- WHEN the content is searched for pipeline phase names
- THEN at least 5 distinct phases are defined
- AND each phase maps to one or more cybersecurity skills from the catalog

#### Scenario: Each phase references specific skills

- GIVEN the gentleman-soc definition
- WHEN the pipeline section is examined
- THEN for each phase, at least one skill from the `soc` category is referenced by name
- AND referenced skill names match their catalog `Name` field (e.g., `malware-triage`, not `malware_triage`)

---

### Requirement: Skill Loading Instructions

The gentleman-soc definition MUST include instructions for WHEN to load each of the 6 SOC skills (`malware-triage`, `specialized-file-analyzer`, `malware-dynamic-analysis`, `malware-report-writer`, `detection-engineer`, `python-security`). The loading triggers SHALL align with the PICERL pipeline phases.

#### Scenario: All 6 SOC skills are referenced with loading conditions

- GIVEN `agents/gentleman-soc.md`
- WHEN all skill references are counted
- THEN all 6 SOC/blue-team skills are referenced by name
- AND each reference includes a trigger condition (e.g., "when user uploads a suspicious file")

#### Scenario: Skill references use correct file paths

- GIVEN the gentleman-soc definition
- WHEN skill file paths are mentioned
- THEN paths follow the pattern `skills/{name}/SKILL.md` (gentle-ai convention)
- AND no paths reference `~/.claude/` or any agent-specific prefix (the installer handles path resolution)

---

### Requirement: MITRE ATT&CK Mapping

The gentleman-soc definition SHOULD include a mapping from incident response phases to relevant MITRE ATT&CK tactics or techniques, providing analysts with a framework for threat classification.

#### Scenario: MITRE ATT&CK references are present

- GIVEN `agents/gentleman-soc.md`
- WHEN the content is searched for MITRE ATT&CK terminology
- THEN at least one MITRE ATT&CK tactic (e.g., `Initial Access`, `Execution`, `Exfiltration`) is referenced
- OR a dedicated MITRE ATT&CK mapping section exists

---

### Requirement: Quality Gates

The gentleman-soc definition MUST define quality gates — conditions that MUST be satisfied before moving from one pipeline phase to the next. Each gate SHALL be testable (verifiable by a human or automated check).

#### Scenario: Quality gates are defined between phases

- GIVEN the gentleman-soc pipeline structure
- WHEN the transition between any two consecutive phases is examined
- THEN at least one quality gate condition is described
- AND the condition is expressed as a checkable statement (not vague prose)

#### Scenario: Quality gates are not implementation instructions

- GIVEN any quality gate definition
- WHEN the text is examined
- THEN it describes WHAT must be verified, not HOW to verify it in code
- AND it does not reference specific APIs, tools, or programming languages

---

### Requirement: Agent Identity

The gentleman-soc definition MUST NOT create a new `AgentID`. It SHALL be an agent augmentation injected via the `StrategyMarkdownSections` or equivalent strategy, not a standalone agent adapter.

#### Scenario: No new AgentID constant for gentleman-soc

- GIVEN `internal/model/types.go`
- WHEN the `AgentID` constants are examined
- THEN no constant named `AgentSOC` or `AgentGentlemanSOC` exists
- AND the total number of AgentID constants matches the gentle-ai upstream count

#### Scenario: Gentleman-soc injects as markdown sections

- GIVEN the cyber preset is installed for an agent using `StrategyMarkdownSections` (e.g., Claude Code)
- WHEN the system prompt file is generated
- THEN gentleman-soc content appears as a markdown section delimited by `<!-- gentle-ai:gentleman-soc -->` markers
- AND it augments, not replaces, the existing agent persona

---

### Requirement: Cross-Agent Availability

The gentleman-soc definition MUST be available to all 13 agents supported by gentle-ai. It SHALL NOT be restricted to a specific agent or agent class.

#### Scenario: Cyber preset installs gentleman-soc for all selected agents

- GIVEN the cyber preset is installed with multiple agents selected
- WHEN the installation completes
- THEN gentleman-soc content is present in each agent's configuration directory
- AND the content is identical across agents (no agent-specific variations)
