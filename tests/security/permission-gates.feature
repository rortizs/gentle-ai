Feature: Permission Gates for Destructive Tools

  As a Gentleman AI user
  I want destructive security tools to require explicit confirmation
  So that no accidental damage occurs to production systems

  Background:
    Given the cyber preset is active
    And kali-mcp manifest declares permission_tier: "destructive"
    And kali-mcp manifest declares a non-empty destructive_tools array

  Scenario: Destructive tool warning is injected for separate-file agents
    Given an agent using StrategySeparateMCPFiles (e.g., Claude Code)
    When kali-mcp overlay JSON is generated
    Then the JSON file is prefixed with a warning comment
    And the warning lists all destructive tools
    And the warning uses the format "/**_WARNING**/"

  Scenario: Destructive tool warning is injected in system prompt
    Given the cyber preset is active
    When the system prompt is generated
    Then a destructive-warning section is injected
    And the section is delimited by <!-- gentle-ai-cyber:destructive-warning --> markers
    And the warning lists all tools from the manifest destructive_tools array

  Scenario: User confirms destructive tool execution
    Given a destructive tool (nmap_scan) is requested
    And the warning has been displayed
    When the user explicitly confirms
    Then the tool executes normally
    And the confirmation is logged

  Scenario: User declines destructive tool execution
    Given a destructive tool (metasploit_run) is requested
    And the warning has been displayed
    When the user does not confirm
    Then the tool is blocked
    And the agent explains why the tool requires confirmation

  Scenario: Unrestricted tools do not require confirmation
    Given shodan-mcp has permission_tier: "unrestricted"
    And virustotal-mcp has permission_tier: "unrestricted"
    When tools from shodan-mcp or virustotal-mcp are requested
    Then no destructive warning is displayed
    And the tools execute without confirmation prompt

  Scenario: Manifest is the single source of truth for permission tier
    Given the permission_tier field in manifest.json
    When the permission tier is checked
    Then the value is read from the manifest, not hardcoded in Go
    And changing the manifest changes the permission behavior without recompilation

  Scenario: All kali-mcp destructive tools are listed
    Given kali-mcp manifest is parsed
    When the destructive_tools array is inspected
    Then it contains: nmap_scan
    And it contains: sqlmap_scan
    And it contains: hydra_attack
    And it contains: metasploit_run
    And it contains: nikto_scan
    And it contains: dirb_scan
    And it contains: gobuster_scan
    And it contains: wpscan_analyze
    And it contains: john_crack
    And it contains: enum4linux_scan
