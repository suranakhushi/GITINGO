package commands

import (
	"fmt"
	"gopract/objects"
	"os"
	"path/filepath"
)

// Checkout switches the working directory to the state of a commit by hash.
func Checkout(repoPath, commitHashStr string) error {
	gitDir := filepath.Join(repoPath, ".git")

	// Ensure the repository exists
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		return fmt.Errorf("not a Git repository: %s", repoPath)
	}

	// Check if the commit hash is valid
	if len(commitHashStr) != 40 {
		return fmt.Errorf("invalid commit hash format: %s", commitHashStr)
	}

	// Traverse the commit history and check out the commit tree
	for commitHashStr != "" {
		objectPath := filepath.Join(gitDir, "objects", commitHashStr[:2], commitHashStr[2:])
		fmt.Printf("Reading commit object from: %s\n", objectPath)

		commit, err := objects.ReadObject(repoPath, commitHashStr)
		if err != nil {
			return fmt.Errorf("failed to read commit object: %w", err)
		}

		commitObj := commit.(*objects.Commit)
		fmt.Printf("Checking out commit %s\n", commitHashStr)

		// Retrieve the tree hash from the commit object
		treeHash := commitObj.Tree

		// Read the tree object from the repository
		treeObj, err := objects.ReadObject(repoPath, treeHash)
		if err != nil {
			return fmt.Errorf("failed to read tree object: %w", err)
		}

		// Check out the tree into the working directory
		if err := checkoutTree(repoPath, treeObj.(*objects.Tree)); err != nil {
			return fmt.Errorf("failed to check out tree: %w", err)
		}

		// Move to the parent commit(s)
		if len(commitObj.Parents) > 0 {
			commitHashStr = commitObj.Parents[0]
		} else {
			commitHashStr = ""
		}
	}

	return nil
}

// checkoutTree writes files from the tree object into the working directory.
func checkoutTree(repoPath string, tree *objects.Tree) error {
	for _, entry := range tree.Entries {
		destPath := filepath.Join(repoPath, entry.Name)

		// If it's a blob, write the file to the working directory
		obj, err := objects.ReadObject(repoPath, entry.Hash)
		if err != nil {
			return fmt.Errorf("failed to read object: %w", err)
		}

		switch obj := obj.(type) {
		case *objects.Blob:
			// Write the file content into the working directory
			if err := os.WriteFile(destPath, obj.Data, 0644); err != nil {
				return fmt.Errorf("failed to write file: %w", err)
			}
		case *objects.Tree:
			// If it's a directory, create the directory and recurse
			if err := os.MkdirAll(destPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
			if err := checkoutTree(repoPath, obj); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown object type: %T", obj)
		}
	}

	return nil
}
