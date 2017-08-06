package parser

type Grammar struct {
	Name string
	Matcher
}

type Scanner interface {
	// Store results in Context, return the number of words to advance.
	Scan(Cursor, *Context) (int, bool)
}

type Matcher interface {
	Scanner
	GetNext() Matcher
}

type Match struct {
	Scanner
	Next Matcher
}

// // Series matches any set of words.
// type Series struct {}

// // Number matches any series of digits
// // FUTURE: commas? and words.
// type Number struct {
// }

// // Reverse queues the next match till later.
// type Reverse struct {
// }

// // Actor handles addressing actors by name. For example: "name, "
// type Actor struct {
// }

func (m *Match) GetNext() Matcher {
	return m.Next
}

// // FIX: return a flag, set a flag, for "greedy" -- i want everything that the next thing cant use.
// func (w *Series) Parse(x *Context) (err error) {
// 	x.Soft = true
// }
