# Proposal: Gentle AI Cybersecurity Fork

## Intent

gentle-ai is a Go 1.25+ TUI installer that configures 13 AI agents with SDD workflows, persistent memory, skills, and MCP servers. It currently targets general-purpose software development. Security professionals (SOC analysts, pentesters, compliance engineers) need the same orchestration quality but with cybersecurity-focused skills, agents, MCP servers, and permission controls.

**gentle-ai-cyber** is a fork that repurposes the installer to ship a cybersecurity edition: security skills (from glitch-ai-toolkit), destructive-tool permission gates, Prowler compliance MCP integration, and a SOC orchestrator agent — while retaining full compatibility with gentle-ai's core (Engram, SDD, persona system, skill infrastructure).

The fork must remain maintainable: downstream from gentle-ai, not a parallel universe. Every feature that works in gentle-ai must work here too — plus security additions on top.

## Scope

### In Scope (v1)

1. **Fork Identity** — Repo, branding, description, positioning
2. **Security Skill Catalog** — 11 security skills from glitch-ai-toolkit, ported to gentle-ai's flat `skills/{name}/SKILL.md` format with frontmatter
3. **Cyber Category in Go Catalog** — New `Category: "red-team" | "blue-team" | "soc" | "compliance"` taxonomy in `internal/catalog/skills.go`
4. **3 Cyber MCP Servers** — kali-mcp, shodan-mcp, virustotal-mcp with `manifest.json` files adapted from glitch-ai-toolkit
5. **Permission Model for Destructive Tools** — Opt-in confirmation gate for tools that execute attacks (nmap, sqlmap, hydra, metasploit)
6. **Prowler MCP Integration** — External MCP reference (no bundling); installer checks for prowler-mcp availability and configures it if present
7. **Skill Format Decision** — Standardize on gentle-ai's simple frontmatter format; Prowler's AgentSkills with `auto_invoke` tables are NOT adopted in v1
8. **gentleman-soc Agent** — SOC orchestrator agent definition, integrated as an optional agent accessible from all 13 existing agents via skill loading
9. **Cyber Preset** — New `PresetCyber` in the TUI alongside full-gentleman, ecosystem-only, and minimal
10. **Installer Cyber Flag** — CLI flag `--cyber` that selects the cybersecurity edition, mapping to `PresetCyber`

### Out of Scope (v2)

- Prowler compliance metadata bundling (1345 check JSONs); v2 can include a sync command
- Prowler AgentSkills format adoption (auto_invoke tables) — requires a skill runtime that doesn't exist yet in gentle-ai
- Custom cybersecurity-specific agent IDs (beyond gentleman-soc which is a persona, not a new AgentID)
- Debian/RPM package distribution
- EULA or legal disclaimer flow for pentest tools (acknowledged in docs only)
- Centralized logging/audit trail for destructive tool usage
- Additional MCP servers beyond kali/shodan/virustotal (e.g., MikroTik, Proxmox — stay in glitch-ai-toolkit)
- Ongoing sync automation with Prowler releases (documented manual process for v1)

## Approach

### 1. Fork Identity

- **Repository**: `github.com/rortizs/gentle-ai-cyber` (or `gentle-ai-cyber` in the same org if preferred)
- **Go module**: `github.com/rortizs/gentle-ai-cyber`
- **Binary name**: `gentle-ai-cyber`
- **Description**: "Cybersecurity edition of Gentle AI — AI agent configurator for SOC analysts, pentesters, and compliance engineers"
- **Positioning**: gentle-ai is the upstream; gentle-ai-cyber imports from it via Go module replace directives or direct copy with attribution. The cyber fork adds security-specific layers without diverging the core.

### 2. Skill Format

**Recommendation: Use gentle-ai's simple frontmatter format universally.**

gentle-ai uses:
```yaml
---
name: skill-name
description: "One-line trigger description"
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---
```

Prowler's AgentSkills standard adds `auto_invoke` tables and `tool_schemas` — these require a runtime that gentle-ai doesn't have. Porting them would mean building a skill runtime first (significant scope creep).

**Hybrid approach**: v1 uses simple frontmatter. When gentle-ai gains a skill runtime (tracked as v2), we can add `auto_invoke` tables to Prowler-derived skills. Until then, skill trigger descriptions serve the same purpose declaratively.

### 3. Catalog Expansion Plan

**v1 Skills (11 security + all existing gentle-ai skills):**

| Category | Skills | Priority |
|----------|--------|----------|
| `red-team` | pentest-orchestrator, ai-pentesting-validation, exploit-chain-patterns, waf-detection-bypass | p0 |
| `blue-team` | detection-engineer, python-security | p0 |
| `soc` | malware-triage, specialized-file-analyzer, malware-dynamic-analysis, malware-report-writer | p0 |
| `compliance` | prowler-sdk-check | p1 (v2) |
| `sdd` (existing) | 10 SDD skills | p0 (inherited) |
| `workflow` (existing) | 7 workflow skills | p0 (inherited) |
| `testing` (existing) | go-testing | p0 (inherited) |

**v2 Skills (nice-to-have):**

| Category | Skills | Priority |
|----------|--------|----------|
| `compliance` | prowler-compliance, prowler-mcp, prowler-attack-paths-query | p1 |
| `red-team` | owasp-security-best-practices | p1 |
| `blue-team` | python-security (promoted from blue-team) | already p0 |

**Go code impact**: Each new skill adds ~3 lines to `internal/model/types.go` (SkillID const) + 1 line to `internal/catalog/skills.go` (Skill struct entry). 11 new skills = ~33 lines in types.go, 11 lines in skills.go. Minimal.

### 4. MCP Integration

**v1 bundles 3 MCP servers:**

| MCP | Type | Config | Security Notes |
|-----|------|--------|----------------|
| kali-mcp | system (external Kali) | Requires `kali-server-mcp` running on a Kali Linux host; MCP client bridges via stdio or SSH | All kali-mcp tools are **destructive** — subject to permission gates |
| shodan-mcp | npx | `npx -y @modelcontextprotocol/server-shodan` | Read-only API queries; no permission gate needed |
| virustotal-mcp | npx | `npx -y @modelcontextprotocol/server-virustotal` | Read-only API queries; no permission gate needed |

**Implementation**: Follow gentle-ai's existing pattern from `internal/components/mcp/context7.go`. Each MCP server gets:
- A `manifest.json` in `mcps/{name}/manifest.json` (adapted from glitch-ai-toolkit)
- A Go file in `internal/components/mcp/{name}.go` with overlay JSON per agent strategy
- Injection logic in `internal/components/mcp/inject.go` extended to handle cyber MCPs

**Prowler MCP**: Not bundled. The installer checks if `prowler-mcp` is installed (PyPI package `prowler-mcp`) and offers to configure it as an MCP server. If not found, it's skipped with a note. This avoids bundling 2.7MB of Prowler checks.

### 5. Permissions Model

**Approach: Opt-in destructive-tool allow-list with confirmation gates.**

The kali-mcp server provides tools like `nmap_scan`, `sqlmap_scan`, `hydra_attack`, `metasploit_run`. These can cause real damage on production systems.

**Design:**

```go
// internal/model/types.go — new types

type PermissionTier string

const (
    PermissionTierUnrestricted PermissionTier = "unrestricted" // shodan, virustotal, context7, engram
    PermissionTierDestructive  PermissionTier = "destructive"  // kali-mcp tools, metasploit
    PermissionTierRestricted   PermissionTier = "restricted"   // future: deletes, resets
)

type ToolPermission struct {
    MCPServer string          // e.g., "kali-mcp"
    Tools     []string        // e.g., ["nmap_scan", "sqlmap_scan", "hydra_attack"]
    Tier      PermissionTier
}
```

**Runtime behavior (in agents, not in the installer):**
1. Each MCP tool is tagged as `unrestricted` or `destructive` in a manifest
2. When the installer writes agent config, destructive MCP servers include a `/**_WARNING**` comment in the JSON listing the dangerous tools
3. Agent system prompts (injected by gentle-ai-cyber) include a block: "Before executing any tool marked DESTRUCTIVE, confirm with the user. Never auto-invoke nmap, sqlmap, hydra, or metasploit without explicit human approval."
4. This is a **soft gate** — it relies on agent prompt compliance, not a hard runtime block. A hard gate (requiring a Go-side proxy that intercepts MCP calls) is v2 scope.

**Manifest addition:**

```json
{
  "name": "kali-mcp",
  "permission_tier": "destructive",
  "destructive_tools": ["nmap_scan", "sqlmap_scan", "hydra_attack", "metasploit_run", "nikto_scan", "dirb_scan", "gobuster_scan", "wpscan_analyze", "john_crack", "enum4linux_scan"],
  "confirmation_required": true
}
```

### 6. Prowler Integration Depth

**Recommendation: External MCP reference. No bundling of Prowler checks.**

| Option | Pros | Cons | Decision |
|--------|------|------|----------|
| Bundle full SDK | Full offline capability | +2.7MB, stale within weeks, Python dependency | Rejected |
| Bundle metadata JSON only | Smaller, check metadata available | Still ~1.5MB, parsing complexity, no ability to run checks | Rejected |
| Reference external MCP | Zero bloat, always current, Prowler team maintains | Requires Prowler installed separately | **Accepted v1** |

**Implementation**: The TUI detects if `prowler-mcp` is available on PATH. If yes, it configures it as an MCP server for all agents. If no, it shows: "Prowler MCP not detected. Install with: `pip install prowler-mcp` or skip." The user can proceed without it.

v2: Add a `gentle-ai-cyber sync prowler` CLI command that pulls the latest Prowler check list and generates a static reference markdown for compliance skill use.

### 7. BDD/TDD Harness for Cyber

Testing strategy for cybersecurity edition:

**Three test categories:**

| Category | What | How |
|----------|------|-----|
| **Validation Pipeline** | Security skills produce non-hallucinated results | Skill SKILL.md includes `## Validation` section with expected output patterns; tests verify skill content parses and declares validation rules |
| **Compliance Mapping** | Prowler checks referenced by skill IDs actually exist | If prowler-mcp is available, integration test runs `prowler_hub_list_checks` and verifies referenced check IDs; if not available, test is skipped |
| **Skill Behavior** | Each skill frontmatter has required fields, trigger matches expected format | Unit tests in `internal/catalog/skills_test.go` extended to validate cyber skill entries |

**Concrete tests:**

```go
// internal/catalog/cyber_skills_test.go

func TestCyberSkillsHaveRequiredFrontmatter(t *testing.T) {
    // Every cyber skill SKILL.md must have: name, description, license, metadata.version
}

func TestDestructiveToolsRequireConfirmation(t *testing.T) {
    // Every MCP manifest with permission_tier=destructive must list destructive_tools
}

func TestCyberSkillCategoriesAreValid(t *testing.T) {
    // Category must be one of: red-team, blue-team, soc, compliance
}
```

**"Passing" definition:**
- All cyber skill SKILL.md files parse with valid frontmatter
- All cyber skill entries appear in `MVPCyberSkills()` catalog
- All destructive-tool manifests declare `permission_tier: destructive`
- No skill references a Prowler check ID format that doesn't match the regex `^[a-z0-9_]+$`

### 8. Agent Additions

**gentleman-soc** is the only new agent in v1.

It's NOT a new `AgentID` (which would require adding a new agent adapter). Instead, it's an **agent definition file** (`agents/gentleman-soc.md`) that gets installed as a system prompt augmentation for existing agents that opt into the cyber preset.

The gentleman-soc definition:
- Is a 583-line markdown file with PICERL pipeline, MITRE ATT&CK mapping, quality gates
- Gets injected into the agent's system prompt (following the `StrategyMarkdownSections` pattern already used by Claude Code)
- Knows WHEN to load each cyber skill (Phase 0→5 pipeline)
- Does NOT replace the existing agent; it augments it

**Implementation**: Add `gentleman-soc.md` to `agents/` directory. The installer copies it to the agent's config directory (e.g., `~/.claude/gentleman-soc.md`) and injects a reference in the system prompt. The skill-loading instructions reference `~/.claude/skills/{name}/SKILL.md` as they do in gentle-ai.

**Future agents** (v2 considerations):
- `pentest-orchestrator` — already exists as a skill; could become an agent definition in v2 if demand exists
- `compliance-auditor` — would orchestrate Prowler checks; requires prowler-mcp integration depth that v1 doesn't have

### 9. Installer Strategy

**Approach: Cybersecurity preset in the existing TUI.**

The gentle-ai TUI currently has 4 presets: `full-gentleman`, `ecosystem-only`, `minimal`, `custom`. We add `PresetCyber`:

```go
// internal/model/types.go
PresetCyber PresetID = "cyber"
```

**Preset definition:**
```
cyber = full-gentleman + 11 cyber skills + 3 cyber MCPs + gentleman-soc agent + destructive-tool warnings
```

The TUI flow:
1. User selects "cyber" preset OR runs `gentle-ai-cyber --cyber`
2. TUI shows: "Cybersecurity Edition: Installs SOC/malware analysis skills, pentest orchestration, offensive security tools, and Prowler compliance integration. ⚠ Destructive tools require manual confirmation."
3. Standard gentle-ai flow continues with cyber additions
4. At the MCP step, kali-mcp, shodan-mcp, virustotal-mcp are added
5. At the skill step, all 11 cyber skills are added
6. System prompt includes gentleman-soc context and destructive-tool warnings

**CLI flag**: `gentle-ai-cyber --cyber` is shorthand for `--preset=cyber`.

**Binary distribution**: Built from the same codebase with a build tag (`-tags=cyber`) that swaps the default preset list to include the cyber option. Alternatively, both presets live in the same binary and the cyber one is gated behind a feature flag. The simplest approach: **same binary, cyber preset always available**, selected via TUI or `--cyber` flag.

### 10. Maintenance & Sync

**Keeping Prowler checks current:**

v1: Manual process. Documentation includes:
```bash
# Update Prowler check reference
pip install --upgrade prowler
prowler-mcp --list-checks > docs/prowler-checks-reference.md
```

v2: Automated sync command:
```bash
gentle-ai-cyber sync prowler
# Downloads latest check metadata from Prowler PyPI package
# Updates local reference markdown
# Validates all skill check ID references
```

**Keeping glitch-ai-toolkit skills current:**

gentle-ai-cyber forks the skill content from glitch-ai-toolkit. Maintenance strategy:
1. Skills are **copied** (not symlinked) into `skills/{name}/SKILL.md`
2. Each skill gets gentle-ai frontmatter added (name, description, license, metadata)
3. A `scripts/sync-glitch-skills.sh` script can pull upstream changes from glitch-ai-toolkit
4. CI checks that skill frontmatter is valid after any sync

**Upstream sync with gentle-ai:**

The fork should rebase on gentle-ai releases, not cherry-pick. Release process:
1. gentle-ai tags a release (e.g., v2.3.0)
2. gentle-ai-cyber rebases its cyber additions on top
3. Conflicts are minimal since cyber additions are in separate files (catalog additions, skill directories, MCP configs)

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `internal/model/types.go` | Modified | Add 11 SkillID constants, PermissionTier type, ToolPermission type, PresetCyber |
| `internal/catalog/skills.go` | Modified | Add 11 cyber skill entries, MVPCyberSkills() function |
| `internal/catalog/skills_test.go` | Modified | Add cyber skill validation tests |
| `internal/catalog/agents.go` | Modified | No change — gentleman-soc is not a new AgentID |
| `internal/tui/screens/preset.go` | Modified | Add PresetCyber option with description |
| `internal/tui/screens/skill_picker.go` | Modified | Show cyber skills when cyber preset selected |
| `internal/components/mcp/` | Modified | Add kali_mcp.go, shodan_mcp.go, virustotal_mcp.go with overlay JSON |
| `internal/components/mcp/inject.go` | Modified | Handle cyber MCP injection per agent strategy |
| `skills/` | New | 11 new skill directories with SKILL.md files |
| `mcps/` | New | 3 MCP manifest directories (kali-mcp, shodan-mcp, virustotal-mcp) |
| `agents/gentleman-soc.md` | New | SOC orchestrator agent definition |
| `internal/catalog/cyber_skills_test.go` | New | Validation and compliance mapping tests |
| `cmd/gentle-ai-cyber/main.go` | Modified | Add --cyber flag (or build from same cmd/) |
| `AGENTS.md` | Modified | Add cybersecurity skill index section |
| `docs/` | New | Permission model docs, cyber preset docs, Prowler setup guide |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| kali-mcp tools cause real damage | High (if misused) | Destructive-tool permission tier with agent prompt warnings; documented disclaimers; opt-in installation |
| Skill content drifts from glitch-ai-toolkit upstream | Medium | sync script + CI validation; copied content is stable text, not generated |
| Prowler MCP API changes break integration | Low | Prowler MCP is separate package with its own versioning; v1 only configures it, doesn't embed |
| Gentle-ai upstream changes conflict with fork | Medium | Rebase strategy; cyber additions are additive (separate files) minimizing conflict surface |
| Agent ignores destructive-tool confirmation prompt | Medium | v1 is soft gate; v2 can add hard proxy. Document that the gate is prompt-based |
| 11 new skills bloat binary size | Low | Skills are markdown files copied at install time, not compiled into binary. Binary size impact: ~0 |
| TUI becomes crowded with preset options | Low | 5 presets is manageable; cyber preset has distinct description |

## Rollback Plan

1. Remove `PresetCyber` from preset list and `--cyber` CLI flag
2. Delete `skills/{red-team,blue-team,soc}/` directories
3. Delete `mcps/{kali-mcp,shodan-mcp,virustotal-mcp}/` directories
4. Revert `internal/model/types.go` to remove cyber SkillIDs and PermissionTier
5. Revert `internal/catalog/skills.go` to remove cyber entries
6. Delete `agents/gentleman-soc.md`
7. Delete `internal/catalog/cyber_skills_test.go`
8. Remove `internal/components/mcp/{kali_mcp,shodan_mcp,virustotal_mcp}.go`

All additions are additive. Rollback = delete the added files and revert the modified Go files to the upstream gentle-ai state.

## Dependencies

- glitch-ai-toolkit: source for security skills and MCP configs (commit-pinned)
- Prowler MCP: optional external dependency (pip package `prowler-mcp`)
- kali-mcp: requires a Kali Linux host running `kali-server-mcp` (documented setup guide)
- shodan-mcp: requires `SHODAN_API_KEY` env var
- virustotal-mcp: requires `VIRUSTOTAL_API_KEY` env var

## Success Criteria

- [ ] `go build ./cmd/gentle-ai-cyber` compiles without errors
- [ ] `go test ./internal/catalog/...` passes with cyber skill validation tests
- [ ] `PresetCyber` appears in the TUI preset selection and selects 11 cyber skills + 3 MCPs
- [ ] `--cyber` CLI flag selects the cyber preset
- [ ] All 11 cyber skill SKILL.md files have valid frontmatter (name, description, license, metadata.version)
- [ ] kali-mcp manifest declares `permission_tier: destructive` with 10 destructive tools listed
- [ ] shodan-mcp and virustotal-mcp manifests declare `permission_tier: unrestricted`
- [ ] Agent system prompt for cyber installations includes destructive-tool confirmation instructions
- [ ] Prowler MCP auto-detection works: skips gracefully if not installed, configures if available
- [ ] `gentleman-soc.md` is valid markdown that references the 6 SOC skills by correct path
- [ ] Rollback procedure succeeds: deleting cyber additions yields a clean gentle-ai build
- [ ] All existing gentle-ai tests continue to pass (no regressions)