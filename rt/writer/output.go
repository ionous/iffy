package writer

// Output defines a sink for text output.
// It's a subset of strings.Builder.
type Output interface {
	Write(p []byte) (int, error)
	WriteByte(c byte) error
	WriteRune(r rune) (int, error)
	WriteString(s string) (int, error)
}

type OutputCloser interface {
	Output
	Close() error
}
