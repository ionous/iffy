package main

import (
	"database/sql"
	"io"
	"strings"

	"github.com/alecthomas/template"
	"github.com/ionous/iffy/dbutil"
)

// some things we could do:
// . an alphabetical index
// . a hierarchical, indented listing
// . headings per kind
//
type Kind struct {
	Kind, Path, Spec string
}

func (k Kind) Parent() string {
	return strings.Split(k.Path, ",")[0]
}

func kinds(w io.Writer, db *sql.DB) error {
	c := make(chan Kind)
	go func() {
		var curr Kind
		if e := dbutil.QueryAll(db, `
		select kind, path, coalesce(spec, '')
		from mdl_kind
		left join mdl_spec 
			on (type='kind' and name=kind)
		order by path, kind`,
			func() (err error) {
				c <- curr
				return
			}, &curr.Kind, &curr.Path, &curr.Spec); e != nil {
			panic(e)
		} else {
			close(c)
		}
	}()
	return kindsTemplate.Execute(w, c)
}

var kindsTemplate = template.Must(template.New("kinds").Parse(`
{{range .}}
<h2 id="{{.Kind}}">{{.Kind}}</h2> 
<div>parent: {{if .Parent}}
<a href="#{{.Parent}}">{{.Parent}}</a>
{{else}}
none
{{end}}</div>
<p class="spec">{{.Spec}}</p>
</h2>
{{end}}
`))
