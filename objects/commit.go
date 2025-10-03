package objects

import (
	"bytes"
	"fmt"
	"strings"
)

// Commit represents a Git commit object.
type Commit struct {
	Tree    string   // SHA-1 hash of the tree object
	Parents []string // SHA-1 hashes of parent commits (if any)
	Author  string   // Author of the commit
	Message string   // Commit message
}

// Serialize converts the commit object into bytes for storage.
func (c *Commit) Serialize() ([]byte, error) {
	var buf bytes.Buffer

	// Write tree hash
	buf.WriteString(fmt.Sprintf("tree %s\n", c.Tree))

	// Write parent hashes
	for _, parent := range c.Parents {
		buf.WriteString(fmt.Sprintf("parent %s\n", parent))
	}

	// Write author information
	buf.WriteString(fmt.Sprintf("author %s\n", c.Author))

	// Write the commit message
	buf.WriteString("\n") // Separate metadata and message with a blank line
	buf.WriteString(c.Message)

	return buf.Bytes(), nil
}

// Deserialize populates the commit object from bytes.
func (c *Commit) Deserialize(data []byte) {
	// Split metadata and message
	parts := bytes.SplitN(data, []byte("\n\n"), 2)
	metadata := string(parts[0])
	message := ""
	if len(parts) > 1 {
		message = string(parts[1])
	}

	// Parse metadata
	lines := strings.Split(metadata, "\n")
	for _, line := range lines {
		// Extract key and value
		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			continue
		}
		key, value := parts[0], parts[1]

		// Handle different keys
		switch key {
		case "tree":
			c.Tree = value
		case "parent":
			c.Parents = append(c.Parents, value)
		case "author":
			c.Author = value
		}
	}

	// Assign the message
	c.Message = message
}

// Type returns the type of the object ("commit").
func (c *Commit) Type() string {
	return "commit"
}
