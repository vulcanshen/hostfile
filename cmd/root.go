package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

var hostsFile string

func defaultHostsPath() string {
	if v := os.Getenv("HOSTFILE__HOSTS_FILE"); v != "" {
		return v
	}
	if runtime.GOOS == "windows" {
		return `C:\Windows\System32\drivers\etc\hosts`
	}
	return "/etc/hosts"
}

var rootCmd = &cobra.Command{
	Use:   "hostfile",
	Short: "Cross-platform hosts file manager",
	Long:  "A CLI tool for managing your hosts file with ease.",
}

// RootCmd returns the root command for documentation generation.
func RootCmd() *cobra.Command {
	return rootCmd
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&hostsFile, "hosts-file", defaultHostsPath(), "path to hosts file")
}
