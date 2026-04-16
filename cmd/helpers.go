package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vulcanshen/hostfile/manager"
	"github.com/vulcanshen/hostfile/parser"
	"github.com/vulcanshen/hostfile/privilege"
)

// readBlock reads the hosts file and returns the parts.
func readBlock() (string, *parser.ManagedBlock, string, error) {
	return manager.ReadHostsFile(hostsFile)
}

// writeBlock writes the hosts file back, using privilege escalation if needed.
func writeBlock(before string, block *parser.ManagedBlock, after string) error {
	// ensure before ends with newline so block marker starts on its own line
	if before != "" && !strings.HasSuffix(before, "\n") {
		before += "\n"
	}

	formatted := parser.FormatBlock(block)
	var content string
	if before == "" && after == "" {
		content = formatted + "\n"
	} else if before == "" {
		content = formatted + "\n" + after
	} else {
		content = before + formatted + "\n" + after
	}
	return privilege.WriteFilePrivileged(hostsFile, []byte(content))
}

// readInput reads content from a file path or stdin (when path is "-").
func readInput(path string) ([]byte, error) {
	if path == "-" {
		return io.ReadAll(os.Stdin)
	}
	return os.ReadFile(path)
}

// parseHostsContent detects JSON or plain hosts format, validates, and returns hosts-format string.
func parseHostsContent(data []byte) (string, error) {
	trimmed := strings.TrimSpace(string(data))
	if len(trimmed) == 0 {
		return "", fmt.Errorf("input is empty")
	}

	// JSON format
	if trimmed[0] == '{' {
		var entries map[string][]string
		if err := json.Unmarshal(data, &entries); err != nil {
			return "", fmt.Errorf("invalid JSON: %w", err)
		}
		if len(entries) == 0 {
			return "", fmt.Errorf("JSON contains no entries")
		}
		for ip, domains := range entries {
			if !parser.ValidIP(ip) {
				return "", fmt.Errorf("invalid IP in JSON: %s", ip)
			}
			if len(domains) == 0 {
				return "", fmt.Errorf("IP %s has no domains", ip)
			}
		}
		var lines []string
		for ip, domains := range entries {
			lines = append(lines, ip+"  "+strings.Join(domains, " "))
		}
		return strings.Join(lines, "\n"), nil
	}

	// plain hosts format — validate that at least one line is parseable
	valid := 0
	invalid := 0
	lines := strings.Split(trimmed, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		entry, err := parser.ParseLine(line)
		if err != nil {
			invalid++
			continue
		}
		if entry != nil {
			valid++
		}
	}
	if valid == 0 {
		if invalid > 0 {
			return "", fmt.Errorf("no valid entries found (%d invalid lines)", invalid)
		}
		return "", fmt.Errorf("no entries found in input")
	}

	return trimmed, nil
}

// exitWithError prints an error and exits.
func exitWithError(err error) {
	fmt.Fprintln(os.Stderr, "error:", err)
	os.Exit(1)
}

// parseOutsideEntries parses hosts entries from text outside the managed block.
// Skips comments and invalid lines silently.
func parseOutsideEntries(text string) []parser.HostEntry {
	var entries []parser.HostEntry
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		entry, err := parser.ParseLine(trimmed)
		if err != nil || entry == nil {
			continue
		}
		entries = append(entries, *entry)
	}
	return entries
}

// completeAllEntries returns all IPs and domains for shell completion.
func completeAllEntries() ([]string, cobra.ShellCompDirective) {
	_, block, _, err := readBlock()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	seen := make(map[string]bool)
	var results []string
	for _, entry := range block.Entries {
		if !seen[entry.IP] {
			seen[entry.IP] = true
			results = append(results, entry.IP)
		}
		for _, d := range entry.Domains {
			if !seen[d] {
				seen[d] = true
				results = append(results, d)
			}
		}
	}
	return results, cobra.ShellCompDirectiveNoFileComp
}

// completeActiveEntries returns active (non-disabled) IPs and domains for shell completion.
func completeActiveEntries() ([]string, cobra.ShellCompDirective) {
	_, block, _, err := readBlock()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	seen := make(map[string]bool)
	var results []string
	for _, entry := range block.Entries {
		switch entry.DisableType {
		case parser.DisableNone:
			if !seen[entry.IP] {
				seen[entry.IP] = true
				results = append(results, entry.IP)
			}
			for _, d := range entry.Domains {
				if !seen[d] {
					seen[d] = true
					results = append(results, d)
				}
			}
		case parser.DisableDomain:
			// IP itself is still active, just this domain is disabled
			if !seen[entry.IP] {
				seen[entry.IP] = true
				results = append(results, entry.IP)
			}
		}
	}
	return results, cobra.ShellCompDirectiveNoFileComp
}

// completeDisabledEntries returns disabled IPs and domains for shell completion.
func completeDisabledEntries() ([]string, cobra.ShellCompDirective) {
	_, block, _, err := readBlock()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	seen := make(map[string]bool)
	var results []string
	for _, entry := range block.Entries {
		switch entry.DisableType {
		case parser.DisableIP:
			if !seen[entry.IP] {
				seen[entry.IP] = true
				results = append(results, entry.IP)
			}
		case parser.DisableDomain:
			for _, d := range entry.Domains {
				if !seen[d] {
					seen[d] = true
					results = append(results, d)
				}
			}
		}
	}
	return results, cobra.ShellCompDirectiveNoFileComp
}
