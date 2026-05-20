package skills

import "github.com/gentleman-programming/gentle-ai/internal/model"

// sddSkills are the SDD orchestrator skills — always included.
var sddSkills = []model.SkillID{
	model.SkillSDDInit,
	model.SkillSDDExplore,
	model.SkillSDDPropose,
	model.SkillSDDSpec,
	model.SkillSDDDesign,
	model.SkillSDDTasks,
	model.SkillSDDApply,
	model.SkillSDDVerify,
	model.SkillSDDArchive,
	model.SkillSDDOnboard,
	model.SkillJudgmentDay,
}

// foundationSkills are baseline learning skills for the "recommended" tier.
var foundationSkills = []model.SkillID{
	model.SkillGoTesting,
	model.SkillCreator,
	model.SkillImprover,
	model.SkillBranchPR,
	model.SkillIssueCreation,
	model.SkillSkillRegistry,
	model.SkillChainedPR,
	model.SkillCognitiveDoc,
	model.SkillCommentWriter,
	model.SkillWorkUnitCommits,
}

// cyberSkills are the v1 cybersecurity skills for the cyber preset.
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

// SkillsForPreset returns which skills should be installed for a given preset.
//
//   - "minimal" / PresetMinimal:       SDD skills only
//   - "ecosystem-only" / PresetEcosystemOnly: SDD + common framework skills
//   - "full-gentleman" / PresetFullGentleman: all available skills
//   - "cyber" / PresetCyber:           all MVP skills + cybersecurity skills
//   - "custom" / PresetCustom:         empty (caller should provide explicit list)
func SkillsForPreset(preset model.PresetID) []model.SkillID {
	switch preset {
	case model.PresetMinimal:
		return copySkills(sddSkills)
	case model.PresetEcosystemOnly:
		return copySkills(append(sddSkills, foundationSkills...))
	case model.PresetFullGentleman:
		all := make([]model.SkillID, 0, len(sddSkills)+len(foundationSkills))
		all = append(all, sddSkills...)
		all = append(all, foundationSkills...)
		return all
	case model.PresetCyber:
		all := make([]model.SkillID, 0, len(sddSkills)+len(foundationSkills)+len(cyberSkills))
		all = append(all, sddSkills...)
		all = append(all, foundationSkills...)
		all = append(all, cyberSkills...)
		return all
	case model.PresetCustom:
		return nil
	default:
		// Unknown preset — default to full.
		all := make([]model.SkillID, 0, len(sddSkills)+len(foundationSkills))
		all = append(all, sddSkills...)
		all = append(all, foundationSkills...)
		return all
	}
}

// AllSkillIDs returns every known skill ID.
func AllSkillIDs() []model.SkillID {
	all := make([]model.SkillID, 0, len(sddSkills)+len(foundationSkills)+len(cyberSkills))
	all = append(all, sddSkills...)
	all = append(all, foundationSkills...)
	all = append(all, cyberSkills...)
	return all
}

func copySkills(src []model.SkillID) []model.SkillID {
	dst := make([]model.SkillID, len(src))
	copy(dst, src)
	return dst
}
