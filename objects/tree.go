package objects

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

// TreeEntry represents a single entry in a tree object.
type TreeEntry struct {
	Mode string // File mode (e.g., "100644" for a file, "040000" for a directory)
	Hash string // SHA-1 hash of the referenced blob or tree
	Name string // Name of the file or directory
}

// Tree represents a Git tree object.
type Tree struct {
	Entries []TreeEntry // List of entries in the tree
}

// Serialize converts the tree object into bytes for storage.
func (t *Tree) Serialize() ([]byte, error) {
	var buf bytes.Buffer

	for _, entry := range t.Entries {
		// Write the mode, name, and null byte separator
		entryData := fmt.Sprintf("%s %s\x00", entry.Mode, entry.Name)
		buf.WriteString(entryData)

		// Write the hash as raw bytes (not as a hex string)
		hashBytes := decodeHex(entry.Hash)
		buf.Write(hashBytes)
	}

	return buf.Bytes(), nil
}

// Deserialize populates the tree object from bytes.
func (t *Tree) Deserialize(data []byte) {
	var entries []TreeEntry
	for len(data) > 0 {
		// Find the null byte separating name and hash
		nullIdx := bytes.IndexByte(data, '\x00')
		if nullIdx < 0 {
			break
		}

		// Parse mode and name
		metadata := string(data[:nullIdx])
		parts := strings.SplitN(metadata, " ", 2)
		if len(parts) != 2 {
			break
		}
		mode, name := parts[0], parts[1]

		// Parse hash
		hashBytes := data[nullIdx+1 : nullIdx+21]
		hash := encodeHex(hashBytes)

		// Create a TreeEntry
		entries = append(entries, TreeEntry{
			Mode: mode,
			Name: name,
			Hash: hash,
		})

		// Move to the next entry
		data = data[nullIdx+21:]
	}
	t.Entries = entries
}

// Type returns the type of the object ("tree").
func (t *Tree) Type() string {
	return "tree"
}

// encodeHex converts raw bytes to a hex string.
func encodeHex(data []byte) string {
	return fmt.Sprintf("%x", data)
}

// decodeHex converts a hex string to raw bytes.
func decodeHex(hexStr string) []byte {
	data := make([]byte, len(hexStr)/2)
	for i := 0; i < len(data); i++ {
		binary.Read(strings.NewReader(hexStr[2*i:2*i+2]), binary.BigEndian, &data[i])
	}
	return data
}
