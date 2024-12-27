package commands

import (
	"fmt"
	"gopract/objects"
	"gopract/repository"
)

func handleLsTree(args []string) {
	// Parse arguments
	if len(args) < 1 {
		fmt.Println("Usage: ls-tree <tree-sha>")
		return
	}

	treeSha := args[0]
	repo, err := repository.Find(".", false)
	if err != nil {
		fmt.Println("Error finding repository:", err)
		return
	}

	// Read the tree object
	tree, err := objects.ReadObject(repo.Worktree, treeSha)
	if err != nil {
		fmt.Println("Error reading tree:", err)
		return
	}

	// Print tree contents
	for _, entry := range tree.(*objects.Tree).Entries {
		fmt.Printf("Mode: %s, Name: %s, SHA: %s\n", entry.Mode, entry.Name, entry.Hash)
	}
}
