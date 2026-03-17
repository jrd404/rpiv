package deploy

import (
	"fmt"
	"os"
	"path/filepath"
)

// ResolveTarget returns the .claude/ directory path based on scope and target.
func ResolveTarget(scope, target string) (string, error) {
	switch scope {
	case "user":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("get home directory: %w", err)
		}
		return filepath.Join(home, ".claude"), nil

	case "project":
		if target == "" {
			var err error
			target, err = os.Getwd()
			if err != nil {
				return "", fmt.Errorf("get working directory: %w", err)
			}
		}
		return filepath.Join(target, ".claude"), nil

	default:
		return "", fmt.Errorf("unknown scope: %s (expected 'user' or 'project')", scope)
	}
}
