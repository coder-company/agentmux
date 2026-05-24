package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	"agentmux/internal/adapters/tmux"

	"github.com/spf13/cobra"
)

func listCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List tmux sessions",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			client := tmux.New()
			if !client.Available() {
				return errTmuxMissing
			}
			sessions, err := client.ListSessions()
			if err != nil {
				return fmt.Errorf("could not list sessions: %w", err)
			}
			if len(sessions) == 0 {
				fmt.Println("No active tmux sessions.")
				return nil
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintf(w, "NAME\tWINDOWS\tDIRECTORY\tSTATUS\n")
			for _, s := range sessions {
				status := ""
				if s.Attached {
					status = "attached"
				}
				fmt.Fprintf(w, "%s\t%d\t%s\t%s\n", s.Name, s.Windows, s.Directory, status)
			}
			return w.Flush()
		},
	}
}
