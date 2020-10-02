package web

import (
	"io"
	"net/http"

	"golang.org/x/net/context"

	"github.com/ionous/errutil"
)

// Resource interfaces with a rest-ish endpoint, there are methods for each major verb.
// See also Wrapper which provides a function-based adapter for creating resources.
type Resource interface {
	// Return the named child resource
	Find(string) Resource
	// Read a resource
	Get(context.Context, http.ResponseWriter) error
	// Alter a resource
	Post(context.Context, io.Reader, http.ResponseWriter) error
	// Receive a resource
	Put(context.Context, io.Reader, http.ResponseWriter) error
}

// Turn one or more Resource compatible functions into a full interface implementation.
type Wrapper struct {
	Finds func(string) Resource
	Gets  func(context.Context, http.ResponseWriter) error
	Posts func(context.Context, io.Reader, http.ResponseWriter) error
	Puts  func(context.Context, io.Reader, http.ResponseWriter) error
}

var _ Resource = (*Wrapper)(nil) // ensure compliance

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

func (n *Wrapper) Post(ctx context.Context, r io.Reader, w http.ResponseWriter) (err error) {
	if post := n.Posts; post != nil {
		err = post(ctx, r, w)
	} else {
		err = errutil.New("unsupported post")
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
