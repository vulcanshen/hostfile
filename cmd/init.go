package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/vulcanshen/hostfile/backup"
	"github.com/vulcanshen/hostfile/manager"
	"github.com/vulcanshen/hostfile/parser"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize hostfile management for the current hosts file",
	Long: `Backup the current hosts file as "origin", then reorganize all entries
into the managed block format and overwrite the hosts file.

If a managed block already exists, entries outside the block are merged in.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// read entire hosts file
		data, err := os.ReadFile(hostsFile)
		if err != nil {
			exitWithError(fmt.Errorf("cannot read hosts file: %w", err))
		}

		content := string(data)
		hasBlock := strings.Contains(content, parser.BlockStart)

		var block *parser.ManagedBlock

		if hasBlock {
			// merge outside entries into the managed block
			before, existingBlock, after, err := manager.ReadHostsFile(hostsFile)
			if err != nil {
				exitWithError(err)
			}

			outside := strings.TrimSpace(before + after)
			if outside == "" {
				exitWithError(fmt.Errorf("all entries are already in the managed block — nothing to do"))
			}

			outsideBlock := parser.ParseBlock(outside)
			for _, w := range outsideBlock.Warnings {
				fmt.Fprintln(os.Stderr, w)
			}

			// merge outside entries into existing block
			block = existingBlock
			for _, entry := range outsideBlock.Entries {
				merged := false
				if entry.DisableType == parser.DisableNone {
					for i, existing := range block.Entries {
						if existing.IP == entry.IP && existing.DisableType == parser.DisableNone {
							for _, d := range entry.Domains {
								if !manager.ContainsDomain(existing.Domains, d) {
									block.Entries[i].Domains = append(block.Entries[i].Domains, d)
								}
							}
							merged = true
							break
						}
					}
				}
				if !merged {
					block.Entries = append(block.Entries, entry)
				}
			}

			fmt.Printf("found %d entries outside the managed block, merging in\n", len(outsideBlock.Entries))
		} else {
			// fresh init — parse everything
			block = parser.ParseBlock(content)
			for _, w := range block.Warnings {
				fmt.Fprintln(os.Stderr, w)
			}
		}

		// determine save name
		saveName := "origin"
		exists, err := backup.Exists("origin")
		if err != nil {
			exitWithError(err)
		}
		if exists {
			saveName = fmt.Sprintf("origin-%s", time.Now().Format("20060102_150405"))
		}

		// confirm with user
		fmt.Println("This will overwrite your current hosts file with hostfile's managed format.")
		fmt.Printf("You can restore the original later with: hostfile load %s\n", saveName)
		if !confirm("Continue?") {
			fmt.Println("aborted")
			return
		}

		// backup original raw content
		if err := backup.CreateRaw(saveName, data); err != nil {
			exitWithError(fmt.Errorf("cannot create backup: %w", err))
		}
		path, err := backup.Path(saveName)
		if err != nil {
			exitWithError(err)
		}
		fmt.Printf("saved current hosts file to: %s\n", path)

		// write back with managed block (no before/after, everything is in the block)
		if err := writeBlock("", block, ""); err != nil {
			exitWithError(err)
		}
		fmt.Printf("initialized (%d entries)\n", len(block.Entries))
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
