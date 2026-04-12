package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/vulcanshen/hostfile/parser"
)

func captureStdout(fn func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestPrintEntries_Alignment(t *testing.T) {
	entries := []parser.HostEntry{
		{IP: "192.168.1.100", Domains: []string{"web.local"}, DisableType: parser.DisableNone},
		{IP: "10.0.0.1", Domains: []string{"db.local"}, DisableType: parser.DisableNone},
	}

	output := captureStdout(func() { printEntries(entries) })
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}

	// both lines should have the same IP column width
	idx1 := strings.Index(lines[0], "web.local")
	idx2 := strings.Index(lines[1], "db.local")
	if idx1 != idx2 {
		t.Errorf("domains not aligned: web.local at %d, db.local at %d", idx1, idx2)
	}
}

func TestPrintEntries_DisabledLabel(t *testing.T) {
	entries := []parser.HostEntry{
		{IP: "10.0.0.1", Domains: []string{"db.local"}, DisableType: parser.DisableIP},
	}

	output := captureStdout(func() { printEntries(entries) })
	if !strings.Contains(output, "[disabled]") {
		t.Errorf("expected [disabled] label, got: %s", output)
	}
}

func TestPrintEntries_DisabledDomainLabel(t *testing.T) {
	entries := []parser.HostEntry{
		{IP: "192.168.1.100", Domains: []string{"api.local"}, DisableType: parser.DisableDomain},
	}

	output := captureStdout(func() { printEntries(entries) })
	if !strings.Contains(output, "[disabled]") {
		t.Errorf("expected [disabled] label, got: %s", output)
	}
}

func TestPrintEntries_Empty(t *testing.T) {
	output := captureStdout(func() { printEntries([]parser.HostEntry{}) })
	if output != "" {
		t.Errorf("expected no output for empty entries, got: %q", output)
	}
}

func TestPrintEntry_SingleEntry(t *testing.T) {
	entry := parser.HostEntry{IP: "10.0.0.1", Domains: []string{"web.local"}, DisableType: parser.DisableNone}
	output := captureStdout(func() { printEntry(entry) })
	if !strings.Contains(output, "10.0.0.1") || !strings.Contains(output, "web.local") {
		t.Errorf("unexpected output: %s", output)
	}
}

func TestIsTTY_Pipe(t *testing.T) {
	// when running in tests, stdout is typically a pipe, not a TTY
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w
	result := isTTY()
	os.Stdout = old
	w.Close()
	r.Close()

	if result {
		t.Error("expected false for pipe")
	}
}
