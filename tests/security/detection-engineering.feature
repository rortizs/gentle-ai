Feature: Blue Team Detection Engineering

  As a blue team operator using Gentleman AI
  I want to create detection rules from analysis findings
  So that threats can be detected and responded to proactively

  Background:
    Given malware analysis findings are available
    And the cyber preset is active

  Scenario: Detection engineer creates Sigma rules
    Given malware analysis findings with behavioral indicators
    When the detection-engineer skill runs
    Then Sigma detection rules are generated
    And the rules target the identified TTPs
    And the rules include proper logsource configuration
    And the rules follow Sigma format conventions

  Scenario: Detection engineer creates Suricata rules
    Given malware analysis findings with network indicators
    When the detection-engineer skill runs
    Then Suricata IDS rules are generated
    And the rules match the identified network patterns
    And the rules include appropriate severity levels
    And the rules include reference metadata

  Scenario: IOCs are defanged for safe sharing
    Given raw IOCs from malware analysis
    When IOCs are prepared for sharing
    Then URLs are defanged (hXXp instead of http)
    And IP addresses are optionally defanged
    And file hashes remain intact (not defanged)
    And the defanged IOCs are safe to share in reports

  Scenario: Detection rules map to MITRE ATT&CK
    Given detection rules generated from analysis
    When rules are validated
    Then each rule maps to at least one MITRE ATT&CK technique
    And the technique IDs are in valid format (T####.###)
    And the mapping is documented in the rule metadata

  Scenario: Complete detection engineering workflow
    Given malware analysis findings
    When the detection-engineer skill creates Sigma rules
    And the detection-engineer skill creates Suricata rules
    And the rules are validated against MITRE ATT&CK
    Then a complete detection rule package is produced
    And the package covers identified TTPs
    And the rules are ready for deployment to SIEM/IDS
