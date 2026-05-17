package config

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"configurator/pkg/config"
)

func NewRootCmd() *cobra.Command {
	var cfgFile string

	rootCmd := &cobra.Command{
		Use:   "configurator",
		Short: "Manage and apply configurations from a directory",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Read(cfgFile)
			if err != nil {
				return err
			}

			cmd.SetContext(config.NewContext(cmd.Context(), cfg))

			return nil
		},
	}

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "configurator.yaml", "path to configurator YAML config")

	rootCmd.AddCommand(newListCmd())
	rootCmd.AddCommand(newApplyCmd())

	return rootCmd
}

func Execute(ctx context.Context) {
	if err := NewRootCmd().ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
