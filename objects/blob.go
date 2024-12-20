package objects

// Blob represents a Git blob object.
type Blob struct {
	Data []byte // Raw data of the file
}

// Serialize converts the blob data into bytes for storage.
func (b *Blob) Serialize() ([]byte, error) {
	return b.Data, nil
}

// Deserialize populates the blob with data from bytes.
func (b *Blob) Deserialize(data []byte) {
	b.Data = data
}

// Type returns the type of the object ("blob").
func (b *Blob) Type() string {
	return "blob"
}
