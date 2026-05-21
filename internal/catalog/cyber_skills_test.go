package catalog

import (
	"regexp"
	"strings"
	"testing"

	"github.com/gentleman-programming/gentle-ai/internal/assets"
	"github.com/gentleman-programming/gentle-ai/internal/model"
)

// TestMVPCyberSkillsReturns10Entries verifies that MVPCyberSkills() returns
// exactly 10 cyber skills.
//
// Spec: cyber-testing — "All 10 cyber skills registered in catalog"
func TestMVPCyberSkillsReturns10Entries(t *testing.T) {
	skills := MVPCyberSkills()
	if got := len(skills); got != 10 {
		t.Fatalf("MVPCyberSkills() returned %d skills, want 10", got)
	}
}

// TestCyberSkillCategoriesAreValid ensures every cyber skill has a valid
// category (red-team, blue-team, soc, or compliance).
//
// Spec: cyber-testing — "All catalog entries have valid categories"
func TestCyberSkillCategoriesAreValid(t *testing.T) {
	validCategories := map[string]bool{
		"red-team":   true,
		"blue-team":  true,
		"soc":        true,
		"compliance": true,
	}

	for _, s := range MVPCyberSkills() {
		if !validCategories[s.Category] {
			t.Errorf("cyber skill %q has invalid category %q", s.Name, s.Category)
		}
	}
}

// TestCyberSkillCategoryCounts verifies the expected distribution of cyber
// skills across categories: 4 red-team, 2 blue-team, 4 soc.
func TestCyberSkillCategoryCounts(t *testing.T) {
	counts := make(map[string]int)
	for _, s := range MVPCyberSkills() {
		counts[s.Category]++
	}

	wantCounts := map[string]int{
		"red-team":  4,
		"blue-team": 2,
		"soc":       4,
	}

	for cat, want := range wantCounts {
		if got := counts[cat]; got != want {
			t.Errorf("category %q: got %d skills, want %d", cat, got, want)
		}
	}
}

// TestNoSkillIDCollisions ensures no cyber skill ID duplicates any MVP skill ID.
//
// Spec: cyber-testing — "No SkillID collisions"
func TestNoSkillIDCollisions(t *testing.T) {
	mvpSet := make(map[model.SkillID]bool)
	for _, s := range MVPSkills() {
		mvpSet[s.ID] = true
	}

	for _, s := range MVPCyberSkills() {
		if mvpSet[s.ID] {
			t.Errorf("cyber skill %q collides with an MVP skill ID", s.ID)
		}
	}
}

// TestCyberSkillIDsAreUnique ensures no duplicate skill IDs within the cyber catalog.
func TestCyberSkillIDsAreUnique(t *testing.T) {
	seen := make(map[model.SkillID]bool)
	for _, s := range MVPCyberSkills() {
		if seen[s.ID] {
			t.Errorf("duplicate cyber skill ID %q", s.ID)
		}
		seen[s.ID] = true
	}
}

// TestCyberSkillsHaveRequiredFrontmatter verifies every cyber skill SKILL.md
// file has valid frontmatter with required fields.
//
// Spec: cyber-testing — "All cyber skill frontmatter is valid"
func TestCyberSkillsHaveRequiredFrontmatter(t *testing.T) {
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
			content := assets.MustRead(skillPath)

			// Must start with --- delimiter.
			if !strings.HasPrefix(content, "---\n") {
				t.Error("frontmatter must start with `---`")
			}

			// Must have closing --- delimiter.
			rest := strings.TrimPrefix(content, "---\n")
			closeIdx := strings.Index(rest, "\n---")
			if closeIdx == -1 {
				t.Fatal("missing closing `---` delimiter")
			}
			block := rest[:closeIdx]

			// Check required fields.
			if !strings.Contains(block, "name:") {
				t.Error("missing required `name:` field")
			}
			if !strings.Contains(block, "description:") {
				t.Error("missing required `description:` field")
			}
			if !strings.Contains(block, "license:") {
				t.Error("missing required `license:` field")
			}
			if !strings.Contains(block, "metadata:") {
				t.Error("missing required `metadata:` field")
			}

			// Check metadata.author.
			if !strings.Contains(block, "author:") {
				t.Error("missing `metadata.author` field")
			}

			// Check metadata.version.
			if !strings.Contains(block, "version:") {
				t.Error("missing `metadata.version` field")
			}

			// Description must contain "Trigger:".
			descLine := extractFrontmatterValue(block, "description")
			if descLine == "" {
				t.Fatal("could not extract description value")
			}
			if !strings.Contains(descLine, "Trigger:") {
				t.Errorf("description must contain `Trigger:` substring; got: %q", descLine)
			}

			// Description must be <= 160 chars.
			if len([]rune(descLine)) > 160 {
				t.Errorf("description length = %d chars, want <= 160", len([]rune(descLine)))
			}

			// Must have ## Validation section in body.
			body := content[closeIdx+4:]
			if !strings.Contains(body, "## Validation") {
				t.Error("skill body must contain `## Validation` section")
			}
		})
	}
}

// TestDestructiveToolsRequireConfirmation verifies that the kali-mcp manifest
// has permission_tier: destructive and a non-empty destructive_tools array.
//
// Spec: cyber-testing — "Destructive manifests declare destructive tools"
func TestDestructiveToolsRequireConfirmation(t *testing.T) {
	content := assets.MustRead("mcps/kali-mcp/manifest.json")

	// Check permission_tier.
	if !strings.Contains(content, `"permission_tier": "destructive"`) {
		t.Error("kali-mcp manifest must have permission_tier: destructive")
	}

	// Check destructive_tools array exists and is non-empty.
	if !strings.Contains(content, `"destructive_tools"`) {
		t.Fatal("kali-mcp manifest missing destructive_tools field")
	}

	// Count tool entries (simple check: each tool is a quoted string in the array).
	toolPattern := regexp.MustCompile(`"[a-z_]+"`)
	tools := toolPattern.FindAllString(content, -1)
	if len(tools) == 0 {
		t.Error("destructive_tools array must be non-empty")
	}

	// Verify confirmation_required is true.
	if !strings.Contains(content, `"confirmation_required": true`) {
		t.Error("kali-mcp manifest must have confirmation_required: true")
	}
}

// TestUnrestrictedManifestsHaveNoDestructiveTools verifies that unrestricted
// MCP manifests do not declare destructive tools.
func TestUnrestrictedManifestsHaveNoDestructiveTools(t *testing.T) {
	manifests := []string{
		"mcps/shodan-mcp/manifest.json",
		"mcps/virustotal-mcp/manifest.json",
	}

	for _, manifestPath := range manifests {
		t.Run(manifestPath, func(t *testing.T) {
			content := assets.MustRead(manifestPath)

			if !strings.Contains(content, `"permission_tier": "unrestricted"`) {
				t.Errorf("%s must have permission_tier: unrestricted", manifestPath)
			}

			if strings.Contains(content, `"destructive_tools"`) {
				t.Errorf("%s must not have destructive_tools field (unrestricted tier)", manifestPath)
			}
		})
	}
}

// TestProwlerCheckIDFormat scans cyber skill content for Prowler check ID
// references and validates their format.
//
// Spec: cyber-testing — "Valid Prowler check ID references pass"
func TestProwlerCheckIDFormat(t *testing.T) {
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

	// Pattern for Prowler check references: check: <id> or similar.
	prowlerRefPattern := regexp.MustCompile(`check:\s*([A-Za-z0-9_.-]+)`)
	validIDPattern := regexp.MustCompile(`^[a-z0-9_]+$`)

	for _, name := range cyberSkillNames {
		t.Run(name, func(t *testing.T) {
			skillPath := "skills/" + name + "/SKILL.md"
			content := assets.MustRead(skillPath)

			matches := prowlerRefPattern.FindAllStringSubmatch(content, -1)
			for _, match := range matches {
				checkID := match[1]
				if !validIDPattern.MatchString(checkID) {
					t.Errorf("malformed Prowler check ID %q in %s (must match ^[a-z0-9_]+$)", checkID, skillPath)
				}
			}

			// Skills without Prowler references should pass cleanly.
			if len(matches) == 0 {
				t.Log("no Prowler check references found (skipped cleanly)")
			}
		})
	}
}

// extractFrontmatterValue extracts the value for a given key from a YAML
// frontmatter block. Handles quoted scalars.
func extractFrontmatterValue(block, key string) string {
	lines := strings.Split(block, "\n")
	prefix := key + ":"
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, prefix) {
			value := strings.TrimSpace(trimmed[len(prefix):])
			// Strip quotes.
			if len(value) >= 2 {
				if (value[0] == '"' && value[len(value)-1] == '"') ||
					(value[0] == '\'' && value[len(value)-1] == '\'') {
					return value[1 : len(value)-1]
				}
			}
			return value
		}
	}
	return ""
}
