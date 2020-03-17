package web

import (
	"io"
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
	// Receive a resource
	Put(io.Reader, http.ResponseWriter) error
}

// Turn one or more Resource compatible functions into a full interface implementation.
type Wrapper struct {
	Finds func(string) Resource
	Gets  func(http.ResponseWriter) error
	Puts  func(io.Reader, http.ResponseWriter) error
}

func (n *Wrapper) Find(child string) (ret Resource) {
	if f := n.Finds; f != nil {
		ret = f(child)
	}
	return
}

func (n *Wrapper) Get(w http.ResponseWriter) (err error) {
	if get := n.Gets; get != nil {
		err = get(w)
	} else {
		err = errutil.New("unsupported get")
	}
	return
}

func (n *Wrapper) Put(r io.Reader, w http.ResponseWriter) (err error) {
	if put := n.Puts; put != nil {
		err = put(r, w)
	} else {
		err = errutil.New("unsupported put")
	}
	return
}
