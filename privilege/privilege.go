package privilege

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// WriteFilePrivileged writes content to a file, using privilege escalation if needed.
func WriteFilePrivileged(path string, content []byte) error {
	// try direct write first
	err := os.WriteFile(path, content, 0644)
	if err == nil {
		return nil
	}

	if !os.IsPermission(err) {
		return err
	}

	// need privilege escalation
	return writeWithEscalation(path, content)
}

func writeWithEscalation(path string, content []byte) error {
	if runtime.GOOS == "windows" {
		return writeWithEscalationWindows(path, content)
	}
	return writeWithEscalationUnix(path, content)
}

func writeWithEscalationUnix(path string, content []byte) error {
	// try sudo
	if _, err := exec.LookPath("sudo"); err == nil {
		cmd := exec.Command("sudo", "tee", path)
		cmd.Stdin = bytes.NewReader(content)
		cmd.Stdout = nil
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	// try doas
	if _, err := exec.LookPath("doas"); err == nil {
		cmd := exec.Command("doas", "tee", path)
		cmd.Stdin = bytes.NewReader(content)
		cmd.Stdout = nil
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	return fmt.Errorf("permission denied: please run with sudo or as root")
}

func writeWithEscalationWindows(path string, content []byte) error {
	// try sudo (Windows 24H2+)
	if _, err := exec.LookPath("sudo"); err == nil {
		cmd := exec.Command("sudo", "tee", path)
		cmd.Stdin = bytes.NewReader(content)
		cmd.Stdout = nil
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	// try gsudo
	if _, err := exec.LookPath("gsudo"); err == nil {
		cmd := exec.Command("gsudo", "tee", path)
		cmd.Stdin = bytes.NewReader(content)
		cmd.Stdout = nil
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	return fmt.Errorf("permission denied: please run this terminal as Administrator")
}
