package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pab",
	Short: "Pharos Advanced Blocking CLI",
	Long:  `A command-line interface to manage, validate, and sync Technitium Advanced Blocking configurations.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to Pharos Advanced Blocking (pab)!")
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
