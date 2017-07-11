package std

//go:generate stringer -type=Verbosity
type Verbosity int

const (
	// "gives short descriptions of locations (even if you haven't been there before)."
	SuperBrief Verbosity = iota
	// "gives long descriptions of places never before visited and short descriptions otherwise."
	Brief
	// "gives long descriptions of locations (even if you've been there before)."
	Verbose
)

type Options struct {
	Verbosity
}
