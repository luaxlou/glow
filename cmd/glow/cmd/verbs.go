package cmd

import "github.com/spf13/cobra"

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Display one or many resources",
}

var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Show details of a specific resource",
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete resources",
}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy or update resources",
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a resource",
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop a resource",
}

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart a resource",
}

var logsCmd = &cobra.Command{
	Use:   "logs [name]",
	Short: "Print the logs for a resource",
	Args:  cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(describeCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(deployCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(restartCmd)
	rootCmd.AddCommand(logsCmd)
}