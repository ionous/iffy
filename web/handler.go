package web

import (
	"log"
	"net/http"
	"strings"

	"github.com/ionous/errutil"
)

// HandleResource turns a Resource into an http.HandlerFunc;
// providing responses to http get and post requests.
func HandleResource(root Resource) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("handling", r.URL.Path, r.Method)
		if e := handleResponse(w, r, root); e != nil {
			log.Println(e)
		}
	}
}

// NOTE: the error, if any, is automatically passed to http.Error
func handleResponse(w http.ResponseWriter, r *http.Request, root Resource) (err error) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// chop off the leading and trailing slash. wise? i dont know.
	if res, e := FindResource(root, r.URL.Path[1:len(r.URL.Path)-1]); e != nil {
		http.NotFound(w, r)
		err = e
	} else if r.Method != "GET" {
		http.Error(w, r.Method, http.StatusMethodNotAllowed)
		err = errutil.Fmt("method %s not allowed", r.Method)
	} else if e := res.Get(w); e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
		err = e
	}
	return
}

// FindResource expands the passed resource, using each element of the passed path in turn.
// Returns an error, PathError, describing the extent of the matched path.
func FindResource(res Resource, path string) (ret Resource, err error) {
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if sub := res.Find(part); sub != nil {
			res = sub // set for next iteration of the loop
		} else {
			err = errutil.Fmt("failed to find resource %d(%s) in %s", i, part, path)
			break
		}
	}
	if err == nil {
		ret = res
	}
	return res, err
}
