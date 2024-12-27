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
	Date    string   // Date when the commit was made (in RFC3339 format)
}

// Serialize converts the Commit object into bytes for storage.
func (c *Commit) Serialize() ([]byte, error) {
	var buf bytes.Buffer

	// Write the tree hash
	buf.WriteString(fmt.Sprintf("tree %s\n", c.Tree))

	// Write parent commit hashes
	for _, parent := range c.Parents {
		buf.WriteString(fmt.Sprintf("parent %s\n", parent))
	}

	// Write author and date
	buf.WriteString(fmt.Sprintf("author %s\n", c.Author))
	buf.WriteString(fmt.Sprintf("date %s\n", c.Date)) // Added date

	// Write commit message
	buf.WriteString("\n")
	buf.WriteString(c.Message)

	return buf.Bytes(), nil
}

// Deserialize populates the Commit object from the given data.
func (c *Commit) Deserialize(data []byte) {
	// Split the data into metadata and message
	parts := bytes.SplitN(data, []byte("\n\n"), 2)
	metadata := string(parts[0])
	message := ""
	if len(parts) > 1 {
		message = string(parts[1])
	}

	// Parse the metadata for tree, parent, author, and date
	lines := strings.Split(metadata, "\n")
	for _, line := range lines {
		// Split key and value (e.g., "author" and "John Doe")
		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			continue
		}
		key, value := parts[0], parts[1]

		// Handle different keys and assign values
		switch key {
		case "tree":
			c.Tree = value
		case "parent":
			c.Parents = append(c.Parents, value)
		case "author":
			c.Author = value
		case "date":
			c.Date = value // Parse the date
		}
	}

	// Set the commit message
	c.Message = message
}

// Type returns the type of the object ("commit").
func (c *Commit) Type() string {
	return "commit"
}
