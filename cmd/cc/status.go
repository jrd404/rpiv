package main

import (
	"fmt"

	"github.com/jarrodchung/cc/internal/assets"
	"github.com/jarrodchung/cc/internal/deploy"
	"github.com/spf13/cobra"
)

func statusCmd() *cobra.Command {
	var (
		scope  string
		target string
	)

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show deployment status and detect drift",
		RunE: func(_ *cobra.Command, _ []string) error {
			targetDir, err := deploy.ResolveTarget(scope, target)
			if err != nil {
				return err
			}

			d := &deploy.Deployer{
				SourceFS: assets.Plugins,
				Plugin:   "research-plan-implement-validate",
				Version:  version,
			}

			assetList, err := d.ListAssets("")
			if err != nil {
				return fmt.Errorf("list assets: %w", err)
			}

			statuses, err := d.Status(targetDir, assetList)
			if err != nil {
				return err
			}

			fmt.Printf("Status for %s\n\n", targetDir)

			counts := map[string]int{}
			for _, s := range statuses {
				counts[s.Status]++
				icon := "✓"
				switch s.Status {
				case "modified":
					icon = "~"
				case "outdated":
					icon = "↑"
				case "missing":
					icon = "✗"
				}
				fmt.Printf("  %s %s (%s)\n", icon, s.RelPath, s.Status)
			}

			fmt.Println()
			if counts["up-to-date"] > 0 {
				fmt.Printf("  %d up-to-date\n", counts["up-to-date"])
			}
			if counts["outdated"] > 0 {
				fmt.Printf("  %d outdated (run 'cc update' to update)\n", counts["outdated"])
			}
			if counts["modified"] > 0 {
				fmt.Printf("  %d locally modified\n", counts["modified"])
			}
			if counts["missing"] > 0 {
				fmt.Printf("  %d missing (run 'cc update' to restore)\n", counts["missing"])
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&scope, "scope", "project", "deployment scope (user or project)")
	cmd.Flags().StringVar(&target, "target", "", "target project root (default: current directory)")

	return cmd
}
