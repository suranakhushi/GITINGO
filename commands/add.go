package commands

import (
	"fmt"
	"gopract/objects"
	"gopract/staging"
	"os"
	"path/filepath"
)

// Add adds a file to the staging area of the repository.
func Add(repoPath, filePath string) error {
	// Ensure the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file %s does not exist", filePath)
	}

	// Ensure the repository exists
	gitDir := filepath.Join(repoPath, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		return fmt.Errorf("not a Git repository: %s", repoPath)
	}

	// Read file contents
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	// Create a Blob object
	blob := &objects.Blob{Data: data}

	// Write the Blob object to `.git/objects`
	blobHash, err := objects.WriteObject(blob, repoPath)
	if err != nil {
		return fmt.Errorf("failed to write blob object for %s: %w", filePath, err)
	}

	// Update the index with the file path and blob hash
	err = staging.UpdateIndex(repoPath, filePath, blobHash)
	if err != nil {
		return fmt.Errorf("failed to update index: %w", err)
	}

	fmt.Printf("Added file %s to staging area\n", filePath)
	return nil
}
