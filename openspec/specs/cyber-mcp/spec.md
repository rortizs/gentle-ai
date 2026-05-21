# Cybersecurity MCP Configuration Specification

## Purpose

Defines the 3 MCP servers bundled in v1 (kali-mcp, shodan-mcp, virustotal-mcp), their manifest format, permission tier declarations, and injection strategy into agent configurations.

## Requirements

### Requirement: MCP Server Manifest Files

Each cybersecurity MCP server MUST have a `manifest.json` file at `mcps/{server-name}/manifest.json` containing at minimum: `name`, `type`, `config`, `permission_tier`, and `description`.

The `permission_tier` field SHALL be one of: `unrestricted`, `destructive`, or `restricted`.

#### Scenario: All 3 MCP manifests exist with required fields

- GIVEN the cybersecurity fork is built
- WHEN MCP manifests are loaded from `mcps/`
- THEN `kali-mcp/manifest.json`, `shodan-mcp/manifest.json`, and `virustotal-mcp/manifest.json` all exist
- AND each manifest contains `name`, `type`, `config`, `permission_tier`, and `description` fields

#### Scenario: Manifests are valid JSON

- GIVEN any cybersecurity MCP manifest file
- WHEN the file is parsed as JSON
- THEN parsing succeeds without error
- AND the root object contains all required top-level keys

---

### Requirement: MCP Permission Tier Classification

The kali-mcp server MUST declare `permission_tier: destructive`. Shodan-mcp and virustotal-mcp MUST declare `permission_tier: unrestricted`.

#### Scenario: Kali-mcp is classified as destructive

- GIVEN the kali-mcp manifest
- WHEN the `permission_tier` field is read
- THEN its value is `"destructive"`

#### Scenario: Shodan-mcp is classified as unrestricted

- GIVEN the shodan-mcp manifest
- WHEN the `permission_tier` field is read
- THEN its value is `"unrestricted"`

#### Scenario: Virustotal-mcp is classified as unrestricted

- GIVEN the virustotal-mcp manifest
- WHEN the `permission_tier` field is read
- THEN its value is `"unrestricted"`

---

### Requirement: Destructive Tool Listing

For any MCP manifest with `permission_tier: destructive`, the manifest MUST include a `destructive_tools` array listing every tool that requires user confirmation before execution.

#### Scenario: Kali-mcp lists all destructive tools

- GIVEN the kali-mcp manifest with `permission_tier: destructive`
- WHEN the `destructive_tools` field is read
- THEN the array contains at minimum: `nmap_scan`, `sqlmap_scan`, `hydra_attack`, `metasploit_run`
- AND each tool name matches an actual tool provided by the kali-mcp server

#### Scenario: Unrestricted manifests have no destructive_tools array

- GIVEN the shodan-mcp manifest with `permission_tier: unrestricted`
- WHEN the `destructive_tools` field is read
- THEN it is either absent or an empty array

#### Scenario: Unrestricted manifests must not list tools as destructive

- GIVEN the virustotal-mcp manifest with `permission_tier: unrestricted`
- WHEN the `destructive_tools` field is examined (if present)
- THEN it is empty
- OR the field is absent entirely

---

### Requirement: MCP Injection Strategy

Each cybersecurity MCP server MUST follow the existing agent-specific MCP injection pattern: a Go file at `internal/components/mcp/{server-name}.go` that provides overlay JSON per agent strategy, wired through `internal/components/mcp/inject.go`.

#### Scenario: Cyber MCP injection files exist

- GIVEN the fork is compiled
- WHEN the `internal/components/mcp/` package is inspected
- THEN `kali_mcp.go`, `shodan_mcp.go`, and `virustotal_mcp.go` exist
- AND each file exports overlay JSON for at least the `StrategySeparateMCPFiles` strategy

#### Scenario: Inject handles cyber MCPs alongside existing MCPs

- GIVEN the cyber preset is selected
- WHEN MCP injection runs
- THEN context7 and engram MCPs are injected as normal
- AND kali-mcp, shodan-mcp, and virustotal-mcp are injected for all agents in the preset

---

### Requirement: Prowler MCP External Detection

The installer MUST detect whether `prowler-mcp` is available on the system PATH. If found, it SHALL offer to configure it as an MCP server. If not found, it SHALL inform the user and continue without error.

The prowler-mcp SHALL NOT be bundled or embedded in the binary.

#### Scenario: Prowler-mcp detected on PATH

- GIVEN `prowler-mcp` is installed and available on PATH
- WHEN the installer reaches the MCP configuration step
- THEN the user is offered to include prowler-mcp in the configuration
- AND the installer shows the command that was used to detect it

#### Scenario: Prowler-mcp not detected on PATH

- GIVEN `prowler-mcp` is NOT installed
- WHEN the installer reaches the MCP configuration step
- THEN a message informs the user: "Prowler MCP not detected. Install with: `pip install prowler-mcp` or skip."
- AND installation continues without error
- AND prowler-mcp is NOT configured

#### Scenario: Prowler-mcp detection failure does not block installation

- GIVEN PATH lookup for `prowler-mcp` fails for any reason (e.g., permission error)
- WHEN the installer attempts detection
- THEN it treats the result as "not found" and continues
- AND a diagnostic message is logged for debugging
