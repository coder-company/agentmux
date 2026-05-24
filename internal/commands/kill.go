package commands

import (
	"fmt"

	"agentmux/internal/adapters/tmux"

	"github.com/spf13/cobra"
)

func killCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "kill <name>",
		Short: "Kill a tmux session",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := tmux.New()
			name := args[0]
			if err := client.KillSession(name); err != nil {
				return fmt.Errorf("failed to kill session %q: %w", name, err)
			}
			fmt.Printf("Killed session: %s\n", name)
			return nil
		},
	}
}
