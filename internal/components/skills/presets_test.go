package skills

import (
	"testing"

	"github.com/gentleman-programming/gentle-ai/internal/model"
)

func TestSkillsForPresetMinimalReturnsSDDOnly(t *testing.T) {
	skills := SkillsForPreset(model.PresetMinimal)
	if len(skills) == 0 {
		t.Fatalf("SkillsForPreset(minimal) returned empty")
	}

	// Orchestration skills that are always bundled with SDD.
	orchestrationSkills := map[model.SkillID]bool{
		model.SkillJudgmentDay: true,
	}

	for _, skill := range skills {
		isSDD := len(skill) >= 4 && skill[:3] == "sdd"
		if !isSDD && !orchestrationSkills[skill] {
			t.Fatalf("minimal preset should only contain SDD/orchestration skills, got %q", skill)
		}
	}
}

func TestSkillsForPresetEcosystemIncludesFrameworks(t *testing.T) {
	skills := SkillsForPreset(model.PresetEcosystemOnly)

	hasGoTesting := false
	hasSkillCreator := false
	hasSDDInit := false
	for _, skill := range skills {
		if skill == model.SkillGoTesting {
			hasGoTesting = true
		}
		if skill == model.SkillCreator {
			hasSkillCreator = true
		}
		if skill == model.SkillSDDInit {
			hasSDDInit = true
		}
	}

	if !hasGoTesting {
		t.Fatalf("ecosystem preset should include go-testing")
	}
	if !hasSDDInit {
		t.Fatalf("ecosystem preset should include sdd-init")
	}
	if !hasSkillCreator {
		t.Fatalf("ecosystem preset should include skill-creator")
	}
}

func TestSkillsForPresetFullIncludesAll(t *testing.T) {
	skills := SkillsForPreset(model.PresetFullGentleman)
	all := AllSkillIDs()

	// PresetFullGentleman does NOT include cyber skills; AllSkillIDs does.
	// So we check that full preset skills are a subset of all skills.
	skillSet := make(map[model.SkillID]struct{}, len(all))
	for _, s := range all {
		skillSet[s] = struct{}{}
	}
	for _, s := range skills {
		if _, ok := skillSet[s]; !ok {
			t.Fatalf("full preset skill %q not in AllSkillIDs", s)
		}
	}
}

func TestSkillsForPresetCyberIncludesAllMVPPlusCyber(t *testing.T) {
	skills := SkillsForPreset(model.PresetCyber)
	mvpCount := len(sddSkills) + len(foundationSkills)
	cyberCount := len(cyberSkills)
	expected := mvpCount + cyberCount

	if len(skills) != expected {
		t.Fatalf("cyber preset skills len = %d, expected %d (mvp=%d + cyber=%d)",
			len(skills), expected, mvpCount, cyberCount)
	}

	// Verify cyber skills are present.
	cyberSet := make(map[model.SkillID]struct{})
	for _, s := range skills {
		cyberSet[s] = struct{}{}
	}
	for _, cs := range cyberSkills {
		if _, ok := cyberSet[cs]; !ok {
			t.Fatalf("cyber preset missing cyber skill %q", cs)
		}
	}
}

func TestSkillsForPresetCustomReturnsNil(t *testing.T) {
	skills := SkillsForPreset(model.PresetCustom)
	if skills != nil {
		t.Fatalf("custom preset should return nil, got %v", skills)
	}
}

func TestAllSkillIDsIncludesEveryKnownSkill(t *testing.T) {
	all := AllSkillIDs()

	required := []model.SkillID{
		model.SkillSDDInit,
		model.SkillGoTesting,
		model.SkillCreator,
		model.SkillJudgmentDay,
	}

	skillSet := make(map[model.SkillID]struct{}, len(all))
	for _, skill := range all {
		skillSet[skill] = struct{}{}
	}

	for _, req := range required {
		if _, ok := skillSet[req]; !ok {
			t.Fatalf("AllSkillIDs() missing %q", req)
		}
	}
}
