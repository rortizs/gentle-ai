# Cybersecurity Permission Model Specification

## Purpose

Defines the soft-gate permission model for destructive cybersecurity tools: how tools are classified, what warnings agents receive, and how the confirmation flow works in the generated agent system prompts.

## Requirements

### Requirement: Permission Tier Type Definitions

The system MUST define a `PermissionTier` string type with exactly three constant values: `unrestricted`, `destructive`, and `restricted`. It SHALL also define a `ToolPermission` struct containing `MCPServer`, `Tools`, and `Tier` fields.

#### Scenario: PermissionTier constants are defined

- GIVEN the `internal/model/types.go` file is compiled
- WHEN the `PermissionTier` type is inspected
- THEN `PermissionTierUnrestricted` has value `"unrestricted"`
- AND `PermissionTierDestructive` has value `"destructive"`
- AND `PermissionTierRestricted` has value `"restricted"`

#### Scenario: ToolPermission struct is defined

- GIVEN the codebase is compiled
- WHEN the `ToolPermission` type is inspected
- THEN it has fields `MCPServer string`, `Tools []string`, and `Tier PermissionTier`

---

### Requirement: Destructive Tool Warning in Agent System Prompts

When the cybersecurity preset is selected, the installer MUST inject a destructive-tool warning block into every agent's system prompt. The block SHALL instruct the agent to require explicit human confirmation before invoking any tool tagged as `destructive`.

#### Scenario: Cyber preset injects destructive tool warning

- GIVEN the user selects the `cyber` preset during installation
- WHEN system prompts are generated for any agent
- THEN each system prompt includes a block containing the text "Before executing any tool marked DESTRUCTIVE"
- AND the block instructs the agent to "confirm with the user" before invoking such tools

#### Scenario: Non-cyber presets do not inject destructive warnings

- GIVEN the user selects `full-gentleman`, `ecosystem-only`, `minimal`, or `custom`
- WHEN system prompts are generated
- THEN no destructive-tool warning block is present
- AND existing gentle-ai prompt behavior is unchanged

#### Scenario: Warning block references specific tools from kali-mcp

- GIVEN the destructive tool warning is injected
- WHEN the block content is examined
- THEN it mentions at minimum: `nmap_scan`, `sqlmap_scan`, `hydra_attack`, and `metasploit_run` as examples requiring confirmation
- AND it states these are from the `kali-mcp` server

---

### Requirement: MCP JSON Warning Comments

For agents that use `StrategySeparateMCPFiles`, the generated MCP JSON file for kali-mcp MUST include a `/**_WARNING**` comment at the top listing the destructive tools and requiring manual confirmation.

#### Scenario: Kali-mcp JSON includes warning comment

- GIVEN the cyber preset is installed for an agent using `StrategySeparateMCPFiles` (e.g., Claude Code)
- WHEN the kali-mcp MCP JSON file is written
- THEN the file begins with a comment containing `_WARNING_`
- AND the comment lists the destructive tools from the manifest
- AND the comment states these tools require human confirmation

#### Scenario: Unrestricted MCP JSON files have no warning comment

- GIVEN the cyber preset is installed
- WHEN shodan-mcp or virustotal-mcp MCP JSON files are written
- THEN no `_WARNING_` comment is present at the top of the file
- OR the comment is informational (non-warning) in nature

---

### Requirement: Confirmation is a Soft Gate

The destructive-tool confirmation MUST be a soft gate — enforced through agent prompt compliance, NOT through a runtime proxy or hard block. The system SHALL NOT intercept, block, or proxy MCP tool calls at the Go runtime level in v1.

#### Scenario: No runtime interception of MCP calls

- GIVEN kali-mcp is configured for an agent
- WHEN the agent invokes `nmap_scan` without prior human confirmation
- THEN the tool call proceeds normally at the MCP transport level
- AND no Go-side error, panic, or block is triggered

#### Scenario: Soft gate is documented in user-facing docs

- GIVEN the cybersecurity fork documentation is generated
- WHEN the permissions documentation is read
- THEN it explicitly states the gate is "prompt-based (soft gate)"
- AND it notes that a hard proxy gate is planned for v2

---

### Requirement: Permission Tier Determination from Manifest

At MCP injection time, the permission tier for each server MUST be read from its `manifest.json` file. The `permission_tier` field SHALL be the single source of truth; it MUST NOT be hardcoded in Go.

#### Scenario: Manifest is the source of truth for permission tier

- GIVEN the kali-mcp `manifest.json` has `"permission_tier": "destructive"`
- WHEN MCP injection reads the permission tier
- THEN the returned tier matches the manifest value
- AND the Go code does not contain a hardcoded `PermissionTierDestructive` for kali-mcp

#### Scenario: Changing manifest changes behavior without recompilation

- GIVEN a user modifies the `permission_tier` in `mcps/kali-mcp/manifest.json` to `"unrestricted"`
- WHEN the agent system prompt is regenerated
- THEN the destructive tool warning is NOT present for kali-mcp
- AND no Go recompilation is required

---

### Requirement: No MCP Servers Are Force-Enabled

Destructive MCP servers (kali-mcp) MUST be opt-in. The installer SHALL list them as available but SHALL NOT enable them by default. The user MUST explicitly approve enabling destructive servers.

#### Scenario: Kali-mcp requires explicit user approval

- GIVEN the user selects the cyber preset
- WHEN the installer reaches the MCP configuration screen
- THEN kali-mcp is listed with a warning label indicating destructive capability
- AND it is NOT pre-selected (requires user toggle)
- AND the user must explicitly confirm before kali-mcp is configured

#### Scenario: Unrestricted MCPs are pre-selected in cyber preset

- GIVEN the user selects the cyber preset
- WHEN the installer reaches the MCP configuration screen
- THEN shodan-mcp and virustotal-mcp are pre-selected (enabled by default)
- AND the user may deselect them if desired
