package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

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
