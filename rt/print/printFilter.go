package print

import "github.com/ionous/iffy/rt/writer"

// implements writer.ChunkWriter and io.Closer
type Filter struct {
	writer.ChunkOutput
	First, Rest func(writer.Chunk) (int, error)
	Last        func() error
	cnt         int
}

func (f *Filter) Close() (err error) {
	if f.Last != nil {
		err = f.Last()
	}
	return
}

func (f *Filter) WriteChunk(c writer.Chunk) (ret int, err error) {
	if f.cnt == 0 && f.First != nil {
		ret, err = f.First(c)
	} else if f.Rest != nil {
		ret, err = f.Rest(c)
	}
	f.cnt += ret
	return
}
