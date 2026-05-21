package catalog

import "github.com/gentleman-programming/gentle-ai/internal/model"

type Skill struct {
	ID       model.SkillID
	Name     string
	Category string
	Priority string
}

var mvpSkills = []Skill{
	// SDD skills
	{ID: model.SkillSDDInit, Name: "sdd-init", Category: "sdd", Priority: "p0"},

	{ID: model.SkillSDDApply, Name: "sdd-apply", Category: "sdd", Priority: "p0"},
	{ID: model.SkillSDDVerify, Name: "sdd-verify", Category: "sdd", Priority: "p0"},
	{ID: model.SkillSDDExplore, Name: "sdd-explore", Category: "sdd", Priority: "p0"},
	{ID: model.SkillSDDPropose, Name: "sdd-propose", Category: "sdd", Priority: "p0"},
	{ID: model.SkillSDDSpec, Name: "sdd-spec", Category: "sdd", Priority: "p0"},
	{ID: model.SkillSDDDesign, Name: "sdd-design", Category: "sdd", Priority: "p0"},
	{ID: model.SkillSDDTasks, Name: "sdd-tasks", Category: "sdd", Priority: "p0"},
	{ID: model.SkillSDDArchive, Name: "sdd-archive", Category: "sdd", Priority: "p0"},
	{ID: model.SkillSDDOnboard, Name: "sdd-onboard", Category: "sdd", Priority: "p0"},
	// Foundation skills
	{ID: model.SkillGoTesting, Name: "go-testing", Category: "testing", Priority: "p0"},
	{ID: model.SkillCreator, Name: "skill-creator", Category: "workflow", Priority: "p0"},
	{ID: model.SkillImprover, Name: "skill-improver", Category: "workflow", Priority: "p0"},
	{ID: model.SkillJudgmentDay, Name: "judgment-day", Category: "workflow", Priority: "p0"},
	{ID: model.SkillBranchPR, Name: "branch-pr", Category: "workflow", Priority: "p0"},
	{ID: model.SkillIssueCreation, Name: "issue-creation", Category: "workflow", Priority: "p0"},
	{ID: model.SkillSkillRegistry, Name: "skill-registry", Category: "workflow", Priority: "p0"},
	// Sustainable review skills
	{ID: model.SkillChainedPR, Name: "chained-pr", Category: "workflow", Priority: "p0"},
	{ID: model.SkillCognitiveDoc, Name: "cognitive-doc-design", Category: "workflow", Priority: "p0"},
	{ID: model.SkillCommentWriter, Name: "comment-writer", Category: "workflow", Priority: "p0"},
	{ID: model.SkillWorkUnitCommits, Name: "work-unit-commits", Category: "workflow", Priority: "p0"},
}

func MVPSkills() []Skill {
	skills := make([]Skill, len(mvpSkills))
	copy(skills, mvpSkills)
	return skills
}

// cyberSkills are the v1 cybersecurity skills for the cyber preset.
var cyberSkills = []Skill{
	// Red-team
	{ID: model.SkillPentestOrchestrator, Name: "pentest-orchestrator", Category: "red-team", Priority: "p0"},
	{ID: model.SkillAIPentestingValidation, Name: "ai-pentesting-validation", Category: "red-team", Priority: "p0"},
	{ID: model.SkillExploitChainPatterns, Name: "exploit-chain-patterns", Category: "red-team", Priority: "p0"},
	{ID: model.SkillWAFDetectionBypass, Name: "waf-detection-bypass", Category: "red-team", Priority: "p0"},
	// Blue-team
	{ID: model.SkillDetectionEngineer, Name: "detection-engineer", Category: "blue-team", Priority: "p0"},
	{ID: model.SkillSecurityPythonScripts, Name: "python-security", Category: "blue-team", Priority: "p0"},
	// SOC
	{ID: model.SkillMalwareTriage, Name: "malware-triage", Category: "soc", Priority: "p0"},
	{ID: model.SkillSpecializedFileAnalyzer, Name: "specialized-file-analyzer", Category: "soc", Priority: "p0"},
	{ID: model.SkillMalwareDynamicAnalysis, Name: "malware-dynamic-analysis", Category: "soc", Priority: "p0"},
	{ID: model.SkillMalwareReportWriter, Name: "malware-report-writer", Category: "soc", Priority: "p0"},
}

// MVPCyberSkills returns the cybersecurity skills for the cyber preset.
func MVPCyberSkills() []Skill {
	skills := make([]Skill, len(cyberSkills))
	copy(skills, cyberSkills)
	return skills
}
