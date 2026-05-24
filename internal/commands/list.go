package commands

import (
	"fmt"

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
			sessions, err := client.ListSessions()
			if err != nil {
				return err
			}
			if len(sessions) == 0 {
				fmt.Println("No tmux sessions running.")
				return nil
			}
			for _, s := range sessions {
				attached := ""
				if s.Attached {
					attached = " (attached)"
				}
				fmt.Printf("%-20s %d windows  %s%s\n", s.Name, s.Windows, s.Directory, attached)
			}
			return nil
		},
	}
}
