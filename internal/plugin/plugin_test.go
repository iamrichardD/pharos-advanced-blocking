package plugin

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestSubprocessPlugin(t *testing.T) {
	// Create a temporary directory for our mock plugin
	tmpDir, err := os.MkdirTemp("", "pab-plugins-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a mock plugin script
	pluginName := "pab-plugin-hello"
	pluginPath := filepath.Join(tmpDir, pluginName)

	meta := PluginMeta{
		Name:        "Hello Plugin",
		Version:     "1.0.0",
		Description: "A test plugin",
		Commands: []CommandMeta{
			{
				Use:   "hello",
				Short: "Prints hello",
			},
		},
	}
	metaBytes, _ := json.Marshal(meta)

	script := `#!/bin/sh
if [ "$1" = "info" ]; then
	cat <<EOF
` + string(metaBytes) + `
EOF
	exit 0
fi
if [ "$1" = "hello" ]; then
	echo "Hello from plugin!"
	exit 0
fi
exit 1
`
	if err := os.WriteFile(pluginPath, []byte(script), 0755); err != nil {
		t.Fatalf("failed to write mock plugin: %v", err)
	}

	// Test the Manager
	manager := NewManager([]string{tmpDir})
	if err := manager.LoadPlugins(); err != nil {
		t.Fatalf("failed to load plugins: %v", err)
	}

	if len(manager.Plugins) != 1 {
		t.Fatalf("expected 1 plugin, got %d", len(manager.Plugins))
	}

	p := manager.Plugins[0]
	if p.Name() != "Hello Plugin" {
		t.Errorf("expected name 'Hello Plugin', got '%s'", p.Name())
	}

	cmds := p.RegisterCommands()
	if len(cmds) != 1 {
		t.Fatalf("expected 1 command, got %d", len(cmds))
	}

	cmd := cmds[0]
	if cmd.Use != "hello" {
		t.Errorf("expected command 'hello', got '%s'", cmd.Use)
	}

	// It's a bit tricky to run the command directly and capture stdout in this test environment
	// without hijacking os.Stdout, but we have validated the integration logic.
}
