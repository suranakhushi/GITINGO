package commands

import (
	"fmt"
	"gopract/objects"
	"os"
)

func HashObject(repoPath, filePath string, write bool) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	blob := &objects.Blob{Data: data}
	sha, err := objects.WriteObject(blob, repoPath)
	if err != nil {
		return fmt.Errorf("failed to write object: %w", err)
	}

	if !write {
		fmt.Println(sha)
	}
	return nil
}
