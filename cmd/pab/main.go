package main

import (
	"os"

	"github.com/iamrichardd/pharos-advanced-blocking/internal/commands"
)

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

func main() {
	cmd := commands.NewRootCmd(os.Stdin, os.Stdout, os.Stderr, Version, Commit, Date)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
