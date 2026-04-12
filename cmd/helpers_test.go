package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

// --- readInput ---

func TestReadInput_File(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "input.txt")
	os.WriteFile(path, []byte("192.168.1.1  web.local\n"), 0644)

	data, err := readInput(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != "192.168.1.1  web.local\n" {
		t.Errorf("unexpected content: %q", string(data))
	}
}

func TestReadInput_FileNotFound(t *testing.T) {
	_, err := readInput("/nonexistent/file")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestReadInput_Stdin(t *testing.T) {
	r, w, _ := os.Pipe()
	origStdin := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = origStdin }()

	go func() {
		w.Write([]byte("10.0.0.1  db.local\n"))
		w.Close()
	}()

	data, err := readInput("-")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != "10.0.0.1  db.local\n" {
		t.Errorf("unexpected content: %q", string(data))
	}
}

// --- parseHostsContent ---

func TestParseHostsContent_ValidHosts(t *testing.T) {
	input := []byte("192.168.1.1  web.local\n10.0.0.1  db.local\n")
	content, err := parseHostsContent(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if content == "" {
		t.Error("expected non-empty content")
	}
}

func TestParseHostsContent_ValidJSON(t *testing.T) {
	input := []byte(`{"192.168.1.1":["web.local","api.local"],"10.0.0.1":["db.local"]}`)
	content, err := parseHostsContent(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if content == "" {
		t.Error("expected non-empty content")
	}
}

func TestParseHostsContent_Empty(t *testing.T) {
	_, err := parseHostsContent([]byte(""))
	if err == nil {
		t.Error("expected error for empty input")
	}
}

func TestParseHostsContent_Whitespace(t *testing.T) {
	_, err := parseHostsContent([]byte("   \n  \n"))
	if err == nil {
		t.Error("expected error for whitespace-only input")
	}
}

func TestParseHostsContent_InvalidJSON(t *testing.T) {
	_, err := parseHostsContent([]byte(`{bad json}`))
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestParseHostsContent_EmptyJSON(t *testing.T) {
	_, err := parseHostsContent([]byte(`{}`))
	if err == nil {
		t.Error("expected error for empty JSON")
	}
}

func TestParseHostsContent_InvalidIPInJSON(t *testing.T) {
	_, err := parseHostsContent([]byte(`{"not-an-ip":["web.local"]}`))
	if err == nil {
		t.Error("expected error for invalid IP in JSON")
	}
}

func TestParseHostsContent_EmptyDomainsInJSON(t *testing.T) {
	_, err := parseHostsContent([]byte(`{"192.168.1.1":[]}`))
	if err == nil {
		t.Error("expected error for empty domains in JSON")
	}
}

func TestParseHostsContent_GarbageText(t *testing.T) {
	_, err := parseHostsContent([]byte("hello world garbage"))
	if err == nil {
		t.Error("expected error for garbage text")
	}
}

func TestParseHostsContent_MixedValidInvalid(t *testing.T) {
	input := []byte("192.168.1.1  web.local\ngarbageline\n")
	content, err := parseHostsContent(input)
	if err != nil {
		t.Fatalf("unexpected error: %v (should pass with at least one valid line)", err)
	}
	if content == "" {
		t.Error("expected non-empty content")
	}
}

func TestParseHostsContent_CommentsOnly(t *testing.T) {
	_, err := parseHostsContent([]byte("# just a comment\n# another\n"))
	if err == nil {
		t.Error("expected error for comments-only input")
	}
}

func TestParseHostsContent_HostsWithComments(t *testing.T) {
	input := []byte("# header\n192.168.1.1  web.local\n# footer\n")
	content, err := parseHostsContent(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if content == "" {
		t.Error("expected non-empty content")
	}
}
