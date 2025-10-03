package objects

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// GitObject is the interface for all Git object types (e.g., blob, tree, commit).
type GitObject interface {
	Type() string               // Returns the object type (e.g., "blob")
	Serialize() ([]byte, error) // Converts the object into bytes for storage
	Deserialize(data []byte)    // Populates the object from bytes
}

// ReadObject reads a Git object from the `.git/objects` directory using its SHA hash.
// ReadObject reads a Git object from the `.git/objects` directory using its SHA hash.
func ReadObject(repoPath, sha string) (GitObject, error) {
	// Construct the path to the object file
	objPath := filepath.Join(repoPath, ".git", "objects", sha[:2], sha[2:])
	file, err := os.Open(objPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open object file: %w", err)
	}
	defer file.Close()

	// Decompress the file
	zr, err := zlib.NewReader(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decompress object: %w", err)
	}
	defer zr.Close()

	// Read decompressed data
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, zr); err != nil {
		return nil, fmt.Errorf("failed to read decompressed data: %w", err)
	}
	raw := buf.Bytes()

	// Extract the object type and data
	spaceIdx := bytes.IndexByte(raw, ' ')
	nullIdx := bytes.IndexByte(raw, '\x00')

	if spaceIdx < 0 || nullIdx < 0 || nullIdx <= spaceIdx {
		return nil, fmt.Errorf("invalid object header")
	}

	objType := string(raw[:spaceIdx])
	objData := raw[nullIdx+1:]

	// Create and return the appropriate GitObject
	var obj GitObject
	switch objType {
	case "blob":
		obj = &Blob{}
	case "tree":
		obj = &Tree{}
	case "commit":
		obj = &Commit{}
	default:
		return nil, fmt.Errorf("unknown object type: %s", objType)
	}

	obj.Deserialize(objData)
	return obj, nil
}

func WriteObject(obj GitObject, repoPath string) (string, error) {
	// Serialize the object data
	data, err := obj.Serialize()
	if err != nil {
		fmt.Printf("DEBUG: Failed to serialize object: %v\n", err)
		return "", fmt.Errorf("failed to serialize object: %w", err)
	}

	// Add the header and prepare the object data
	header := fmt.Sprintf("%s %d\x00", obj.Type(), len(data))
	storeData := append([]byte(header), data...)

	// Compute the SHA-1 hash
	sha := fmt.Sprintf("%x", sha1.Sum(storeData))
	fmt.Printf("DEBUG: Computed SHA-1: %s\n", sha)

	// Determine the object path
	objDir := filepath.Join(repoPath, ".git", "objects", sha[:2])
	objPath := filepath.Join(objDir, sha[2:])
	fmt.Printf("DEBUG: Object directory: %s\n", objDir)
	fmt.Printf("DEBUG: Object path: %s\n", objPath)

	// Ensure the directory exists
	if err := os.MkdirAll(objDir, 0755); err != nil {
		fmt.Printf("DEBUG: Failed to create directory %s: %v\n", objDir, err)
		return "", fmt.Errorf("failed to create object directory: %w", err)
	}
	fmt.Printf("DEBUG: Directory created: %s\n", objDir)

	// Write the compressed object to the file
	file, err := os.Create(objPath)
	if err != nil {
		fmt.Printf("DEBUG: Failed to create object file %s: %v\n", objPath, err)
		return "", fmt.Errorf("failed to create object file: %w", err)
	}
	defer file.Close()
	fmt.Printf("DEBUG: File created: %s\n", objPath)

	zw := zlib.NewWriter(file)
	if _, err := zw.Write(storeData); err != nil {
		fmt.Printf("DEBUG: Failed to write compressed data: %v\n", err)
		return "", fmt.Errorf("failed to write compressed data: %w", err)
	}
	if err := zw.Close(); err != nil {
		fmt.Printf("DEBUG: Failed to close zlib writer: %v\n", err)
		return "", fmt.Errorf("failed to close zlib writer: %w", err)
	}
	fmt.Printf("DEBUG: Object successfully written: %s\n", objPath)

	return sha, nil
}
