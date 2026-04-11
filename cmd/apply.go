package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vulcanshen/hostfile/manager"
)

var applyCmd = &cobra.Command{
	Use:   "apply <file>",
	Short: "Replace the managed block with content from a file",
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

		manager.Apply(block, string(data))

		if err := writeBlock(before, block, after); err != nil {
			exitWithError(err)
		}
		fmt.Printf("applied %s to managed block (%d entries)\n", filePath, len(block.Entries))
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
}
