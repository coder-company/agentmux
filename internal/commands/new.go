package commands

import (
	"fmt"

	"agentmux/internal/adapters/tmux"

	"github.com/spf13/cobra"
)

func newCmd() *cobra.Command {
	var dir string

	cmd := &cobra.Command{
		Use:   "new <name>",
		Short: "Create a new tmux session",
		Long:  `Create a detached tmux session with the given name. Optionally set the working directory with --dir.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := tmux.New()
			if !client.Available() {
				return errTmuxMissing
			}
			name := args[0]
			if err := tmux.ValidateSessionName(name); err != nil {
				return err
			}
			if err := client.NewSession(name, dir); err != nil {
				return fmt.Errorf("could not create session %q: %w", name, err)
			}
			fmt.Printf("Created session: %s\n", name)
			return nil
		},
	}

	cmd.Flags().StringVarP(&dir, "dir", "d", "", "working directory for the new session")
	return cmd
}
