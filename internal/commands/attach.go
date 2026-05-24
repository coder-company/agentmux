package commands

import (
	"fmt"
	"os"
	"os/exec"

	"agentmux/internal/adapters/tmux"
	"agentmux/internal/app"

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
			bin, err := exec.LookPath("tmux")
			if err != nil {
				return fmt.Errorf("could not find tmux: %w", err)
			}
			return app.Execvp(bin, []string{"tmux", "attach-session", "-t", name}, os.Environ())
		},
	}
}
