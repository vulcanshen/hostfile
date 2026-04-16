package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vulcanshen/hostfile/manager"
	"github.com/vulcanshen/hostfile/parser"
)

var searchAll bool

var searchCmd = &cobra.Command{
	Use:           "search <ip|domain>",
	Short:         "Search the managed block for an IP or domain",
	SilenceUsage:  true,
	SilenceErrors: true,
	Args: cobra.ExactArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return completeAllEntries()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]

		before, block, after, err := readBlock()
		if err != nil {
			exitWithError(err)
		}

		results := manager.Search(block, query)

		if searchAll {
			outside := parseOutsideEntries(before + "\n" + after)
			outsideBlock := &parser.ManagedBlock{Entries: outside}
			outsideResults := manager.Search(outsideBlock, query)
			results = append(outsideResults, results...)
		}

		if len(results) == 0 {
			return fmt.Errorf("no entries found for %q", query)
		}

		printEntriesHighlight(results, query)
		return nil
	},
}

const (
	colorReset   = "\033[0m"
	colorGreen   = "\033[32m"
	colorGray    = "\033[90m"
	colorYellow  = "\033[33m"
	colorBold    = "\033[1m"
	colorCyan    = "\033[36m"
	colorReverse = "\033[7m"
)

func isTTY() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

func highlightSubstring(s, query string) string {
	lower := strings.ToLower(s)
	lowerQuery := strings.ToLower(query)
	idx := strings.Index(lower, lowerQuery)
	if idx < 0 {
		return s
	}
	return s[:idx] + colorYellow + s[idx:idx+len(query)] + colorReset + s[idx+len(query):]
}

func highlightSubstringReverse(s, query string) string {
	lower := strings.ToLower(s)
	lowerQuery := strings.ToLower(query)
	idx := strings.Index(lower, lowerQuery)
	if idx < 0 {
		return s
	}
	return s[:idx] + colorReverse + s[idx:idx+len(query)] + colorReset + colorGray + s[idx+len(query):]
}

func maxIPWidth(entries []parser.HostEntry) int {
	w := len("IP")
	for _, e := range entries {
		if len(e.IP) > w {
			w = len(e.IP)
		}
	}
	return w
}

func printHeader(ipWidth int, tty bool) {
	ip := fmt.Sprintf("%-*s", ipWidth, "IP")
	if tty {
		fmt.Printf("%s%s%s  %s%s%s\n", colorBold+colorCyan, ip, colorReset, colorBold+colorCyan, "DOMAIN", colorReset)
		ipUnderline := fmt.Sprintf("%-*s", ipWidth, strings.Repeat("─", len("IP")))
		fmt.Printf("%s%s%s  %s%s%s\n", colorGray, ipUnderline, colorReset, colorGray, strings.Repeat("─", len("DOMAIN")), colorReset)
	} else {
		fmt.Printf("%s  %s\n", ip, "DOMAIN")
	}
}

func printEntries(entries []parser.HostEntry) {
	if len(entries) == 0 {
		return
	}
	tty := isTTY()
	w := maxIPWidth(entries)
	printHeader(w, tty)
	for _, e := range entries {
		printEntryRows(e, w, "", tty)
	}
}

func printEntriesHighlight(entries []parser.HostEntry, query string) {
	if len(entries) == 0 {
		return
	}
	tty := isTTY()
	w := maxIPWidth(entries)
	printHeader(w, tty)
	for _, e := range entries {
		printEntryRows(e, w, query, tty)
	}
}

func printEntryRows(entry parser.HostEntry, ipWidth int, query string, tty bool) {
	blank := fmt.Sprintf("%-*s", ipWidth, "")

	switch entry.DisableType {
	case parser.DisableIP:
		for i, d := range entry.Domains {
			ip := blank
			if i == 0 {
				ip = fmt.Sprintf("%-*s", ipWidth, entry.IP)
			}
			if tty {
				if query != "" {
					ip = highlightSubstringReverse(ip, query)
					d = highlightSubstringReverse(d, query)
				}
				fmt.Printf("%s%s  %s  [disabled]%s\n", colorGray, ip, d, colorReset)
			} else {
				fmt.Printf("%s  %s  [disabled]\n", ip, d)
			}
		}
	case parser.DisableDomain:
		ip := fmt.Sprintf("%-*s", ipWidth, entry.IP)
		d := ""
		if len(entry.Domains) > 0 {
			d = entry.Domains[0]
		}
		if tty {
			if query != "" {
				ip = highlightSubstringReverse(ip, query)
				d = highlightSubstringReverse(d, query)
			}
			fmt.Printf("%s%s  %s  [disabled]%s\n", colorGray, ip, d, colorReset)
		} else {
			fmt.Printf("%s  %s  [disabled]\n", ip, d)
		}
	default:
		for i, d := range entry.Domains {
			ip := blank
			if i == 0 {
				ip = fmt.Sprintf("%-*s", ipWidth, entry.IP)
			}
			if tty {
				if query != "" {
					ip = highlightSubstring(ip, query)
					d = highlightSubstring(d, query)
					fmt.Printf("%s  %s\n", ip, d)
				} else {
					fmt.Printf("%s  %s%s%s\n", ip, colorGreen, d, colorReset)
				}
			} else {
				fmt.Printf("%s  %s\n", ip, d)
			}
		}
	}
}

func init() {
	searchCmd.Flags().BoolVar(&searchAll, "all", false, "include entries outside the managed block")
	rootCmd.AddCommand(searchCmd)
}
