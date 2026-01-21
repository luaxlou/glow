package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	version   string
	commit    string
	buildDate string
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run: func(cmd *cobra.Command, args []string) {
		if version == "" {
			version = "dev"
		}
		if commit == "" {
			commit = "unknown"
		}
		if buildDate == "" {
			buildDate = "unknown"
		}
		fmt.Printf("glow-server version %s\n", version)
		fmt.Printf("commit: %s\n", commit)
		fmt.Printf("built at: %s\n", buildDate)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
