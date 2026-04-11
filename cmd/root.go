package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

var hostsFile string

func defaultHostsPath() string {
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

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&hostsFile, "hosts-file", defaultHostsPath(), "path to hosts file")
}
