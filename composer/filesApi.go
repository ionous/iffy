package composer

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/web"
	"golang.org/x/net/context"
)

func FilesApi(cfg *Config) web.Resource {
	return &web.Wrapper{
		Finds: func(name string) (ret web.Resource) {
			switch name {
			case "files":
				ret = dir(cfg.Root)
			}
			return
		},
	}
}

type dir string

// Return the named sub-resource
func (d dir) Find(sub string) (ret web.Resource) {
	if !strings.HasPrefix(sub, ".") {
		ret = dir(path.Join(string(d), sub))
	}
	return
}

// Write the resource
func (d dir) Get(ctx context.Context, w http.ResponseWriter) (err error) {
	files := []string{}
	start := string(d)
	if e := filepath.Walk(start, func(path string, info os.FileInfo, e error) (err error) {
		if path != start && e == nil {
			isDir := info.IsDir()
			if isDir {
				err = filepath.SkipDir
			}
			if !strings.HasPrefix(info.Name(), ".") {
				if isDir || strings.HasSuffix(path, ".if.js") {
					chop := len(d)
					if !isDir {
						chop++
					}
					files = append(files, path[chop:])
				}
			}
		}
		return
	}); e != nil {
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
func (d dir) Put(context.Context, io.Reader, http.ResponseWriter) error {
	return errutil.New("unsupported put")
}
