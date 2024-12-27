package objects

type Blob struct {
	Data []byte // Raw data of the file
}
// basically convert the commit object into bytes 
func (b *Blob) Serialize() ([]byte, error) {
	return b.Data, nil
}

func (b *Blob) Deserialize(data []byte) {
	b.Data = data
}

func (b *Blob) Type() string {
	return "blob"
}