package deploy

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const ManifestFile = ".rpiv-manifest.json"

type ManifestItem struct {
	SHA256        string `json:"sha256"`
	SourceVersion string `json:"source_version"`
}

type Manifest struct {
	Version     string                  `json:"version"`
	InstalledAt string                  `json:"installed_at"`
	UpdatedAt   string                  `json:"updated_at"`
	Scope       string                  `json:"scope"`
	Plugin      string                  `json:"plugin"`
	Items       map[string]ManifestItem `json:"items"`
}

func NewManifest(scope, plugin, version string) *Manifest {
	now := time.Now().UTC().Format(time.RFC3339)
	return &Manifest{
		Version:     version,
		InstalledAt: now,
		UpdatedAt:   now,
		Scope:       scope,
		Plugin:      plugin,
		Items:       make(map[string]ManifestItem),
	}
}

func LoadManifest(dir string) (*Manifest, error) {
	path := filepath.Join(dir, ManifestFile)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read manifest: %w", err)
	}

	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parse manifest: %w", err)
	}
	return &m, nil
}

func (m *Manifest) Save(dir string) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal manifest: %w", err)
	}

	path := filepath.Join(dir, ManifestFile)
	if err := os.WriteFile(path, append(data, '\n'), 0644); err != nil {
		return fmt.Errorf("write manifest: %w", err)
	}
	return nil
}

func (m *Manifest) AddItem(relPath string, content []byte, version string) {
	m.Items[relPath] = ManifestItem{
		SHA256:        hashContent(content),
		SourceVersion: version,
	}
	m.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
}

func hashContent(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}

// FileStatus represents the state of a deployed file.
type FileStatus struct {
	RelPath       string
	Status        string // "up-to-date", "outdated", "modified", "missing"
	InstalledHash string
	CurrentHash   string
	SourceHash    string
}

// CheckFile compares a manifest item against the filesystem and source.
func CheckFile(targetDir, relPath string, item ManifestItem, sourceContent []byte) FileStatus {
	fs := FileStatus{
		RelPath:       relPath,
		InstalledHash: item.SHA256,
		SourceHash:    hashContent(sourceContent),
	}

	filePath := filepath.Join(targetDir, relPath)
	data, err := os.ReadFile(filePath)
	if err != nil {
		fs.Status = "missing"
		return fs
	}

	fs.CurrentHash = hashContent(data)

	if fs.CurrentHash != fs.InstalledHash {
		fs.Status = "modified"
	} else if fs.InstalledHash != fs.SourceHash {
		fs.Status = "outdated"
	} else {
		fs.Status = "up-to-date"
	}

	return fs
}
