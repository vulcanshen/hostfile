package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vulcanshen/hostfile/manager"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all entries in the managed block",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		_, block, _, err := readBlock()
		if err != nil {
			exitWithError(err)
		}

		entries := manager.List(block)
		if len(entries) == 0 {
			fmt.Println("managed block is empty")
			return
		}

		for _, entry := range entries {
			printEntry(entry)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
