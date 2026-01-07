package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication and connection settings",
}

var authViewCmd = &cobra.Command{
	Use:   "view",
	Short: "View current connection settings",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := loadConfig()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}
		
		var currentCtx *Context
		for _, ctx := range cfg.Contexts {
			if ctx.Name == cfg.CurrentContext {
				currentCtx = &ctx
				break
			}
		}
		
		if currentCtx == nil {
			fmt.Println("No active context found.")
			return
		}

		keyMasked := "********"
		if len(currentCtx.APIKey) > 4 {
			keyMasked = currentCtx.APIKey[:4] + "****"
		}
		fmt.Printf("Context:    %s\n", currentCtx.Name)
		fmt.Printf("Server URL: %s\n", currentCtx.ServerURL)
		fmt.Printf("API Key:    %s\n", keyMasked)
	},
}

var authResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset ALL connection settings (clears all contexts)",
	Run: func(cmd *cobra.Command, args []string) {
		if err := os.Remove(cfgFile); err != nil && !os.IsNotExist(err) {
			fmt.Printf("Error removing config: %v\n", err)
			return
		}
		fmt.Println("Configuration reset.")
		ensureConfig() // Trigger bootstrap
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(authViewCmd)
	authCmd.AddCommand(authResetCmd)
}