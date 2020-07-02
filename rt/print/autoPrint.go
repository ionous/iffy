package print

import (
	"unicode"

	"github.com/ionous/iffy/rt/writer"
)

type AutoWriter struct {
	writer.ChunkOutput
	Target writer.Output
}

func NewAutoWriter(w writer.Output) *AutoWriter {
	a := &AutoWriter{Target: w}
	writer.InitChunks(a)
	return a
}

func (w *AutoWriter) WriteChunk(c writer.Chunk) (int, error) {
	n, e := c.WriteTo(w.Target)
	if last, _ := c.DecodeLastRune(); unicode.Is(unicode.Terminal_Punctuation, last) {
		w.Target.WriteRune('\n')
	}
	return n, e
}
