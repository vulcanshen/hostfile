package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// stdin is the input source for interactive prompts (injectable for testing).
var stdin io.Reader = os.Stdin

// confirm asks the user a yes/no question. Returns true if the user answers yes.
func confirm(message string) bool {
	fmt.Printf("%s [y/N]: ", message)
	scanner := bufio.NewScanner(stdin)
	if scanner.Scan() {
		answer := strings.TrimSpace(strings.ToLower(scanner.Text()))
		return answer == "y" || answer == "yes"
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error reading input: %v\n", err)
	}
	return false
}
