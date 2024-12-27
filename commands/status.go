package commands

import (
	"crypto/sha1"
	"fmt"
	"gopract/staging"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Status(repoPath string) error {
	gitDir := filepath.Join(repoPath, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		return fmt.Errorf("not a Git repository: %s", repoPath)
	}

	headPath := filepath.Join(gitDir, "HEAD")
	headData, err := os.ReadFile(headPath)
	if err != nil {
		return fmt.Errorf("error reading HEAD: %w", err)
	}

	currentBranch := strings.TrimPrefix(string(headData), "ref: refs/heads/")
	fmt.Printf("On branch %s\n", currentBranch)

	index, err := staging.ReadIndex(repoPath)
	if err != nil {
		return fmt.Errorf("failed to read index: %w", err)
	}

	var modifiedFiles []string
	var untrackedFiles []string

	for filePath, _ := range index {
		modified, err := checkFileModified(repoPath, filePath, index)
		if err != nil {
			return err
		}
		if modified {
			modifiedFiles = append(modifiedFiles, filePath)
		}
	}

	untrackedFiles, err = getUntrackedFiles(repoPath)
	if err != nil {
		return err
	}

	if len(modifiedFiles) > 0 {
		fmt.Println("Changes not staged for commit:")
		for _, file := range modifiedFiles {
			fmt.Printf("  modified: %s\n", file)
		}
	}

	if len(untrackedFiles) > 0 {
		fmt.Println("Untracked files:")
		for _, file := range untrackedFiles {
			fmt.Printf("  %s\n", file)
		}
	}
	fmt.Println("Changes to be committed:")
	for filePath := range index {
		fmt.Printf("  new file:   %s\n", filePath)
	}

	return nil
}

func checkFileModified(repoPath, filePath string, index map[string]string) (bool, error) {

	stagedFileHash, ok := index[filePath]
	if !ok {
		return false, nil
	}

	filePathFull := filepath.Join(repoPath, filePath)
	fileHash, err := calculateFileHash(filePathFull)
	if err != nil {
		return false, err
	}
	if stagedFileHash != fileHash {
		return true, nil
	}

	return false, nil
}

func calculateFileHash(filePath string) (string, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	hash := sha1.New()

	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to calculate hash for file %s: %w", filePath, err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func getUntrackedFiles(repoPath string) ([]string, error) {
	var untrackedFiles []string
	err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.Contains(path, ".git") {
			return nil
		}

		// Check if the file is untracked
		// Logic for checking if the file is untracked goes here
		untrackedFiles = append(untrackedFiles, path)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk the repository: %w", err)
	}

	return untrackedFiles, nil
}
