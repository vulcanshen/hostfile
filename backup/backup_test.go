package backup

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vulcanshen/hostfile/parser"
)

// override backupBasePath for tests
func setupTestDir(t *testing.T) (string, func()) {
	t.Helper()
	dir := t.TempDir()

	origFn := backupBasePathFn
	backupBasePathFn = func() (string, error) {
		return dir, nil
	}
	return dir, func() { backupBasePathFn = origFn }
}

func TestCreateAndRestore(t *testing.T) {
	_, cleanup := setupTestDir(t)
	defer cleanup()

	block := &parser.ManagedBlock{
		Entries: []parser.HostEntry{
			{IP: "10.0.0.1", Domains: []string{"web.local"}, DisableType: parser.DisableNone},
			{IP: "10.0.0.2", Domains: []string{"api.local"}, DisableType: parser.DisableIP},
		},
	}

	err := Create("test-backup", block)
	if err != nil {
		t.Fatalf("create error: %v", err)
	}

	restored, err := Restore("test-backup")
	if err != nil {
		t.Fatalf("restore error: %v", err)
	}
	if len(restored.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(restored.Entries))
	}
	if restored.Entries[0].IP != "10.0.0.1" {
		t.Errorf("expected 10.0.0.1, got %s", restored.Entries[0].IP)
	}
	if restored.Entries[1].DisableType != parser.DisableIP {
		t.Errorf("expected DisableIP, got %v", restored.Entries[1].DisableType)
	}
}

func TestList(t *testing.T) {
	dir, cleanup := setupTestDir(t)
	defer cleanup()

	// create some backup files
	os.WriteFile(filepath.Join(dir, "a.hostfile"), []byte("10.0.0.1  web.local\n"), 0644)
	os.WriteFile(filepath.Join(dir, "b.hostfile"), []byte("10.0.0.2  api.local\n"), 0644)
	os.WriteFile(filepath.Join(dir, "not-a-backup.txt"), []byte("ignored"), 0644)

	backups, err := List()
	if err != nil {
		t.Fatalf("list error: %v", err)
	}
	if len(backups) != 2 {
		t.Fatalf("expected 2 backups, got %d", len(backups))
	}
}

func TestDelete(t *testing.T) {
	dir, cleanup := setupTestDir(t)
	defer cleanup()

	os.WriteFile(filepath.Join(dir, "todel.hostfile"), []byte("data"), 0644)

	err := Delete("todel")
	if err != nil {
		t.Fatalf("delete error: %v", err)
	}

	// verify deleted
	_, err = os.Stat(filepath.Join(dir, "todel.hostfile"))
	if !os.IsNotExist(err) {
		t.Error("expected file to be deleted")
	}
}

func TestDelete_NotFound(t *testing.T) {
	_, cleanup := setupTestDir(t)
	defer cleanup()

	err := Delete("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent backup")
	}
}

func TestRestore_NotFound(t *testing.T) {
	_, cleanup := setupTestDir(t)
	defer cleanup()

	_, err := Restore("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent backup")
	}
}
