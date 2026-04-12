package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/vulcanshen/hostfile/parser"
)

func TestPrintJSON_ActiveOnly(t *testing.T) {
	entries := []parser.HostEntry{
		{IP: "192.168.1.100", Domains: []string{"web.local", "api.local"}, DisableType: parser.DisableNone},
		{IP: "10.0.0.1", Domains: []string{"db.local"}, DisableType: parser.DisableIP},
		{IP: "192.168.1.100", Domains: []string{"old.local"}, DisableType: parser.DisableDomain},
	}

	// capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printJSON(entries)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)

	var result map[string][]string
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}

	// only active entry should be present
	if len(result) != 1 {
		t.Fatalf("expected 1 IP in JSON, got %d: %v", len(result), result)
	}
	domains, ok := result["192.168.1.100"]
	if !ok {
		t.Fatal("expected 192.168.1.100 in JSON")
	}
	if len(domains) != 2 || domains[0] != "web.local" || domains[1] != "api.local" {
		t.Errorf("unexpected domains: %v", domains)
	}
}

func TestPrintJSON_Empty(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printJSON([]parser.HostEntry{})

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)

	var result map[string][]string
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty JSON, got %v", result)
	}
}

func TestPrintJSON_AllDisabled(t *testing.T) {
	entries := []parser.HostEntry{
		{IP: "10.0.0.1", Domains: []string{"db.local"}, DisableType: parser.DisableIP},
	}

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printJSON(entries)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)

	var result map[string][]string
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty JSON for all-disabled, got %v", result)
	}
}

func TestPrintJSON_MergeSameIP(t *testing.T) {
	entries := []parser.HostEntry{
		{IP: "10.0.0.1", Domains: []string{"web.local"}, DisableType: parser.DisableNone},
		{IP: "10.0.0.1", Domains: []string{"api.local"}, DisableType: parser.DisableNone},
	}

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printJSON(entries)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)

	var result map[string][]string
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	domains := result["10.0.0.1"]
	if len(domains) != 2 {
		t.Errorf("expected 2 domains merged, got %v", domains)
	}
}
