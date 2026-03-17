package deploy

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Deployer handles installing plugin assets to a target .claude/ directory.
type Deployer struct {
	SourceFS fs.FS
	Plugin   string
	Version  string
}

// Asset represents a file to be deployed.
type Asset struct {
	// RelPath is relative to the .claude/ directory (e.g., "commands/research.md").
	RelPath string
	Content []byte
}

// ListAssets returns all assets for the given plugin and category.
// category should be "commands", "agents", or "" for all.
func (d *Deployer) ListAssets(category string) ([]Asset, error) {
	var assets []Asset
	pluginRoot := filepath.Join("plugins", d.Plugin)

	walkDir := pluginRoot
	if category != "" {
		walkDir = filepath.Join(pluginRoot, category)
	}

	err := fs.WalkDir(d.SourceFS, walkDir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		if !strings.HasSuffix(entry.Name(), ".md") {
			return nil
		}
		// Skip plugin-level metadata files (only deploy commands/ and agents/ content)
		if entry.Name() == "PLUGIN.md" {
			return nil
		}

		content, err := fs.ReadFile(d.SourceFS, path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}

		// Convert path relative to plugin root into path relative to .claude/
		// e.g., "plugins/research-plan-implement-validate/commands/research.md" → "commands/research.md"
		relPath, err := filepath.Rel(pluginRoot, path)
		if err != nil {
			return fmt.Errorf("rel path: %w", err)
		}

		assets = append(assets, Asset{
			RelPath: relPath,
			Content: content,
		})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk %s: %w", walkDir, err)
	}

	return assets, nil
}

// Install deploys assets to the target directory and writes a manifest.
func (d *Deployer) Install(targetDir string, scope string, assets []Asset, dryRun bool) error {
	manifest := NewManifest(scope, d.Plugin, d.Version)

	for _, a := range assets {
		destPath := filepath.Join(targetDir, a.RelPath)

		if dryRun {
			fmt.Printf("  would write: %s\n", destPath)
			manifest.AddItem(a.RelPath, a.Content, d.Version)
			continue
		}

		dir := filepath.Dir(destPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("mkdir %s: %w", dir, err)
		}

		if err := os.WriteFile(destPath, a.Content, 0644); err != nil {
			return fmt.Errorf("write %s: %w", destPath, err)
		}

		manifest.AddItem(a.RelPath, a.Content, d.Version)
		fmt.Printf("  installed: %s\n", a.RelPath)
	}

	if dryRun {
		fmt.Printf("  would write: %s\n", filepath.Join(targetDir, ManifestFile))
		return nil
	}

	return manifest.Save(targetDir)
}

// Update re-deploys assets, respecting local modifications.
func (d *Deployer) Update(targetDir string, assets []Asset, force bool, showDiff bool) error {
	manifest, err := LoadManifest(targetDir)
	if err != nil {
		return fmt.Errorf("no existing installation found: %w", err)
	}

	updated := 0
	skipped := 0

	for _, a := range assets {
		item, exists := manifest.Items[a.RelPath]
		if exists {
			status := CheckFile(targetDir, a.RelPath, item, a.Content)
			switch status.Status {
			case "up-to-date":
				continue
			case "modified":
				if !force {
					fmt.Printf("  skipped (locally modified): %s\n", a.RelPath)
					skipped++
					continue
				}
				fmt.Printf("  overwriting (locally modified): %s\n", a.RelPath)
			case "outdated":
				fmt.Printf("  updating: %s\n", a.RelPath)
			case "missing":
				fmt.Printf("  restoring: %s\n", a.RelPath)
			}
		} else {
			fmt.Printf("  adding: %s\n", a.RelPath)
		}

		destPath := filepath.Join(targetDir, a.RelPath)
		dir := filepath.Dir(destPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("mkdir %s: %w", dir, err)
		}

		if err := os.WriteFile(destPath, a.Content, 0644); err != nil {
			return fmt.Errorf("write %s: %w", destPath, err)
		}

		manifest.AddItem(a.RelPath, a.Content, d.Version)
		updated++
	}

	if updated == 0 && skipped == 0 {
		fmt.Println("  everything up-to-date")
		return nil
	}

	if skipped > 0 {
		fmt.Printf("  %d file(s) skipped (use --force to overwrite)\n", skipped)
	}

	return manifest.Save(targetDir)
}

// Status checks all installed files against the manifest and source.
func (d *Deployer) Status(targetDir string, assets []Asset) ([]FileStatus, error) {
	manifest, err := LoadManifest(targetDir)
	if err != nil {
		return nil, fmt.Errorf("no existing installation found: %w", err)
	}

	// Build source content map
	sourceMap := make(map[string][]byte)
	for _, a := range assets {
		sourceMap[a.RelPath] = a.Content
	}

	var statuses []FileStatus
	for relPath, item := range manifest.Items {
		source, ok := sourceMap[relPath]
		if !ok {
			source = []byte{} // file removed from source
		}
		statuses = append(statuses, CheckFile(targetDir, relPath, item, source))
	}

	return statuses, nil
}

// Uninstall removes all files tracked by the manifest.
func (d *Deployer) Uninstall(targetDir string) error {
	manifest, err := LoadManifest(targetDir)
	if err != nil {
		return fmt.Errorf("no existing installation found: %w", err)
	}

	for relPath := range manifest.Items {
		filePath := filepath.Join(targetDir, relPath)
		if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
			fmt.Printf("  warning: could not remove %s: %v\n", relPath, err)
			continue
		}
		fmt.Printf("  removed: %s\n", relPath)
	}

	// Remove manifest
	manifestPath := filepath.Join(targetDir, ManifestFile)
	if err := os.Remove(manifestPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove manifest: %w", err)
	}
	fmt.Println("  removed: manifest")

	// Try to clean up empty directories
	for _, subdir := range []string{"commands", "agents"} {
		dirPath := filepath.Join(targetDir, subdir)
		entries, err := os.ReadDir(dirPath)
		if err == nil && len(entries) == 0 {
			os.Remove(dirPath)
		}
	}

	return nil
}
