package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vulcanshen/hostfile/manager"
	"github.com/vulcanshen/hostfile/parser"
)

var searchCmd = &cobra.Command{
	Use:   "search <ip|domain>",
	Short: "Search the managed block for an IP or domain",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := args[0]

		_, block, _, err := readBlock()
		if err != nil {
			exitWithError(err)
		}

		results := manager.Search(block, query)
		if len(results) == 0 {
			fmt.Printf("no entries found for %q\n", query)
			return
		}

		for _, entry := range results {
			printEntry(entry)
		}
	},
}

func printEntry(entry parser.HostEntry) {
	line := parser.FormatEntry(&entry)
	status := ""
	switch entry.DisableType {
	case parser.DisableIP:
		status = " (disabled-ip)"
	case parser.DisableDomain:
		status = " (disabled-domain)"
	}
	fmt.Printf("%s%s\n", line, status)
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
