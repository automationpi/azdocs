package main

import (
	"fmt"
	"os"

	"github.com/automationpi/azdocs/cmd/azdoc/commands"
)

var (
	// Set via ldflags during build
	Version   = "dev"
	GitCommit = "unknown"
	BuildDate = "unknown"
)

func main() {
	commands.SetVersion(Version, GitCommit, BuildDate)

	if err := commands.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
