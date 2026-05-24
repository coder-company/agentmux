package commands

import (
	"fmt"

	"agentmux/internal/config"

	"github.com/spf13/cobra"
)

func initCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Create a default config file",
		Long:  `Create a default config file at ~/.config/agentmux/config.toml if one does not already exist.`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			path := config.DefaultPath()
			if err := config.InitDefault(path); err != nil {
				return fmt.Errorf("could not create config: %w", err)
			}
			fmt.Printf("Config ready: %s\n", path)
			return nil
		},
	}
}
