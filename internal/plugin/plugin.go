package plugin

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// Plugin defines the interface for pab extensions.
type Plugin interface {
	Name() string
	Version() string
	Description() string
	RegisterCommands() []*cobra.Command
}

// SubprocessPlugin implements the Plugin interface for executable sub-process plugins.
// It follows the hashicorp/go-plugin / git-style plugin pattern.
type SubprocessPlugin struct {
	path        string
	name        string
	version     string
	description string
	commands    []CommandMeta
}

// CommandMeta describes a subcommand provided by the plugin.
type CommandMeta struct {
	Use         string `json:"use"`
	Short       string `json:"short"`
	Long        string `json:"long"`
	Example     string `json:"example"`
}

// PluginMeta describes the metadata returned by the plugin's 'info' command.
type PluginMeta struct {
	Name        string        `json:"name"`
	Version     string        `json:"version"`
	Description string        `json:"description"`
	Commands    []CommandMeta `json:"commands"`
}

func (p *SubprocessPlugin) Name() string {
	return p.name
}

func (p *SubprocessPlugin) Version() string {
	return p.version
}

func (p *SubprocessPlugin) Description() string {
	return p.description
}

func (p *SubprocessPlugin) RegisterCommands() []*cobra.Command {
	var cmds []*cobra.Command
	for _, meta := range p.commands {
		metaCopy := meta
		cmd := &cobra.Command{
			Use:     metaCopy.Use,
			Short:   metaCopy.Short,
			Long:    metaCopy.Long,
			Example: metaCopy.Example,
			RunE: func(cmd *cobra.Command, args []string) error {
				// Execute the plugin with the command name and arguments
				pluginArgs := append([]string{metaCopy.Use}, args...)
				execCmd := exec.Command(p.path, pluginArgs...)
				execCmd.Stdout = os.Stdout
				execCmd.Stderr = os.Stderr
				execCmd.Stdin = os.Stdin
				return execCmd.Run()
			},
			// Skip parsing flags so they are passed to the plugin
			DisableFlagParsing: true,
		}
		cmds = append(cmds, cmd)
	}
	return cmds
}

// Manager handles finding and loading plugins.
type Manager struct {
	pluginDirs []string
	Plugins    []Plugin
}

// NewManager creates a new plugin manager looking in the specified directories.
func NewManager(dirs []string) *Manager {
	return &Manager{
		pluginDirs: dirs,
	}
}

// LoadPlugins scans the configured directories for executables starting with 'pab-plugin-'
// and queries them for their metadata.
func (m *Manager) LoadPlugins() error {
	for _, dir := range m.pluginDirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			continue
		}

		entries, err := os.ReadDir(dir)
		if err != nil {
			return fmt.Errorf("failed to read plugin directory %s: %w", dir, err)
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			name := entry.Name()
			if strings.HasPrefix(name, "pab-plugin-") {
				path := filepath.Join(dir, name)
				
				// Ensure it's executable
				info, err := entry.Info()
				if err != nil {
					continue
				}
				if info.Mode()&0111 == 0 {
					continue
				}

				plugin, err := loadSubprocessPlugin(path)
				if err == nil {
					m.Plugins = append(m.Plugins, plugin)
				}
			}
		}
	}
	return nil
}

func loadSubprocessPlugin(path string) (*SubprocessPlugin, error) {
	// Query the plugin for its metadata using the 'info' command
	cmd := exec.Command(path, "info")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to query plugin info: %w", err)
	}

	var meta PluginMeta
	if err := json.Unmarshal(output, &meta); err != nil {
		return nil, fmt.Errorf("invalid plugin metadata: %w", err)
	}

	return &SubprocessPlugin{
		path:        path,
		name:        meta.Name,
		version:     meta.Version,
		description: meta.Description,
		commands:    meta.Commands,
	}, nil
}
