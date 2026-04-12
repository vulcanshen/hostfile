package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vulcanshen/hostfile/backup"
	"github.com/vulcanshen/hostfile/privilege"
)

var loadCmd = &cobra.Command{
	Use:   "load <name>",
	Short: "Load a saved snapshot into the managed block",
	Args:  cobra.ExactArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return completeSaveNames()
	},
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		// check if this is a raw backup (e.g. from init)
		raw, err := backup.IsRaw(name)
		if err != nil {
			exitWithError(err)
		}

		if raw {
			data, err := backup.RestoreRaw(name)
			if err != nil {
				exitWithError(err)
			}
			if err := privilege.WriteFilePrivileged(hostsFile, data); err != nil {
				exitWithError(err)
			}
			fmt.Printf("loaded %q (raw)\n", name)
			return
		}

		restored, err := backup.Restore(name)
		if err != nil {
			exitWithError(err)
		}

		before, _, after, err := readBlock()
		if err != nil {
			exitWithError(err)
		}

		if err := writeBlock(before, restored, after); err != nil {
			exitWithError(err)
		}
		fmt.Printf("loaded %q (%d entries)\n", name, len(restored.Entries))
	},
}

func init() {
	rootCmd.AddCommand(loadCmd)
}
