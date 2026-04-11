package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vulcanshen/hostfile/backup"
	"github.com/vulcanshen/hostfile/manager"
	"github.com/vulcanshen/hostfile/privilege"
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
	Use:               "restore <name>",
	Short:             "Restore the managed block from a backup",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeBackupNames,
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
			fmt.Printf("restored raw backup %q\n", name)
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

		if err := manager.WriteHostsFile(hostsFile, before, restored, after); err != nil {
			exitWithError(err)
		}
		fmt.Printf("restored backup %q (%d entries)\n", name, len(restored.Entries))
	},
}

var backupDeleteCmd = &cobra.Command{
	Use:               "delete <name>",
	Short:             "Delete a backup",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeBackupNames,
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		if err := backup.Delete(name); err != nil {
			exitWithError(err)
		}
		fmt.Printf("backup %q deleted\n", name)
	},
}

func completeBackupNames(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
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
	backupCmd.AddCommand(backupListCmd)
	backupCmd.AddCommand(backupRestoreCmd)
	backupCmd.AddCommand(backupDeleteCmd)
	rootCmd.AddCommand(backupCmd)
}
