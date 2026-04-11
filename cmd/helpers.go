package cmd

import (
	"fmt"
	"os"

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

// exitWithError prints an error and exits.
func exitWithError(err error) {
	fmt.Fprintln(os.Stderr, "error:", err)
	os.Exit(1)
}
