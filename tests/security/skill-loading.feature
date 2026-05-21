Feature: Cyber Skill Loading and Categorization

  As a Gentleman AI user selecting the cyber preset
  I want all cybersecurity skills to be properly loaded and categorized
  So that the right skills are available for security assessments

  Background:
    Given the cyber preset is selected

  Scenario: Cyber preset includes all 10 cyber skills
    When SkillsForPreset(PresetCyber) is called
    Then all MVP skills are included
    And all 10 cyber skills are included
    And the total skill count is deterministic

  Scenario: Red-team skills are correctly categorized
    Given MVPCyberSkills() returns the cyber skill catalog
    When skill categories are checked
    Then pentest-orchestrator has category "red-team"
    And ai-pentesting-validation has category "red-team"
    And exploit-chain-patterns has category "red-team"
    And waf-detection-bypass has category "red-team"
    And there are exactly 4 red-team skills

  Scenario: Blue-team skills are correctly categorized
    Given MVPCyberSkills() returns the cyber skill catalog
    When skill categories are checked
    Then detection-engineer has category "blue-team"
    And python-security has category "blue-team"
    And there are exactly 2 blue-team skills

  Scenario: SOC skills are correctly categorized
    Given MVPCyberSkills() returns the cyber skill catalog
    When skill categories are checked
    Then malware-triage has category "soc"
    And specialized-file-analyzer has category "soc"
    And malware-dynamic-analysis has category "soc"
    And malware-report-writer has category "soc"
    And there are exactly 4 SOC skills

  Scenario: No skill ID collisions between MVP and cyber skills
    Given all MVP skill IDs
    And all cyber skill IDs
    When skill IDs are compared
    Then no cyber skill ID duplicates any MVP skill ID
    And all skill IDs are unique across the entire catalog

  Scenario: Cyber skill frontmatter is valid
    Given all cyber skill SKILL.md files exist
    When frontmatter is parsed
    Then each skill has a valid name field matching its directory
    And each skill has a quoted description field
    And each skill has a license field
    And each skill has metadata with author and version
    And each description contains "Trigger:" substring
    And no skill has auto_invoke or tool_schemas in frontmatter

  Scenario: Non-cyber presets are unchanged
    Given the full-gentleman preset
    When SkillsForPreset(PresetFullGentleman) is called
    Then no cyber skills are included
    And the skill count matches the original MVP count

  Scenario: Cyber preset skill count is deterministic
    Given PresetCyber is used to compute the component set
    When the skill list is computed twice
    Then both results have the same skill count
    And both results contain the same skill IDs
    And the order is deterministic
