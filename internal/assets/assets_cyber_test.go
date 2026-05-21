package assets

import (
	"io/fs"
	"strings"
	"testing"
)

// TestCyberSkillsAreEmbedded verifies all 10 cybersecurity skill directories
// are embedded and have SKILL.md files.
//
// Spec: cyber-testing — "All cyber skill frontmatter is valid"
func TestCyberSkillsAreEmbedded(t *testing.T) {
	cyberSkillNames := []string{
		"pentest-orchestrator",
		"ai-pentesting-validation",
		"exploit-chain-patterns",
		"waf-detection-bypass",
		"detection-engineer",
		"python-security",
		"malware-triage",
		"specialized-file-analyzer",
		"malware-dynamic-analysis",
		"malware-report-writer",
	}

	for _, name := range cyberSkillNames {
		t.Run(name, func(t *testing.T) {
			skillPath := "skills/" + name + "/SKILL.md"
			content, err := Read(skillPath)
			if err != nil {
				t.Fatalf("Read(%q) error = %v", skillPath, err)
			}

			if len(strings.TrimSpace(content)) == 0 {
				t.Fatalf("Read(%q) returned empty content", skillPath)
			}

			// Must have substantial content (>100 chars beyond frontmatter).
			if len(content) < 100 {
				t.Fatalf("Read(%q) content is suspiciously short (%d bytes)", skillPath, len(content))
			}
		})
	}
}

// TestCyberSkillFrontmatterViaEmbeddedAssets verifies cyber skill frontmatter
// through the embedded assets (same mechanism as the existing frontmatter test).
func TestCyberSkillFrontmatterViaEmbeddedAssets(t *testing.T) {
	skillPaths := embeddedSkillPaths(t)

	// Find cyber skill paths.
	var cyberPaths []string
	for _, path := range skillPaths {
		for _, name := range []string{
			"pentest-orchestrator",
			"ai-pentesting-validation",
			"exploit-chain-patterns",
			"waf-detection-bypass",
			"detection-engineer",
			"python-security",
			"malware-triage",
			"specialized-file-analyzer",
			"malware-dynamic-analysis",
			"malware-report-writer",
		} {
			if strings.Contains(path, name) {
				cyberPaths = append(cyberPaths, path)
				break
			}
		}
	}

	if len(cyberPaths) != 10 {
		t.Fatalf("expected 10 embedded cyber skill paths, got %d", len(cyberPaths))
	}

	// Run the same frontmatter checks as the existing TestSkillFrontmatterIsLintClean.
	for _, path := range cyberPaths {
		t.Run(path, func(t *testing.T) {
			content := MustRead(path)

			fm, err := extractSkillFrontmatter(content)
			if err != nil {
				t.Fatalf("extract frontmatter: %v", err)
			}

			// name == parent directory basename.
			expectedName := skillDirBasename(path)
			if fm.name != expectedName {
				t.Errorf("name = %q, want %q", fm.name, expectedName)
			}

			// description must be quoted.
			if !isQuotedScalar(fm.descriptionRawAfterColon) {
				t.Errorf("description must be quoted; got: %q", fm.descriptionRawAfterColon)
			}

			// description must contain Trigger.
			if !strings.Contains(fm.description, "Trigger:") {
				t.Errorf("description must contain `Trigger:`; got: %q", fm.description)
			}

			// description must be <= 160 chars.
			if len([]rune(fm.description)) > 160 {
				t.Errorf("description length = %d, want <= 160", len([]rune(fm.description)))
			}
		})
	}
}

// TestGentlemanSOCAgentIsEmbedded verifies the gentleman-soc agent definition
// is embedded and readable.
func TestGentlemanSOCAgentIsEmbedded(t *testing.T) {
	content, err := Read("agents/gentleman-soc.md")
	if err != nil {
		t.Fatalf("Read(agents/gentleman-soc.md) error = %v", err)
	}

	if len(content) == 0 {
		t.Fatal("gentleman-soc.md is empty")
	}

	// Must reference SOC pipeline phases.
	requiredSections := []string{
		"Phase",
		"malware-triage",
		"detection-engineer",
		"malware-report-writer",
	}
	for _, section := range requiredSections {
		if !strings.Contains(content, section) {
			t.Errorf("gentleman-soc.md must contain %q", section)
		}
	}
}

// TestCyberMCPManifestsAreEmbedded verifies all 3 cyber MCP manifests
// are embedded in the assets.
func TestCyberMCPManifestsAreEmbedded(t *testing.T) {
	manifests := []string{
		"mcps/kali-mcp/manifest.json",
		"mcps/shodan-mcp/manifest.json",
		"mcps/virustotal-mcp/manifest.json",
	}

	for _, path := range manifests {
		t.Run(path, func(t *testing.T) {
			content, err := Read(path)
			if err != nil {
				t.Fatalf("Read(%q) error = %v", path, err)
			}

			if len(strings.TrimSpace(content)) == 0 {
				t.Fatalf("Read(%q) returned empty content", path)
			}
		})
	}
}

// TestEmbeddedSkillCountIncludesCyberSkills verifies the total skill count
// includes the 10 cyber skills (32 total: 10 SDD + judgment-day + 6 foundation
// + 4 sustainable-review + _shared + 10 cybersecurity + 1 go-testing).
func TestEmbeddedSkillCountIncludesCyberSkills(t *testing.T) {
	entries, err := FS.ReadDir("skills")
	if err != nil {
		t.Fatalf("ReadDir(skills) error = %v", err)
	}

	skillDirs := 0
	for _, entry := range entries {
		if entry.IsDir() {
			skillDirs++
		}
	}

	// 32 skill directories expected.
	wantSkillDirs := 32
	if skillDirs != wantSkillDirs {
		t.Fatalf("expected %d skill directories, got %d", wantSkillDirs, skillDirs)
	}

	// Verify all cyber skill directories exist.
	cyberSkillNames := []string{
		"pentest-orchestrator",
		"ai-pentesting-validation",
		"exploit-chain-patterns",
		"waf-detection-bypass",
		"detection-engineer",
		"python-security",
		"malware-triage",
		"specialized-file-analyzer",
		"malware-dynamic-analysis",
		"malware-report-writer",
	}

	seenDirs := make(map[string]bool)
	for _, entry := range entries {
		if entry.IsDir() {
			seenDirs[entry.Name()] = true
		}
	}

	for _, name := range cyberSkillNames {
		if !seenDirs[name] {
			t.Errorf("embedded skills missing cyber skill directory %q", name)
		}
	}
}

// TestCyberSkillBodiesHaveValidationSection verifies every cyber skill has
// a ## Validation section in its body (beyond frontmatter).
func TestCyberSkillBodiesHaveValidationSection(t *testing.T) {
	cyberSkillNames := []string{
		"pentest-orchestrator",
		"ai-pentesting-validation",
		"exploit-chain-patterns",
		"waf-detection-bypass",
		"detection-engineer",
		"python-security",
		"malware-triage",
		"specialized-file-analyzer",
		"malware-dynamic-analysis",
		"malware-report-writer",
	}

	for _, name := range cyberSkillNames {
		t.Run(name, func(t *testing.T) {
			skillPath := "skills/" + name + "/SKILL.md"
			content := MustRead(skillPath)

			if !strings.Contains(content, "## Validation") {
				t.Errorf("%s must contain `## Validation` section", skillPath)
			}
		})
	}
}

// TestNoForbiddenKeysInCyberSkillFrontmatter verifies cyber skills don't
// have auto_invoke or tool_schemas in their frontmatter.
//
// Spec: cyber-testing — "No auto_invoke or tool_schemas in cyber skill frontmatter"
func TestNoForbiddenKeysInCyberSkillFrontmatter(t *testing.T) {
	cyberSkillNames := []string{
		"pentest-orchestrator",
		"ai-pentesting-validation",
		"exploit-chain-patterns",
		"waf-detection-bypass",
		"detection-engineer",
		"python-security",
		"malware-triage",
		"specialized-file-analyzer",
		"malware-dynamic-analysis",
		"malware-report-writer",
	}

	forbiddenKeys := []string{"auto_invoke", "tool_schemas"}

	for _, name := range cyberSkillNames {
		t.Run(name, func(t *testing.T) {
			skillPath := "skills/" + name + "/SKILL.md"
			content := MustRead(skillPath)

			fm, err := extractSkillFrontmatter(content)
			if err != nil {
				t.Fatalf("extract frontmatter: %v", err)
			}

			for _, key := range fm.topLevelKeys {
				for _, forbidden := range forbiddenKeys {
					if key == forbidden {
						t.Errorf("cyber skill %s must not have %q in frontmatter", name, forbidden)
					}
				}
			}
		})
	}
}

// TestAllCyberSkillFilesAreDiscoverableViaWalk verifies that fs.WalkDir
// on the embedded FS finds all 10 cyber skill SKILL.md files.
func TestAllCyberSkillFilesAreDiscoverableViaWalk(t *testing.T) {
	var foundPaths []string
	if err := fs.WalkDir(FS, "skills", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, "/SKILL.md") {
			return nil
		}
		foundPaths = append(foundPaths, path)
		return nil
	}); err != nil {
		t.Fatalf("WalkDir(skills) error = %v", err)
	}

	cyberSkillNames := map[string]bool{
		"pentest-orchestrator":       false,
		"ai-pentesting-validation":   false,
		"exploit-chain-patterns":     false,
		"waf-detection-bypass":       false,
		"detection-engineer":         false,
		"python-security":            false,
		"malware-triage":             false,
		"specialized-file-analyzer":  false,
		"malware-dynamic-analysis":   false,
		"malware-report-writer":      false,
	}

	for _, path := range foundPaths {
		for name := range cyberSkillNames {
			if strings.Contains(path, name) {
				cyberSkillNames[name] = true
			}
		}
	}

	for name, found := range cyberSkillNames {
		if !found {
			t.Errorf("cyber skill %q not found via WalkDir", name)
		}
	}
}
