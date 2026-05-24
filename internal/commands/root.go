package commands

import (
	"agentmux/internal/app"

	"github.com/spf13/cobra"
)

// Version is set at build time via ldflags.
var Version = "dev"

// Root returns the root cobra command.
func Root() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agentmux",
		Short: "Terminal workspace control plane over tmux",
		Long: `agentmux is a keyboard-first TUI that wraps tmux.

Browse sessions, preview pane output, launch workspaces, and manage
your terminal workflow — all without leaving the terminal.

Run without arguments to launch the interactive TUI.
Run with a subcommand for quick CLI access.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			a, err := app.New()
			if err != nil {
				return err
			}
			defer a.Close()
			return a.RunTUI()
		},
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       Version,
	}

	cmd.AddCommand(listCmd())
	cmd.AddCommand(newCmd())
	cmd.AddCommand(attachCmd())
	cmd.AddCommand(killCmd())
	cmd.AddCommand(initCmd())

	return cmd
}
