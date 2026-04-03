package utils

import (
	"os"
	"path/filepath"
	"strings"
)

func DefaultAppDirs() []string {
	home, _ := os.UserHomeDir()

	return append([]string{
		// user local apps
		os.Getenv("XDG_DATA_HOME"),
		filepath.Join(home, ".local/share"),

		// flatpak
		filepath.Join(home, ".local/share/flatpak/exports/share"),
		"/var/lib/flatpak/exports/share",

		// system
		"/usr/local/share",
		"/usr/share",

		// xdg
	}, strings.Split(os.Getenv("XDG_DATA_DIRS"), ":")...)
}

func ProcessDirectories(paths []string) []string {
	uniquePaths := make(map[string]bool)
	var result []string

	for _, path := range paths {
		if path == "" {
			continue
		}

		path = filepath.Clean(path)
		path = filepath.Join(path, "applications")

		if !filepath.IsAbs(path) {
			continue
		}

		fileInfo, err := os.Stat(path)
		if err != nil {
			continue
		}

		if !fileInfo.IsDir() {
			continue
		}

		if !uniquePaths[path] {
			uniquePaths[path] = true
			result = append(result, path)
		}
	}

	return result
}
