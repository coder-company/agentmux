package commands

import (
	"agentmux/internal/app"

	"github.com/spf13/cobra"
)

// Root returns the root cobra command.
func Root() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agentmux",
		Short: "A modern terminal workspace control plane over tmux",
		RunE: func(cmd *cobra.Command, args []string) error {
			a, err := app.New()
			if err != nil {
				return err
			}
			defer a.Close()
			return a.RunTUI()
		},
		SilenceUsage: true,
	}

	cmd.AddCommand(listCmd())
	cmd.AddCommand(newCmd())
	cmd.AddCommand(attachCmd())
	cmd.AddCommand(killCmd())

	return cmd
}
