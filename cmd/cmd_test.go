package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// runCmd executes a hostfile command with the given args and temp hosts file.
// Returns stdout output and any error.
func runCmd(t *testing.T, hostsPath string, args ...string) (string, error) {
	t.Helper()

	// reset and configure
	cmd := RootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(append(args, "--hosts-file", hostsPath))

	// override hostsFile via flag
	err := cmd.Execute()

	return buf.String(), err
}

// runCmdWithStdin executes a command with injected stdin.
func runCmdWithStdin(t *testing.T, hostsPath string, stdinContent string, args ...string) (string, error) {
	t.Helper()

	origStdin := stdin
	defer func() { stdin = origStdin }()
	stdin = strings.NewReader(stdinContent)

	return runCmd(t, hostsPath, args...)
}

// setupHosts creates a temp hosts file with optional initial content.
func setupHosts(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "hosts")
	if content != "" {
		os.WriteFile(path, []byte(content), 0644)
	}
	return path
}

// readHosts reads the content of a temp hosts file.
func readHosts(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("cannot read hosts file: %v", err)
	}
	return string(data)
}

// --- Integration tests ---

func TestCmd_AddAndShow(t *testing.T) {
	path := setupHosts(t, "")

	runCmd(t, path, "add", "192.168.1.100", "web.local", "api.local")
	runCmd(t, path, "add", "10.0.0.1", "db.local")

	content := readHosts(t, path)
	if !strings.Contains(content, "192.168.1.100  web.local api.local") {
		t.Errorf("expected web.local entry, got:\n%s", content)
	}
	if !strings.Contains(content, "10.0.0.1  db.local") {
		t.Errorf("expected db.local entry, got:\n%s", content)
	}
}

func TestCmd_AddDuplicate(t *testing.T) {
	path := setupHosts(t, "")

	runCmd(t, path, "add", "192.168.1.100", "web.local")
	out, _ := captureOutput(t, func() {
		runCmd(t, path, "add", "192.168.1.100", "web.local")
	})

	if !strings.Contains(out, "already exists") {
		t.Errorf("expected 'already exists' message, got: %s", out)
	}
}

func TestCmd_Remove(t *testing.T) {
	path := setupHosts(t, "")

	runCmd(t, path, "add", "192.168.1.100", "web.local", "api.local")
	runCmd(t, path, "remove", "api.local")

	content := readHosts(t, path)
	if strings.Contains(content, "api.local") {
		t.Error("api.local should have been removed")
	}
	if !strings.Contains(content, "web.local") {
		t.Error("web.local should still exist")
	}
}

func TestCmd_DisableEnable(t *testing.T) {
	path := setupHosts(t, "")

	runCmd(t, path, "add", "192.168.1.100", "web.local", "api.local")
	runCmd(t, path, "disable", "api.local")

	content := readHosts(t, path)
	if !strings.Contains(content, "#[disable-domain]") {
		t.Errorf("expected disable-domain marker, got:\n%s", content)
	}

	runCmd(t, path, "enable", "api.local")
	content = readHosts(t, path)
	if strings.Contains(content, "#[disable-domain]") {
		t.Error("disable-domain marker should be gone after enable")
	}
}

func TestCmd_DisableIP(t *testing.T) {
	path := setupHosts(t, "")

	runCmd(t, path, "add", "10.0.0.1", "db.local")
	runCmd(t, path, "disable", "10.0.0.1")

	content := readHosts(t, path)
	if !strings.Contains(content, "#[disable-ip]") {
		t.Errorf("expected disable-ip marker, got:\n%s", content)
	}
}

func TestCmd_Search(t *testing.T) {
	path := setupHosts(t, "")
	runCmd(t, path, "add", "192.168.1.100", "web.local", "api.local")

	out, _ := captureOutput(t, func() {
		runCmd(t, path, "search", "web.local")
	})
	if !strings.Contains(out, "192.168.1.100") {
		t.Errorf("expected IP in search result, got: %s", out)
	}

	out, _ = captureOutput(t, func() {
		runCmd(t, path, "search", "192.168.1.100")
	})
	if !strings.Contains(out, "web.local") {
		t.Errorf("expected domain in search result, got: %s", out)
	}
}

func TestCmd_ShowJSON(t *testing.T) {
	path := setupHosts(t, "")
	runCmd(t, path, "add", "192.168.1.100", "web.local")

	out, _ := captureOutput(t, func() {
		runCmd(t, path, "show", "--json")
	})
	if !strings.Contains(out, `"192.168.1.100"`) {
		t.Errorf("expected JSON with IP, got: %s", out)
	}
	if !strings.Contains(out, `"web.local"`) {
		t.Errorf("expected JSON with domain, got: %s", out)
	}
}

func TestCmd_ApplyFile(t *testing.T) {
	path := setupHosts(t, "")
	runCmd(t, path, "add", "10.0.0.1", "old.local")

	// create input file
	dir := t.TempDir()
	inputPath := filepath.Join(dir, "new.hosts")
	os.WriteFile(inputPath, []byte("192.168.1.1  new.local\n"), 0644)

	runCmd(t, path, "apply", inputPath)

	content := readHosts(t, path)
	if strings.Contains(content, "old.local") {
		t.Error("old entry should be replaced")
	}
	if !strings.Contains(content, "new.local") {
		t.Error("new entry should exist")
	}
}

func TestCmd_ApplyJSON(t *testing.T) {
	path := setupHosts(t, "")

	dir := t.TempDir()
	inputPath := filepath.Join(dir, "config.json")
	os.WriteFile(inputPath, []byte(`{"192.168.1.100":["web.local","api.local"]}`), 0644)

	runCmd(t, path, "apply", inputPath)

	content := readHosts(t, path)
	if !strings.Contains(content, "192.168.1.100") || !strings.Contains(content, "web.local") {
		t.Errorf("expected JSON entries applied, got:\n%s", content)
	}
}

// TestCmd_ApplyInvalid is tested via parseHostsContent unit tests in helpers_test.go.
// Cannot integration-test here because exitWithError calls os.Exit(1).

func TestCmd_MergeFile(t *testing.T) {
	path := setupHosts(t, "")
	runCmd(t, path, "add", "10.0.0.1", "existing.local")

	dir := t.TempDir()
	inputPath := filepath.Join(dir, "merge.hosts")
	os.WriteFile(inputPath, []byte("192.168.1.1  new.local\n"), 0644)

	runCmd(t, path, "merge", inputPath)

	content := readHosts(t, path)
	if !strings.Contains(content, "existing.local") {
		t.Error("existing entry should be preserved")
	}
	if !strings.Contains(content, "new.local") {
		t.Error("merged entry should exist")
	}
}

func TestCmd_Clean(t *testing.T) {
	path := setupHosts(t, "")
	runCmd(t, path, "add", "10.0.0.1", "web.local")

	runCmdWithStdin(t, path, "y\n", "clean")

	content := readHosts(t, path)
	if strings.Contains(content, "web.local") {
		t.Error("entries should be cleaned")
	}
}

func TestCmd_Clean_Abort(t *testing.T) {
	path := setupHosts(t, "")
	runCmd(t, path, "add", "10.0.0.1", "web.local")

	runCmdWithStdin(t, path, "n\n", "clean")

	content := readHosts(t, path)
	if !strings.Contains(content, "web.local") {
		t.Error("entries should be preserved after abort")
	}
}

func TestCmd_ShowEmpty(t *testing.T) {
	path := setupHosts(t, "")

	out, _ := captureOutput(t, func() {
		runCmd(t, path, "show")
	})
	if !strings.Contains(out, "managed block is empty") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestCmd_SearchNotFound(t *testing.T) {
	path := setupHosts(t, "")

	_, err := runCmd(t, path, "search", "nonexistent.local")
	if err == nil {
		t.Error("expected error for search with no results")
	}
	if !strings.Contains(err.Error(), "no entries found") {
		t.Errorf("expected not-found in error, got: %s", err)
	}
}

func TestCmd_PreservesOutsideContent(t *testing.T) {
	// hosts file with existing content and a managed block
	initial := "127.0.0.1 localhost\n::1 localhost\n"
	path := setupHosts(t, initial)

	runCmd(t, path, "add", "10.0.0.1", "web.local")

	content := readHosts(t, path)
	if !strings.Contains(content, "127.0.0.1 localhost") {
		t.Error("original content should be preserved")
	}
	if !strings.Contains(content, "10.0.0.1  web.local") {
		t.Error("new entry should exist")
	}
}

// captureOutput captures os.Stdout during fn execution.
func captureOutput(t *testing.T, fn func()) (string, error) {
	t.Helper()
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String(), nil
}
