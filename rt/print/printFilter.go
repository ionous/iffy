package print

// implements writer.Output and io.Closer
type Filter struct {
	First, Rest func(Chunk) (int, error)
	Last        func() error
	cnt         int
}

func (f *Filter) Close() (err error) {
	if f.Last != nil {
		err = f.Last()
	}
	return
}

func (f *Filter) Write(p []byte) (int, error) {
	return f.write(Chunk{p})
}
func (f *Filter) WriteByte(c byte) error {
	_, e := f.write(Chunk{c})
	return e
}
func (f *Filter) WriteRune(r rune) (int, error) {
	return f.write(Chunk{r})
}
func (f *Filter) WriteString(s string) (int, error) {
	return f.write(Chunk{s})
}

func (f *Filter) write(c Chunk) (ret int, err error) {
	if f.cnt == 0 && f.First != nil {
		ret, err = f.First(c)
	} else if f.Rest != nil {
		ret, err = f.Rest(c)
	}
	f.cnt += ret
	return
}
