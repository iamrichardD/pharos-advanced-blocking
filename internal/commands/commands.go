package commands

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"

	"github.com/iamrichardd/pharos-advanced-blocking/internal/client"
	"github.com/iamrichardd/pharos-advanced-blocking/internal/config"
	"github.com/iamrichardd/pharos-advanced-blocking/internal/plugin"
	"github.com/spf13/cobra"
)

// SecretsConfig defines the structure for local secrets file ~/.config/pab/secrets.json
type SecretsConfig struct {
	Nodes map[string]NodeConfig `json:"nodes"`
}

// NodeConfig defines the URL and token for a Technitium node.
type NodeConfig struct {
	URL   string `json:"url"`
	Token string `json:"token"`
}

// GlobalFlags stores global CLI flag values.
type GlobalFlags struct {
	ConfigFile string
}

// NewRootCmd constructs and registers all CLI commands.
// It allows passing stdin, stdout, and stderr for flexible CLI testing.
func NewRootCmd(stdin io.Reader, stdout, stderr io.Writer, version, commit, date string) *cobra.Command {
	global := &GlobalFlags{}

	rootCmd := &cobra.Command{
		Use:     "pab",
		Short:   "Pharos Advanced Blocking CLI",
		Long:    `A command-line interface to manage, validate, and sync Technitium Advanced Blocking configurations.`,
		Version: fmt.Sprintf("%s, commit %s, built at %s", version, commit, date),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Strict Credentials Guard check on startup.
			// Abort immediately with a high-priority warning if permissions are weaker than 0600.
			configDir, err := os.UserConfigDir()
			if err == nil {
				secretsPath := filepath.Join(configDir, "pab", "secrets.json")
				err = config.VerifyCredentialsFile(secretsPath)
				if err != nil {
					if errors.Is(err, config.ErrWeakerPermissions) {
						// Print warning and return error to abort command run
						fmt.Fprintf(stderr, "SECURITY WARNING: %v\n", err)
						return err
					}
					// ErrCredentialsNotFound is ok, we fall back to environment variables
				}
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(stdout, "Welcome to Pharos Advanced Blocking (pab)!")
		},
	}

	rootCmd.SetIn(stdin)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)

	rootCmd.PersistentFlags().StringVarP(&global.ConfigFile, "config", "c", "dnsApp.config", "Path to Advanced Blocking configuration file")

	// Register subcommands
	rootCmd.AddCommand(newMapCmd(global, stdout, stderr))
	rootCmd.AddCommand(newUnmapCmd(global, stdout, stderr))
	rootCmd.AddCommand(newDeployCmd(global, stdin, stdout, stderr))

	// Load Plugins
	configDir, err := os.UserConfigDir()
	if err == nil {
		pluginDir := filepath.Join(configDir, "pab", "plugins")
		manager := plugin.NewManager([]string{pluginDir, "./plugins"})
		// Best-effort loading, ignore errors
		_ = manager.LoadPlugins()
		for _, p := range manager.Plugins {
			for _, pCmd := range p.RegisterCommands() {
				// Set output streams for plugin commands if applicable
				pCmd.SetIn(stdin)
				pCmd.SetOut(stdout)
				pCmd.SetErr(stderr)
				rootCmd.AddCommand(pCmd)
			}
		}
	}

	return rootCmd
}

// newMapCmd defines the "pab map --ip <IP> --group <GROUP>" command.
func newMapCmd(global *GlobalFlags, stdout, stderr io.Writer) *cobra.Command {
	var ip string
	var group string

	cmd := &cobra.Command{
		Use:   "map",
		Short: "Map/reassign a client IP to a blocking group",
		Long:  `Adds or updates a client IP mapping in the configuration.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if ip == "" {
				return errors.New("missing required flag: --ip")
			}
			if group == "" {
				return errors.New("missing required flag: --group")
			}

			// Load existing configuration from disk
			cfg, err := loadConfig(global.ConfigFile)
			if err != nil {
				return err
			}

			// Initialize NetworkGroupMap if nil
			if cfg.NetworkGroupMap == nil {
				cfg.NetworkGroupMap = make(map[string]string)
			}

			// Update the client IP mapping
			cfg.NetworkGroupMap[ip] = group

			// Run validation engine on the updated configuration
			if err := config.ValidateConfig(cfg); err != nil {
				return fmt.Errorf("configuration validation failed: %w", err)
			}

			// Save back to disk
			if err := saveConfig(global.ConfigFile, cfg); err != nil {
				return err
			}

			fmt.Fprintf(stdout, "Successfully mapped IP %s to group %q\n", ip, group)
			return nil
		},
	}

	cmd.Flags().StringVar(&ip, "ip", "", "Client IP or CIDR range")
	cmd.Flags().StringVar(&group, "group", "", "Target blocking group name")

	return cmd
}

// newUnmapCmd defines the "pab unmap --ip <IP>" command.
func newUnmapCmd(global *GlobalFlags, stdout, stderr io.Writer) *cobra.Command {
	var ip string

	cmd := &cobra.Command{
		Use:   "unmap",
		Short: "Delete a client IP mapping from the configuration",
		Long:  `Deletes a client IP mapping from the configuration.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if ip == "" {
				return errors.New("missing required flag: --ip")
			}

			cfg, err := loadConfig(global.ConfigFile)
			if err != nil {
				return err
			}

			if cfg.NetworkGroupMap == nil {
				return fmt.Errorf("client IP %s is not mapped (no network group map exists)", ip)
			}

			// Check if mapping exists
			if _, exists := cfg.NetworkGroupMap[ip]; !exists {
				return fmt.Errorf("client IP %s is not mapped in the configuration", ip)
			}

			// Remove mapping
			delete(cfg.NetworkGroupMap, ip)

			// Run validation
			if err := config.ValidateConfig(cfg); err != nil {
				return fmt.Errorf("configuration validation failed: %w", err)
			}

			// Save config
			if err := saveConfig(global.ConfigFile, cfg); err != nil {
				return err
			}

			fmt.Fprintf(stdout, "Successfully unmapped IP %s\n", ip)
			return nil
		},
	}

	cmd.Flags().StringVar(&ip, "ip", "", "Client IP or CIDR range")

	return cmd
}

// newDeployCmd defines the "pab deploy [--dry-run] [--yes] [--node <name>]" command.
func newDeployCmd(global *GlobalFlags, stdin io.Reader, stdout, stderr io.Writer) *cobra.Command {
	var dryRun bool
	var yes bool
	var targetNode string

	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Sync configuration to Technitium server API nodes",
		Long:  `Syncs the validated configuration to Technitium server API nodes.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load and validate local configuration
			cfg, err := loadConfig(global.ConfigFile)
			if err != nil {
				return err
			}

			if err := config.ValidateConfig(cfg); err != nil {
				return fmt.Errorf("local configuration is invalid: %w", err)
			}

			// Marshal configuration to JSON
			configJSONBytes, err := json.Marshal(cfg)
			if err != nil {
				return fmt.Errorf("failed to marshal local configuration: %w", err)
			}
			configJSON := string(configJSONBytes)

			// Resolve Technitium target nodes and credentials
			nodes, err := resolveNodes()
			if err != nil {
				return err
			}

			if len(nodes) == 0 {
				return errors.New("no Technitium nodes configured. Set environment variables or define them in ~/.config/pab/secrets.json")
			}

			// If target node is specified, verify it exists and narrow down targets
			targets := make(map[string]NodeConfig)
			if targetNode != "" {
				node, exists := nodes[targetNode]
				if !exists {
					return fmt.Errorf("target node %q is not defined in the configuration or environment", targetNode)
				}
				targets[targetNode] = node
			} else {
				targets = nodes
			}

			// Process nodes
			for name, node := range targets {
				c := client.NewClient(node.URL, node.Token)

				if dryRun {
					fmt.Fprintf(stdout, "Checking configuration on node %q...\n", name)
					remoteConfigJSON, err := c.GetAppConfig()
					if err != nil {
						fmt.Fprintf(stderr, "Warning: failed to fetch remote configuration for dry-run comparison: %v\n", err)
						continue
					}

					var remoteCfg config.Config
					if remoteConfigJSON != "" {
						_ = json.Unmarshal([]byte(remoteConfigJSON), &remoteCfg)
					}

					// Print structural diff
					printStructuralDiff(stdout, name, &remoteCfg, cfg)
					continue
				}

				// If not dry-run, ask for confirmation unless --yes is passed
				if !yes {
					fmt.Fprintf(stdout, "Deploy local configuration to node %q? (y/N): ", name)
					reader := bufio.NewReader(stdin)
					text, err := reader.ReadString('\n')
					if err != nil {
						return fmt.Errorf("failed to read confirmation: %w", err)
					}
					text = strings.ToLower(strings.TrimSpace(text))
					if text != "y" && text != "yes" {
						fmt.Fprintf(stdout, "Deployment to node %q cancelled.\n", name)
						continue
					}
				}

				fmt.Fprintf(stdout, "Syncing configuration to node %q...\n", name)
				if err := c.SetAppConfig(configJSON); err != nil {
					return fmt.Errorf("failed to deploy to node %q: %w", name, err)
				}
				fmt.Fprintf(stdout, "Successfully deployed configuration to node %q\n", name)
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Runs validation check and structural diff without writing to API")
	cmd.Flags().BoolVar(&yes, "yes", false, "Bypass confirmation prompts")
	cmd.Flags().StringVar(&targetNode, "node", "", "Override target to a specific node name")

	return cmd
}

// loadConfig reads the config file from disk.
func loadConfig(path string) (*config.Config, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("configuration file %q not found. Please verify the file path", path)
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg config.Config
	if err := json.Unmarshal(bytes, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse JSON from config file: %w", err)
	}

	return &cfg, nil
}

// saveConfig writes the config file to disk with pretty-print.
func saveConfig(path string, cfg *config.Config) error {
	bytes, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to format config: %w", err)
	}

	if err := os.WriteFile(path, bytes, 0644); err != nil {
		return fmt.Errorf("failed to write config to file: %w", err)
	}

	return nil
}

// resolveNodes discovers Technitium node configurations from environment variables and secrets.json.
func resolveNodes() (map[string]NodeConfig, error) {
	nodes := make(map[string]NodeConfig)

	// 1. Try loading from environment variables: TECHNITIUM_URL_<suffix> and TECHNITIUM_TOKEN_<suffix>
	// or TECHNITIUM_NODE_<suffix>_URL/TOKEN, plus fallback single node envs: TECHNITIUM_URL and TECHNITIUM_TOKEN.
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := parts[0]
		val := parts[1]

		if strings.HasPrefix(key, "TECHNITIUM_URL_") {
			suffix := strings.TrimPrefix(key, "TECHNITIUM_URL_")
			tokenKey := "TECHNITIUM_TOKEN_" + suffix
			if tokenVal := os.Getenv(tokenKey); tokenVal != "" {
				nodes["node-"+suffix] = NodeConfig{
					URL:   val,
					Token: tokenVal,
				}
			}
		} else if strings.HasPrefix(key, "TECHNITIUM_NODE_") && strings.HasSuffix(key, "_URL") {
			nodePart := strings.TrimPrefix(key, "TECHNITIUM_NODE_")
			nodePart = strings.TrimSuffix(nodePart, "_URL")
			tokenKey := fmt.Sprintf("TECHNITIUM_NODE_%s_TOKEN", nodePart)
			if tokenVal := os.Getenv(tokenKey); tokenVal != "" {
				name := strings.ToLower(nodePart)
				nodes[name] = NodeConfig{
					URL:   val,
					Token: tokenVal,
				}
			}
		}
	}

	// Fallback to single/default node env variable
	if defaultURL := os.Getenv("TECHNITIUM_URL"); defaultURL != "" {
		if defaultToken := os.Getenv("TECHNITIUM_TOKEN"); defaultToken != "" {
			nodes["default"] = NodeConfig{
				URL:   defaultURL,
				Token: defaultToken,
			}
		}
	}

	// 2. Read ~/.config/pab/secrets.json
	configDir, err := os.UserConfigDir()
	if err == nil {
		secretsPath := filepath.Join(configDir, "pab", "secrets.json")
		if _, statErr := os.Stat(secretsPath); statErr == nil {
			// Permission is verified in PersistentPreRunE, so we just read here
			fileBytes, readErr := os.ReadFile(secretsPath)
			if readErr == nil {
				var sc SecretsConfig
				if jsonErr := json.Unmarshal(fileBytes, &sc); jsonErr == nil {
					for name, node := range sc.Nodes {
						nodes[name] = node
					}
				}
			}
		}
	}

	return nodes, nil
}

// printStructuralDiff renders a simplified, human-readable structural diff.
func printStructuralDiff(w io.Writer, nodeName string, remote, local *config.Config) {
	fmt.Fprintf(w, "Structural configuration diff for node %q:\n", nodeName)

	if remote.EnableBlocking != local.EnableBlocking {
		fmt.Fprintf(w, "  ~ EnableBlocking: %t -> %t\n", remote.EnableBlocking, local.EnableBlocking)
	}
	if remote.BlockingAnswerTtl != local.BlockingAnswerTtl {
		fmt.Fprintf(w, "  ~ BlockingAnswerTtl: %d -> %d\n", remote.BlockingAnswerTtl, local.BlockingAnswerTtl)
	}

	// Diff Groups
	remoteGroups := make(map[string]config.Group)
	for _, g := range remote.Groups {
		remoteGroups[g.Name] = g
	}

	localGroups := make(map[string]config.Group)
	for _, g := range local.Groups {
		localGroups[g.Name] = g
	}

	// Sorted list of all group names
	allGroupNames := make(map[string]bool)
	for k := range remoteGroups {
		allGroupNames[k] = true
	}
	for k := range localGroups {
		allGroupNames[k] = true
	}

	var groupNames []string
	for k := range allGroupNames {
		groupNames = append(groupNames, k)
	}
	slices.Sort(groupNames)

	for _, gName := range groupNames {
		rg, inRemote := remoteGroups[gName]
		lg, inLocal := localGroups[gName]

		if !inRemote {
			fmt.Fprintf(w, "  + Group %q (Added)\n", gName)
		} else if !inLocal {
			fmt.Fprintf(w, "  - Group %q (Removed)\n", gName)
		} else {
			// Check if modified (simplified check)
			if rg.EnableBlocking != lg.EnableBlocking || !reflect.DeepEqual(rg.Blocked, lg.Blocked) || !reflect.DeepEqual(rg.Allowed, lg.Allowed) {
				fmt.Fprintf(w, "  ~ Group %q (Modified)\n", gName)
			}
		}
	}

	// Diff NetworkGroupMap
	allIPs := make(map[string]bool)
	for k := range remote.NetworkGroupMap {
		allIPs[k] = true
	}
	for k := range local.NetworkGroupMap {
		allIPs[k] = true
	}

	var ips []string
	for k := range allIPs {
		ips = append(ips, k)
	}
	slices.Sort(ips)

	hasMapDiff := false
	for _, ip := range ips {
		rv, inRemote := remote.NetworkGroupMap[ip]
		lv, inLocal := local.NetworkGroupMap[ip]

		if !inRemote {
			if !hasMapDiff {
				fmt.Fprintln(w, "  Client Mappings:")
				hasMapDiff = true
			}
			fmt.Fprintf(w, "    + %s -> %s\n", ip, lv)
		} else if !inLocal {
			if !hasMapDiff {
				fmt.Fprintln(w, "  Client Mappings:")
				hasMapDiff = true
			}
			fmt.Fprintf(w, "    - %s -> %s\n", ip, rv)
		} else if rv != lv {
			if !hasMapDiff {
				fmt.Fprintln(w, "  Client Mappings:")
				hasMapDiff = true
			}
			fmt.Fprintf(w, "    ~ %s: %s -> %s\n", ip, rv, lv)
		}
	}
	fmt.Fprintln(w)
}
