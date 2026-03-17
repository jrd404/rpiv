package main

import (
	"fmt"

	"github.com/jarrodchung/cc/internal/assets"
	"github.com/jarrodchung/cc/internal/deploy"
	"github.com/spf13/cobra"
)

func installCmd() *cobra.Command {
	var (
		scope  string
		target string
		dryRun bool
		plugin string
	)

	cmd := &cobra.Command{
		Use:       "install [commands|agents|all]",
		Short:     "Deploy plugin commands and agents to a target project",
		Args:      cobra.MaximumNArgs(1),
		ValidArgs: []string{"commands", "agents", "all"},
		RunE: func(_ *cobra.Command, args []string) error {
			category := ""
			if len(args) > 0 && args[0] != "all" {
				category = args[0]
			}

			if plugin == "" {
				plugin = "research-plan-implement-validate"
			}

			targetDir, err := deploy.ResolveTarget(scope, target)
			if err != nil {
				return err
			}

			d := &deploy.Deployer{
				SourceFS: assets.Plugins,
				Plugin:   plugin,
				Version:  version,
			}

			assetList, err := d.ListAssets(category)
			if err != nil {
				return fmt.Errorf("list assets: %w", err)
			}

			if len(assetList) == 0 {
				return fmt.Errorf("no assets found for plugin %q category %q", plugin, category)
			}

			action := "Installing"
			if dryRun {
				action = "Would install"
			}
			fmt.Printf("%s %d files to %s\n", action, len(assetList), targetDir)

			return d.Install(targetDir, scope, assetList, dryRun)
		},
	}

	cmd.Flags().StringVar(&scope, "scope", "project", "deployment scope (user or project)")
	cmd.Flags().StringVar(&target, "target", "", "target project root (default: current directory)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "show what would be installed without writing")
	cmd.Flags().StringVar(&plugin, "plugin", "research-plan-implement-validate", "plugin to install")

	return cmd
}
