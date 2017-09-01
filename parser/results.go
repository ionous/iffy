package parser

// Results used by the parser include, a list of results, a resolved object, a resolved action, etc. On success, the parser generally returns a ResultList as its primary result.
type Result interface {
	// the number of words used to match this result.
	WordsMatched() int
}

type ResolvedAction struct {
	Name string
}

// ResolvedActor
// ResolvedNumber
// ResolvedWords

type ResolvedMulti struct {
	Nouns     []NounInstance
	WordCount int
}
type ResolvedObject struct {
	NounInstance NounInstance
	Words        []string // what the user said to identify the object
}
type ResolvedWord struct {
	Word string
}

func (f ResolvedAction) WordsMatched() int {
	return 0
}
func (f ResolvedMulti) WordsMatched() int {
	return f.WordCount
}
func (f ResolvedObject) WordsMatched() int {
	return len(f.Words)
}
func (f ResolvedWord) WordsMatched() int {
	return 1
}
