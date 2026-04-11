package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vulcanshen/hostfile/manager"
)

var removeCmd = &cobra.Command{
	Use:   "remove <ip|domain>",
	Short: "Remove an IP or domain from the managed block",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]

		before, block, after, err := readBlock()
		if err != nil {
			exitWithError(err)
		}

		result := manager.Remove(block, target)

		if result.LastDomain {
			msg := fmt.Sprintf("%q is the last domain on %s. Remove the entire line?", target, result.LastDomainIP)
			if confirm(msg) {
				manager.RemoveConfirmed(block, result.LastDomainIP)
			} else {
				fmt.Println("aborted")
				return
			}
		} else if !result.Removed {
			fmt.Printf("%q not found in managed block\n", target)
			return
		}

		if err := writeBlock(before, block, after); err != nil {
			exitWithError(err)
		}
		fmt.Printf("removed %s\n", target)
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
