# Gentle AI — Cybersecurity Edition

> **gentle-ai-cyber** is an additive cybersecurity layer on top of Gentle AI. It adds offensive and defensive security skills, MCP server integrations, a SOC orchestrator agent, and a permission model for destructive tools — without changing any existing functionality.

## What Is It

The cybersecurity edition installs alongside the standard Gentle AI setup. When you select the **cyber preset** during installation (or pass `--preset cyber` / `--cyber` on the CLI), you get:

- **10 security-focused coding skills** organized by team role (Red Team, Blue Team, SOC)
- **3 MCP server integrations** (Kali Linux, VirusTotal, Shodan) with permission tiers
- **gentleman-soc** — a SOC orchestrator agent definition that guides incident response through the PICERL pipeline
- **Destructive tool warnings** — soft gates that require human confirmation before executing potentially harmful commands

Everything is additive. The standard `full-gentleman` preset is unchanged.

## Skills Taxonomy

### Red Team (4 skills)

| Skill | Purpose |
|-------|---------|
| `pentest-orchestrator` | AI-driven multi-phase penetration testing orchestration |
| `ai-pentesting-validation` | Anti-hallucination validation pipeline for security assessments |
| `exploit-chain-patterns` | Combining vulnerabilities into high-impact attack paths |
| `waf-detection-bypass` | WAF identification and bypass techniques for authorized testing |

### Blue Team (2 skills)

| Skill | Purpose |
|-------|---------|
| `detection-engineer` | Sigma/Suricata detection rules and hunting queries from malware findings |
| `python-security` | Secure Python patterns and defensive coding best practices |

### SOC (4 skills)

| Skill | Purpose |
|-------|---------|
| `malware-triage` | Systematic malware triage and initial assessment workflow |
| `specialized-file-analyzer` | Analysis of non-PE files: Office macros, PDFs, scripts, archives, ELF |
| `malware-dynamic-analysis` | Safe sandbox execution with Procmon, Wireshark, Process Hacker, Sysmon |
| `malware-report-writer` | Professional malware analysis report creation for enterprise IR |

## MCP Integrations

### Kali Linux MCP (`kali-mcp`)

Connects to a running Kali Linux server to execute security tools programmatically:

- **Tools**: nmap, sqlmap, hydra, metasploit, nikto, dirb, gobuster, wpscan, john, enum4linux
- **Permission tier**: `destructive` — requires explicit human confirmation before execution
- **Setup**: Requires a running Kali Linux instance with the MCP server installed

### VirusTotal MCP (`virustotal-mcp`)

- **Permission tier**: `unrestricted` — read-only threat intelligence queries
- **Setup**: Requires a VirusTotal API key

### Shodan MCP (`shodan-mcp`)

- **Permission tier**: `unrestricted` — read-only internet-wide scanning data
- **Setup**: Requires a Shodan API key

## Permission Model

The cyber edition uses a **soft gate** permission model — warnings and confirmation prompts in the agent's system prompt, not hard runtime blocks.

| Tier | Behavior | Example |
|------|----------|---------|
| `unrestricted` | Agent may invoke without confirmation | VirusTotal lookups, Shodan queries |
| `destructive` | Agent must ask user before invoking | nmap scans, metasploit exploits, hydra attacks |
| `restricted` | Agent should not invoke (reserved) | Future use |

Destructive tools are listed in two places:
1. **System prompt warning** — injected as `<!-- gentle-ai-cyber:destructive-warning -->` section
2. **MCP JSON comment** — `/**_WARNING**/` block prefixed to the MCP config for agents that support separate MCP files

> **Important**: These are soft gates. The agent is instructed to ask for confirmation, but the model ultimately decides compliance. Never use the cyber preset on production systems without human oversight.

## How to Install

### Via TUI

```bash
gentle-ai
# Navigate to "Preset" → Select "cyber" → Follow installation
```

### Via CLI

```bash
gentle-ai --preset cyber
# or
gentle-ai --cyber
```

### After Install

The cyber preset installs:
- 10 skill directories under `skills/` (or embedded, depending on agent)
- 3 MCP server configurations per agent
- `gentleman-soc.md` agent definition (for agents that support markdown augmentation)
- Destructive tool warnings in system prompts

## Using the gentleman-soc Agent

`gentleman-soc` is not a standalone agent — it is an **orchestration persona** that works through your existing AI agent (Claude Code, OpenCode, etc.). It guides you through the **PICERL** incident response pipeline:

| Phase | Name | Focus |
|-------|------|-------|
| 0 | Preparation | Assess situation, plan pipeline |
| 1 | Identification | Triage, classify, prioritize |
| 2 | Identification | Deep static/specialized analysis |
| 3 | Identification + Containment | Dynamic analysis, IOC extraction |
| 4 | Containment + Recovery | Detection rules, hunting queries |
| 5 | Lessons Learned | Professional report |

The SOC agent definition tells your AI agent **when** to load each skill based on the current phase. It does not load all skills upfront.

### Workflow Example

1. **Triage** → `malware-triage` skill loads → classify the sample
2. **Deep Analysis** → `specialized-file-analyzer` → extract strings, sections, behavior
3. **Dynamic Analysis** → `malware-dynamic-analysis` → sandbox execution, IOC capture
4. **Detection** → `detection-engineer` → write Sigma/Suricata rules
5. **Reporting** → `malware-report-writer` → generate professional report

Throughout, the agent maintains a MITRE ATT&CK mapping table and quality gates between phases.

## v1 vs v2 Scope

### v1 (Current)

- 10 skills (4 red-team, 2 blue-team, 4 SOC)
- 3 MCP servers (Kali, VirusTotal, Shodan)
- Soft gate permission model (prompt-based)
- `gentleman-soc` agent definition (markdown augmentation)
- Destructive tool warnings in system prompts and MCP configs

### v2 (Planned)

- Additional MCP servers (Prowler, additional threat intel)
- Hard runtime permission gates (not just prompt-based)
- Expanded skill catalog from upstream glitch-ai-toolkit
- Compliance category skills
- Cross-agent SOC collaboration patterns

## Security Considerations

1. **Authorized use only** — These tools are for authorized security assessments. Always have explicit permission before testing any system you do not own.
2. **Destructive tools require confirmation** — The soft gate model relies on agent compliance. Always verify before execution.
3. **API keys are sensitive** — VirusTotal and Shodan API keys are stored in your agent's MCP config. Protect them accordingly.
4. **Kali server is external** — The kali-mcp connects to an external Kali Linux instance. Ensure it is properly isolated from production networks.

## Related Documentation

- [AGENTS.md](../AGENTS.md) — Full skill index with paths
- [README.md](../README.md) — Gentle AI overview and quick start
- [Intended Usage](./intended-usage.md) — Gentle AI mental model
