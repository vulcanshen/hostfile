package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open",
	Short: "Open the hosts file in the default editor",
	Long:  "Open the hosts file in the system default editor with appropriate privileges.",
	Args:  cobra.NoArgs,
	ValidArgsFunction: cobra.NoFileCompletions,
	Run: func(cmd *cobra.Command, args []string) {
		path := hostsFile
		fmt.Printf("opening %s\n", path)

		switch runtime.GOOS {
		case "windows":
			openWindows(path)
		case "darwin":
			openDarwin(path)
		default:
			openUnix(path)
		}
	},
}

func init() {
	rootCmd.AddCommand(openCmd)
}

// needsEscalation checks if the file requires elevated privileges to write.
func needsEscalation(path string) bool {
	f, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		return os.IsPermission(err)
	}
	f.Close()
	return false
}

// openDarwin opens the hosts file on macOS.
// Uses $EDITOR if set, falls back to nano.
func openDarwin(path string) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}
	openWithTerminalEditor(editor, path)
}

// openUnix opens the hosts file on Linux.
// Uses $EDITOR if set, falls back to vi.
func openUnix(path string) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}
	openWithTerminalEditor(editor, path)
}

// openWindows opens the hosts file on Windows.
// Uses notepad with elevation if needed.
func openWindows(path string) {
	if !needsEscalation(path) {
		runGUI(exec.Command("notepad", path))
		return
	}

	// try sudo (Windows 24H2+)
	if _, err := exec.LookPath("sudo"); err == nil {
		runGUI(exec.Command("sudo", "notepad", path))
		return
	}

	// try gsudo
	if _, err := exec.LookPath("gsudo"); err == nil {
		runGUI(exec.Command("gsudo", "notepad", path))
		return
	}

	fmt.Fprintln(os.Stderr, "error: permission denied — sudo and gsudo not found, please install gsudo or run as Administrator")
	os.Exit(1)
}

// openWithTerminalEditor launches a terminal editor, with sudo/doas if needed.
func openWithTerminalEditor(editor, path string) {
	if !needsEscalation(path) {
		runInteractive(exec.Command(editor, path))
		return
	}

	// try sudo
	if _, err := exec.LookPath("sudo"); err == nil {
		runInteractive(exec.Command("sudo", editor, path))
		return
	}

	// try doas
	if _, err := exec.LookPath("doas"); err == nil {
		runInteractive(exec.Command("doas", editor, path))
		return
	}

	fmt.Fprintln(os.Stderr, "error: permission denied — sudo and doas not found, please install sudo or run as root")
	os.Exit(1)
}

// runInteractive runs a command attached to the current terminal (stdin/stdout/stderr).
func runInteractive(cmd *exec.Cmd) {
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		exitWithError(fmt.Errorf("editor exited with error: %w", err))
	}
}

// runGUI runs a command with stdout/stderr attached (no stdin — for GUI editors).
func runGUI(cmd *exec.Cmd) {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		exitWithError(fmt.Errorf("failed to open editor: %w", err))
	}
}
