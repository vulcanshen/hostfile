package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vulcanshen/hostfile/manager"
)

var disableCmd = &cobra.Command{
	Use:   "disable <ip|domain>",
	Short: "Disable an entry without removing it",
	Args:  cobra.ExactArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return completeActiveEntries()
	},
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]

		before, block, after, err := readBlock()
		if err != nil {
			exitWithError(err)
		}

		result, err := manager.Disable(block, target)
		if err != nil {
			exitWithError(err)
		}

		if result.LastDomain {
			msg := fmt.Sprintf("%q is the last active domain on %s. Disable the entire IP?", target, result.LastDomainIP)
			if confirm(msg) {
				manager.DisableIPConfirmed(block, result.LastDomainIP)
			} else {
				fmt.Println("aborted")
				return
			}
		}

		if err := writeBlock(before, block, after); err != nil {
			exitWithError(err)
		}
		fmt.Printf("disabled %s\n", target)
	},
}

func init() {
	rootCmd.AddCommand(disableCmd)
}
