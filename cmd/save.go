package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vulcanshen/hostfile/backup"
)

var saveCmd = &cobra.Command{
	Use:   "save <name>",
	Short: "Save the managed block as a snapshot",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: cobra.NoFileCompletions,
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		_, block, _, err := readBlock()
		if err != nil {
			exitWithError(err)
		}

		if err := backup.Create(name, block); err != nil {
			exitWithError(err)
		}
		path, err := backup.Path(name)
		if err != nil {
			exitWithError(err)
		}
		fmt.Printf("saved %q (%d entries) to: %s\n", name, len(block.Entries), path)
	},
}

func completeSaveNames() ([]string, cobra.ShellCompDirective) {
	backups, err := backup.List()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	var names []string
	for _, b := range backups {
		names = append(names, b.Name)
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}

func init() {
	rootCmd.AddCommand(saveCmd)
}
