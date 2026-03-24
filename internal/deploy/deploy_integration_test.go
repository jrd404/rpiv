package deploy_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jarrodchung/rpiv/internal/assets"
	"github.com/jarrodchung/rpiv/internal/deploy"
)

// TestInstallToFakeProject simulates installing rpiv into a fresh project and
// validates the resulting directory structure, manifest, and file contents.
func TestInstallToFakeProject(t *testing.T) {
	// Create a temporary directory to act as a fake project's .claude/ dir
	tmpDir := t.TempDir()
	targetDir := filepath.Join(tmpDir, ".claude")

	d := &deploy.Deployer{
		SourceFS: assets.Plugins,
		Plugin:   "rpiv",
		Version:  "test",
	}

	// List all assets
	allAssets, err := d.ListAssets("")
	if err != nil {
		t.Fatalf("ListAssets: %v", err)
	}
	if len(allAssets) == 0 {
		t.Fatal("ListAssets returned no assets")
	}

	// Install
	if err := d.Install(targetDir, "project", allAssets, false); err != nil {
		t.Fatalf("Install: %v", err)
	}

	// --- Validate directory structure ---

	expectedSkills := []string{
		"rpiv.commit",
		"rpiv.debug",
		"rpiv.handoff",
		"rpiv.implement",
		"rpiv.init",
		"rpiv.iterate",
		"rpiv.oneshot",
		"rpiv.plan",
		"rpiv.pr",
		"rpiv.research",
		"rpiv.resume",
		"rpiv.status",
		"rpiv.validate",
	}

	for _, skill := range expectedSkills {
		skillFile := filepath.Join(targetDir, "skills", skill, "SKILL.md")
		if _, err := os.Stat(skillFile); os.IsNotExist(err) {
			t.Errorf("missing skill file: %s", skillFile)
		}
	}

	expectedAgents := []string{
		"codebase-analyzer.md",
		"codebase-locator.md",
		"codebase-pattern-finder.md",
		"web-researcher.md",
	}

	for _, agent := range expectedAgents {
		agentFile := filepath.Join(targetDir, "agents", agent)
		if _, err := os.Stat(agentFile); os.IsNotExist(err) {
			t.Errorf("missing agent file: %s", agentFile)
		}
	}

	// Hooks
	hooksFile := filepath.Join(targetDir, "hooks", "hooks.json")
	if _, err := os.Stat(hooksFile); os.IsNotExist(err) {
		t.Error("missing hooks/hooks.json")
	}

	// Scripts
	expectedScripts := []string{"update-tracker.sh", "log-commit.sh"}
	for _, script := range expectedScripts {
		scriptFile := filepath.Join(targetDir, "scripts", script)
		if _, err := os.Stat(scriptFile); os.IsNotExist(err) {
			t.Errorf("missing script: %s", scriptFile)
		}
	}

	// --- Validate manifest ---

	manifestPath := filepath.Join(targetDir, deploy.ManifestFile)
	manifestData, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}

	var manifest struct {
		Version     string                     `json:"version"`
		Scope       string                     `json:"scope"`
		Plugin      string                     `json:"plugin"`
		InstalledAt string                     `json:"installed_at"`
		Items       map[string]json.RawMessage `json:"items"`
	}
	if err := json.Unmarshal(manifestData, &manifest); err != nil {
		t.Fatalf("parse manifest: %v", err)
	}

	if manifest.Version != "test" {
		t.Errorf("manifest version = %q, want %q", manifest.Version, "test")
	}
	if manifest.Scope != "project" {
		t.Errorf("manifest scope = %q, want %q", manifest.Scope, "project")
	}
	if manifest.Plugin != "rpiv" {
		t.Errorf("manifest plugin = %q, want %q", manifest.Plugin, "rpiv")
	}
	if len(manifest.Items) != len(allAssets) {
		t.Errorf("manifest has %d items, want %d", len(manifest.Items), len(allAssets))
	}

	// --- Validate SKILL.md frontmatter names match directory names ---

	for _, skill := range expectedSkills {
		skillFile := filepath.Join(targetDir, "skills", skill, "SKILL.md")
		content, err := os.ReadFile(skillFile)
		if err != nil {
			continue // already reported above
		}

		// Check that name: field in frontmatter matches directory name
		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "name: ") {
				name := strings.TrimPrefix(line, "name: ")
				if name != skill {
					t.Errorf("skill %s: frontmatter name = %q, want %q", skill, name, skill)
				}
				break
			}
		}
	}

	// --- Validate no stale /rpiv: references in skill files ---

	for _, skill := range expectedSkills {
		skillFile := filepath.Join(targetDir, "skills", skill, "SKILL.md")
		content, err := os.ReadFile(skillFile)
		if err != nil {
			continue
		}
		if strings.Contains(string(content), "/rpiv:") {
			t.Errorf("skill %s: contains stale /rpiv: reference", skill)
		}
	}

	// --- Validate status check works ---

	statuses, err := d.Status(targetDir, allAssets)
	if err != nil {
		t.Fatalf("Status: %v", err)
	}

	for _, s := range statuses {
		if s.Status != "up-to-date" {
			t.Errorf("file %s: status = %q after fresh install, want %q", s.RelPath, s.Status, "up-to-date")
		}
	}

	// --- Validate uninstall cleans up ---

	if err := d.Uninstall(targetDir); err != nil {
		t.Fatalf("Uninstall: %v", err)
	}

	// Manifest should be gone
	if _, err := os.Stat(manifestPath); !os.IsNotExist(err) {
		t.Error("manifest still exists after uninstall")
	}

	// Skill files should be gone
	for _, skill := range expectedSkills {
		skillFile := filepath.Join(targetDir, "skills", skill, "SKILL.md")
		if _, err := os.Stat(skillFile); !os.IsNotExist(err) {
			t.Errorf("skill file still exists after uninstall: %s", skillFile)
		}
	}
}
