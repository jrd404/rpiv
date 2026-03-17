package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "dev"

func main() {
	rootCmd := &cobra.Command{
		Use:   "cc",
		Short: "Deploy and manage Claude Code plugins across projects",
		Long: `cc is a CLI for deploying the research-plan-implement-validate
framework and other Claude Code plugins to any project.

It embeds all plugin assets (commands, agents, skills) into a single
binary, making it easy to install on any machine via go install.`,
	}

	rootCmd.AddCommand(
		installCmd(),
		updateCmd(),
		statusCmd(),
		listCmd(),
		uninstallCmd(),
		versionCmd(),
		completionCmd(rootCmd),
		completerCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Println("cc", version)
		},
	}
}
