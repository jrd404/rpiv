package main

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/jarrodchung/rpiv/internal/assets"
	"github.com/spf13/cobra"
)

func listCmd() *cobra.Command {
	return &cobra.Command{
		Use:       "list [skills|agents|plugins]",
		Short:     "List available plugin assets",
		Args:      cobra.MaximumNArgs(1),
		ValidArgs: []string{"skills", "agents", "plugins"},
		RunE: func(_ *cobra.Command, args []string) error {
			category := ""
			if len(args) > 0 {
				category = args[0]
			}

			if category == "plugins" {
				return listPlugins()
			}

			return listAssets(category)
		},
	}
}

func listPlugins() error {
	fmt.Println("Available plugins:")
	entries, err := fs.ReadDir(assets.Plugins, "plugins")
	if err != nil {
		return fmt.Errorf("read plugins: %w", err)
	}
	for _, e := range entries {
		if e.IsDir() {
			fmt.Printf("  %s\n", e.Name())
		}
	}
	return nil
}

func listAssets(category string) error {
	plugin := "rpiv"
	pluginRoot := filepath.Join("plugins", plugin)

	walkRoot := pluginRoot
	if category != "" {
		walkRoot = filepath.Join(pluginRoot, category)
	}

	var skills []string
	var agents []string

	err := fs.WalkDir(assets.Plugins, walkRoot, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Skip .claude-plugin metadata directory
		if entry.IsDir() && entry.Name() == ".claude-plugin" {
			return fs.SkipDir
		}
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			return nil
		}

		relPath, _ := filepath.Rel(pluginRoot, path)
		if strings.HasPrefix(relPath, "skills/") {
			// Extract skill name from skills/<name>/SKILL.md
			parts := strings.Split(relPath, string(filepath.Separator))
			if len(parts) >= 2 {
				name := parts[1]
				// Deduplicate (only add once per skill directory)
				found := false
				for _, s := range skills {
					if s == name {
						found = true
						break
					}
				}
				if !found {
					skills = append(skills, name)
				}
			}
		} else if strings.HasPrefix(relPath, "agents/") {
			name := strings.TrimSuffix(strings.TrimPrefix(relPath, "agents/"), ".md")
			agents = append(agents, name)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("walk assets: %w", err)
	}

	if (category == "" || category == "skills") && len(skills) > 0 {
		fmt.Println("Skills:")
		for _, s := range skills {
			fmt.Printf("  /rpiv:%s\n", s)
		}
	}

	if (category == "" || category == "agents") && len(agents) > 0 {
		if category == "" && len(skills) > 0 {
			fmt.Println()
		}
		fmt.Println("Agents:")
		for _, a := range agents {
			fmt.Printf("  %s\n", a)
		}
	}

	return nil
}
