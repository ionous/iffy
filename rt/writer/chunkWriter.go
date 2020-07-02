package writer

// ChunkWriter adapts a WriteChunk method into writer friendly output.
type ChunkWriter interface {
	WriteChunk(Chunk) (int, error)
	init(alson ChunkWriter)
}

type ChunkOutput struct {
	target ChunkWriter
}

// InitChunks must be called before the first write
func InitChunks(n ChunkWriter) {
	n.init(n)
}

func (n *ChunkOutput) init(alson ChunkWriter) {
	n.target = alson
}

func (n ChunkOutput) Write(p []byte) (int, error) {
	return n.target.WriteChunk(Chunk{p})
}
func (n ChunkOutput) WriteByte(c byte) error {
	_, e := n.target.WriteChunk(Chunk{c})
	return e
}
func (n ChunkOutput) WriteRune(r rune) (int, error) {
	return n.target.WriteChunk(Chunk{r})
}
func (n ChunkOutput) WriteString(s string) (int, error) {
	return n.target.WriteChunk(Chunk{s})
}
