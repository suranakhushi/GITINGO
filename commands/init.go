package commands

import (
	"fmt"
	"os"
	"path/filepath"
)

// Init initializes a new Git repository.
func Init(repoPath string) error {
	// Ensure the target path exists
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		return fmt.Errorf("failed to create repository path %s: %w", repoPath, err)
	}

	// Create the .git directory
	gitDir := filepath.Join(repoPath, ".git")
	if _, err := os.Stat(gitDir); !os.IsNotExist(err) {
		return fmt.Errorf("repository already exists at %s", gitDir)
	}
	if err := os.MkdirAll(gitDir, 0755); err != nil {
		return fmt.Errorf("failed to create .git directory: %w", err)
	}

	// Create required subdirectories
	subDirs := []string{"objects", "refs/heads", "refs/tags"}
	for _, dir := range subDirs {
		path := filepath.Join(gitDir, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", path, err)
		}
	}

	// Create default files
	if err := createDefaultFiles(gitDir); err != nil {
		return fmt.Errorf("failed to create default files: %w", err)
	}

	fmt.Printf("Initialized empty Git repository in %s\n", repoPath)
	return nil
}

// createDefaultFiles creates the default files required in a new repository.
func createDefaultFiles(gitDir string) error {
	// Create HEAD file
	headPath := filepath.Join(gitDir, "HEAD")
	if err := os.WriteFile(headPath, []byte("ref: refs/heads/master\n"), 0644); err != nil {
		return fmt.Errorf("failed to write HEAD file: %w", err)
	}

	// Create an empty config file
	configPath := filepath.Join(gitDir, "config")
	configContent := `[core]
	repositoryformatversion = 0
	filemode = true
	bare = false
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	// Create description file
	descriptionPath := filepath.Join(gitDir, "description")
	descriptionContent := `Unnamed repository; edit this file to name the repository.
`
	if err := os.WriteFile(descriptionPath, []byte(descriptionContent), 0644); err != nil {
		return fmt.Errorf("failed to write description file: %w", err)
	}

	return nil
}
