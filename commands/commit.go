package commands

import (
	"fmt"
	"gopract/objects"
	"gopract/repository"
	"gopract/staging"
	"os"
	"time"
)

// Commit creates a commit object and writes it to the repository.
func Commit(repoPath, message string) error {
	gitDir := fmt.Sprintf("%s/.git", repoPath)

	// Check if the directory exists
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

	// Write the tree object to .git/objects
	treeHash, err := objects.WriteObject(tree, repoPath)
	if err != nil {
		return fmt.Errorf("failed to write tree object: %w", err)
	}

	// Read the current HEAD, to determine parent commit (if any)
	var parents []string
	headPath := fmt.Sprintf("%s/HEAD", gitDir)
	head, err := os.ReadFile(headPath)
	if err == nil {
		// If HEAD exists and is a reference to a branch (like refs/heads/master)
		ref := string(head)
		refPath := fmt.Sprintf("%s/%s", gitDir, ref)
		parentHash, err := os.ReadFile(refPath)
		if err == nil {
			parents = append(parents, string(parentHash))
		}
	}

	// Create a commit object
	commit := &objects.Commit{
		Tree:    treeHash,
		Parents: parents,
		Author:  fmt.Sprintf("Your Name <your.email@example.com> %d +0000", time.Now().Unix()), // Example author
		Message: message,
	}

	// Write the commit object to .git/objects
	commitHash, err := objects.WriteObject(commit, repoPath)
	if err != nil {
		return fmt.Errorf("failed to write commit object: %w", err)
	}

	// Update the HEAD to point to the new commit
	refPath := fmt.Sprintf("%s/refs/heads/master", gitDir)
	err = os.WriteFile(refPath, []byte(commitHash), 0644)
	if err != nil {
		return fmt.Errorf("failed to update HEAD: %w", err)
	}

	fmt.Printf("Committed with hash %s\n", commitHash)
	return nil
}

// GetCommit retrieves a commit object by its hash.
func GetCommit(repo *repository.Repository, commitHash string) (*objects.Commit, error) {
	// Read the object from the .git/objects directory
	commitObj, err := objects.ReadObject(repo.Worktree, commitHash)
	if err != nil {
		return nil, fmt.Errorf("failed to read commit object: %w", err)
	}

	// Ensure that the object type is "commit"
	if commitObj.Type() != "commit" {
		return nil, fmt.Errorf("object is not a commit: %s", commitHash)
	}

	// Cast and return the commit object
	return commitObj.(*objects.Commit), nil
}

// Log shows the commit history (log).
