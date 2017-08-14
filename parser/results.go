package parser

type Result interface {
	// Len (number) of words used to match this result.
	ResultLen() int
}

type ResolvedAction struct {
	Name string
}

// ResolvedActor
// ResolvedNumber
// ResolvedWords

type ResolvedMulti struct {
	Nouns     []Noun
	WordCount int
}
type ResolvedObject struct {
	Noun  Noun
	Words []string // what the user said to identify the object
}
type ResolvedWord struct {
	Word string
}

func (f ResolvedAction) ResultLen() int {
	return 0
}
func (f ResolvedMulti) ResultLen() int {
	return f.WordCount
}
func (f ResolvedObject) ResultLen() int {
	return len(f.Words)
}
func (f ResolvedWord) ResultLen() int {
	return 1
}
