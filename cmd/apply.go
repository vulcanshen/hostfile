package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vulcanshen/hostfile/manager"
)

var applyCmd = &cobra.Command{
	Use:   "apply <file | ->",
	Short: "Replace the managed block with content from a file or stdin",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]

		data, err := readInput(filePath)
		if err != nil {
			exitWithError(fmt.Errorf("cannot read input: %w", err))
		}

		content, err := parseHostsContent(data)
		if err != nil {
			exitWithError(err)
		}

		before, block, after, err := readBlock()
		if err != nil {
			exitWithError(err)
		}

		manager.Apply(block, content)

		if err := writeBlock(before, block, after); err != nil {
			exitWithError(err)
		}
		fmt.Printf("applied (%d entries)\n", len(block.Entries))
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
}
