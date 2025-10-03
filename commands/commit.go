package commands

import (
	"fmt"
	"gopract/objects"
	"gopract/staging"
	"os"
	"time"
)

// Commit creates a new commit object and updates the repository state.
func Commit(repoPath, message string) error {
	// Ensure the repository exists
	gitDir := fmt.Sprintf("%s/.git", repoPath)
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		return fmt.Errorf("not a Git repository: %s", repoPath)
	}

	// Read the staging area (index)
	index, err := staging.ReadIndex(repoPath)
	if err != nil {
		return fmt.Errorf("failed to read index: %w", err)
	}

	// Create a tree object from the index
	tree := &objects.Tree{}
	for filePath, blobHash := range index {
		tree.Entries = append(tree.Entries, objects.TreeEntry{
			Mode: "100644", // Regular file mode
			Hash: blobHash,
			Name: filePath,
		})
	}

	// Write the tree object to `.git/objects`
	treeHash, err := objects.WriteObject(tree, repoPath)
	if err != nil {
		return fmt.Errorf("failed to write tree object: %w", err)
	}

	// Determine the parent commit (if any)
	var parents []string
	headPath := fmt.Sprintf("%s/HEAD", gitDir)
	head, err := os.ReadFile(headPath)
	if err == nil {
		// Read the current HEAD, if it exists
		ref := string(head)
		refPath := fmt.Sprintf("%s/%s", gitDir, ref)
		if parentHash, err := os.ReadFile(refPath); err == nil {
			parents = append(parents, string(parentHash))
		}
	}

	// Create a new commit object
	commit := &objects.Commit{
		Tree:    treeHash,
		Parents: parents,
		Author:  fmt.Sprintf("Your Name <your.email@example.com> %d +0000", time.Now().Unix()), // Example author
		Message: message,
	}

	// Write the commit object to `.git/objects`
	commitHash, err := objects.WriteObject(commit, repoPath)
	if err != nil {
		return fmt.Errorf("failed to write commit object: %w", err)
	}

	// Update the HEAD to point to the new commit
	refPath := fmt.Sprintf("%s/refs/heads/master", gitDir)
	if err := os.WriteFile(refPath, []byte(commitHash), 0644); err != nil {
		return fmt.Errorf("failed to update HEAD: %w", err)
	}

	fmt.Printf("Committed with hash %s\n", commitHash)
	return nil
}
