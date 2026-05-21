package skills

import (
	"sort"
	"testing"

	"github.com/gentleman-programming/gentle-ai/internal/model"
)

// TestPresetCyberResolvesToMVPPPlusCyberSkills verifies that
// SkillsForPreset(PresetCyber) returns all MVP skills plus all cyber skills.
//
// Spec: cyber-testing — "Cyber preset includes expected components"
func TestPresetCyberResolvesToMVPPPlusCyberSkills(t *testing.T) {
	skills := SkillsForPreset(model.PresetCyber)

	// Count MVP + cyber skills.
	mvpCount := len(sddSkills) + len(foundationSkills)
	cyberCount := len(cyberSkills)
	wantCount := mvpCount + cyberCount

	if got := len(skills); got != wantCount {
		t.Fatalf("SkillsForPreset(PresetCyber) returned %d skills, want %d (%d MVP + %d cyber)",
			got, wantCount, mvpCount, cyberCount)
	}

	// Verify all cyber skills are present.
	skillSet := make(map[model.SkillID]bool)
	for _, s := range skills {
		skillSet[s] = true
	}

	for _, cyberSkill := range cyberSkills {
		if !skillSet[cyberSkill] {
			t.Errorf("cyber skill %q missing from PresetCyber result", cyberSkill)
		}
	}

	// Verify all MVP skills are present.
	for _, sddSkill := range sddSkills {
		if !skillSet[sddSkill] {
			t.Errorf("SDD skill %q missing from PresetCyber result", sddSkill)
		}
	}
	for _, foundSkill := range foundationSkills {
		if !skillSet[foundSkill] {
			t.Errorf("foundation skill %q missing from PresetCyber result", foundSkill)
		}
	}
}

// TestPresetCyberIsDeterministic verifies that calling SkillsForPreset
// twice returns identical results.
//
// Spec: cyber-testing — "Cyber preset component count is deterministic"
func TestPresetCyberIsDeterministic(t *testing.T) {
	first := SkillsForPreset(model.PresetCyber)
	second := SkillsForPreset(model.PresetCyber)

	if len(first) != len(second) {
		t.Fatalf("length mismatch: first=%d, second=%d", len(first), len(second))
	}

	// Sort both slices for comparison.
	sort.Slice(first, func(i, j int) bool { return first[i] < first[j] })
	sort.Slice(second, func(i, j int) bool { return second[i] < second[j] })

	for i := range first {
		if first[i] != second[i] {
			t.Errorf("result mismatch at index %d: first=%q, second=%q", i, first[i], second[i])
		}
	}
}

// TestNonCyberPresetsUnchanged verifies that existing presets do not
// include cyber skills.
//
// Spec: cyber-testing — "Non-cyber presets are unchanged"
func TestNonCyberPresetsUnchanged(t *testing.T) {
	nonCyberPresets := []model.PresetID{
		model.PresetFullGentleman,
		model.PresetEcosystemOnly,
		model.PresetMinimal,
	}

	for _, preset := range nonCyberPresets {
		t.Run(string(preset), func(t *testing.T) {
			skills := SkillsForPreset(preset)

			for _, s := range skills {
				for _, cyberSkill := range cyberSkills {
					if s == cyberSkill {
						t.Errorf("non-cyber preset %q should not include cyber skill %q", preset, cyberSkill)
					}
				}
			}
		})
	}
}

// TestAllSkillIDsIncludesCyberSkills verifies that AllSkillIDs() returns
// all cyber skill IDs.
func TestAllSkillIDsIncludesCyberSkills(t *testing.T) {
	allIDs := AllSkillIDs()

	idSet := make(map[model.SkillID]bool)
	for _, id := range allIDs {
		idSet[id] = true
	}

	for _, cyberSkill := range cyberSkills {
		if !idSet[cyberSkill] {
			t.Errorf("cyber skill %q missing from AllSkillIDs()", cyberSkill)
		}
	}
}

// TestCyberSkillsSliceHas10Entries verifies the cyberSkills slice in
// presets.go has exactly 10 entries.
func TestCyberSkillsSliceHas10Entries(t *testing.T) {
	if got := len(cyberSkills); got != 10 {
		t.Fatalf("cyberSkills slice has %d entries, want 10", got)
	}
}

// TestCyberSkillsInPresetsAreUnique verifies no duplicate skill IDs in
// the presets.go cyberSkills slice.
func TestCyberSkillsInPresetsAreUnique(t *testing.T) {
	seen := make(map[model.SkillID]bool)
	for _, id := range cyberSkills {
		if seen[id] {
			t.Errorf("duplicate cyber skill ID %q in presets.go cyberSkills", id)
		}
		seen[id] = true
	}
}

// TestNoSkillIDCollisionsBetweenPresetsAndCatalog verifies that the
// cyberSkills in presets.go match the cyberSkills in catalog.go.
func TestNoSkillIDCollisionsBetweenPresetsAndCatalog(t *testing.T) {
	// Build set from presets.go cyberSkills.
	presetSet := make(map[model.SkillID]bool)
	for _, id := range cyberSkills {
		presetSet[id] = true
	}

	// Verify each catalog cyber skill is in presets.
	// We import catalog indirectly by checking the known cyber skill IDs.
	catalogCyberSkills := []model.SkillID{
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

	for _, id := range catalogCyberSkills {
		if !presetSet[id] {
			t.Errorf("catalog cyber skill %q missing from presets.go cyberSkills", id)
		}
	}
}
