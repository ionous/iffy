package composer

import (
	"io"
	"net/http"
	"os"

	"github.com/ionous/iffy/web"
	"golang.org/x/net/context"
)

type storyFile string

// returns nil. files have no sub-resources.
func (d storyFile) Find(sub string) (ret web.Resource) {
	return
}

// Write the resource
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

// Receive a resource
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
