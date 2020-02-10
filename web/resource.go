package web

import (
	"net/http"

	"github.com/ionous/errutil"
)

// Resource interfaces with a rest-ish endpoint.
// See also, Wrapper, which provides a function-based adapter.
type Resource interface {
	// Return the named sub-resource
	Find(string) Resource
	// Write the resource
	Get(http.ResponseWriter) error
}

// Turn one or more Resource compatible functions into a full interface implementation.
type Wrapper struct {
	Finds func(string) Resource
	Gets  func(http.ResponseWriter) error
}

func (n *Wrapper) Find(child string) (ret Resource) {
	if f := n.Finds; f != nil {
		ret = f(child)
	}
	return
}

func (n *Wrapper) Get(w http.ResponseWriter) (err error) {
	if q := n.Gets; q != nil {
		err = q(w)
	} else {
		err = errutil.New("unsupported get")
	}
	return
}
