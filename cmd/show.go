package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vulcanshen/hostfile/backup"
	"github.com/vulcanshen/hostfile/manager"
	"github.com/vulcanshen/hostfile/parser"
)

var showJSON bool

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

	if showJSON {
		printJSON(entries)
		return
	}
	printEntries(entries)
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

	if showJSON {
		printJSON(restored.Entries)
		return
	}
	printEntries(restored.Entries)
}

func printJSON(entries []parser.HostEntry) {
	result := make(map[string][]string)
	for _, e := range entries {
		if e.DisableType != parser.DisableNone {
			continue
		}
		result[e.IP] = append(result[e.IP], e.Domains...)
	}
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		exitWithError(err)
	}
	fmt.Println(string(data))
}

func init() {
	showCmd.Flags().BoolVar(&showJSON, "json", false, "output active entries as JSON")
	rootCmd.AddCommand(showCmd)
}
