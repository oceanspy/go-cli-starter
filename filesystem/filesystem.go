package filesystem

import (
	"os"
	"path/filepath"
	"strings"
)

// ExpandHome replaces a leading "~" (alone or followed by "/") with the
// current user's home directory. Other paths are returned unchanged.
func ExpandHome(path string) (string, error) {
	if path == "~" || strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		return filepath.Join(home, path[1:]), nil
	}

	return path, nil
}
