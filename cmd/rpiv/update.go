package main

import (
	"fmt"

	"github.com/jarrodchung/rpiv/internal/assets"
	"github.com/jarrodchung/rpiv/internal/deploy"
	"github.com/spf13/cobra"
)

func updateCmd() *cobra.Command {
	var (
		scope  string
		target string
		diff   bool
		force  bool
	)

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update deployed plugin to latest embedded version",
		RunE: func(_ *cobra.Command, _ []string) error {
			targetDir, err := deploy.ResolveTarget(scope, target)
			if err != nil {
				return err
			}

			d := &deploy.Deployer{
				SourceFS: assets.Plugins,
				Plugin:   "rpiv",
				Version:  version,
			}

			assetList, err := d.ListAssets("")
			if err != nil {
				return fmt.Errorf("list assets: %w", err)
			}

			fmt.Printf("Updating %s\n", targetDir)
			return d.Update(targetDir, assetList, force, diff)
		},
	}

	cmd.Flags().StringVar(&scope, "scope", "project", "deployment scope (user or project)")
	cmd.Flags().StringVar(&target, "target", "", "target project root (default: current directory)")
	cmd.Flags().BoolVar(&diff, "diff", false, "show changes before applying")
	cmd.Flags().BoolVar(&force, "force", false, "overwrite locally modified files")

	return cmd
}
