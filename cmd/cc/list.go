package main

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/jarrodchung/cc/internal/assets"
	"github.com/spf13/cobra"
)

func listCmd() *cobra.Command {
	return &cobra.Command{
		Use:       "list [commands|agents|plugins]",
		Short:     "List available plugin assets",
		Args:      cobra.MaximumNArgs(1),
		ValidArgs: []string{"commands", "agents", "plugins"},
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
	plugin := "research-plan-implement-validate"
	pluginRoot := filepath.Join("plugins", plugin)

	walkRoot := pluginRoot
	if category != "" {
		walkRoot = filepath.Join(pluginRoot, category)
	}

	var commands []string
	var agents []string

	err := fs.WalkDir(assets.Plugins, walkRoot, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			return nil
		}

		relPath, _ := filepath.Rel(pluginRoot, path)
		if strings.HasPrefix(relPath, "commands/") {
			name := strings.TrimSuffix(strings.TrimPrefix(relPath, "commands/"), ".md")
			commands = append(commands, name)
		} else if strings.HasPrefix(relPath, "agents/") {
			name := strings.TrimSuffix(strings.TrimPrefix(relPath, "agents/"), ".md")
			agents = append(agents, name)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("walk assets: %w", err)
	}

	if (category == "" || category == "commands") && len(commands) > 0 {
		fmt.Println("Commands:")
		for _, c := range commands {
			fmt.Printf("  /%s\n", c)
		}
	}

	if (category == "" || category == "agents") && len(agents) > 0 {
		if category == "" && len(commands) > 0 {
			fmt.Println()
		}
		fmt.Println("Agents:")
		for _, a := range agents {
			fmt.Printf("  %s\n", a)
		}
	}

	return nil
}
