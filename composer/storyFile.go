package composer

import (
	"io"
	"net/http"
	"os"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/web"
	"golang.org/x/net/context"
)

// path of a local .if file
type storyFile string

// String name of the file.
func (d storyFile) String() string {
	return string(d)
}

// Find actions for individual files
func (d storyFile) Find(sub string) (ret web.Resource) {
	switch sub {
	case "check":
		ret = &web.Wrapper{
			Posts: func(ctx context.Context, in io.Reader, out http.ResponseWriter) (err error) {
				if e := tempTest(ctx, d.String(), in); e != nil {
					err = e
				}
				return
			},
		}
	}
	return
}

// Get the contents of this resource.
func (d storyFile) Get(ctx context.Context, w http.ResponseWriter) (err error) {
	if f, e := os.Open(string(d)); e != nil {
		err = e
	} else {
		defer f.Close()
		w.Header().Set("Content-Type", "application/json")
		_, err = io.Copy(w, f)
	}
	return
}

// Post a modification to this resource
func (d storyFile) Post(ctx context.Context, r io.Reader, w http.ResponseWriter) (err error) {
	return errutil.New("unsupported post", d)
}

// Put new resource data in our place
func (d storyFile) Put(ctx context.Context, r io.Reader, w http.ResponseWriter) (err error) {
	// its okay to use Create because storyFolder.Get() ensures it already exists.
	if f, e := os.Create(string(d)); e != nil {
		err = e
	} else {
		defer f.Close()
		_, err = io.Copy(f, r)
	}
	return
}
