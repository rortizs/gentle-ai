# Tasks: Gentle AI Cybersecurity Fork

## Phase 1: Foundation (Infrastructure + Types)

- [x] 1.1 Add `PermissionTier` type and 3 constants (`unrestricted`, `destructive`, `restricted`) to `internal/model/types.go`
- [x] 1.2 Add `ToolPermission` struct (`MCPServer`, `Tools`, `Tier`) to `internal/model/types.go`
- [x] 1.3 Add `PresetCyber PresetID = "cyber"` constant to `internal/model/types.go`
- [x] 1.4 Add 10 cyber SkillID constants to `internal/model/types.go`: `SkillPentestOrchestrator`, `SkillAIPentestingValidation`, `SkillExploitChainPatterns`, `SkillWAFDetectionBypass`, `SkillDetectionEngineer`, `SkillSecurityPythonScripts`, `SkillMalwareTriage`, `SkillSpecializedFileAnalyzer`, `SkillMalwareDynamicAnalysis`, `SkillMalwareReportWriter`
- [x] 1.5 Add `cyberSkills` slice and `MVPCyberSkills()` function to `internal/catalog/skills.go`
- [x] 1.6 Add `cyberSkills` variable and extend `SkillsForPreset()` with `PresetCyber` case in `internal/components/skills/presets.go`
- [x] 1.7 Extend `AllSkillIDs()` in `internal/components/skills/presets.go` to include cyber skills
- [x] 1.8 Add `PresetCyber` to `PresetOptions()` and `presetDescriptions` map in `internal/tui/screens/preset.go`
- [x] 1.9 Add `normalizePreset("cyber")` support and cyber skill ID acceptance in `internal/cli/validate.go`

**Verification for Phase 1**:
```bash
go build ./internal/model/... && echo "types OK"
go build ./internal/catalog/... && echo "catalog OK"
go build ./internal/components/skills/... && echo "presets OK"
go test ./internal/model/... ./internal/catalog/... -v -run "Skill|Preset" 2>&1 | head -50
```

---

## Phase 2: MCP Infrastructure

- [x] 2.1 Create `mcps/kali-mcp/manifest.json` with `permission_tier: destructive`, `destructive_tools` array (nmap_scan, sqlmap_scan, hydra_attack, metasploit_run, nikto_scan, dirb_scan, gobuster_scan, wpscan_analyze, john_crack, enum4linux_scan)
- [x] 2.2 Create `mcps/shodan-mcp/manifest.json` with `permission_tier: unrestricted`
- [x] 2.3 Create `mcps/virustotal-mcp/manifest.json` with `permission_tier: unrestricted`
- [x] 2.4 Create `internal/components/mcp/kali_mcp.go` with overlay JSON per agent strategy (following context7.go pattern)
- [ ] 2.5 Create `internal/components/mcp/shodan_mcp.go` with overlay JSON per agent strategy
- [ ] 2.6 Create `internal/components/mcp/virustotal_mcp.go` with overlay JSON per agent strategy
- [ ] 2.7 Extend `internal/components/mcp/inject.go` with `InjectCyber()` function and extend `Inject()` to handle cyber MCPs when preset is cyber
- [ ] 2.8 Add Prowler MCP external detection logic in `internal/components/mcp/inject.go` (exec.LookPath for prowler-mcp)

> **Note**: Tasks 2.5–2.8 were scoped for implementation in the Go-side MCP infrastructure but were not completed in the v1 PRs. The 5 PRs focused on: (1) Foundation types/manifests, (2) Skill content, (3) Agent definition/warning, (4) Testing harness, (5) Documentation. The MCP Go injection layer (shodan_mcp.go, virustotal_mcp.go, inject.go extensions) and Prowler detection remain as work-in-progress. These are tracked as deferred items for the next SDD cycle.

**Verification for Phase 2**:
```bash
cat mcps/kali-mcp/manifest.json | python3 -c "import sys,json; m=json.load(sys.stdin); assert m['permission_tier']=='destructive'; assert len(m['destructive_tools'])>0; print('kali manifest OK')"
cat mcps/shodan-mcp/manifest.json | python3 -c "import sys,json; m=json.load(sys.stdin); assert m['permission_tier']=='unrestricted'; print('shodan manifest OK')"
cat mcps/virustotal-mcp/manifest.json | python3 -c "import sys,json; m=json.load(sys.stdin); assert m['permission_tier']=='unrestricted'; print('virustotal manifest OK')"
go build ./internal/components/mcp/... && echo "MCP Go files OK"
```

---

## Phase 3: Skill Files (Content)

- [x] 3.1 Create `skills/pentest-orchestrator/SKILL.md` (red-team, ~300-500 lines, frontmatter + PICERL workflow + ## Validation section)
- [x] 3.2 Create `skills/ai-pentesting-validation/SKILL.md` (red-team, ~200-400 lines, anti-hallucination validation patterns)
- [x] 3.3 Create `skills/exploit-chain-patterns/SKILL.md` (red-team, ~300-400 lines, vulnerability chaining)
- [x] 3.4 Create `skills/waf-detection-bypass/SKILL.md` (red-team, ~250-400 lines, WAF detection/bypass)
- [x] 3.5 Create `skills/detection-engineer/SKILL.md` (blue-team, ~300-500 lines, Sigma/Suricata rules)
- [x] 3.6 Create `skills/python-security/SKILL.md` (blue-team, ~200-400 lines, secure python patterns)
- [x] 3.7 Create `skills/malware-triage/SKILL.md` (soc, ~300-500 lines, initial assessment)
- [x] 3.8 Create `skills/specialized-file-analyzer/SKILL.md` (soc, ~300-600 lines, non-PE analysis)
- [x] 3.9 Create `skills/malware-dynamic-analysis/SKILL.md` (soc, ~300-500 lines, sandbox execution)
- [x] 3.10 Create `skills/malware-report-writer/SKILL.md` (soc, ~300-500 lines, report authoring)
- [x] 3.11 Embed skill files in `internal/assets/skills/` (add to existing embed directive or add new embed)

**Verification for Phase 3**:
```bash
# Check frontmatter format on all skill files
for f in skills/*/SKILL.md; do
  python3 -c "
import sys, re
content = open('$f'.replace('$f',f)).read()
assert re.search(r'^name: .+$', content, re.M)
assert 'description:' in content
assert 'license: Apache-2.0' in content
assert 'metadata:' in content
assert 'version:' in content
assert '## Validation' in content
print('OK: $f')
" || echo "FAIL: $f"
done

# Check total skill content lines
wc -l skills/*/SKILL.md | tail -1
```

---

## Phase 4: Agent Definition + Destructive Warning

- [x] 4.1 Create `agents/gentleman-soc.md` (SOC orchestrator agent definition, ~500-600 lines, PICERL pipeline, MITRE ATT&CK mapping, quality gates, skill loading instructions for 6 SOC/blue-team skills)
- [x] 4.2 Add `gentleman-soc.md` to Go embed in `internal/assets/`
- [x] 4.3 Implement destructive-tool warning block injection in agent system prompts (add `<!-- gentle-ai-cyber:destructive-warning -->` section to prompt generation)
- [x] 4.4 Implement `/**_WARNING**/` comment in kali-mcp JSON for `StrategySeparateMCPFiles` agents
- [x] 4.5 Extend `internal/cli/validate.go` to handle `--cyber` CLI flag and conflict detection with `--preset`

**Verification for Phase 4**:
```bash
# Check gentleman-soc.md content
grep -c "Phase" agents/gentleman-soc.md   # should have >=5 phases
grep -c "skill" agents/gentleman-soc.md   # should reference skills
grep "PICERL\|Prepare\|Identify\|Contain\|Eradicate\|Recover" agents/gentleman-soc.md | head -5

# Check destructive warning in generated prompts (unit test)
go test ./internal/cli/... -v -run "Destructive|Cyber" 2>&1 | head -30

# Check --cyber flag behavior
go run ./cmd/gentle-ai --help 2>&1 | grep -i cyber || echo "no --cyber in help"
```

---

## Phase 5: Testing

- [x] 5.1 Create `internal/catalog/cyber_skills_test.go` with `TestCyberSkillsHaveRequiredFrontmatter`, `TestCyberSkillCategoriesAreValid`, `TestDestructiveToolsRequireConfirmation`, `TestProwlerCheckIDFormat`
- [x] 5.2 Create `internal/components/mcp/kali_mcp_test.go` for MCP manifest validation and destructive warning injection
- [x] 5.3 Create `internal/components/skills/presets_cyber_test.go` to verify `PresetCyber` resolves to 30 skills (20 MVP + 10 cyber)
- [x] 5.4 Extend `internal/catalog/skills_test.go` with `TestMVPCyberSkillsReturns10Entries`, `TestNoSkillIDCollisions`
- [x] 5.5 Create `internal/assets/assets_cyber_test.go` to verify cyber skill embedding and frontmatter

**Verification for Phase 5**:
```bash
go test ./internal/catalog/... -v -run "Cyber" 2>&1
go test ./internal/components/mcp/... -v -run "Manifest" 2>&1
go test ./internal/components/skills/... -v -run "PresetCyber" 2>&1
go test ./... 2>&1 | tail -20  # full suite, check for regressions
```

---

## Phase 6: Documentation

- [ ] 6.1 Create `docs/cyber-permissions.md` (soft gate model, destructive tool list, warning semantics)
- [ ] 6.2 Create `docs/cyber-setup.md` (Kali MCP setup, shodan/virustotal API keys, Prowler MCP install)
- [ ] 6.3 Create `scripts/sync-glitch-skills.sh` (pull upstream from glitch-ai-toolkit)

> **Note**: Tasks 6.1–6.3 were planned as supplementary documentation and sync tooling but were not included in the 5 PR scope. Documentation was consolidated into `docs/cybersecurity-edition.md` (created) and `AGENTS.md` (updated). The permissions model is documented in the `kali_mcp.go` source code and `gentleman-soc.md` agent definition.
- [x] 6.4 Update `AGENTS.md` with cyber skill taxonomy (red-team, blue-team, soc categories)

**Verification for Phase 6**:
```bash
# Documentation checks
ls -la docs/cyber-*.md
grep "soft gate\|prompt-based" docs/cyber-permissions.md
grep "pip install prowler-mcp" docs/cyber-setup.md

# AGENTS.md has cyber skills
grep -A5 "cyber\|red-team\|blue-team\|soc" AGENTS.md | head -20
```

---

## Dependencies Summary

```
Phase 1 (Foundation)
  └── must complete before: Phase 5 (tests that import cyber SkillIDs)

Phase 2 (MCP Infrastructure)
  └── depends on: Phase 1 (types for manifest parsing)
  └── can run parallel with: Phase 3 (no dependency)

Phase 3 (Skill Files)
  └── can run parallel with: Phase 2
  └── must complete before: Phase 5 (skill content tests)

Phase 4 (Agent + Warning)
  └── depends on: Phase 1, Phase 2 (types and MCP files exist)

Phase 5 (Testing)
  └── depends on: Phase 1, Phase 2, Phase 3, Phase 4

Phase 6 (Documentation)
  └── depends on: Phase 4 (agent definition exists)
```

---

## Work Units / Commit Grouping

### Commit 1: Core Type & Catalog Additions
```
internal/model/types.go                    (+25 lines)
internal/catalog/skills.go                (+15 lines)
internal/components/skills/presets.go     (+20 lines)
internal/tui/screens/preset.go            (+5 lines)
internal/cli/validate.go                  (+10 lines)
```
Small, foundational. Reviewable in one pass.

### Commit 2: MCP Infrastructure
```
mcps/kali-mcp/manifest.json                (new)
mcps/shodan-mcp/manifest.json              (new)
mcps/virustotal-mcp/manifest.json          (new)
internal/components/mcp/kali_mcp.go        (new)
internal/components/mcp/shodan_mcp.go      (new)
internal/components/mcp/virustotal_mcp.go  (new)
internal/components/mcp/inject.go          (+30 lines)
```
Medium. Clear boundaries.

### Commit 3: Skill Files (Content)
```
skills/pentest-orchestrator/SKILL.md
skills/ai-pentesting-validation/SKILL.md
skills/exploit-chain-patterns/SKILL.md
skills/waf-detection-bypass/SKILL.md
skills/detection-engineer/SKILL.md
skills/python-security/SKILL.md
skills/malware-triage/SKILL.md
skills/specialized-file-analyzer/SKILL.md
skills/malware-dynamic-analysis/SKILL.md
skills/malware-report-writer/SKILL.md
internal/assets/skills/...                 (embed updates)
```
Large. ~2500-5000 lines of skill content. CHAINED PR recommended.

### Commit 4: Agent Definition + Destructive Warning
```
agents/gentleman-soc.md                    (new, ~580 lines)
internal/cli/cyber_flag.go                 (new or extend validate.go)
internal/components/prompt/warning.go     (new, destructive warning block)
```
Medium. Agent definition is large but self-contained.

### Commit 5: Testing
```
internal/catalog/cyber_skills_test.go      (new)
internal/components/mcp/manifest_test.go   (new)
internal/components/skills/presets_cyber_test.go (new)
```
Small-Medium. Test files.

### Commit 6: Documentation
```
docs/cyber-permissions.md                  (new)
docs/cyber-setup.md                        (new)
scripts/sync-glitch-skills.sh              (new)
AGENTS.md                                  (update)
```
Small. Can be combined with other commits or done last.

---

## Review Workload Forecast

| Metric | Value |
|--------|-------|
| New files | 24 |
| Modified files | 8 |
| Deleted files | 0 |
| **Total changed files** | **32** |
| Go code additions (est.) | ~200 lines |
| Go code deletions (est.) | ~0 lines |
| **Go net additions** | **~200 lines** |
| Skill SKILL.md content (est.) | ~3000-5000 lines |
| Agent definition (gentleman-soc.md) | ~580 lines |
| Documentation content | ~400 lines |
| **Total lines (additions only)** | **~4200-6200 lines** |

### Exceeds 400 lines? ✅ YES
Estimated total: **4200-6200 lines** of additions across all phases.

### Chained PRs Recommended? ✅ YES
The skill files phase alone (Phase 3) contains 10 new SKILL.md files with ~300-500 lines each. Even at 300 lines average, that's 3000 lines — more than 7x the 400-line threshold.

**Recommended chain:**
1. **PR 1**: Types + Catalog + Presets + TUI + MCP manifests + MCP Go files (Foundation: ~150 lines Go + 3 JSON manifests)
2. **PR 2**: Skill files (Content: ~3000-5000 lines — MUST be separate)
3. **PR 3**: Agent definition + destructive warning injection (Middleware: ~600 lines)
4. **PR 4**: Testing (Tests: ~200 lines)
5. **PR 5**: Documentation + AGENTS.md update (Docs: ~400 lines + updates)

### Risk Level: 🟡 **MEDIUM-HIGH**

| Risk Factor | Assessment |
|------------|------------|
| Scope | 10 skill files + 3 MCP configs + agent definition = high content volume |
| Skill content accuracy | Each SKILL.md must have valid frontmatter + ## Validation section; upstream drift possible |
| Frontmatter lint | Existing `TestSkillFrontmatterIsLintClean` will catch malformed skills — must ensure all 10 pass |
| Soft gate reliability | Agent prompt compliance is the gate — if model ignores, no runtime enforcement |
| No breaking changes | All additive; existing presets unchanged. Low regression risk on Go side. |
| Prowler MCP detection | exec.LookPath may vary by OS; needs cross-platform verification |
| kali-mcp external dep | Requires Kali host running `kali-server-mcp` — can't fully test in CI |

### Implementation Risk
- **Content creation** (10 skill files): High effort, medium risk. Each skill needs careful frontmatter + valid content. Use existing skill patterns from gentle-ai.
- **MCP overlay JSON per agent**: 13 agents × 3 MCPs = 39 JSON combinations. context7.go pattern handles this. Medium risk.
- **Destructive warning injection**: Changes system prompt generation — need to ensure non-cyber presets remain clean. Medium risk.

### Confidence
- **Go infrastructure**: High confidence (follows existing patterns exactly)
- **Skill content**: Medium confidence (depends on glitch-ai-toolkit source quality)
- **Agent definition**: Medium confidence (gentleman-soc.md is new, no upstream to validate against)