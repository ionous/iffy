package parser

type Scanner interface {
	// Store results in Context, return the number of words to advance.
	Scan(Cursor, *Context) (int, bool)
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

// // FIX: return a flag, set a flag, for "greedy" -- i want everything that the next thing cant use.
// func (w *Series) Parse(x *Context) (err error) {
// 	x.Soft = true
// }

// Action terminates a matcher sequence by setting the context to the desired action.
type Action struct {
	Name string
}

// Scan matches only if the cursor has finished with all words.
func (a *Action) Scan(scan Cursor, ctx *Context) (ret int, okay bool) {
	_, exists := scan.GetNext()
	// println("trying", a.Name, scan.Word, !exists)
	if !exists {
		// FIX: can we make this a return?
		ctx.Results.Action = a.Name
		okay = true
	}
	return
}
