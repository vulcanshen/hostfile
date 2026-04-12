package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vulcanshen/hostfile/backup"
	"github.com/vulcanshen/hostfile/parser"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize hostfile management for the current hosts file",
	Long: `Backup the current hosts file as "origin", then reorganize all entries
into the managed block format and overwrite the hosts file.

This is a one-time setup command to start managing your hosts file with hostfile.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// read entire hosts file
		data, err := os.ReadFile(hostsFile)
		if err != nil {
			exitWithError(fmt.Errorf("cannot read hosts file: %w", err))
		}

		content := string(data)

		// check if already managed
		if strings.Contains(content, parser.BlockStart) {
			exitWithError(fmt.Errorf("hosts file already contains a managed block — no need to init"))
		}

		// parse all lines into entries
		block := parser.ParseBlock(content)
		for _, w := range block.Warnings {
			fmt.Fprintln(os.Stderr, w)
		}

		// confirm with user
		fmt.Println("This will overwrite your current hosts file with hostfile's managed format.")
		fmt.Println("You can restore the original later with: hostfile load origin")
		if !confirm("Continue?") {
			fmt.Println("aborted")
			return
		}

		// backup original raw content as "origin"
		if err := backup.CreateRaw("origin", data); err != nil {
			exitWithError(fmt.Errorf("cannot create backup: %w", err))
		}
		fmt.Println("backed up current hosts file as \"origin\"")

		// write back with managed block (no before/after, everything is in the block)
		if err := writeBlock("", block, ""); err != nil {
			exitWithError(err)
		}
		fmt.Printf("initialized hostfile management (%d entries)\n", len(block.Entries))
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
