package commands

import (
	"fmt"

	"agentmux/internal/adapters/tmux"

	"github.com/spf13/cobra"
)

func attachCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "attach <name>",
		Short: "Attach to a tmux session",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := tmux.New()
			name := args[0]
			if err := client.AttachSession(name); err != nil {
				return fmt.Errorf("failed to attach to %q: %w", name, err)
			}
			return nil
		},
	}
}
