package parser

type ErrorDepth interface {
	ErrorDepth() int
}

type Depth int

func DepthOf(e interface{}) (ret int) {
	if d, ok := e.(ErrorDepth); ok {
		ret = d.ErrorDepth()
	}
	return
}

type AmbiguousObject struct {
	Nouns []Noun
	Depth
}

type MismatchedWord struct {
	Word string
	Depth
}

type MissingObject struct {
	Depth
}

// NoSuchObjects after asking for multiple items, and finding none.
type NoSuchObjects struct {
	Depth
}

// Overflow when we expect to be done, but input tokens remain.
type Overflow struct {
	Depth
}

// Underflow when we expect a word, but the input is empty
type Underflow struct {
	Depth
}

type UnknownObject struct {
	Depth
}

func (d Depth) ErrorDepth() int {
	return int(d)
}

func (AmbiguousObject) Error() string {
	return "couldnt determine object"
}
func (MismatchedWord) Error() string {
	return "too few words"
}
func (MissingObject) Error() string {
	return "expected an object"
}
func (NoSuchObjects) Error() string {
	return "you cant see any such things"
}
func (Overflow) Error() string {
	return "too many words"
}
func (Underflow) Error() string {
	return "too few words"
}
func (UnknownObject) Error() string {
	return "you cant see any such thing"
}
