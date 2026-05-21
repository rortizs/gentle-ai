# Cybersecurity Skills Catalog Specification

## Purpose

Defines the 11 v1 cybersecurity skills that MUST be ported from glitch-ai-toolkit into gentle-ai's skill catalog, their category taxonomy, frontmatter format, and registration in the Go catalog.

## Requirements

### Requirement: Skill Identity and Categorization

The system MUST define exactly 11 cybersecurity skills across 4 security categories. Each skill SHALL be identified by a unique `SkillID` constant in `internal/model/types.go` and registered as a catalog entry in `internal/catalog/skills.go`.

The security category taxonomy SHALL consist of: `red-team`, `blue-team`, `soc`, and `compliance`.

#### Scenario: All 11 v1 skills are registered in the catalog

- GIVEN the cybersecurity fork is built
- WHEN `MVPCyberSkills()` is called
- THEN exactly 11 skill entries are returned
- AND each entry has a non-empty ID, Name, Category, and Priority

#### Scenario: Red team skills exist with correct categorization

- GIVEN the catalog is loaded
- WHEN the `red-team` category is queried
- THEN `pentest-orchestrator`, `ai-pentesting-validation`, `exploit-chain-patterns`, and `waf-detection-bypass` are present
- AND each is tagged with Category `red-team` and Priority `p0`

#### Scenario: Blue team skills exist with correct categorization

- GIVEN the catalog is loaded
- WHEN the `blue-team` category is queried
- THEN `detection-engineer` and `python-security` are present
- AND each is tagged with Category `blue-team` and Priority `p0`

#### Scenario: SOC skills exist with correct categorization

- GIVEN the catalog is loaded
- WHEN the `soc` category is queried
- THEN `malware-triage`, `specialized-file-analyzer`, `malware-dynamic-analysis`, and `malware-report-writer` are present
- AND each is tagged with Category `soc` and Priority `p0`

#### Scenario: Each skill has a unique SkillID that does not collide with existing IDs

- GIVEN the existing 20 skill IDs in `internal/model/types.go`
- WHEN the 11 cyber SkillIDs are added
- THEN no cyber SkillID string duplicates an existing SkillID string
- AND each cyber SkillID has a corresponding Go constant in the `SkillID` type

---

### Requirement: Skill Frontmatter Format

Every cybersecurity skill SKILL.md file MUST contain valid YAML frontmatter with exactly four required fields: `name`, `description`, `license`, and `metadata` (containing `author` and `version`).

The frontmatter SHALL follow gentle-ai's simple format — no `auto_invoke` tables, no `tool_schemas`, no execution directives.

#### Scenario: Valid frontmatter exists on every cyber skill file

- GIVEN a cybersecurity skill file at `skills/{name}/SKILL.md`
- WHEN the file is parsed
- THEN the frontmatter contains a `name` field matching the skill name
- AND the frontmatter contains a `description` field with a one-line trigger description
- AND the frontmatter contains `license: Apache-2.0`
- AND `metadata.author` is `gentleman-programming`
- AND `metadata.version` is a semantic version string (e.g., `"1.0"`)

#### Scenario: Skill file has no auto_invoke or tool_schemas blocks

- GIVEN any cybersecurity skill SKILL.md
- WHEN the file content is examined
- THEN it contains no `auto_invoke` YAML key
- AND it contains no `tool_schemas` YAML key

#### Scenario: Skill frontmatter parsing fails gracefully

- GIVEN a cybersecurity skill SKILL.md with malformed YAML frontmatter
- WHEN the catalog validation runs
- THEN an error is returned indicating the file path and parse failure reason

---

### Requirement: Skill Content Completeness

Each cybersecurity skill SKILL.md MUST include a `## Validation` section that declares expected output patterns for testing purposes. The validation section SHALL define at least one expected behavior or output pattern the skill should produce when loaded by an agent.

#### Scenario: Skill file includes a validation section

- GIVEN any cybersecurity skill SKILL.md
- WHEN the file is read
- THEN a `## Validation` section header exists
- AND the section contains at least one declarative statement describing expected skill behavior

#### Scenario: Skill file has minimum viable content

- GIVEN a cybersecurity skill SKILL.md
- WHEN the file content beyond frontmatter is examined
- THEN at least one markdown heading exists (e.g., `## Overview` or `## Workflow`)
- AND the total content length (excluding frontmatter) exceeds 100 characters

---

### Requirement: Owasp Security Best Practices Deferred to v2

The `owasp-security-best-practices` skill SHALL NOT be included in v1. It MUST be tracked as a v2 skill and MUST NOT appear in the `MVPCyberSkills()` catalog.

#### Scenario: Owasp security best practices is absent from v1 catalog

- GIVEN the v1 cybersecurity catalog is built
- WHEN `MVPCyberSkills()` is called
- THEN `owasp-security-best-practices` is NOT in the returned list

---

### Requirement: No Regression on Existing Skills

The addition of cybersecurity skills MUST NOT alter or remove any existing skill entry in `internal/catalog/skills.go`. The existing `MVPSkills()` function SHALL continue to return all 20 existing skill entries unchanged.

#### Scenario: Existing MVP skills are unchanged

- GIVEN the fork is built with cybersecurity additions
- WHEN `MVPSkills()` is called
- THEN exactly 20 skill entries are returned
- AND all existing SkillIDs (sdd-init through work-unit-commits) are present
- AND their Category and Priority values match the upstream gentle-ai definitions
