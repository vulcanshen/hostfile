package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vulcanshen/hostfile/manager"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clear all entries from the managed block",
	Args:              cobra.NoArgs,
	ValidArgsFunction: cobra.NoFileCompletions,
	Run: func(cmd *cobra.Command, args []string) {
		before, block, after, err := readBlock()
		if err != nil {
			exitWithError(err)
		}

		if len(block.Entries) == 0 {
			fmt.Println("managed block is already empty")
			return
		}

		if !confirm("Remove all entries from the managed block?") {
			fmt.Println("aborted")
			return
		}

		manager.Clean(block)

		if err := writeBlock(before, block, after); err != nil {
			exitWithError(err)
		}
		fmt.Println("cleared managed block")
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}
