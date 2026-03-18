package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var (
	completionCommands    = []string{"install", "update", "status", "list", "uninstall", "version", "completion"}
	completionGlobalFlags = []string{"--help"}
)

func completionCmd(root *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:   "completion [bash|zsh]",
		Short: "Generate shell completion scripts",
		Long: `Generate shell completion scripts for bash or zsh.

To load completions:

Bash (recommended):
  $ complete -C 'rpiv __completer' rpiv

Bash (traditional):
  $ source <(rpiv completion bash)

Zsh:
  $ source <(rpiv completion zsh)
  # To load completions for each session:
  $ rpiv completion zsh > ~/.zsh/completions/_rpiv
`,
		Args:      cobra.ExactArgs(1),
		ValidArgs: []string{"bash", "zsh"},
		RunE: func(_ *cobra.Command, args []string) error {
			switch args[0] {
			case "bash":
				return root.GenBashCompletion(os.Stdout)
			case "zsh":
				return root.GenZshCompletion(os.Stdout)
			default:
				return fmt.Errorf("unsupported shell: %s", args[0])
			}
		},
	}
}

func completerCmd() *cobra.Command {
	return &cobra.Command{
		Use:                "__completer",
		Hidden:             true,
		DisableFlagParsing: true,
		Run: func(_ *cobra.Command, _ []string) {
			line := os.Getenv("COMP_LINE")
			pointStr := os.Getenv("COMP_POINT")

			point, _ := strconv.Atoi(pointStr)
			if point > len(line) {
				point = len(line)
			}

			lineUpToCursor := line[:point]
			words := parseCompletionWords(lineUpToCursor)
			completions := getCompletions(words)

			for _, c := range completions {
				fmt.Println(c)
			}
		},
	}
}

func parseCompletionWords(line string) []string {
	var words []string
	var current strings.Builder
	inQuote := false
	quoteChar := rune(0)

	for _, r := range line {
		switch {
		case r == '"' || r == '\'':
			if inQuote && r == quoteChar {
				inQuote = false
			} else if !inQuote {
				inQuote = true
				quoteChar = r
			}
			current.WriteRune(r)
		case r == ' ' && !inQuote:
			if current.Len() > 0 {
				words = append(words, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(r)
		}
	}

	words = append(words, current.String())
	return words
}

func getCompletions(words []string) []string {
	subcommandFlags := map[string][]string{
		"install":    {"--scope", "--target", "--dry-run", "--plugin", "--help"},
		"update":     {"--scope", "--target", "--diff", "--force", "--help"},
		"status":     {"--scope", "--target", "--help"},
		"list":       {"--help"},
		"uninstall":  {"--scope", "--target", "--help"},
		"version":    {"--help"},
		"completion": {"--help"},
	}

	installArgs := []string{"skills", "agents", "all"}
	listArgs := []string{"skills", "agents", "plugins"}
	scopes := []string{"user", "project"}
	shells := []string{"bash", "zsh"}
	plugins := []string{"rpiv"}

	if len(words) <= 1 {
		return append(completionCommands, completionGlobalFlags...)
	}

	current := words[len(words)-1]
	prev := ""
	if len(words) >= 2 {
		prev = words[len(words)-2]
	}

	// Flag value completions
	switch prev {
	case "--scope":
		return filterPrefix(scopes, current)
	case "--target":
		return nil // directory completion handled by shell
	case "--plugin":
		return filterPrefix(plugins, current)
	}

	// Flag name completions
	if strings.HasPrefix(current, "-") {
		for _, w := range words[1:] {
			if !strings.HasPrefix(w, "-") {
				if flags, ok := subcommandFlags[w]; ok {
					return filterPrefix(flags, current)
				}
			}
		}
		return filterPrefix(completionGlobalFlags, current)
	}

	// Detect subcommand
	subCmd := ""
	for _, w := range words[1 : len(words)-1] {
		if !strings.HasPrefix(w, "-") {
			for _, cmd := range completionCommands {
				if w == cmd {
					subCmd = w
					break
				}
			}
			if subCmd != "" {
				break
			}
		}
	}

	if subCmd == "" {
		return filterPrefix(completionCommands, current)
	}

	// Positional argument completions
	switch subCmd {
	case "install":
		posArgs := countPositionalArgs(words, subCmd)
		if posArgs == 0 {
			return filterPrefix(installArgs, current)
		}
		return nil

	case "list":
		posArgs := countPositionalArgs(words, subCmd)
		if posArgs == 0 {
			return filterPrefix(listArgs, current)
		}
		return nil

	case "completion":
		posArgs := countPositionalArgs(words, subCmd)
		if posArgs == 0 {
			return filterPrefix(shells, current)
		}
		return nil

	default:
		if current == "" {
			return subcommandFlags[subCmd]
		}
		return nil
	}
}

func countPositionalArgs(words []string, subCmd string) int {
	count := 0
	foundCmd := false
	skipNext := false
	for _, w := range words[:len(words)-1] {
		if skipNext {
			skipNext = false
			continue
		}
		if foundCmd {
			if strings.HasPrefix(w, "-") {
				// Skip flags that take values
				if w == "--scope" || w == "--target" || w == "--plugin" {
					skipNext = true
				}
				continue
			}
			count++
		}
		if w == subCmd {
			foundCmd = true
		}
	}
	return count
}

func filterPrefix(items []string, prefix string) []string {
	if prefix == "" {
		return items
	}
	var result []string
	for _, item := range items {
		if strings.HasPrefix(item, prefix) {
			result = append(result, item)
		}
	}
	return result
}
