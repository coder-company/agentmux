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
		Long:  `Destroy a tmux session by name. This terminates all processes in the session.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := tmux.New()
			if !client.Available() {
				return errTmuxMissing
			}
			name := args[0]
			if err := client.KillSession(name); err != nil {
				return fmt.Errorf("could not kill session %q: %w", name, err)
			}
			fmt.Printf("Killed session: %s\n", name)
			return nil
		},
	}
}
