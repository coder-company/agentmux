package main

import (
	"fmt"
	"os"

	"agentmux/internal/commands"
)

func main() {
	if err := commands.Root().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
