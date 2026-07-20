package commands

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/iamrichardd/pharos-advanced-blocking/internal/config"
)

// helperWriteConfig creates a mock config file.
func helperWriteConfig(t *testing.T, path string, cfg *config.Config) {
	t.Helper()
	bytes, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("failed to marshal config: %v", err)
	}
	err = os.WriteFile(path, bytes, 0644)
	if err != nil {
		t.Fatalf("failed to write config: %v", err)
	}
}

// helperReadConfig reads a config file.
func helperReadConfig(t *testing.T, path string) *config.Config {
	t.Helper()
	bytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read config: %v", err)
	}
	var cfg config.Config
	err = json.Unmarshal(bytes, &cfg)
	if err != nil {
		t.Fatalf("failed to unmarshal config: %v", err)
	}
	return &cfg
}

func TestMapCommand(t *testing.T) {
	t.Run("successfully maps an IP to an existing group", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "dnsApp.config")
		initialConfig := &config.Config{
			Groups: []config.Group{
				{Name: "Default", EnableBlocking: true},
			},
			NetworkGroupMap: map[string]string{
				"192.168.1.10": "Default",
			},
		}
		helperWriteConfig(t, configPath, initialConfig)

		var outBuf, errBuf bytes.Buffer
		rootCmd := NewRootCmd(nil, &outBuf, &errBuf, "dev", "none", "unknown")
		rootCmd.SetArgs([]string{"--config", configPath, "map", "--ip", "192.168.1.20", "--group", "Default"})

		if err := rootCmd.Execute(); err != nil {
			t.Fatalf("unexpected error executing map command: %v", err)
		}

		// Verify output
		if !strings.Contains(outBuf.String(), "Successfully mapped IP 192.168.1.20 to group \"Default\"") {
			t.Errorf("unexpected output: %s", outBuf.String())
		}

		// Verify on-disk config updated
		cfg := helperReadConfig(t, configPath)
		if cfg.NetworkGroupMap["192.168.1.20"] != "Default" {
			t.Errorf("expected 192.168.1.20 to be mapped to Default, got %s", cfg.NetworkGroupMap["192.168.1.20"])
		}
	})

	t.Run("fails when mapping to a non-existent group", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "dnsApp.config")
		initialConfig := &config.Config{
			Groups: []config.Group{
				{Name: "Default", EnableBlocking: true},
			},
		}
		helperWriteConfig(t, configPath, initialConfig)

		var outBuf, errBuf bytes.Buffer
		rootCmd := NewRootCmd(nil, &outBuf, &errBuf, "dev", "none", "unknown")
		rootCmd.SetArgs([]string{"--config", configPath, "map", "--ip", "192.168.1.20", "--group", "NonExistent"})

		err := rootCmd.Execute()
		if err == nil {
			t.Fatal("expected error due to invalid group name, got nil")
		}
		if !strings.Contains(err.Error(), "target group \"NonExistent\" does not exist") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("fails when mapping an invalid IP address", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "dnsApp.config")
		initialConfig := &config.Config{
			Groups: []config.Group{
				{Name: "Default", EnableBlocking: true},
			},
		}
		helperWriteConfig(t, configPath, initialConfig)

		var outBuf, errBuf bytes.Buffer
		rootCmd := NewRootCmd(nil, &outBuf, &errBuf, "dev", "none", "unknown")
		rootCmd.SetArgs([]string{"--config", configPath, "map", "--ip", "invalid-ip", "--group", "Default"})

		err := rootCmd.Execute()
		if err == nil {
			t.Fatal("expected error due to invalid IP format, got nil")
		}
		if !strings.Contains(err.Error(), "must be a valid IPv4 or IPv6 address") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("fails when missing required flags", func(t *testing.T) {
		var outBuf, errBuf bytes.Buffer
		rootCmd := NewRootCmd(nil, &outBuf, &errBuf, "dev", "none", "unknown")
		rootCmd.SetArgs([]string{"map", "--ip", "192.168.1.10"}) // missing --group

		err := rootCmd.Execute()
		if err == nil || !strings.Contains(err.Error(), "missing required flag: --group") {
			t.Errorf("expected missing flag error, got: %v", err)
		}
	})
}

func TestUnmapCommand(t *testing.T) {
	t.Run("successfully removes an existing mapping", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "dnsApp.config")
		initialConfig := &config.Config{
			Groups: []config.Group{
				{Name: "Default", EnableBlocking: true},
			},
			NetworkGroupMap: map[string]string{
				"192.168.1.10": "Default",
			},
		}
		helperWriteConfig(t, configPath, initialConfig)

		var outBuf, errBuf bytes.Buffer
		rootCmd := NewRootCmd(nil, &outBuf, &errBuf, "dev", "none", "unknown")
		rootCmd.SetArgs([]string{"--config", configPath, "unmap", "--ip", "192.168.1.10"})

		if err := rootCmd.Execute(); err != nil {
			t.Fatalf("unexpected error executing unmap: %v", err)
		}

		if !strings.Contains(outBuf.String(), "Successfully unmapped IP 192.168.1.10") {
			t.Errorf("unexpected output: %s", outBuf.String())
		}

		cfg := helperReadConfig(t, configPath)
		if _, exists := cfg.NetworkGroupMap["192.168.1.10"]; exists {
			t.Fatal("expected 192.168.1.10 to be removed from NetworkGroupMap")
		}
	})

	t.Run("fails when IP is not mapped", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "dnsApp.config")
		initialConfig := &config.Config{
			Groups: []config.Group{
				{Name: "Default", EnableBlocking: true},
			},
			NetworkGroupMap: map[string]string{
				"192.168.1.10": "Default",
			},
		}
		helperWriteConfig(t, configPath, initialConfig)

		var outBuf, errBuf bytes.Buffer
		rootCmd := NewRootCmd(nil, &outBuf, &errBuf, "dev", "none", "unknown")
		rootCmd.SetArgs([]string{"--config", configPath, "unmap", "--ip", "192.168.1.99"})

		err := rootCmd.Execute()
		if err == nil {
			t.Fatal("expected error unmapping non-existent IP, got nil")
		}
		if !strings.Contains(err.Error(), "client IP 192.168.1.99 is not mapped") {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestDeployCommand(t *testing.T) {
	t.Run("successfully deploys configuration with --yes flag", func(t *testing.T) {
		// Mock Technitium API Server
		serverCalled := false
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost && r.URL.Path == "/api/apps/config/set" {
				serverCalled = true
				w.Write([]byte(`{"status": "ok"}`))
				return
			}
			w.WriteHeader(http.StatusBadRequest)
		}))
		defer server.Close()

		t.Setenv("TECHNITIUM_URL", server.URL)
		t.Setenv("TECHNITIUM_TOKEN", "mock-token")

		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "dnsApp.config")
		cfg := &config.Config{
			Groups: []config.Group{
				{Name: "Default", EnableBlocking: true},
			},
		}
		helperWriteConfig(t, configPath, cfg)

		var outBuf, errBuf bytes.Buffer
		rootCmd := NewRootCmd(nil, &outBuf, &errBuf, "dev", "none", "unknown")
		rootCmd.SetArgs([]string{"--config", configPath, "deploy", "--yes"})

		if err := rootCmd.Execute(); err != nil {
			t.Fatalf("unexpected deploy error: %v", err)
		}

		if !serverCalled {
			t.Error("expected mock server to be called during deployment")
		}
		if !strings.Contains(outBuf.String(), "Successfully deployed configuration to node") {
			t.Errorf("unexpected output: %s", outBuf.String())
		}
	})

	t.Run("successfully performs dry-run structural diff comparison", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet && r.URL.Path == "/api/apps/config/get" {
				// Return a remote config with a deleted IP and a missing group
				w.Write([]byte(`{"status": "ok", "response": {"config": "{\"groups\":[{\"name\":\"Default\"},{\"name\":\"OldGroup\"}],\"networkGroupMap\":{\"192.168.1.15\":\"Default\"}}"}}`))
				return
			}
			w.WriteHeader(http.StatusBadRequest)
		}))
		defer server.Close()

		t.Setenv("TECHNITIUM_URL", server.URL)
		t.Setenv("TECHNITIUM_TOKEN", "mock-token")

		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "dnsApp.config")
		// Local config has NewGroup instead of OldGroup, and IP 192.168.1.20 instead of 192.168.1.15
		cfg := &config.Config{
			Groups: []config.Group{
				{Name: "Default", EnableBlocking: true},
				{Name: "NewGroup", EnableBlocking: true},
			},
			NetworkGroupMap: map[string]string{
				"192.168.1.20": "NewGroup",
			},
		}
		helperWriteConfig(t, configPath, cfg)

		var outBuf, errBuf bytes.Buffer
		rootCmd := NewRootCmd(nil, &outBuf, &errBuf, "dev", "none", "unknown")
		rootCmd.SetArgs([]string{"--config", configPath, "deploy", "--dry-run"})

		if err := rootCmd.Execute(); err != nil {
			t.Fatalf("unexpected deploy dry-run error: %v", err)
		}

		output := outBuf.String()
		if !strings.Contains(output, `+ Group "NewGroup" (Added)`) {
			t.Errorf("missing added group output in diff: %s", output)
		}
		if !strings.Contains(output, `- Group "OldGroup" (Removed)`) {
			t.Errorf("missing removed group output in diff: %s", output)
		}
		if !strings.Contains(output, `+ 192.168.1.20 -> NewGroup`) {
			t.Errorf("missing added IP output in diff: %s", output)
		}
		if !strings.Contains(output, `- 192.168.1.15 -> Default`) {
			t.Errorf("missing removed IP output in diff: %s", output)
		}
	})

	t.Run("aborts deploy when confirmation is rejected", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "dnsApp.config")
		cfg := &config.Config{
			Groups: []config.Group{
				{Name: "Default", EnableBlocking: true},
			},
		}
		helperWriteConfig(t, configPath, cfg)

		t.Setenv("TECHNITIUM_URL", "http://localhost:5380")
		t.Setenv("TECHNITIUM_TOKEN", "mock-token")

		// Pass 'n' to stdin to reject confirmation
		inBuf := bytes.NewBufferString("n\n")
		var outBuf, errBuf bytes.Buffer
		rootCmd := NewRootCmd(inBuf, &outBuf, &errBuf, "dev", "none", "unknown")
		rootCmd.SetArgs([]string{"--config", configPath, "deploy"})

		if err := rootCmd.Execute(); err != nil {
			t.Fatalf("unexpected deploy error: %v", err)
		}

		if !strings.Contains(outBuf.String(), "Deployment to node \"default\" cancelled.") {
			t.Errorf("unexpected output: %s", outBuf.String())
		}
	})

	t.Run("aborts deploy when secrets file has weaker permissions", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "dnsApp.config")
		cfg := &config.Config{
			Groups: []config.Group{
				{Name: "Default", EnableBlocking: true},
			},
		}
		helperWriteConfig(t, configPath, cfg)

		// Set XDG_CONFIG_HOME to check user config secrets file permission check
		t.Setenv("XDG_CONFIG_HOME", tmpDir)
		pabConfigDir := filepath.Join(tmpDir, "pab")
		err := os.MkdirAll(pabConfigDir, 0755)
		if err != nil {
			t.Fatalf("failed to create mock config dir: %v", err)
		}

		secretsPath := filepath.Join(pabConfigDir, "secrets.json")
		// Write secrets file with weak permissions (0644 instead of 0600)
		err = os.WriteFile(secretsPath, []byte(`{"nodes":{"tech-01":{"url":"http://localhost:5380","token":"123"}}}`), 0644)
		if err != nil {
			t.Fatalf("failed to write secrets file: %v", err)
		}

		// Explicitly chmod to ensure permissions are weaker
		_ = os.Chmod(secretsPath, 0644)

		var outBuf, errBuf bytes.Buffer
		rootCmd := NewRootCmd(nil, &outBuf, &errBuf, "dev", "none", "unknown")
		rootCmd.SetArgs([]string{"--config", configPath, "deploy"})

		err = rootCmd.Execute()
		if err == nil {
			t.Fatal("expected error due to weak secrets file permissions, got nil")
		}

		if !errors.Is(err, config.ErrWeakerPermissions) {
			t.Errorf("expected ErrWeakerPermissions, got: %v", err)
		}
		if !strings.Contains(errBuf.String(), "SECURITY WARNING:") {
			t.Errorf("expected warning in stderr, got: %s", errBuf.String())
		}
	})
}
