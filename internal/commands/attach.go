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
		Long:  `Attach to an existing tmux session by name. Replaces the current process.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := tmux.New()
			if !client.Available() {
				return errTmuxMissing
			}
			name := args[0]
			if err := client.AttachSession(name); err != nil {
				return fmt.Errorf("could not attach to session %q: %w", name, err)
			}
			return nil
		},
	}
}
