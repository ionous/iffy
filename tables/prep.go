package tables

import (
	"database/sql"
	"io"

	"github.com/ionous/errutil"
)

type Prep []*prepErr

type prepErr struct {
	q string
	e error
}

// Error implements the error interface.
func (ps prepErr) Error() (ret string) {
	return errutil.Sprint("error", ps.e, "preparing", ps.q)
}

// unfortunately, this returns nil on error
// pkg sql doesnt provide a way to prepare an empty statement with a "sticky error" from outside the package.
func (ps *Prep) Prep(db *sql.DB, q string) (ret *sql.Stmt) {
	if stmt, e := db.Prepare(q); e != nil {
		*ps = append(*ps, &prepErr{q, e})
	} else {
		ret = stmt
	}
	return
}

// Reset resets an Prep to no errors.
func (ps *Prep) Reset() {
	*ps = (*ps)[0:0]
}

// Prep implements the error interface.
func (ps Prep) Error() (ret string) {
	switch len(ps) {
	case 0:
		ret = "no issues"
	case 1:
		ret = ps[0].Error()
	default:
		ret = errutil.Sprintf("%s (and %d more issues)", ps[0], len(ps)-1)
	}
	return
}

// Err returns an error equivalent to this error list.
// If the list is empty, Err returns nil.
func (ps Prep) Err() (ret error) {
	if len(ps) > 0 {
		ret = ps
	}
	return
}

// Prep is a utility function that prints a list of errors to w,
// one error per line, if the err parameter is an Prep. Otherwise
// it prints the err string.
func PrintPrep(w io.Writer, err error) {
	if list, ok := err.(Prep); ok {
		for _, e := range list {
			errutil.Fprintf(w, "%s\n", e)
		}
	} else if err != nil {
		errutil.Fprintf(w, "%s\n", err)
	}
}
