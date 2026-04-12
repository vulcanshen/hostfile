package privilege

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteFilePrivileged_DirectWrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "hosts")

	content := []byte("127.0.0.1  localhost\n")
	err := WriteFilePrivileged(path, content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("cannot read file: %v", err)
	}
	if string(data) != string(content) {
		t.Errorf("expected %q, got %q", string(content), string(data))
	}
}

func TestWriteFilePrivileged_Overwrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "hosts")

	os.WriteFile(path, []byte("old content\n"), 0644)

	newContent := []byte("new content\n")
	err := WriteFilePrivileged(path, newContent)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(path)
	if string(data) != string(newContent) {
		t.Errorf("expected %q, got %q", string(newContent), string(data))
	}
}

func TestWriteFilePrivileged_NonExistentDir(t *testing.T) {
	path := "/nonexistent/dir/hosts"
	err := WriteFilePrivileged(path, []byte("data"))
	if err == nil {
		t.Error("expected error for nonexistent directory")
	}
}
