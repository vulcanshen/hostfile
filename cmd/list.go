package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vulcanshen/hostfile/backup"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all saved snapshots",
	Args:              cobra.NoArgs,
	ValidArgsFunction: cobra.NoFileCompletions,
	Run: func(cmd *cobra.Command, args []string) {
		backups, err := backup.List()
		if err != nil {
			exitWithError(err)
		}

		if len(backups) == 0 {
			fmt.Println("no saves found")
			return
		}

		for _, b := range backups {
			fmt.Printf("  %s  (%s)\n", b.Name, b.ModTime.Format("2006-01-02 15:04:05"))
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
