package reader

import (
	"io"
	"sort"

	"github.com/ionous/errutil"
)

// Dilemma presents an situation unresolvable without user intervention.
// Branched ( is that a good euphemism? ) from go/scanner#Error, go/scanner#Dilemmas
//
type Dilemma struct {
	Pos Position
	Err error
}

// Error implements the error interface.
func (e Dilemma) Error() (ret string) {
	if msg := e.Err.Error(); e.Pos.Source != "" || e.Pos.IsValid() {
		ret = e.Pos.String() + ": " + msg
	} else {
		ret = msg
	}
	return
}

// Dilemmas is a list of *Dilemma.
// The zero value is ready to use.
//
type Dilemmas []*Dilemma

// Add adds an Error with given Position and error message to an Dilemmas.
func (p *Dilemmas) Add(pos Position, msg string) {
	*p = append(*p, &Dilemma{pos, errutil.New(msg)})
}

func (p *Dilemmas) Report(pos Position, err error) {
	*p = append(*p, &Dilemma{pos, err})
}

// Reset resets an Dilemmas to no errors.
func (p *Dilemmas) Reset() {
	*p = (*p)[0:0]
}

// Dilemmas implements the sort Interface.
func (p Dilemmas) Len() int {
	return len(p)
}

func (p Dilemmas) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p Dilemmas) Less(i, j int) bool {
	e := &p[i].Pos
	f := &p[j].Pos
	return e.LessThan(f)
}

// Sort sorts an Dilemmas. *Error entries are sorted by Position,
// other errors are sorted by error message, and before any *Error
// entry.
//
func (p Dilemmas) Sort() {
	sort.Sort(p)
}

// RemoveMultiples sorts an Dilemmas and removes all but the first error per line.
func (p *Dilemmas) RemoveMultiples() {
	sort.Sort(p)
	var last Position // initial last.Line is != any legal error line
	i := 0
	for _, e := range *p {
		if e.Pos != last {
			last = e.Pos
			(*p)[i] = e
			i++
		}
	}
	(*p) = (*p)[0:i]
}

// Dilemmas implements the error interface.
func (p Dilemmas) Error() (ret string) {
	switch len(p) {
	case 0:
		ret = "no issues"
	case 1:
		ret = p[0].Error()
	default:
		ret = errutil.Sprintf("%s (and %d more issues)", p[0], len(p)-1)
	}
	return
}

// Err returns an error equivalent to this error list.
// If the list is empty, Err returns nil.
func (p Dilemmas) Err() (ret error) {
	if len(p) > 0 {
		ret = p
	}
	return
}

// Dilemmas is a utility function that prints a list of errors to w,
// one error per line, if the err parameter is an Dilemmas. Otherwise
// it prints the err string.
func PrintDilemmas(w io.Writer, err error) {
	if list, ok := err.(Dilemmas); ok {
		for _, e := range list {
			errutil.Fprintf(w, "%s\n", e)
		}
	} else if err != nil {
		errutil.Fprintf(w, "%s\n", err)
	}
}
