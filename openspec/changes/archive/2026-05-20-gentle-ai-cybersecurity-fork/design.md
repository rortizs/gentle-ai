# Design: Gentle AI Cybersecurity Fork

## Technical Approach

The cybersecurity fork adds a **additive layer** on top of gentle-ai's existing installer: new SkillID constants, a cyber preset, security-specific skill files, three MCP server configs, destructive-tool soft gates, a SOC orchestrator agent definition, and validation tests. The design follows gentle-ai's established patterns exactly — flat `skills/{name}/SKILL.md` directories, `internal/components/mcp/{name}.go` overlay JSON files, `internal/catalog/skills.go` catalog entries, and `internal/components/skills/presets.go` preset resolution. No existing code paths are modified; all additions are parallel to the existing MVP infrastructure.

References: `cyber-skills/spec.md` (skill catalog), `cyber-mcp/spec.md` (MCP manifests and injection), `cyber-permissions/spec.md` (soft gate model).

---

## Architecture Decisions

### Decision: Flat skill directory layout (not hierarchical)

**Choice**: Security skills live in flat `skills/{name}/SKILL.md` directories, matching gentle-ai's existing layout.
**Alternatives considered**: Hierarchical grouping like `skills/red-team/pentest-orchestrator/SKILL.md` to visually separate categories.
**Rationale**: gentle-ai's asset embedding (`internal/assets/skills/`) and the `Inject()` function both expect flat directories where the skill name equals the directory name. The skill catalog `Category` field handles logical grouping. Adding hierarchy would require changes to the embed directive, the walk logic in `inject.go`, and the frontmatter validation test — all for no functional benefit. The flat layout also keeps cross-references (e.g., `malware-triage` referencing `specialized-file-analyzer`) simpler since every skill is at the same path depth.

### Decision: Separate `mvpcyberSkills` function (not extending `mvpSkills`)

**Choice**: A new `MVPCyberSkills()` function in `internal/catalog/skills.go` returns only the 11 cyber skills. A `CyberSkillsForPreset()` helper composes them with `MVPSkills()` when the cyber preset is selected.
**Alternatives considered**: Adding all 31 skills (20 MVP + 11 cyber) to the existing `mvpSkills` slice and gating with a feature flag.
**Rationale**: The existing `MVPSkills()` function is called by `normalizeSkills()` in `validate.go` and the TUI skill picker. Both rely on `MVPSkills()` returning exactly the skills available for the selected preset. Merging cyber skills into `mvpSkills` would break the standard `full-gentleman` preset — it would suddenly include pentest skills. Keeping them separate means: (1) existing presets are unchanged, (2) `PresetCyber` composes MVP + cyber, (3) the TUI skill picker shows cyber skills only when the cyber preset is active, (4) validation test `TestMVPSkillsCoverAllPresetSkills` continues to pass unchanged.

### Decision: Permission tier is manifest-driven, not hardcoded

**Choice**: Each MCP manifest JSON (`mcps/{name}/manifest.json`) contains a `permission_tier` field. The Go runtime reads this at injection time to determine whether to inject the destructive-tool warning into the system prompt.
**Alternatives considered**: Hardcoding `PermissionTier` in Go constants per MCP server name (e.g., a switch statement `case "kali-mcp": return PermissionTierDestructive`).
**Rationale**: The cyber-permissions spec requirement 5 explicitly states: "The `permission_tier` field SHALL be the single source of truth; it MUST NOT be hardcoded in Go." This allows users to reclassify MCP servers without recompilation — changing the manifest file is sufficient. Hardcoding would be simpler code but would violate the spec and prevent runtime customization.

### Decision: gentleman-soc is an agent definition file, not a new AgentID

**Choice**: `gentleman-soc.md` is a markdown file stored in `agents/gentleman-soc.md` and installed alongside system prompts as an augmentation file. It does NOT add `AgentGentlemanSOC` to `internal/catalog/agents.go`.
**Alternatives considered**: Adding a new `AgentGentlemanSOC` constant and adapter, treating it as a full agent with its own config paths.
**Rationale**: gentleman-soc is an orchestration persona that works *through* existing agents (Claude Code, OpenCode, etc.), not a standalone agent. Making it an AgentID would require an adapter implementation, config paths, MCP strategy, system prompt strategy — all for something that is conceptually a skill-loading instruction set. The definition file approach follows the same pattern gentle-ai uses for system prompt augmentation (`StrategyMarkdownSections`). The SOC pipeline phases (Recon → Triage → Analysis → Report) are skill-loading recipes, not agent behaviors.

### Decision: Soft gate via agent prompt (not runtime proxy)

**Choice**: Destructive-tool warnings are injected into the agent's system prompt as a `<!-- gentle-ai-cyber:destructive-warning -->` section. No Go runtime intercepts, blocks, or proxies MCP calls.
**Alternatives considered**: A Go-side MCP proxy that intercepts `tools/call` requests for destructive tools and requires interactive confirmation via the TUI.
**Rationale**: v1 scope explicitly excludes a hard proxy. A proxy would require: (1) a local MCP gateway server between the agent and kali-mcp, (2) a TUI confirmation prompt that pauses the agent, (3) handling of timeouts and agent retries. This is a significant feature that deserves its own design cycle. The soft gate approach: LLMs reliably follow instructions in their system prompt (especially Claude, GPT-4, and similar models), the warning is visible in the system prompt for user auditing, and the spec explicitly requires this approach.

### Decision: Prowler MCP as external detection (not bundled)

**Choice**: The installer runs `exec.LookPath("prowler-mcp")` at MCP configuration time. If found, it offers to configure it. If not found, it shows an installation hint and continues.
**Alternatives considered**: Bundling `prowler-mcp` as a Python wheel inside the Go binary, embedding its check metadata as static JSON.
**Rationale**: Prowler's checks are a living dataset (1345+ checks) that update frequently. Bundling would add 2.7MB+ of stale data. Embedding would create a sync maintenance burden. External detection means: (1) always-current Prowler, (2) zero binary bloat, (3) Prowler team maintains their own MCP. This aligns with the proposal's v1 scope.

---

## Data Flow

### User Request Through Agent → Skill → MCP → Result

```
┌──────────┐    ┌──────────────┐    ┌──────────────┐    ┌───────────┐
│  User    │───→│ Agent System  │───→│ Skill Load   │───→│ Skill     │
│ Request  │    │ Prompt+Warning│    │ (malware-    │    │ SKILL.md  │
│          │    │              │    │  triage)      │    │ Execution │
└──────────┘    └──────┬───────┘    └──────────────┘    └─────┬─────┘
                       │                                       │
                       │  ┌──────────────────────────┐        │
                       │  │ Destructive-Tool Warning  │        │
                       │  │ (soft gate: confirm before│        │
                       │  │  calling kali-mcp tools)  │        │
                       │  └──────────────────────────┘        │
                       │                                       ▼
                ┌──────┴───────┐    ┌──────────────┐    ┌───────────┐
                │ MCP Client   │───→│ kali-mcp /   │───→│ Tool      │
                │ (Agent side) │    │ shodan-mcp / │    │ Execution │
                │              │    │ virustotal   │    │ (nmap,etc)│
                └──────────────┘    └──────────────┘    └───────────┘
```

### gentleman-soc Orchestration Pipeline

```
gentleman-soc.md (583-line agent definition)
    │
    ├── Phase 0: Reconnaissance ──→ Load pentest-orchestrator skill
    ├── Phase 1: Triage ──────────→ Load malware-triage skill
    ├── Phase 2: Analysis ────────→ Load specialized-file-analyzer
    │                                  └── malware-dynamic-analysis
    ├── Phase 3: Detection ───────→ Load detection-engineer skill
    └── Phase 4: Reporting ───────→ Load malware-report-writer skill
```

The gentleman-soc definition file tells the agent WHEN to load each skill based on the PICERL (Prepare → Identify → Control → Eradicate → Recover → Lessons) phase. It does not load all skills upfront — it uses the existing skill-loading mechanism referenced in the AGENTS.md `available_skills` block.

### Permission Gate Data Flow

```
MCP Manifest (manifest.json)
    │
    │  permission_tier: "destructive"
    │  destructive_tools: ["nmap_scan", "sqlmap_scan", ...]
    │
    ▼
┌──────────────────┐     ┌───────────────────────────┐
│ mcp.Injection    │────→│ Agent System Prompt        │
│ reads manifest   │     │ <!-- gentle-ai-cyber:       │
│ at inject time   │     │   destructive-warning --> │
│                  │     │ "Before executing any tool  │
│                  │     │  marked DESTRUCTIVE..."     │
└──────────────────┘     └───────────────────────────┘
```

The manifest is the single source of truth. The Go code reads it once during installation, injects the appropriate warning into the system prompt, and then has no further runtime involvement.

---

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `internal/model/types.go` | Modify | Add 11 SkillID constants (`SkillPentestOrchestrator` through `SkillMalwareReportWriter`), `PresetCyber` constant, `PermissionTier` type with 3 constants, `ToolPermission` struct |
| `internal/catalog/skills.go` | Modify | Add `cyberSkills` slice and `MVPCyberSkills()` function returning 11 entries |
| `internal/catalog/cyber_skills_test.go` | Create | Validation tests: frontmatter parsing, category validity, destructive tool manifest checks |
| `internal/components/skills/presets.go` | Modify | Add `cyberSkills` var, extend `SkillsForPreset()` with `PresetCyber` case, extend `AllSkillIDs()` |
| `internal/tui/screens/preset.go` | Modify | Add `PresetCyber` to `PresetOptions()`, add description in `presetDescriptions` |
| `internal/tui/screens/skill_picker.go` | Modify | Add `cyberSkillIDs` group and `cyberSkillLabels`, render "Security Skills" section |
| `internal/components/mcp/kali_mcp.go` | Create | Overlay JSON per agent strategy for kali-mcp (follows `context7.go` pattern) |
| `internal/components/mcp/shodan_mcp.go` | Create | Overlay JSON per agent strategy for shodan-mcp (npx-based, similar to context7) |
| `internal/components/mcp/virustotal_mcp.go` | Create | Overlay JSON per agent strategy for virustotal-mcp (npx-based) |
| `internal/components/mcp/inject.go` | Modify | Add `InjectCyber()` function and extend existing `Inject()` to handle cyber MCPs when preset is cyber |
| `mcps/kali-mcp/manifest.json` | Create | kali-mcp manifest: name, type stdio, command, args, `permission_tier: destructive`, `destructive_tools` list |
| `mcps/shodan-mcp/manifest.json` | Create | shodan-mcp manifest: name, type npx, `permission_tier: unrestricted` |
| `mcps/virustotal-mcp/manifest.json` | Create | virustotal-mcp manifest: name, type npx, `permission_tier: unrestricted` |
| `skills/pentest-orchestrator/SKILL.md` | Create | Red-team skill: pentest orchestration (from glitch-ai-toolkit, adapted to gentle-ai frontmatter) |
| `skills/ai-pentesting-validation/SKILL.md` | Create | Red-team skill: anti-hallucination validation for pentest findings |
| `skills/exploit-chain-patterns/SKILL.md` | Create | Red-team skill: vulnerability chaining methodology |
| `skills/waf-detection-bypass/SKILL.md` | Create | Red-team skill: WAF detection and bypass techniques |
| `skills/detection-engineer/SKILL.md` | Create | Blue-team skill: Sigma/Suricata rule creation from analysis |
| `skills/malware-triage/SKILL.md` | Create | SOC skill: initial malware assessment and classification |
| `skills/specialized-file-analyzer/SKILL.md` | Create | SOC skill: non-PE file analysis (.NET, Office, PDF, ELF) |
| `skills/malware-dynamic-analysis/SKILL.md` | Create | SOC skill: sandbox execution and behavioral observation |
| `skills/malware-report-writer/SKILL.md` | Create | SOC skill: professional malware analysis report authoring |
| `agents/gentleman-soc.md` | Create | SOC orchestrator agent definition (PICERL pipeline, skill-loading instructions) |
| `internal/assets/skills/_shared/SKILL.md` | Modify | Add security skill cross-reference documentation in shared skill content |
| `internal/cli/validate.go` | Modify | Extend `normalizePreset()` to accept `PresetCyber`, extend `normalizeSkills()` to accept cyber skill IDs |
| `internal/components/skills/inject.go` | Modify | No change needed — cyber skills use the same embed-based injection path; `assets.FS` already walks by skill name |
| `docs/cyber-permissions.md` | Create | User-facing documentation of the soft gate model, destructive tool list, warning semantics |
| `docs/cyber-setup.md` | Create | Kali Linux MCP setup guide, shodan/virustotal API key configuration, Prowler MCP installation |
| `scripts/sync-glitch-skills.sh` | Create | Shell script to pull upstream skill content from glitch-ai-toolkit and reformat frontmatter |

---

## Interfaces / Contracts

### New SkillID Constants (`internal/model/types.go`)

```go
// Cybersecurity skills (v1)
SkillPentestOrchestrator    SkillID = "pentest-orchestrator"
SkillAIPentestingValidation SkillID = "ai-pentesting-validation"
SkillExploitChainPatterns   SkillID = "exploit-chain-patterns"
SkillWAFDetectionBypass     SkillID = "waf-detection-bypass"
SkillDetectionEngineer      SkillID = "detection-engineer"
SkillSecurityPythonScripts  SkillID = "python-security" // reserved; src name TBD post-content review
SkillMalwareTriage          SkillID = "malware-triage"
SkillSpecializedFileAnalyzer SkillID = "specialized-file-analyzer"
SkillMalwareDynamicAnalysis SkillID = "malware-dynamic-analysis"
SkillMalwareReportWriter    SkillID = "malware-report-writer"
```

Note: 10 constant names are listed here. The 11th skill in the proposal, `prowler-sdk-check`, is deferred to v2 (per spec `cyber-skills/spec.md` requirement 4). The `python-security` skill may rename depending on content review; the SkillID constant matches the directory name.

### Cyber Skill Catalog Entries (`internal/catalog/skills.go`)

```go
var cyberSkills = []Skill{
    {ID: model.SkillPentestOrchestrator, Name: "pentest-orchestrator", Category: "red-team", Priority: "p0"},
    {ID: model.SkillAIPentestingValidation, Name: "ai-pentesting-validation", Category: "red-team", Priority: "p0"},
    {ID: model.SkillExploitChainPatterns, Name: "exploit-chain-patterns", Category: "red-team", Priority: "p0"},
    {ID: model.SkillWAFDetectionBypass, Name: "waf-detection-bypass", Category: "red-team", Priority: "p0"},
    {ID: model.SkillDetectionEngineer, Name: "detection-engineer", Category: "blue-team", Priority: "p0"},
    {ID: model.SkillSecurityPythonScripts, Name: "python-security", Category: "blue-team", Priority: "p0"},
    {ID: model.SkillMalwareTriage, Name: "malware-triage", Category: "soc", Priority: "p0"},
    {ID: model.SkillSpecializedFileAnalyzer, Name: "specialized-file-analyzer", Category: "soc", Priority: "p0"},
    {ID: model.SkillMalwareDynamicAnalysis, Name: "malware-dynamic-analysis", Category: "soc", Priority: "p0"},
    {ID: model.SkillMalwareReportWriter, Name: "malware-report-writer", Category: "soc", Priority: "p0"},
}

func MVPCyberSkills() []Skill {
    skills := make([]Skill, len(cyberSkills))
    copy(skills, cyberSkills)
    return skills
}
```

### Permission Tier Types (`internal/model/types.go`)

```go
type PermissionTier string

const (
    PermissionTierUnrestricted PermissionTier = "unrestricted"
    PermissionTierDestructive  PermissionTier = "destructive"
    PermissionTierRestricted  PermissionTier = "restricted"
)

type ToolPermission struct {
    MCPServer string
    Tools     []string
    Tier      PermissionTier
}
```

### Cyber Preset (`internal/model/types.go`)

```go
PresetCyber PresetID = "cyber"
```

### MCP Manifest Contract (`mcps/{name}/manifest.json`)

```json
{
  "name": "kali-mcp",
  "type": "stdio",
  "description": "Kali Linux security tools via MCP",
  "config": {
    "command": "kali-server-mcp",
    "args": []
  },
  "permission_tier": "destructive",
  "destructive_tools": [
    "nmap_scan",
    "sqlmap_scan",
    "hydra_attack",
    "metasploit_run",
    "nikto_scan",
    "dirb_scan",
    "gobuster_scan",
    "wpscan_analyze",
    "john_crack",
    "enum4linux_scan"
  ],
  "confirmation_required": true,
  "env": {
    "KALI_HOST": "<user-configured>",
    "KALI_SSH_KEY": "<user-configured>"
  }
}
```

### Skill Frontmatter Contract (Cyber skills)

```yaml
---
name: malware-triage
description: "Systematic malware triage and initial assessment. Trigger: When you need to perform initial malware assessment, classify samples, or determine analysis priority."
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---
```

Key rules from `internal/assets/skills_frontmatter_test.go`:
- `description` must be a quoted scalar (not `>` or `|` block)
- `description` must contain `Trigger:` substring (except `_shared`)
- `description` must be ≤160 chars
- `name` must match the parent directory basename
- No `auto_invoke` or `tool_schemas` keys (blocked by allowed-keys test)

### Preset Composition (`internal/components/skills/presets.go`)

```go
// cyberSkills are the 10 v1 cybersecurity skills.
var cyberSkills = []model.SkillID{
    model.SkillPentestOrchestrator,
    model.SkillAIPentestingValidation,
    model.SkillExploitChainPatterns,
    model.SkillWAFDetectionBypass,
    model.SkillDetectionEngineer,
    model.SkillSecurityPythonScripts,
    model.SkillMalwareTriage,
    model.SkillSpecializedFileAnalyzer,
    model.SkillMalwareDynamicAnalysis,
    model.SkillMalwareReportWriter,
}
```

With `SkillsForPreset()` gaining a new case:

```go
case model.PresetCyber:
    all := make([]model.SkillID, 0, len(sddSkills)+len(foundationSkills)+len(cyberSkills))
    all = append(all, sddSkills...)
    all = append(all, foundationSkills...)
    all = append(all, cyberSkills...)
    return all
```

And `AllSkillIDs()` updated to include cyber skills.

---

## Testing Strategy

| Layer | What to Test | Approach |
|-------|-------------|----------|
| Unit | SkillID constants compile and are unique | `go vet` + duplicate check test |
| Unit | `MVPCyberSkills()` returns exactly 10 entries with valid categories | Test that each entry has Category in `{"red-team","blue-team","soc"}` and non-empty Name |
| Unit | `SkillsForPreset(PresetCyber)` returns 30 skills (20 MVP + 10 cyber) | Verify SDD + foundation + cyber skill IDs are all present |
| Unit | `normalizePreset("cyber")` returns `PresetCyber` | Extend existing `normalizePreset` test table |
| Unit | `normalizeSkills()` accepts all 10 cyber skill IDs | Extend existing skill normalization test |
| Unit | Cyber skill SKILL.md files parse with valid frontmatter | Reuse existing `TestSkillFrontmatterIsLintClean` — no changes needed since it walks embedded assets |
| Unit | Cyber skill minimal content (>100 chars beyond frontmatter, has `## Validation` section) | New test in `cyber_skills_test.go` |
| Unit | MCP manifest JSON validates: required fields, `permission_tier` is one of the 3 valid values, `destructive_tools` is non-empty when tier is `destructive` | New test, possibly in `internal/components/mcp/manifest_test.go` |
| Unit | `PresetCyber` appears in `PresetOptions()` and has a description | New test in preset_test |
| Unit | `InjectCyber()` correctly writes kali-mcp, shodan-mcp, virustotal-mcp overlays per agent strategy | New test with mock adapters |
| Integration | Full cyber preset TUI flow: select preset → see skills → install succeeds | Extend existing TUI integration test pattern |
| Integration | Prowler MCP detection: `exec.LookPath` mock returning found/not found | New test in `internal/components/mcp/prowler_detect_test.go` |
| E2E | `go build ./cmd/gentle-ai` compiles with all new code | CI build check |
| E2E | All existing `go test ./...` pass with no regressions | CI test check |

### Priority Test Scenarios (BDD)

These map to the spec requirements in `cyber-skills/spec.md`, `cyber-mcp/spec.md`, and `cyber-permissions/spec.md`:

1. **All 10 cyber skills registered in catalog** — `MVPCyberSkills()` returns exactly 10 entries
2. **Red-team skills have correct category** — 4 skills with `Category: "red-team"`
3. **No SkillID collisions** — 10 new constants don't duplicate any of the 20 existing ones
4. **Skill frontmatter is lint-clean** — Existing asset test covers this automatically
5. **Kali-mcp manifest has `permission_tier: destructive`** — JSON parsing test
6. **Shodan/virustotal manifests have `permission_tier: unrestricted`** — JSON parsing test
7. **Destructive tools list is non-empty for kali-mcp** — Array length check
8. **PresetCyber resolves skills correctly** — 10 cyber + 20 MVP = 30 total
9. **Non-cyber presets are unchanged** — `MVPSkills()` still returns 20 entries
10. **No auto_invoke or tool_schemas in cyber skill frontmatter** — New test

---

## Migration / Rollout

No data migration required. The change is purely additive:

1. **Build step**: `go build ./cmd/gentle-ai` compiles all new code along with existing code
2. **Asset embedding**: New skill directories under `internal/assets/skills/` are picked up automatically by the `//go:embed all:skills` directive in `assets.go`
3. **Backward compatibility**: Selecting any existing preset (`full-gentleman`, `ecosystem-only`, `minimal`, `custom`) produces the same behavior as before. The `PresetCyber` is opt-in.

**Rollback plan** (from proposal): Delete the added files and revert the modified Go files. All additions are additive, with no breaking changes to existing data structures or function signatures.

---

## Risks & Mitigations

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| kali-mcp tools cause real damage on production systems | High if misused | Critical | Soft gate via agent prompt warnings (v1); hard proxy planned for v2; kali-mcp requires external Kali host (not local); destructive MCPs are opt-in (not pre-selected) |
| Skill content drifts from glitch-ai-toolkit upstream | Medium | Medium | `scripts/sync-glitch-skills.sh` pulls upstream content; frontmatter is stable (name/description/license/metadata); CI validates frontmatter on every build |
| Agent ignores destructive-tool confirmation prompt | Medium | High | Documented as soft gate; v2 adds hard proxy; warning is visible in system prompt for auditing; kali-mcp requires explicit user consent during installation |
| Prowler MCP API changes break integration | Low | Low | Prowler MCP is external dependency (pip package); v1 only configures it if detected, doesn't embed it; breaking changes in Prowler are Prowler's responsibility |
| gentle-ai upstream changes conflict with fork | Medium | Medium | Rebase strategy; cyber additions are in separate files (catalog additions, skill directories, MCP configs) minimizing conflict surface; 11 SkillIDs and 1 PresetID are additive constants |
| 10 new skills bloat binary size | Low | Low | Skills are markdown files copied at install time from embedded assets (`internal/assets/skills/`), not compiled. Each SKILL.md is ~5-15KB. Total: ~100-150KB added to binary |
| Frontmatter lint test catches non-compliant upstream skill content | High | High | The `skills_frontmatter_test.go` enforces `Trigger:` substring, ≤160 chars, quoted description. Sync script must validate content before embedding. Test failures block the build |

---

## Open Questions

- [ ] Should `python-security` skill keep this name or be renamed to `secure-python` or `python-security-scripts`? The glitch-ai-toolkit source needs to be checked for what this skill actually covers.
- [ ] The 11th skill from the proposal (`prowler-sdk-check`) is deferred to v2 per the spec. Should its SkillID constant be reserved as a comment in `types.go` to avoid future naming collisions?
- [ ] Should kali-mcp overlay JSON include `KALI_HOST` and `KALI_SSH_KEY` as template env vars that the user fills in, or should the TUI collect these during installation?
- [ ] The `agents/gentleman-soc.md` definition file references PICERL phases that load skills by name. Should the SOC agent definition also inject `gentle-ai-cyber:destructive-warning` as a standalone section, or should it reference the same warning block that the MCP injection produces? (Currently designed as a single source of truth from the manifest.)
- [ ] Should the `mcps/` directory live at the repo root (alongside `skills/`) or inside `internal/assets/mcps/` for embedding? The current pattern has MCP overlays in Go source (`internal/components/mcp/`), not as embedded assets. Manifests are a new concept — they need a home.