# Cybersecurity BDD/TDD Testing Specification

## Purpose

Defines the security scenario test harness: validation pipeline tests, compliance mapping tests, skill behavior tests, and the definition of "passing" for the cybersecurity edition.

## Requirements

### Requirement: Skill Frontmatter Validation Tests

The system MUST include automated tests that validate every cybersecurity skill SKILL.md file has complete and valid YAML frontmatter. Tests SHALL be in `internal/catalog/cyber_skills_test.go`.

#### Scenario: All cyber skill frontmatter is valid

- GIVEN all cybersecurity skill files exist under `skills/`
- WHEN `TestCyberSkillsHaveRequiredFrontmatter` runs
- THEN every skill file parses with valid YAML frontmatter
- AND each frontmatter contains `name`, `description`, `license`, and `metadata` fields
- AND `metadata.author` is `gentleman-programming`
- AND `metadata.version` is a non-empty string

#### Scenario: Malformed frontmatter causes test failure

- GIVEN a cybersecurity skill SKILL.md has malformed YAML (missing closing `---`)
- WHEN `TestCyberSkillsHaveRequiredFrontmatter` runs
- THEN the test fails
- AND the failure message includes the file path and parse error

#### Scenario: Missing required field causes test failure

- GIVEN a cybersecurity skill SKILL.md has frontmatter without a `description` field
- WHEN `TestCyberSkillsHaveRequiredFrontmatter` runs
- THEN the test fails
- AND the failure message indicates the missing field and the file path

---

### Requirement: Skill Category Validation Tests

The system MUST include a test that verifies every cybersecurity skill entry in the catalog has a valid category value: one of `red-team`, `blue-team`, `soc`, or `compliance`.

#### Scenario: All catalog entries have valid categories

- GIVEN `MVPCyberSkills()` returns the cyber skill catalog
- WHEN `TestCyberSkillCategoriesAreValid` runs
- THEN every entry's `Category` field is one of: `red-team`, `blue-team`, `soc`, `compliance`
- AND no entry has an empty or unrecognized category

#### Scenario: Invalid category is detected

- GIVEN a cyber skill catalog entry has `Category: "pentest"` (not a valid category)
- WHEN `TestCyberSkillCategoriesAreValid` runs
- THEN the test fails
- AND the failure message identifies the skill with the invalid category

---

### Requirement: Destructive Tool Manifest Tests

The system MUST include a test that verifies every MCP manifest with `permission_tier: destructive` declares a non-empty `destructive_tools` array.

#### Scenario: Destructive manifests declare destructive tools

- GIVEN all MCP manifests under `mcps/`
- WHEN `TestDestructiveToolsRequireConfirmation` runs
- THEN every manifest with `permission_tier: destructive` has a `destructive_tools` array
- AND the array is non-empty
- AND every tool name in the array is a non-empty string

#### Scenario: Manifests without destructive_tools field fail the test

- GIVEN a manifest has `permission_tier: destructive` but no `destructive_tools` field
- WHEN `TestDestructiveToolsRequireConfirmation` runs
- THEN the test fails
- AND the failure message names the manifest file

---

### Requirement: Prowler Check Reference Format Tests

Any skill that references Prowler check IDs (in v1 or deferred for v2) MUST use check IDs matching the regex `^[a-z0-9_]+$`. The system SHALL include a test that scans skill content for Prowler check references and validates their format.

#### Scenario: Valid Prowler check ID references pass

- GIVEN a skill SKILL.md contains a reference like `check: iam_user_no_policies`
- WHEN `TestProwlerCheckIDFormat` runs and scans the content
- THEN the check ID `iam_user_no_policies` matches the expected format
- AND the test passes

#### Scenario: Invalid Prowler check ID format fails

- GIVEN a skill SKILL.md contains a reference like `check: IAM-User-No-Policies` (uppercase, hyphens)
- WHEN `TestProwlerCheckIDFormat` runs
- THEN the test fails
- AND the failure message identifies the malformed check ID and the file

#### Scenario: Skills without Prowler check references skip cleanly

- GIVEN a cybersecurity skill has no Prowler check ID references at all
- WHEN `TestProwlerCheckIDFormat` scans the content
- THEN the skill is skipped (not a failure)
- AND no false positive is reported

---

### Requirement: Cyber Preset Test

The system MUST include a test that verifies `PresetCyber` resolves to the correct set of components: full-gentleman base PLUS 11 cyber skills PLUS 3 cyber MCPs PLUS gentleman-soc agent augmentation.

#### Scenario: Cyber preset includes expected components

- GIVEN `PresetCyber` is used to compute the component set
- WHEN `componentsForPreset(model.PresetCyber)` is called
- THEN all full-gentleman components are present
- AND 11 cyber skills are in the skill list
- AND kali-mcp, shodan-mcp, and virustotal-mcp are in the MCP list
- AND gentleman-soc agent content is flagged for injection

#### Scenario: Cyber preset component count is deterministic

- GIVEN `PresetCyber` is used twice in succession
- WHEN the component set is computed each time
- THEN both results are identical (same components, same order)
- AND no random or environment-dependent components appear

---

### Requirement: No Regression on Existing Tests

All existing gentle-ai tests MUST continue to pass when the cybersecurity additions are present. The test suite SHALL be run with `go test ./...` and SHALL produce zero failures.

#### Scenario: Full test suite passes with cyber additions

- GIVEN the cybersecurity fork with all additions compiled
- WHEN `go test ./...` is executed
- THEN zero test failures are reported
- AND zero test panics occur
- AND all existing test files compile without modification

#### Scenario: Existing catalog tests still pass

- GIVEN the cybersecurity additions to `internal/catalog/`
- WHEN `go test ./internal/catalog/...` is executed
- THEN all pre-existing tests pass
- AND new cyber tests pass independently

---

### Requirement: BDD Scenario Traceability

Every Gherkin-style scenario in these specs SHALL be traceable to at least one automated test in the test suite. The mapping is documented in test comments referencing the spec requirement.

#### Scenario: Each spec scenario has a corresponding test

- GIVEN the spec document for any cybersecurity domain
- WHEN the test suite is audited
- THEN for each `#### Scenario:` heading in the specs, at least one test function exists
- AND the test function has a comment referencing the spec requirement name
