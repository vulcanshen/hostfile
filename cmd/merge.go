package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vulcanshen/hostfile/manager"
)

var mergeCmd = &cobra.Command{
	Use:   "merge <file>",
	Short: "Merge content from a file into the managed block",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]

		data, err := os.ReadFile(filePath)
		if err != nil {
			exitWithError(fmt.Errorf("cannot read file %s: %w", filePath, err))
		}

		before, block, after, err := readBlock()
		if err != nil {
			exitWithError(err)
		}

		manager.Merge(block, string(data))

		if err := writeBlock(before, block, after); err != nil {
			exitWithError(err)
		}
		fmt.Printf("merged %q (%d entries)\n", filePath, len(block.Entries))
	},
}

func init() {
	rootCmd.AddCommand(mergeCmd)
}
