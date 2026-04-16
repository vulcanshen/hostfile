package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vulcanshen/hostfile/manager"
)

var enableCmd = &cobra.Command{
	Use:   "enable <ip|domain>",
	Short: "Enable a disabled entry",
	Args:  cobra.ExactArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return completeDisabledEntries()
	},
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]

		before, block, after, err := readBlock()
		if err != nil {
			exitWithError(err)
		}

		if err := manager.Enable(block, target); err != nil {
			exitWithError(err)
		}

		if err := writeBlock(before, block, after); err != nil {
			exitWithError(err)
		}
		fmt.Printf("enabled %s\n", target)
	},
}

func init() {
	rootCmd.AddCommand(enableCmd)
}
