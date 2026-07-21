package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// CommandMeta matches the structure expected by the host's plugin manager.
type CommandMeta struct {
	Use         string `json:"use"`
	Short       string `json:"short"`
	Long        string `json:"long"`
	Example     string `json:"example"`
}

// PluginMeta matches the structure expected by the host's plugin manager.
type PluginMeta struct {
	Name        string        `json:"name"`
	Version     string        `json:"version"`
	Description string        `json:"description"`
	Commands    []CommandMeta `json:"commands"`
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "info" {
		// Output the plugin metadata in JSON format
		meta := PluginMeta{
			Name:        "Sample Plugin",
			Version:     "0.1.0",
			Description: "A sample plugin for pab that demonstrates extensibility",
			Commands: []CommandMeta{
				{
					Use:   "sample",
					Short: "Run the sample plugin command",
					Long:  "Run the sample plugin command which prints a greeting.",
				},
			},
		}
		
		encoder := json.NewEncoder(os.Stdout)
		if err := encoder.Encode(meta); err != nil {
			fmt.Fprintf(os.Stderr, "failed to encode metadata: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// For actual execution of commands:
	var rootCmd = &cobra.Command{
		Use: "pab-plugin-sample",
	}

	var sampleCmd = &cobra.Command{
		Use:   "sample",
		Short: "Run the sample plugin command",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hello from the sample plugin!")
			if len(args) > 0 {
				fmt.Printf("Received arguments: %v\n", args)
			}
		},
	}

	rootCmd.AddCommand(sampleCmd)
	
	// Execute without "info" logic interference
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
