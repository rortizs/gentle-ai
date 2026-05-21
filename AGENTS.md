# Gentle AI — Agent Skills Index

When working on this project, load the relevant skill(s) BEFORE writing any code.

Naming convention: `gentle-ai-*` skills are repo-specific workflow skills. Unprefixed skills are portable writing or work-unit skills and intentionally keep their canonical names.

## How to Use

1. Check the trigger column to find skills that match your current task
2. Load the skill by reading the SKILL.md file at the listed path
3. Follow ALL patterns and rules from the loaded skill
4. Multiple skills can apply simultaneously

## Skills

| Skill | Trigger | Path |
|-------|---------|------|
| `gentle-ai-issue-creation` | When creating a GitHub issue, reporting a bug, or requesting a feature. | [`skills/issue-creation/SKILL.md`](skills/issue-creation/SKILL.md) |
| `gentle-ai-branch-pr` | When creating a pull request, opening a PR, or preparing changes for review. | [`skills/branch-pr/SKILL.md`](skills/branch-pr/SKILL.md) |
| `gentle-ai-chained-pr` | When a change is too large for one review, or when creating chained/stacked pull requests. | [`skills/chained-pr/SKILL.md`](skills/chained-pr/SKILL.md) |
| `cognitive-doc-design` | When writing docs that must reduce cognitive load for readers or reviewers. | [`skills/cognitive-doc-design/SKILL.md`](skills/cognitive-doc-design/SKILL.md) |
| `comment-writer` | When drafting human comments, PR feedback, issue replies, or async updates. | [`skills/comment-writer/SKILL.md`](skills/comment-writer/SKILL.md) |
| `work-unit-commits` | When splitting implementation work into deliverable commits or chained PRs. | [`skills/work-unit-commits/SKILL.md`](skills/work-unit-commits/SKILL.md) |

## Cybersecurity Skills (gentle-ai-cyber)

Cybersecurity skills are organized into three categories: red-team, blue-team, and SOC.
Select the `cyber` preset during installation (`gentle-ai --preset cyber`) to activate these skills.

Full reference: [Cybersecurity Edition](docs/cybersecurity-edition.md)

### Red-Team

| Skill | Trigger | Path |
|-------|---------|------|
| `pentest-orchestrator` | AI-driven penetration testing orchestration for multi-phase assessments. | [`skills/pentest-orchestrator/SKILL.md`](skills/pentest-orchestrator/SKILL.md) |
| `ai-pentesting-validation` | Anti-hallucination validation pipeline for AI-driven security assessments. | [`skills/ai-pentesting-validation/SKILL.md`](skills/ai-pentesting-validation/SKILL.md) |
| `exploit-chain-patterns` | Exploit chaining methodology for combining vulnerabilities into high-impact attack paths. | [`skills/exploit-chain-patterns/SKILL.md`](skills/exploit-chain-patterns/SKILL.md) |
| `waf-detection-bypass` | Web Application Firewall detection and bypass techniques for authorized testing. | [`skills/waf-detection-bypass/SKILL.md`](skills/waf-detection-bypass/SKILL.md) |

### Blue-Team

| Skill | Trigger | Path |
|-------|---------|------|
| `detection-engineer` | Create detection rules and hunting queries from malware analysis findings. | [`skills/detection-engineer/SKILL.md`](skills/detection-engineer/SKILL.md) |
| `python-security` | Secure Python patterns and best practices for defensive coding. | [`skills/python-security/SKILL.md`](skills/python-security/SKILL.md) |

### SOC

| Skill | Trigger | Path |
|-------|---------|------|
| `malware-triage` | Systematic malware triage and initial assessment workflow. | [`skills/malware-triage/SKILL.md`](skills/malware-triage/SKILL.md) |
| `specialized-file-analyzer` | Analyze specialized file types beyond standard PE executables. | [`skills/specialized-file-analyzer/SKILL.md`](skills/specialized-file-analyzer/SKILL.md) |
| `malware-dynamic-analysis` | Execute and monitor malware in controlled sandbox environments. | [`skills/malware-dynamic-analysis/SKILL.md`](skills/malware-dynamic-analysis/SKILL.md) |
| `malware-report-writer` | Professional malware analysis report creation for enterprise malware analysis and incident response. | [`skills/malware-report-writer/SKILL.md`](skills/malware-report-writer/SKILL.md) |

### SOC Agent: gentleman-soc

`gentleman-soc` is a SOC orchestrator agent definition (not a standalone AgentID). It is installed as a markdown augmentation file that guides incident response through the PICERL pipeline (Prepare → Identify → Contain → Eradicate → Recover → Lessons). It references the 6 SOC/blue-team skills above and maintains MITRE ATT&CK mapping throughout the investigation.

### MCP Integrations

The cyber preset includes three MCP server integrations:

| MCP Server | Permission Tier | Purpose |
|------------|----------------|---------|
| `kali-mcp` | destructive | Offensive security tools (nmap, metasploit, hydra, etc.) — requires human confirmation |
| `virustotal-mcp` | unrestricted | Threat intelligence lookups (read-only) |
| `shodan-mcp` | unrestricted | Internet-wide scanning data (read-only) |

> **Warning**: Destructive tools (kali-mcp) require explicit human confirmation before execution. The permission model is a soft gate (prompt-based), not a hard runtime block.
