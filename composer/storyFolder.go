package composer

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/web"
	"golang.org/x/net/context"
)

type storyFolder string

// Return the named sub-resource
func (d storyFolder) Find(sub string) (ret web.Resource) {
	base := string(d)
	path := filepath.Join(base, sub)
	// join cleans the elements; removing .. paths
	// it helps let us make sure we're still in the right spot.
	if strings.HasPrefix(path, base) {
		if info, e := os.Lstat(path); e != nil {
			// we could return an erroring resource if we really wanted i suppose...
			println("error reading", sub, e)
		} else if info.IsDir() {
			ret = storyFolder(path)
		} else {
			ret = storyFile(path)
		}
	}
	return
}

// Write the resource
func (d storyFolder) Get(ctx context.Context, w http.ResponseWriter) (err error) {
	if files, e := listDirectory(string(d)); e != nil {
		err = e
	} else if b, e := json.Marshal(files); e != nil {
		err = e
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
	return
}

// Receive a resource
func (d storyFolder) Put(context.Context, io.Reader, http.ResponseWriter) error {
	return errutil.New("unsupported put")
}

// based on filepath.Walk
func listDirectory(path string) (ret []string, err error) {
	if f, e := os.Open(path); e != nil {
		err = e
	} else if names, e := f.Readdirnames(-1); e != nil {
		err = e
	} else {
		for _, name := range names {
			filename := filepath.Join(path, name)
			if info, e := os.Lstat(filename); e != nil {
				err = e
				break
			} else {
				isDir := info.IsDir()
				if isDir || strings.HasSuffix(name, ".if") {
					if name[0] != '_' && name[0] != '.' {
						if isDir {
							name = "/" + name
						}
						ret = append(ret, name)
					}
				}
			}
		}
	}
	return
}
