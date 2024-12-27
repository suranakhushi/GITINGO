package commands

import (
	"fmt"
	"gopract/objects"
	"os"
	"path/filepath"
	"strings"
)

// Log shows the commit history (log).
func Log(repoPath string) error {
	// Ensure the repository exists
	gitDir := filepath.Join(repoPath, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		return fmt.Errorf("not a Git repository: %s", repoPath)
	}

	// Read the HEAD to get the current branch
	headPath := filepath.Join(gitDir, "HEAD")
	head, err := os.ReadFile(headPath)
	if err != nil {
		return fmt.Errorf("failed to read HEAD: %w", err)
	}

	// Trim and check for 'ref: refs/heads/master' format in HEAD
	refLine := strings.TrimSpace(string(head))
	if !strings.HasPrefix(refLine, "ref: ") {
		return fmt.Errorf("unexpected format in HEAD file: %s", refLine)
	}

	// Extract branch reference (e.g., refs/heads/master)
	ref := strings.TrimPrefix(refLine, "ref: ")
	refPath := filepath.Join(gitDir, ref)

	// Read the commit hash from the reference file (e.g., refs/heads/master)
	commitHash, err := os.ReadFile(refPath)
	if err != nil {
		return fmt.Errorf("failed to read branch reference %s: %w", ref, err)
	}

	// Trim any leading/trailing spaces from the commit hash and sanitize
	commitHashStr := strings.TrimSpace(string(commitHash))
	// Remove any non-ASCII characters if necessary
	commitHashStr = strings.ReplaceAll(commitHashStr, "'", "")
	commitHashStr = strings.ReplaceAll(commitHashStr, "��", "")

	// Debugging: Show the commit hash after sanitization
	fmt.Printf("Commit Hash: '%s'\n", commitHashStr)

	// Check if the commit hash seems valid (length of 40 characters)
	if len(commitHashStr) != 40 {
		return fmt.Errorf("invalid commit hash format: %s", commitHashStr)
	}

	// Traverse commit history
	for commitHashStr != "" {
		// Debugging: Show the current object path being accessed
		objectPath := filepath.Join(gitDir, "objects", commitHashStr[:2], commitHashStr[2:])
		fmt.Printf("Reading commit object from: %s\n", objectPath)

		commit, err := objects.ReadObject(repoPath, commitHashStr)
		if err != nil {
			return fmt.Errorf("failed to read commit object: %w", err)
		}

		commitObj := commit.(*objects.Commit)
		fmt.Printf("commit %s\n", commitHashStr)
		fmt.Printf("Author: %s\n", commitObj.Author)
		// Removed Timestamp field here
		fmt.Printf("\n    %s\n\n", commitObj.Message)

		// Move to the parent commit(s)
		if len(commitObj.Parents) > 0 {
			commitHashStr = commitObj.Parents[0] // Follow the first parent for simplicity
		} else {
			commitHashStr = ""
		}
	}

	return nil
}
