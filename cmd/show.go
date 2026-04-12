package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vulcanshen/hostfile/backup"
	"github.com/vulcanshen/hostfile/manager"
)

var showCmd = &cobra.Command{
	Use:   "show [name]",
	Short: "Show entries in the managed block, or show a saved snapshot",
	Args:  cobra.MaximumNArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return completeSaveNames()
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			showCurrentBlock()
			return
		}
		showSavedSnapshot(args[0])
	},
}

func showCurrentBlock() {
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
}

func showSavedSnapshot(name string) {
	raw, err := backup.IsRaw(name)
	if err != nil {
		exitWithError(err)
	}

	if raw {
		data, err := backup.RestoreRaw(name)
		if err != nil {
			exitWithError(err)
		}
		fmt.Print(string(data))
		return
	}

	restored, err := backup.Restore(name)
	if err != nil {
		exitWithError(err)
	}

	if len(restored.Entries) == 0 {
		fmt.Printf("save %q is empty\n", name)
		return
	}

	for _, entry := range restored.Entries {
		printEntry(entry)
	}
}

func init() {
	rootCmd.AddCommand(showCmd)
}
