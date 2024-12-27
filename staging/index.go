package staging

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
)

// ReadIndex reads and decodes the index file.
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

	// Debugging: print raw index file data length
	fmt.Printf("Raw index data length: %d bytes\n", len(data))

	// Decode the index file
	decoder := gob.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(&index); err != nil {
		return nil, fmt.Errorf("failed to decode index file: %w", err)
	}

	// Debugging: print decoded index
	fmt.Println("Decoded index:")
	for filePath, hash := range index {
		fmt.Printf("  File: %s, Hash: %s\n", filePath, hash)
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

	// Debugging: print encoded index data length
	fmt.Printf("Encoded index data length: %d bytes\n", buffer.Len())

	// Write the encoded index to the file
	if err := os.WriteFile(indexPath, buffer.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write index file: %w", err)
	}

	// Debugging: confirm index file written
	fmt.Printf("Index written successfully to %s\n", indexPath)

	return nil
}

// UpdateIndex updates the staging index with the new file path and blob hash.
func UpdateIndex(repoPath, filePath, blobHash string) error {
	// Read the existing index
	index, err := ReadIndex(repoPath)
	if err != nil {
		return fmt.Errorf("failed to read index: %w", err)
	}

	// Print the index before updating
	fmt.Println("Index before update:")
	for path, hash := range index {
		fmt.Printf("  File: %s, Hash: %s\n", path, hash)
	}

	// Update the index
	index[filePath] = blobHash

	// Print the index after updating
	fmt.Println("Index after update:")
	for path, hash := range index {
		fmt.Printf("  File: %s, Hash: %s\n", path, hash)
	}

	// Write the updated index back to the file
	if err := WriteIndex(repoPath, index); err != nil {
		return fmt.Errorf("failed to update index: %w", err)
	}

	return nil
}

// ResetIndex resets the index file by removing it, useful for fixing index corruption.
func ResetIndex(repoPath string) error {
	indexPath := filepath.Join(repoPath, ".git", "index")
	if err := os.Remove(indexPath); err != nil {
		return fmt.Errorf("failed to remove index file: %w", err)
	}

	// Debugging: confirm reset
	fmt.Printf("Index file at %s has been removed.\n", indexPath)

	return nil
}
