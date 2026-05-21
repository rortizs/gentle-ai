package mcp

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gentleman-programming/gentle-ai/internal/assets"
	"github.com/gentleman-programming/gentle-ai/internal/components/filemerge"
	"github.com/gentleman-programming/gentle-ai/internal/model"
)

// kaliMCPManifest holds the parsed kali-mcp manifest.json.
type kaliMCPManifest struct {
	Name                 string   `json:"name"`
	Type                 string   `json:"type"`
	Description          string   `json:"description"`
	PermissionTier       string   `json:"permission_tier"`
	DestructiveTools     []string `json:"destructive_tools"`
	ConfirmationRequired bool     `json:"confirmation_required"`
}

// loadKaliManifest reads and parses the kali-mcp manifest.json from embedded assets.
func loadKaliManifest() (kaliMCPManifest, error) {
	content, err := assets.FS.ReadFile("mcps/kali-mcp/manifest.json")
	if err != nil {
		return kaliMCPManifest{}, fmt.Errorf("read kali-mcp manifest: %w", err)
	}

	var manifest kaliMCPManifest
	if err := json.Unmarshal(content, &manifest); err != nil {
		return kaliMCPManifest{}, fmt.Errorf("parse kali-mcp manifest: %w", err)
	}

	return manifest, nil
}

// kaliWarningComment returns the warning comment block for destructive tools.
// The comment is prefixed to MCP JSON files for agents using StrategySeparateMCPFiles.
func kaliWarningComment(tools []string) string {
	comment := "/**_WARNING**/\n"
	comment += " * DESTRUCTIVE TOOLS: The following kali-mcp tools can cause real damage\n"
	comment += " * to production systems. NEVER auto-invoke these tools without explicit\n"
	comment += " * human confirmation:\n"
	for _, tool := range tools {
		comment += fmt.Sprintf(" *   - %s\n", tool)
	}
	comment += " * Always confirm with the user before executing any of these tools.\n"
	comment += " */\n"
	return comment
}

// destructiveWarningBlock returns the markdown block injected into agent system prompts
// when the cyber preset is selected.
func destructiveWarningBlock(tools []string) string {
	block := "<!-- gentle-ai-cyber:destructive-warning -->\n\n"
	block += "## Destructive Tool Warning\n\n"
	block += "Before executing any tool marked **DESTRUCTIVE**, confirm with the user. "
	block += "Never auto-invoke these tools without explicit human approval.\n\n"
	block += "The following kali-mcp tools require manual confirmation:\n\n"
	for _, tool := range tools {
		block += fmt.Sprintf("- `%s`\n", tool)
	}
	block += "\n<!-- /gentle-ai-cyber:destructive-warning -->\n"
	return block
}

// KaliMCPOverlayJSON returns the overlay JSON for kali-mcp with the appropriate
// strategy for the given agent. For StrategySeparateMCPFiles, the JSON is prefixed
// with a destructive tool warning comment.
func KaliMCPOverlayJSON(agent model.AgentID) ([]byte, error) {
	manifest, err := loadKaliManifest()
	if err != nil {
		return nil, err
	}

	// Build the base server config.
	serverJSON := map[string]any{
		"command": "kali-server-mcp",
		"args":    []string{},
	}

	// For StrategySeparateMCPFiles, prefix with warning comment.
	adapter, err := newAgentAdapter(agent)
	if err != nil {
		return nil, err
	}

	if adapter != nil && adapter.MCPStrategy() == model.StrategySeparateMCPFiles {
		warning := kaliWarningComment(manifest.DestructiveTools)
		jsonBytes, err := json.MarshalIndent(serverJSON, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("marshal kali-mcp server JSON: %w", err)
		}
		commented := warning + string(jsonBytes) + "\n"
		return []byte(commented), nil
	}

	// For other strategies, return standard overlay.
	overlay := map[string]any{
		"mcpServers": map[string]any{
			"kali-mcp": serverJSON,
		},
	}

	jsonBytes, err := json.MarshalIndent(overlay, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal kali-mcp overlay JSON: %w", err)
	}

	return append(jsonBytes, '\n'), nil
}

// KaliMCPServerJSON returns the raw server config JSON for kali-mcp (without mcpServers wrapper).
func KaliMCPServerJSON() ([]byte, error) {
	manifest, err := loadKaliManifest()
	if err != nil {
		return nil, err
	}

	serverJSON := map[string]any{
		"command": "kali-server-mcp",
		"args":    []string{},
	}

	jsonBytes, err := json.MarshalIndent(serverJSON, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal kali-mcp server JSON: %w", err)
	}

	// For StrategySeparateMCPFiles, prefix with warning comment.
	warning := kaliWarningComment(manifest.DestructiveTools)
	return []byte(warning + string(jsonBytes) + "\n"), nil
}

// KaliDestructiveWarningBlock returns the markdown warning block for system prompt injection.
func KaliDestructiveWarningBlock() (string, error) {
	manifest, err := loadKaliManifest()
	if err != nil {
		return "", err
	}
	return destructiveWarningBlock(manifest.DestructiveTools), nil
}

// newAgentAdapter creates a minimal adapter for MCP strategy detection.
// Returns nil if the agent is not supported.
func newAgentAdapter(agent model.AgentID) (*minimalAdapter, error) {
	// MCP strategy map based on known agent configurations.
	strategies := map[model.AgentID]model.MCPStrategy{
		model.AgentClaudeCode:    model.StrategySeparateMCPFiles,
		model.AgentOpenCode:      model.StrategyMergeIntoSettings,
		model.AgentKilocode:      model.StrategyMergeIntoSettings,
		model.AgentGeminiCLI:     model.StrategyMergeIntoSettings,
		model.AgentCursor:        model.StrategyMCPConfigFile,
		model.AgentVSCodeCopilot: model.StrategyMCPConfigFile,
		model.AgentCodex:         model.StrategyTOMLFile,
		model.AgentAntigravity:   model.StrategyMCPConfigFile,
		model.AgentWindsurf:      model.StrategyMergeIntoSettings,
		model.AgentKimi:          model.StrategyMergeIntoSettings,
		model.AgentQwenCode:      model.StrategyMergeIntoSettings,
		model.AgentKiroIDE:       model.StrategyMergeIntoSettings,
		model.AgentOpenClaw:      model.StrategyMergeIntoSettings,
		model.AgentPi:            model.StrategyMergeIntoSettings,
	}

	s, ok := strategies[agent]
	if !ok {
		return nil, fmt.Errorf("unknown agent %q", agent)
	}

	return &minimalAdapter{strategy: s}, nil
}

type minimalAdapter struct {
	strategy model.MCPStrategy
}

func (m *minimalAdapter) MCPStrategy() model.MCPStrategy {
	return m.strategy
}

// InjectKaliMCP injects kali-mcp configuration for the given agent.
// Returns the injection result with changed status and affected files.
func InjectKaliMCP(homeDir string, adapter interface {
	MCPStrategy() model.MCPStrategy
	MCPConfigPath(homeDir, name string) string
	SettingsPath(homeDir string) string
}) (InjectionResult, error) {
	manifest, err := loadKaliManifest()
	if err != nil {
		return InjectionResult{}, err
	}

	serverJSON := map[string]any{
		"command": "kali-server-mcp",
		"args":    []string{},
	}

	switch adapter.MCPStrategy() {
	case model.StrategySeparateMCPFiles:
		return injectKaliSeparateFile(homeDir, adapter, serverJSON, manifest.DestructiveTools)
	case model.StrategyMergeIntoSettings:
		return injectKaliMergeIntoSettings(homeDir, adapter, serverJSON)
	default:
		return InjectionResult{}, nil
	}
}

func injectKaliSeparateFile(homeDir string, adapter interface {
	MCPConfigPath(homeDir, name string) string
}, serverJSON map[string]any, destructiveTools []string) (InjectionResult, error) {
	path := adapter.MCPConfigPath(homeDir, "kali-mcp")

	jsonBytes, err := json.MarshalIndent(serverJSON, "", "  ")
	if err != nil {
		return InjectionResult{}, fmt.Errorf("marshal kali-mcp server JSON: %w", err)
	}

	warning := kaliWarningComment(destructiveTools)
	content := []byte(warning + string(jsonBytes) + "\n")

	writeResult, err := filemerge.WriteFileAtomic(path, content, 0o644)
	if err != nil {
		return InjectionResult{}, err
	}

	return InjectionResult{Changed: writeResult.Changed, Files: []string{path}}, nil
}

func injectKaliMergeIntoSettings(homeDir string, adapter interface {
	SettingsPath(homeDir string) string
}, serverJSON map[string]any) (InjectionResult, error) {
	settingsPath := adapter.SettingsPath(homeDir)
	if settingsPath == "" {
		return InjectionResult{}, nil
	}

	overlay := map[string]any{
		"mcpServers": map[string]any{
			"kali-mcp": serverJSON,
		},
	}

	overlayBytes, err := json.MarshalIndent(overlay, "", "  ")
	if err != nil {
		return InjectionResult{}, fmt.Errorf("marshal kali-mcp overlay: %w", err)
	}
	overlayBytes = append(overlayBytes, '\n')

	writeResult, err := mergeJSONFile(settingsPath, overlayBytes)
	if err != nil {
		return InjectionResult{}, err
	}

	return InjectionResult{Changed: writeResult.Changed, Files: []string{settingsPath}}, nil
}

// InjectShodanMCP injects shodan-mcp configuration for the given agent.
func InjectShodanMCP(homeDir string, adapter interface {
	MCPStrategy() model.MCPStrategy
	MCPConfigPath(homeDir, name string) string
	SettingsPath(homeDir string) string
}) (InjectionResult, error) {
	serverJSON := map[string]any{
		"command": "npx",
		"args":    []string{"-y", "@modelcontextprotocol/server-shodan"},
	}

	switch adapter.MCPStrategy() {
	case model.StrategySeparateMCPFiles:
		path := adapter.MCPConfigPath(homeDir, "shodan-mcp")
		jsonBytes, err := json.MarshalIndent(serverJSON, "", "  ")
		if err != nil {
			return InjectionResult{}, err
		}
		writeResult, err := filemerge.WriteFileAtomic(path, append(jsonBytes, '\n'), 0o644)
		if err != nil {
			return InjectionResult{}, err
		}
		return InjectionResult{Changed: writeResult.Changed, Files: []string{path}}, nil
	case model.StrategyMergeIntoSettings:
		settingsPath := adapter.SettingsPath(homeDir)
		if settingsPath == "" {
			return InjectionResult{}, nil
		}
		overlay := map[string]any{
			"mcpServers": map[string]any{
				"shodan-mcp": serverJSON,
			},
		}
		overlayBytes, err := json.MarshalIndent(overlay, "", "  ")
		if err != nil {
			return InjectionResult{}, err
		}
		writeResult, err := mergeJSONFile(settingsPath, append(overlayBytes, '\n'))
		if err != nil {
			return InjectionResult{}, err
		}
		return InjectionResult{Changed: writeResult.Changed, Files: []string{settingsPath}}, nil
	default:
		return InjectionResult{}, nil
	}
}

// InjectVirusTotalMCP injects virustotal-mcp configuration for the given agent.
func InjectVirusTotalMCP(homeDir string, adapter interface {
	MCPStrategy() model.MCPStrategy
	MCPConfigPath(homeDir, name string) string
	SettingsPath(homeDir string) string
}) (InjectionResult, error) {
	serverJSON := map[string]any{
		"command": "npx",
		"args":    []string{"-y", "@modelcontextprotocol/server-virustotal"},
	}

	switch adapter.MCPStrategy() {
	case model.StrategySeparateMCPFiles:
		path := adapter.MCPConfigPath(homeDir, "virustotal-mcp")
		jsonBytes, err := json.MarshalIndent(serverJSON, "", "  ")
		if err != nil {
			return InjectionResult{}, err
		}
		writeResult, err := filemerge.WriteFileAtomic(path, append(jsonBytes, '\n'), 0o644)
		if err != nil {
			return InjectionResult{}, err
		}
		return InjectionResult{Changed: writeResult.Changed, Files: []string{path}}, nil
	case model.StrategyMergeIntoSettings:
		settingsPath := adapter.SettingsPath(homeDir)
		if settingsPath == "" {
			return InjectionResult{}, nil
		}
		overlay := map[string]any{
			"mcpServers": map[string]any{
				"virustotal-mcp": serverJSON,
			},
		}
		overlayBytes, err := json.MarshalIndent(overlay, "", "  ")
		if err != nil {
			return InjectionResult{}, err
		}
		writeResult, err := mergeJSONFile(settingsPath, append(overlayBytes, '\n'))
		if err != nil {
			return InjectionResult{}, err
		}
		return InjectionResult{Changed: writeResult.Changed, Files: []string{settingsPath}}, nil
	default:
		return InjectionResult{}, nil
	}
}

// GentlemanSOCContent returns the gentleman-soc agent definition content from embedded assets.
func GentlemanSOCContent() (string, error) {
	content, err := assets.FS.ReadFile("agents/gentleman-soc.md")
	if err != nil {
		return "", fmt.Errorf("read gentleman-soc agent definition: %w", err)
	}
	return string(content), nil
}

// InjectGentlemanSOC injects the gentleman-soc agent definition into the system prompt.
// Only applies to agents using StrategyMarkdownSections (e.g., Claude Code).
func InjectGentlemanSOC(homeDir string, adapter interface {
	SystemPromptFile(homeDir string) string
}) (InjectionResult, error) {
	content, err := GentlemanSOCContent()
	if err != nil {
		return InjectionResult{}, err
	}

	promptPath := adapter.SystemPromptFile(homeDir)
	existing, err := osReadFileForString(promptPath)
	if err != nil {
		return InjectionResult{}, err
	}

	updated := filemerge.InjectMarkdownSection(existing, "gentleman-soc", content)
	writeResult, err := filemerge.WriteFileAtomic(promptPath, []byte(updated), 0o644)
	if err != nil {
		return InjectionResult{}, err
	}

	return InjectionResult{Changed: writeResult.Changed, Files: []string{promptPath}}, nil
}

// InjectDestructiveWarning injects the destructive tool warning into the system prompt.
// Only applies to agents using StrategyMarkdownSections (e.g., Claude Code).
func InjectDestructiveWarning(homeDir string, adapter interface {
	SystemPromptFile(homeDir string) string
}) (InjectionResult, error) {
	warning, err := KaliDestructiveWarningBlock()
	if err != nil {
		return InjectionResult{}, err
	}

	promptPath := adapter.SystemPromptFile(homeDir)
	existing, err := osReadFileForString(promptPath)
	if err != nil {
		return InjectionResult{}, err
	}

	updated := filemerge.InjectMarkdownSection(existing, "gentle-ai-cyber:destructive-warning", warning)
	writeResult, err := filemerge.WriteFileAtomic(promptPath, []byte(updated), 0o644)
	if err != nil {
		return InjectionResult{}, err
	}

	return InjectionResult{Changed: writeResult.Changed, Files: []string{promptPath}}, nil
}

func osReadFileForString(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", fmt.Errorf("read file %q: %w", path, err)
	}
	return string(content), nil
}
