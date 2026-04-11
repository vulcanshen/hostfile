package backup

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/vulcanshen/hostfile/parser"
)

const (
	backupDir = ".hostfile"
	backupExt = ".hostfile"
)

// BackupInfo holds metadata about a backup.
type BackupInfo struct {
	Name    string
	Path    string
	ModTime time.Time
}

// backupBasePathFn is a function variable for testing override.
var backupBasePathFn = defaultBackupBasePath

func defaultBackupBasePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine home directory: %w", err)
	}
	return filepath.Join(home, backupDir), nil
}

// Create saves the managed block content to a backup file.
func Create(name string, block *parser.ManagedBlock) error {
	base, err := backupBasePathFn()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(base, 0755); err != nil {
		return fmt.Errorf("cannot create backup directory: %w", err)
	}

	path := filepath.Join(base, name+backupExt)
	content := parser.FormatBlock(block)
	return os.WriteFile(path, []byte(content+"\n"), 0644)
}

// List returns all backups.
func List() ([]BackupInfo, error) {
	base, err := backupBasePathFn()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(base)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var backups []BackupInfo
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), backupExt) {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		name := strings.TrimSuffix(entry.Name(), backupExt)
		backups = append(backups, BackupInfo{
			Name:    name,
			Path:    filepath.Join(base, entry.Name()),
			ModTime: info.ModTime(),
		})
	}
	return backups, nil
}

// Restore reads a backup and returns the managed block.
func Restore(name string) (*parser.ManagedBlock, error) {
	base, err := backupBasePathFn()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(base, name+backupExt)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("backup %q not found: %w", name, err)
	}

	content := string(data)
	// strip block markers if present
	content = strings.TrimPrefix(content, parser.BlockStart+"\n")
	content = strings.TrimSuffix(content, "\n"+parser.BlockEnd+"\n")
	content = strings.TrimSuffix(content, "\n"+parser.BlockEnd)

	return parser.ParseBlock(content), nil
}

// CreateRaw saves raw file content as a backup (preserving original formatting).
func CreateRaw(name string, data []byte) error {
	base, err := backupBasePathFn()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(base, 0755); err != nil {
		return fmt.Errorf("cannot create backup directory: %w", err)
	}

	path := filepath.Join(base, name+backupExt)
	return os.WriteFile(path, data, 0644)
}

// IsRaw checks if a backup is a raw file (no managed block markers).
func IsRaw(name string) (bool, error) {
	base, err := backupBasePathFn()
	if err != nil {
		return false, err
	}

	path := filepath.Join(base, name+backupExt)
	data, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}

	return !strings.Contains(string(data), parser.BlockStart), nil
}

// RestoreRaw reads a raw backup and returns the raw content.
func RestoreRaw(name string) ([]byte, error) {
	base, err := backupBasePathFn()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(base, name+backupExt)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("backup %q not found: %w", name, err)
	}
	return data, nil
}

// Delete removes a backup file.
func Delete(name string) error {
	base, err := backupBasePathFn()
	if err != nil {
		return err
	}

	path := filepath.Join(base, name+backupExt)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("backup %q not found", name)
	}
	return os.Remove(path)
}
