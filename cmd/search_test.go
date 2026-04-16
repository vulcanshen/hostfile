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
	// filter out empty lines (separators)
	var lines []string
	for _, l := range strings.Split(strings.TrimSpace(output), "\n") {
		if strings.TrimSpace(l) != "" {
			lines = append(lines, l)
		}
	}
	// header + 2 domain lines
	if len(lines) != 3 {
		t.Fatalf("expected 3 non-empty lines (header + 2 entries), got %d: %q", len(lines), output)
	}

	// header should contain IP and DOMAIN
	if !strings.Contains(lines[0], "IP") || !strings.Contains(lines[0], "DOMAIN") {
		t.Errorf("expected header with IP and DOMAIN, got: %s", lines[0])
	}

	// data lines should have the same domain column offset
	idx1 := strings.Index(lines[1], "web.local")
	idx2 := strings.Index(lines[2], "db.local")
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

func TestPrintEntries_MultiDomainPerLine(t *testing.T) {
	entries := []parser.HostEntry{
		{IP: "192.168.1.100", Domains: []string{"web.local", "api.local"}, DisableType: parser.DisableNone},
	}

	output := captureStdout(func() { printEntries(entries) })
	lines := strings.Split(strings.TrimSpace(output), "\n")
	// header + 2 domain lines (one per domain)
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines (header + 2 domains), got %d: %q", len(lines), output)
	}

	// first data line has IP, second has blank IP
	if !strings.Contains(lines[1], "192.168.1.100") {
		t.Errorf("first domain line should contain IP, got: %s", lines[1])
	}
	if !strings.Contains(lines[2], "api.local") {
		t.Errorf("second domain line should contain api.local, got: %s", lines[2])
	}
	// second line should NOT contain the IP
	if strings.Contains(lines[2], "192.168.1.100") {
		t.Errorf("second domain line should not repeat IP, got: %s", lines[2])
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
