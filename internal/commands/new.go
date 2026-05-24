package commands

import (
	"fmt"

	"agentmux/internal/adapters/tmux"

	"github.com/spf13/cobra"
)

func newCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "new <name>",
		Short: "Create a new tmux session",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := tmux.New()
			name := args[0]
			if err := client.NewSession(name, ""); err != nil {
				return fmt.Errorf("failed to create session %q: %w", name, err)
			}
			fmt.Printf("Created session: %s\n", name)
			return nil
		},
	}
}
