package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vulcanshen/hostfile/backup"
	"github.com/vulcanshen/hostfile/manager"
)

var backupCmd = &cobra.Command{
	Use:   "backup <name>",
	Short: "Backup the managed block",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		_, block, _, err := readBlock()
		if err != nil {
			exitWithError(err)
		}

		if err := backup.Create(name, block); err != nil {
			exitWithError(err)
		}
		fmt.Printf("backup %q created (%d entries)\n", name, len(block.Entries))
	},
}

var backupListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all backups",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		backups, err := backup.List()
		if err != nil {
			exitWithError(err)
		}

		if len(backups) == 0 {
			fmt.Println("no backups found")
			return
		}

		for _, b := range backups {
			fmt.Printf("  %s  (%s)\n", b.Name, b.ModTime.Format("2006-01-02 15:04:05"))
		}
	},
}

var backupRestoreCmd = &cobra.Command{
	Use:   "restore <name>",
	Short: "Restore the managed block from a backup",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		restored, err := backup.Restore(name)
		if err != nil {
			exitWithError(err)
		}

		before, _, after, err := readBlock()
		if err != nil {
			exitWithError(err)
		}

		if err := manager.WriteHostsFile(hostsFile, before, restored, after); err != nil {
			exitWithError(err)
		}
		fmt.Printf("restored backup %q (%d entries)\n", name, len(restored.Entries))
	},
}

var backupDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete a backup",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		if err := backup.Delete(name); err != nil {
			exitWithError(err)
		}
		fmt.Printf("backup %q deleted\n", name)
	},
}

func init() {
	backupCmd.AddCommand(backupListCmd)
	backupCmd.AddCommand(backupRestoreCmd)
	backupCmd.AddCommand(backupDeleteCmd)
	rootCmd.AddCommand(backupCmd)
}
