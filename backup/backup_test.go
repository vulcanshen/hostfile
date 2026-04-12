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

func TestCreateRawAndRestoreRaw(t *testing.T) {
	_, cleanup := setupTestDir(t)
	defer cleanup()

	raw := []byte("127.0.0.1 localhost\n::1 localhost\n")
	if err := CreateRaw("origin", raw); err != nil {
		t.Fatalf("create raw error: %v", err)
	}

	data, err := RestoreRaw("origin")
	if err != nil {
		t.Fatalf("restore raw error: %v", err)
	}
	if string(data) != string(raw) {
		t.Errorf("expected %q, got %q", string(raw), string(data))
	}
}

func TestIsRaw(t *testing.T) {
	_, cleanup := setupTestDir(t)
	defer cleanup()

	// raw backup (no block markers)
	CreateRaw("raw-test", []byte("127.0.0.1 localhost\n"))
	raw, err := IsRaw("raw-test")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !raw {
		t.Error("expected raw=true for file without block markers")
	}

	// managed backup (has block markers)
	block := &parser.ManagedBlock{
		Entries: []parser.HostEntry{
			{IP: "10.0.0.1", Domains: []string{"web.local"}, DisableType: parser.DisableNone},
		},
	}
	Create("managed-test", block)
	raw, err = IsRaw("managed-test")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if raw {
		t.Error("expected raw=false for file with block markers")
	}
}

func TestExists(t *testing.T) {
	_, cleanup := setupTestDir(t)
	defer cleanup()

	exists, err := Exists("nope")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if exists {
		t.Error("expected false for nonexistent backup")
	}

	CreateRaw("exists-test", []byte("data"))
	exists, err = Exists("exists-test")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !exists {
		t.Error("expected true for existing backup")
	}
}

func TestPath(t *testing.T) {
	dir, cleanup := setupTestDir(t)
	defer cleanup()

	path, err := Path("my-save")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	expected := filepath.Join(dir, "my-save.hostfile")
	if path != expected {
		t.Errorf("expected %q, got %q", expected, path)
	}
}

func TestRestoreRaw_NotFound(t *testing.T) {
	_, cleanup := setupTestDir(t)
	defer cleanup()

	_, err := RestoreRaw("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent raw backup")
	}
}
