package staging

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
)

// IndexEntry represents a single entry in the staging area.
type IndexEntry struct {
	FilePath string // Relative path of the file
	BlobHash string // SHA-1 hash of the blob
}

// ReadIndex reads the contents of the `.git/index` file.
func ReadIndex(repoPath string) (map[string]string, error) {
	indexPath := filepath.Join(repoPath, ".git", "index")
	index := make(map[string]string)

	// Check if the index file exists
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		return index, nil // Return an empty index if the file doesn't exist
	}

	// Read the index file
	data, err := os.ReadFile(indexPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read index file: %w", err)
	}

	// Decode the index file
	decoder := gob.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(&index); err != nil {
		return nil, fmt.Errorf("failed to decode index file: %w", err)
	}

	return index, nil
}

// WriteIndex writes the given index to the `.git/index` file.
func WriteIndex(repoPath string, index map[string]string) error {
	indexPath := filepath.Join(repoPath, ".git", "index")

	// Encode the index into bytes
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(index); err != nil {
		return fmt.Errorf("failed to encode index: %w", err)
	}

	// Write the encoded index to the file
	if err := os.WriteFile(indexPath, buffer.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write index file: %w", err)
	}

	return nil
}

// UpdateIndex adds or updates a file entry in the `.git/index` file.
func UpdateIndex(repoPath, filePath, blobHash string) error {
	// Read the existing index
	index, err := ReadIndex(repoPath)
	if err != nil {
		return fmt.Errorf("failed to read index: %w", err)
	}

	// Update the index
	index[filePath] = blobHash

	// Write the updated index back to the file
	if err := WriteIndex(repoPath, index); err != nil {
		return fmt.Errorf("failed to update index: %w", err)
	}

	return nil
}
