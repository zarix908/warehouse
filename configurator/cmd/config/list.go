package config

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"configurator/pkg/config"
)

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List available configurations",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.FromContext(cmd.Context())

			entries, err := os.ReadDir(cfg.ConfigsDir)
			if err != nil {
				return fmt.Errorf("read configs dir: %w", err)
			}

			for _, e := range entries {
				fmt.Printf("%s\n", e.Name())
			}

			return nil
		},
	}
}
