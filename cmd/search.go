package cmd

import (
	"fmt"
	"os"
	"strings"

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

		printEntries(results)
	},
}

const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorGray   = "\033[90m"
	colorYellow = "\033[33m"
)

func isTTY() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

func printEntries(entries []parser.HostEntry) {
	if len(entries) == 0 {
		return
	}

	// calculate max IP width for alignment
	maxIPLen := 0
	for _, e := range entries {
		if len(e.IP) > maxIPLen {
			maxIPLen = len(e.IP)
		}
	}

	for _, e := range entries {
		printEntryAligned(e, maxIPLen)
	}
}

func printEntryAligned(entry parser.HostEntry, ipWidth int) {
	tty := isTTY()
	domains := strings.Join(entry.Domains, " ")

	switch entry.DisableType {
	case parser.DisableIP:
		ip := fmt.Sprintf("%-*s", ipWidth, entry.IP)
		if tty {
			fmt.Printf("%s%s  %s  [disabled]%s\n", colorGray, ip, domains, colorReset)
		} else {
			fmt.Printf("%s  %s  [disabled]\n", ip, domains)
		}
	case parser.DisableDomain:
		ip := fmt.Sprintf("%-*s", ipWidth, entry.IP)
		domain := ""
		if len(entry.Domains) > 0 {
			domain = entry.Domains[0]
		}
		if tty {
			fmt.Printf("%s%s  %s  [disabled]%s\n", colorGray, ip, domain, colorReset)
		} else {
			fmt.Printf("%s  %s  [disabled]\n", ip, domain)
		}
	default:
		ip := fmt.Sprintf("%-*s", ipWidth, entry.IP)
		if tty {
			fmt.Printf("%s  %s%s%s\n", ip, colorGreen, domains, colorReset)
		} else {
			fmt.Printf("%s  %s\n", ip, domains)
		}
	}
}

func printEntry(entry parser.HostEntry) {
	printEntryAligned(entry, len(entry.IP))
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
