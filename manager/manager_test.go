package manager

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/vulcanshen/hostfile/parser"
)

func newBlock(entries ...parser.HostEntry) *parser.ManagedBlock {
	return &parser.ManagedBlock{Entries: entries}
}

func normalEntry(ip string, domains ...string) parser.HostEntry {
	return parser.HostEntry{IP: ip, Domains: domains, DisableType: parser.DisableNone}
}

func disabledIPEntry(ip string, domains ...string) parser.HostEntry {
	return parser.HostEntry{IP: ip, Domains: domains, DisableType: parser.DisableIP}
}

func disabledDomainEntry(ip string, domain string) parser.HostEntry {
	return parser.HostEntry{IP: ip, Domains: []string{domain}, DisableType: parser.DisableDomain}
}

// --- ReadHostsFile / WriteHostsFile ---

func TestReadWriteHostsFile_NewFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "hosts")

	before, block, after, err := ReadHostsFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if before != "" || after != "" {
		t.Error("expected empty before/after for new file")
	}
	if len(block.Entries) != 0 {
		t.Error("expected empty block for new file")
	}

	block.Entries = append(block.Entries, normalEntry("10.0.0.1", "web.local"))
	err = WriteHostsFile(path, before, block, after)
	if err != nil {
		t.Fatalf("write error: %v", err)
	}

	data, _ := os.ReadFile(path)
	content := string(data)
	if !strings.Contains(content, "10.0.0.1  web.local") {
		t.Errorf("expected entry in output, got:\n%s", content)
	}
	if !strings.HasPrefix(content, parser.BlockStart) {
		t.Errorf("expected block at top, got:\n%s", content)
	}
}

func TestReadHostsFile_ExistingContentNoBlock(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "hosts")
	os.WriteFile(path, []byte("127.0.0.1 localhost\n::1 localhost\n"), 0644)

	before, block, after, err := ReadHostsFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if before != "" {
		t.Error("expected empty before")
	}
	if after != "127.0.0.1 localhost\n::1 localhost\n" {
		t.Errorf("expected existing content in after, got: %q", after)
	}
	if len(block.Entries) != 0 {
		t.Error("expected empty block")
	}
}

func TestReadHostsFile_WithBlock(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "hosts")
	content := "127.0.0.1 localhost\n" +
		parser.BlockStart + "\n" +
		"10.0.0.1  web.local\n" +
		parser.BlockEnd + "\n" +
		"# trailing\n"
	os.WriteFile(path, []byte(content), 0644)

	before, block, after, err := ReadHostsFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if before != "127.0.0.1 localhost\n" {
		t.Errorf("unexpected before: %q", before)
	}
	if after != "\n# trailing\n" {
		t.Errorf("unexpected after: %q", after)
	}
	if len(block.Entries) != 1 || block.Entries[0].IP != "10.0.0.1" {
		t.Errorf("unexpected block: %+v", block.Entries)
	}
}

// --- Add ---

func TestAdd_NewIP(t *testing.T) {
	block := newBlock()
	conflicts, err := Add(block, "10.0.0.1", []string{"web.local", "api.local"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(conflicts) != 0 {
		t.Fatalf("unexpected conflicts: %v", conflicts)
	}
	if len(block.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(block.Entries))
	}
	if len(block.Entries[0].Domains) != 2 {
		t.Errorf("expected 2 domains, got %v", block.Entries[0].Domains)
	}
}

func TestAdd_MergeExistingIP(t *testing.T) {
	block := newBlock(normalEntry("10.0.0.1", "web.local"))
	conflicts, err := Add(block, "10.0.0.1", []string{"api.local"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(conflicts) != 0 {
		t.Fatalf("unexpected conflicts: %v", conflicts)
	}
	if len(block.Entries[0].Domains) != 2 {
		t.Errorf("expected 2 domains, got %v", block.Entries[0].Domains)
	}
}

func TestAdd_DuplicateDomain(t *testing.T) {
	block := newBlock(normalEntry("10.0.0.1", "web.local"))
	Add(block, "10.0.0.1", []string{"web.local"})
	if len(block.Entries[0].Domains) != 1 {
		t.Errorf("expected no duplicate, got %v", block.Entries[0].Domains)
	}
}

func TestAdd_Conflict(t *testing.T) {
	block := newBlock(normalEntry("10.0.0.1", "web.local"))
	conflicts, err := Add(block, "10.0.0.2", []string{"web.local"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(conflicts) != 1 {
		t.Fatalf("expected 1 conflict, got %d", len(conflicts))
	}
	if conflicts[0].Domain != "web.local" || conflicts[0].CurrentIP != "10.0.0.1" {
		t.Errorf("unexpected conflict: %+v", conflicts[0])
	}
}

func TestAdd_InvalidIP(t *testing.T) {
	block := newBlock()
	_, err := Add(block, "not-an-ip", []string{"web.local"})
	if err == nil {
		t.Error("expected error for invalid IP")
	}
}

func TestAddForce(t *testing.T) {
	block := newBlock(normalEntry("10.0.0.1", "web.local", "api.local"))
	err := AddForce(block, "10.0.0.2", []string{"web.local"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// web.local should be moved from 10.0.0.1 to 10.0.0.2
	if len(block.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(block.Entries))
	}
	if containsDomain(block.Entries[0].Domains, "web.local") {
		t.Error("web.local should have been removed from 10.0.0.1")
	}
	if !containsDomain(block.Entries[1].Domains, "web.local") {
		t.Error("web.local should be on 10.0.0.2")
	}
}

// --- Remove ---

func TestRemove_IP(t *testing.T) {
	block := newBlock(
		normalEntry("10.0.0.1", "web.local"),
		normalEntry("10.0.0.2", "api.local"),
	)
	result := Remove(block, "10.0.0.1")
	if !result.Removed {
		t.Error("expected removed=true")
	}
	if len(block.Entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(block.Entries))
	}
}

func TestRemove_Domain(t *testing.T) {
	block := newBlock(normalEntry("10.0.0.1", "web.local", "api.local"))
	result := Remove(block, "web.local")
	if !result.Removed {
		t.Error("expected removed=true")
	}
	if len(block.Entries[0].Domains) != 1 || block.Entries[0].Domains[0] != "api.local" {
		t.Errorf("expected [api.local], got %v", block.Entries[0].Domains)
	}
}

func TestRemove_LastDomain(t *testing.T) {
	block := newBlock(normalEntry("10.0.0.1", "web.local"))
	result := Remove(block, "web.local")
	if !result.LastDomain {
		t.Error("expected LastDomain=true")
	}
	if result.LastDomainIP != "10.0.0.1" {
		t.Errorf("expected LastDomainIP=10.0.0.1, got %s", result.LastDomainIP)
	}
	// entry should still exist (not removed yet)
	if len(block.Entries) != 1 {
		t.Error("expected entry to still exist before confirmation")
	}
}

func TestRemove_DisabledDomain(t *testing.T) {
	block := newBlock(disabledDomainEntry("10.0.0.1", "web.local"))
	result := Remove(block, "web.local")
	if !result.Removed {
		t.Error("expected removed=true")
	}
	if len(block.Entries) != 0 {
		t.Error("expected empty block")
	}
}

func TestRemoveConfirmed(t *testing.T) {
	block := newBlock(normalEntry("10.0.0.1", "web.local"))
	RemoveConfirmed(block, "10.0.0.1")
	if len(block.Entries) != 0 {
		t.Error("expected empty block after confirmed removal")
	}
}

// --- Search ---

func TestSearch_ByIP(t *testing.T) {
	block := newBlock(
		normalEntry("10.0.0.1", "web.local"),
		normalEntry("10.0.0.2", "api.local"),
		disabledDomainEntry("10.0.0.1", "old.local"),
	)
	results := Search(block, "10.0.0.1")
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}

func TestSearch_ByDomain(t *testing.T) {
	block := newBlock(
		normalEntry("10.0.0.1", "web.local", "api.local"),
		normalEntry("10.0.0.2", "db.local"),
	)
	results := Search(block, "api.local")
	if len(results) != 1 || results[0].IP != "10.0.0.1" {
		t.Errorf("unexpected results: %+v", results)
	}
}

func TestSearch_DisabledDomain(t *testing.T) {
	block := newBlock(disabledDomainEntry("10.0.0.1", "old.local"))
	results := Search(block, "old.local")
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
}

// --- Enable ---

func TestEnable_IP(t *testing.T) {
	block := newBlock(disabledIPEntry("10.0.0.1", "web.local"))
	err := Enable(block, "10.0.0.1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if block.Entries[0].DisableType != parser.DisableNone {
		t.Error("expected entry to be enabled")
	}
}

func TestEnable_Domain(t *testing.T) {
	block := newBlock(
		normalEntry("10.0.0.1", "api.local"),
		disabledDomainEntry("10.0.0.1", "web.local"),
	)
	err := Enable(block, "web.local")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// web.local should be added back to the 10.0.0.1 line
	if len(block.Entries) != 1 {
		t.Fatalf("expected 1 entry (merged), got %d", len(block.Entries))
	}
	if !containsDomain(block.Entries[0].Domains, "web.local") {
		t.Error("expected web.local to be re-added")
	}
}

func TestEnable_DomainNewLine(t *testing.T) {
	// when the IP line doesn't exist, enabling should create a new one
	block := newBlock(disabledDomainEntry("10.0.0.1", "web.local"))
	err := Enable(block, "web.local")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(block.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(block.Entries))
	}
	if block.Entries[0].DisableType != parser.DisableNone {
		t.Error("expected new line to be enabled")
	}
}

func TestEnable_NotFound(t *testing.T) {
	block := newBlock(normalEntry("10.0.0.1", "web.local"))
	err := Enable(block, "10.0.0.2")
	if err == nil {
		t.Error("expected error for not-found IP")
	}
}

// --- Disable ---

func TestDisable_IP(t *testing.T) {
	block := newBlock(normalEntry("10.0.0.1", "web.local"))
	result, err := Disable(block, "10.0.0.1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Done {
		t.Error("expected Done=true")
	}
	if block.Entries[0].DisableType != parser.DisableIP {
		t.Error("expected entry to be disabled")
	}
}

func TestDisable_Domain(t *testing.T) {
	block := newBlock(normalEntry("10.0.0.1", "web.local", "api.local"))
	result, err := Disable(block, "web.local")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Done {
		t.Error("expected Done=true")
	}
	// web.local should be removed from the entry and added as disable-domain
	if containsDomain(block.Entries[0].Domains, "web.local") {
		t.Error("web.local should have been removed from normal entry")
	}
	if len(block.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(block.Entries))
	}
	if block.Entries[1].DisableType != parser.DisableDomain {
		t.Error("expected disable-domain entry")
	}
}

func TestDisable_LastDomain(t *testing.T) {
	block := newBlock(normalEntry("10.0.0.1", "web.local"))
	result, err := Disable(block, "web.local")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.LastDomain {
		t.Error("expected LastDomain=true")
	}
	if result.LastDomainIP != "10.0.0.1" {
		t.Errorf("expected LastDomainIP=10.0.0.1, got %s", result.LastDomainIP)
	}
}

func TestDisable_NotFound(t *testing.T) {
	block := newBlock(normalEntry("10.0.0.1", "web.local"))
	_, err := Disable(block, "10.0.0.2")
	if err == nil {
		t.Error("expected error for not-found IP")
	}
}

// --- Clean ---

func TestClean(t *testing.T) {
	block := newBlock(normalEntry("10.0.0.1", "web.local"), normalEntry("10.0.0.2", "api.local"))
	Clean(block)
	if len(block.Entries) != 0 {
		t.Error("expected empty block after clean")
	}
}

// --- Apply ---

func TestApply(t *testing.T) {
	block := newBlock(normalEntry("10.0.0.1", "old.local"))
	Apply(block, "10.0.0.2  new.local\n10.0.0.3  other.local")
	if len(block.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(block.Entries))
	}
	if block.Entries[0].IP != "10.0.0.2" {
		t.Errorf("expected 10.0.0.2, got %s", block.Entries[0].IP)
	}
}

// --- Merge ---

func TestMerge_NewIP(t *testing.T) {
	block := newBlock(normalEntry("10.0.0.1", "web.local"))
	Merge(block, "10.0.0.2  api.local")
	if len(block.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(block.Entries))
	}
}

func TestMerge_ExistingIP(t *testing.T) {
	block := newBlock(normalEntry("10.0.0.1", "web.local"))
	Merge(block, "10.0.0.1  api.local")
	if len(block.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(block.Entries))
	}
	if len(block.Entries[0].Domains) != 2 {
		t.Errorf("expected 2 domains, got %v", block.Entries[0].Domains)
	}
}

func TestMerge_NoDuplicate(t *testing.T) {
	block := newBlock(normalEntry("10.0.0.1", "web.local"))
	Merge(block, "10.0.0.1  web.local")
	if len(block.Entries[0].Domains) != 1 {
		t.Errorf("expected no duplicate, got %v", block.Entries[0].Domains)
	}
}

// --- WriteHostsFile preserves surrounding content ---

func TestWriteHostsFile_Prepend(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "hosts")

	block := newBlock(normalEntry("10.0.0.1", "web.local"))
	after := "127.0.0.1 localhost\n"

	err := WriteHostsFile(path, "", block, after)
	if err != nil {
		t.Fatalf("write error: %v", err)
	}

	data, _ := os.ReadFile(path)
	content := string(data)
	// block should come before localhost
	blockIdx := strings.Index(content, parser.BlockStart)
	localhostIdx := strings.Index(content, "127.0.0.1 localhost")
	if blockIdx > localhostIdx {
		t.Errorf("expected block before localhost:\n%s", content)
	}
}
