package commands

import (
	"fmt"
	"gopract/objects"
	"os"
	"path/filepath"
)

// CatFile retrieves and displays the raw content of a Git object.
func CatFile(repoPath, sha string) error {
	// Ensure the repository exists
	gitDir := filepath.Join(repoPath, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		return fmt.Errorf("not a Git repository: %s", repoPath)
	}

	// Read the object from the `.git/objects` directory
	obj, err := objects.ReadObject(repoPath, sha)
	if err != nil {
		return fmt.Errorf("failed to read object %s: %w", sha, err)
	}

	// Print the object content based on its type
	switch obj.Type() {
	case "blob":
		fmt.Print(string(obj.(*objects.Blob).Data))
	case "tree":
		for _, entry := range obj.(*objects.Tree).Entries {
			fmt.Printf("%s %s %s\n", entry.Mode, entry.Hash, entry.Name)
		}
	case "commit":
		commit := obj.(*objects.Commit)
		fmt.Printf("tree %s\n", commit.Tree)
		for _, parent := range commit.Parents {
			fmt.Printf("parent %s\n", parent)
		}
		fmt.Printf("author %s\n", commit.Author)
		fmt.Printf("\n%s\n", commit.Message)
	default:
		return fmt.Errorf("unknown object type: %s", obj.Type())
	}

	return nil
}
