package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vulcanshen/hostfile/manager"
)

var mergeDryRun bool

var mergeCmd = &cobra.Command{
	Use:   "merge <file | ->",
	Short: "Merge content from a file or stdin into the managed block",
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

		manager.Merge(block, content)

		if mergeDryRun {
			fmt.Printf("(dry-run) would result in %d entries:\n", len(block.Entries))
			printEntries(block.Entries)
			return
		}

		if err := writeBlock(before, block, after); err != nil {
			exitWithError(err)
		}
		fmt.Printf("merged (%d entries)\n", len(block.Entries))
	},
}

func init() {
	mergeCmd.Flags().BoolVar(&mergeDryRun, "dry-run", false, "preview changes without writing")
	rootCmd.AddCommand(mergeCmd)
}
