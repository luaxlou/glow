package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a resource from the system",
}

var removeYes bool

var removeMysqlCmd = &cobra.Command{
	Use:   "mysql",
	Short: "Remove MySQL integration",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRemoveResource(os.Stdin, os.Stdout, "mysql", removeYes)
	},
}

var removeRedisCmd = &cobra.Command{
	Use:   "redis",
	Short: "Remove Redis integration",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRemoveResource(os.Stdin, os.Stdout, "redis", removeYes)
	},
}

var removeNginxCmd = &cobra.Command{
	Use:   "nginx",
	Short: "Remove Nginx integration",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRemoveResource(os.Stdin, os.Stdout, "nginx", removeYes)
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)

	removeCmd.PersistentFlags().BoolVarP(&removeYes, "yes", "y", false, "Skip confirmation")
	removeCmd.AddCommand(removeMysqlCmd)
	removeCmd.AddCommand(removeRedisCmd)
	removeCmd.AddCommand(removeNginxCmd)
}

func runRemoveResource(in io.Reader, out io.Writer, resource string, yes bool) error {
	var keys []string
	switch resource {
	case "mysql":
		keys = []string{"mysql_info", "mysql_users"}
	case "redis":
		keys = []string{"redis_info", "redis_users"}
	case "nginx":
		keys = []string{"nginx_info"}
	default:
		return fmt.Errorf("unsupported resource: %s", resource)
	}

	existed, err := anySystemConfigExists(keys)
	if err != nil {
		return err
	}
	if !existed {
		fmt.Fprintf(out, "%s: [NOT CONFIGURED]\n", strings.ToUpper(resource))
		return nil
	}

	if !yes {
		ok, err := confirm(in, out, fmt.Sprintf("Remove %s integration and related cached credentials", strings.ToUpper(resource)))
		if err != nil {
			return err
		}
		if !ok {
			fmt.Fprintln(out, "aborted")
			return nil
		}
	}

	for _, k := range keys {
		if err := deleteSystemConfig(k); err != nil {
			return err
		}
	}

	fmt.Fprintf(out, "%s removed\n", strings.ToUpper(resource))
	return nil
}

func anySystemConfigExists(keys []string) (bool, error) {
	for _, k := range keys {
		val, err := getSystemConfig(k)
		if err != nil {
			return false, err
		}
		if strings.TrimSpace(val) != "" {
			return true, nil
		}
	}
	return false, nil
}

func confirm(in io.Reader, out io.Writer, prompt string) (bool, error) {
	fmt.Fprintf(out, "%s? [y/N]: ", prompt)
	reader := bufio.NewReader(in)
	line, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return false, err
	}
	answer := strings.TrimSpace(strings.ToLower(line))
	return answer == "y" || answer == "yes", nil
}
