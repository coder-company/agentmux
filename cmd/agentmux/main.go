package main

import (
	"fmt"
	"os"

	"agentmux/internal/commands"
)

func main() {
	cmd := commands.Root()
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "agentmux: %v\n", err)
		os.Exit(1)
	}
}
