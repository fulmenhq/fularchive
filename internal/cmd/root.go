package cmd

import (
	"github.com/fulmenhq/refbolt/internal/config"
	"github.com/spf13/cobra"
)

var (
	verbose    bool
	configFlag string
)

var rootCmd = &cobra.Command{
	Use:   "refbolt",
	Short: "Archive web docs into clean, versioned Markdown trees",
	Long: `refbolt snapshots documentation sites (especially LLM APIs)
into date-versioned Markdown + JSON archives for offline use.`,
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// init and validate handle their own config loading.
		if cmd.Name() == "init" || cmd.Name() == "version" {
			return nil
		}

		strict := cmd.Name() == "validate"
		resolved := config.ResolveConfigPath(configFlag)

		return config.Load(config.LoadOptions{
			ConfigPath:  resolved,
			Strict:      strict,
			UseEmbedded: resolved == "",
		})
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().StringVar(&configFlag, "config", "", "Path to providers config file")
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}
