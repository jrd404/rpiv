package main

import (
	"fmt"

	"github.com/jarrodchung/rpiv/internal/assets"
	"github.com/jarrodchung/rpiv/internal/deploy"
	"github.com/spf13/cobra"
)

func uninstallCmd() *cobra.Command {
	var (
		scope  string
		target string
	)

	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Remove deployed plugin skills and agents",
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

			fmt.Printf("Uninstalling from %s\n", targetDir)
			return d.Uninstall(targetDir)
		},
	}

	cmd.Flags().StringVar(&scope, "scope", "project", "deployment scope (user or project)")
	cmd.Flags().StringVar(&target, "target", "", "target project root (default: current directory)")

	return cmd
}
