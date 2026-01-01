package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// FindProjectRoot finds the root directory of the git repository
// by walking up the directory tree looking for a .git directory
func FindProjectRoot(startPath string) (string, error) {
	if startPath == "" {
		var err error
		startPath, err = os.Getwd()
		if err != nil {
			return "", err
		}
	}

	currentPath := startPath
	for {
		gitPath := filepath.Join(currentPath, ".git")
		if _, err := os.Stat(gitPath); err == nil {
			return currentPath, nil
		}

		parentPath := filepath.Dir(currentPath)
		if parentPath == currentPath {
			// Reached filesystem root without finding .git
			return "", os.ErrNotExist
		}
		currentPath = parentPath
	}
}

// IsGitRepository checks if the given directory is within a git repository
func IsGitRepository(path string) bool {
	_, err := FindProjectRoot(path)
	return err == nil
}

// GetDkitDataDir returns the .dkit data directory path for the project
// If projectRoot is empty, it will attempt to find it automatically
func GetDkitDataDir(projectRoot string) (string, error) {
	if projectRoot == "" {
		var err error
		projectRoot, err = FindProjectRoot("")
		if err != nil {
			// If not in git repo, use current directory
			projectRoot, err = os.Getwd()
			if err != nil {
				return "", err
			}
		}
	}

	return filepath.Join(projectRoot, ".dkit"), nil
}

// EnsureDkitDataDir creates the .dkit directory if it doesn't exist
func EnsureDkitDataDir(projectRoot string) (string, error) {
	dataDir, err := GetDkitDataDir(projectRoot)
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return "", err
	}

	return dataDir, nil
}

// SanitizePathForFilename converts a path or command into a safe filename
func SanitizePathForFilename(path string) string {
	// Replace path separators and other unsafe characters
	replacer := strings.NewReplacer(
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
		" ", "_",
	)
	return replacer.Replace(path)
}
