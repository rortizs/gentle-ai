package mcp

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/gentleman-programming/gentle-ai/internal/model"
)

// TestKaliMCPOverlayJSONHasWarningForSeparateFileAgents verifies that
// KaliMCPOverlayJSON prefixes the JSON with a destructive tool warning
// comment for agents using StrategySeparateMCPFiles.
//
// Spec: cyber-mcp — "Destructive tool warning in overlay JSON"
func TestKaliMCPOverlayJSONHasWarningForSeparateFileAgents(t *testing.T) {
	// Claude Code uses StrategySeparateMCPFiles.
	jsonBytes, err := KaliMCPOverlayJSON(model.AgentClaudeCode)
	if err != nil {
		t.Fatalf("KaliMCPOverlayJSON(claude-code) error = %v", err)
	}

	content := string(jsonBytes)

	// Must have the warning comment prefix.
	if !strings.Contains(content, "/**_WARNING**/") {
		t.Error("kali-mcp overlay for Claude Code must contain /**_WARNING**/ comment")
	}

	// Must list destructive tools in the warning.
	destructiveTools := []string{
		"nmap_scan",
		"sqlmap_scan",
		"metasploit_run",
		"hydra_attack",
	}
	for _, tool := range destructiveTools {
		if !strings.Contains(content, tool) {
			t.Errorf("kali-mcp overlay warning must list tool %q", tool)
		}
	}

	// The JSON portion after the comment must still be valid JSON.
	jsonStart := strings.Index(content, "{")
	if jsonStart == -1 {
		t.Fatal("kali-mcp overlay must contain valid JSON after warning comment")
	}
	jsonPart := content[jsonStart:]
	var parsed map[string]any
	if err := json.Unmarshal([]byte(jsonPart), &parsed); err != nil {
		t.Fatalf("JSON portion after warning must be valid JSON: %v", err)
	}

	// Must have command field.
	if parsed["command"] != "kali-server-mcp" {
		t.Errorf("kali-mcp overlay command = %v, want kali-server-mcp", parsed["command"])
	}
}

// TestKaliMCPOverlayJSONNoWarningForMergeAgents verifies that agents using
// StrategyMergeIntoSettings do NOT get the warning comment prefix (the
// warning goes in the system prompt instead).
func TestKaliMCPOverlayJSONNoWarningForMergeAgents(t *testing.T) {
	// OpenCode uses StrategyMergeIntoSettings.
	jsonBytes, err := KaliMCPOverlayJSON(model.AgentOpenCode)
	if err != nil {
		t.Fatalf("KaliMCPOverlayJSON(opencode) error = %v", err)
	}

	content := string(jsonBytes)

	// Must NOT have the warning comment prefix for merge-strategy agents.
	if strings.Contains(content, "/**_WARNING**/") {
		t.Error("kali-mcp overlay for OpenCode must NOT contain warning comment (warning goes in system prompt)")
	}

	// Must have mcpServers wrapper.
	if !strings.Contains(content, `"mcpServers"`) {
		t.Error("kali-mcp overlay for OpenCode must have mcpServers wrapper")
	}

	// Must be valid JSON.
	var parsed map[string]any
	if err := json.Unmarshal(jsonBytes, &parsed); err != nil {
		t.Fatalf("kali-mcp overlay for OpenCode must be valid JSON: %v", err)
	}
}

// TestKaliWarningCommentListsAllDestructiveTools verifies the warning
// comment includes every tool from the manifest.
func TestKaliWarningCommentListsAllDestructiveTools(t *testing.T) {
	manifest, err := loadKaliManifest()
	if err != nil {
		t.Fatalf("loadKaliManifest() error = %v", err)
	}

	comment := kaliWarningComment(manifest.DestructiveTools)

	for _, tool := range manifest.DestructiveTools {
		if !strings.Contains(comment, tool) {
			t.Errorf("warning comment must list tool %q", tool)
		}
	}
}

// TestDestructiveWarningBlockHasCorrectFormat verifies the markdown block
// injected into system prompts has the correct delimiters and content.
//
// Spec: cyber-permissions — "Soft gate via agent prompt"
func TestDestructiveWarningBlockHasCorrectFormat(t *testing.T) {
	block, err := KaliDestructiveWarningBlock()
	if err != nil {
		t.Fatalf("KaliDestructiveWarningBlock() error = %v", err)
	}

	// Must have the opening marker.
	if !strings.Contains(block, "<!-- gentle-ai-cyber:destructive-warning -->") {
		t.Error("warning block must have opening <!-- gentle-ai-cyber:destructive-warning --> marker")
	}

	// Must have the closing marker.
	if !strings.Contains(block, "<!-- /gentle-ai-cyber:destructive-warning -->") {
		t.Error("warning block must have closing <!-- /gentle-ai-cyber:destructive-warning --> marker")
	}

	// Must have the heading.
	if !strings.Contains(block, "## Destructive Tool Warning") {
		t.Error("warning block must have `## Destructive Tool Warning` heading")
	}

	// Must mention confirmation.
	if !strings.Contains(block, "confirm") {
		t.Error("warning block must mention user confirmation")
	}

	// Must list destructive tools as code-formatted items.
	manifest, err := loadKaliManifest()
	if err != nil {
		t.Fatalf("loadKaliManifest() error = %v", err)
	}
	for _, tool := range manifest.DestructiveTools {
		if !strings.Contains(block, "`"+tool+"`") {
			t.Errorf("warning block must list tool %q as code-formatted item", tool)
		}
	}
}

// TestKaliMCPServerJSONIncludesWarning verifies that KaliMCPServerJSON
// returns the server config prefixed with the warning comment.
func TestKaliMCPServerJSONIncludesWarning(t *testing.T) {
	jsonBytes, err := KaliMCPServerJSON()
	if err != nil {
		t.Fatalf("KaliMCPServerJSON() error = %v", err)
	}

	content := string(jsonBytes)

	if !strings.Contains(content, "/**_WARNING**/") {
		t.Error("KaliMCPServerJSON must include warning comment")
	}

	// Must contain the server command.
	if !strings.Contains(content, "kali-server-mcp") {
		t.Error("KaliMCPServerJSON must contain kali-server-mcp command")
	}
}

// TestLoadKaliManifestParsesAllFields verifies the manifest is parsed
// correctly with all required fields.
func TestLoadKaliManifestParsesAllFields(t *testing.T) {
	manifest, err := loadKaliManifest()
	if err != nil {
		t.Fatalf("loadKaliManifest() error = %v", err)
	}

	if manifest.Name != "kali-mcp" {
		t.Errorf("manifest name = %q, want %q", manifest.Name, "kali-mcp")
	}

	if manifest.PermissionTier != "destructive" {
		t.Errorf("manifest permission_tier = %q, want %q", manifest.PermissionTier, "destructive")
	}

	if len(manifest.DestructiveTools) == 0 {
		t.Error("manifest destructive_tools must be non-empty")
	}

	if !manifest.ConfirmationRequired {
		t.Error("manifest confirmation_required must be true")
	}

	// Verify all 10 destructive tools are present.
	wantTools := []string{
		"nmap_scan",
		"sqlmap_scan",
		"hydra_attack",
		"metasploit_run",
		"nikto_scan",
		"dirb_scan",
		"gobuster_scan",
		"wpscan_analyze",
		"john_crack",
		"enum4linux_scan",
	}

	for _, want := range wantTools {
		found := false
		for _, got := range manifest.DestructiveTools {
			if got == want {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("manifest missing destructive tool %q", want)
		}
	}
}

// TestGentlemanSOCContentIsReadable verifies the gentleman-soc agent
// definition can be loaded from embedded assets.
func TestGentlemanSOCContentIsReadable(t *testing.T) {
	content, err := GentlemanSOCContent()
	if err != nil {
		t.Fatalf("GentlemanSOCContent() error = %v", err)
	}

	if len(content) == 0 {
		t.Fatal("GentlemanSOCContent() returned empty content")
	}

	// Must reference PICERL phases (uses "Preparation", "Identification", "Containment", "Recovery").
	picerlPhases := []string{"PICERL", "Preparation", "Identification", "Containment", "Recovery"}
	for _, phase := range picerlPhases {
		if !strings.Contains(content, phase) {
			t.Errorf("gentleman-soc.md must reference PICERL phase %q", phase)
		}
	}

	// Must reference cyber skills.
	cyberSkillRefs := []string{
		"malware-triage",
		"specialized-file-analyzer",
		"detection-engineer",
		"malware-report-writer",
	}
	for _, ref := range cyberSkillRefs {
		if !strings.Contains(content, ref) {
			t.Errorf("gentleman-soc.md must reference skill %q", ref)
		}
	}
}

// TestNewAgentAdapterReturnsStrategyForAllKnownAgents verifies that
// newAgentAdapter returns a valid strategy for every known agent.
func TestNewAgentAdapterReturnsStrategyForAllKnownAgents(t *testing.T) {
	knownAgents := []model.AgentID{
		model.AgentClaudeCode,
		model.AgentOpenCode,
		model.AgentKilocode,
		model.AgentGeminiCLI,
		model.AgentCursor,
		model.AgentVSCodeCopilot,
		model.AgentCodex,
		model.AgentAntigravity,
		model.AgentWindsurf,
		model.AgentKimi,
		model.AgentQwenCode,
		model.AgentKiroIDE,
		model.AgentOpenClaw,
		model.AgentPi,
	}

	for _, agent := range knownAgents {
		t.Run(string(agent), func(t *testing.T) {
			adapter, err := newAgentAdapter(agent)
			if err != nil {
				t.Fatalf("newAgentAdapter(%q) error = %v", agent, err)
			}
			if adapter == nil {
				t.Fatal("newAgentAdapter returned nil adapter for known agent")
			}
			// Strategy must be a valid value.
			strategy := adapter.MCPStrategy()
			if strategy < model.StrategySeparateMCPFiles || strategy > model.StrategyTOMLFile {
				t.Errorf("invalid MCP strategy %d for agent %q", strategy, agent)
			}
		})
	}
}

// TestNewAgentAdapterErrorsForUnknownAgent verifies that an unknown agent
// returns an error.
func TestNewAgentAdapterErrorsForUnknownAgent(t *testing.T) {
	_, err := newAgentAdapter(model.AgentID("unknown-agent"))
	if err == nil {
		t.Fatal("newAgentAdapter(unknown-agent) should return error, got nil")
	}
}
