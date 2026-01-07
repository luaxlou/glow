package cmd

import (
	"os"
	"path/filepath"

	"github.com/luaxlou/glow/starter/glowsqlite"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "glow-server",
	Short: "Glow Server - Application Lifecycle Management",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		home, err := os.UserHomeDir()
		if err != nil {
			// Fallback if home dir cannot be found
			home = "."
		}
		dbPath := filepath.Join(home, ".glow-server", "glow.db")
		glowsqlite.Init(dbPath)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
