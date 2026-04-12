package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vulcanshen/hostfile/backup"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete a saved snapshot",
	Args:  cobra.ExactArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return completeSaveNames()
	},
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		if err := backup.Delete(name); err != nil {
			exitWithError(err)
		}
		fmt.Printf("deleted %q\n", name)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
