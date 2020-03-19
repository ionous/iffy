package web

import (
	"io"
	"net/http"

	"golang.org/x/net/context"

	"github.com/ionous/errutil"
)

// Resource interfaces with a rest-ish endpoint.
// See also, Wrapper, which provides a function-based adapter.
type Resource interface {
	// Return the named sub-resource
	Find(string) Resource
	// Write the resource
	Get(context.Context, http.ResponseWriter) error
	// Receive a resource
	Put(context.Context, io.Reader, http.ResponseWriter) error
}

// Turn one or more Resource compatible functions into a full interface implementation.
type Wrapper struct {
	Finds func(string) Resource
	Gets  func(context.Context, http.ResponseWriter) error
	Puts  func(context.Context, io.Reader, http.ResponseWriter) error
}

func (n *Wrapper) Find(child string) (ret Resource) {
	if f := n.Finds; f != nil {
		ret = f(child)
	}
	return
}

func (n *Wrapper) Get(ctx context.Context, w http.ResponseWriter) (err error) {
	if get := n.Gets; get != nil {
		err = get(ctx, w)
	} else {
		err = errutil.New("unsupported get")
	}
	return
}

func (n *Wrapper) Put(ctx context.Context, r io.Reader, w http.ResponseWriter) (err error) {
	if put := n.Puts; put != nil {
		err = put(ctx, r, w)
	} else {
		err = errutil.New("unsupported put")
	}
	return
}
