package main

import (
	"database/sql"
	"io"
	"strings"
	"text/template"

	"github.com/ionous/iffy/dbutil"
)

func listOfRelations(w io.Writer, db *sql.DB) (err error) {
	// 	// originally used a channel, but the template iterates over the same elements multiple times
	var rels []*Relation
	var rel Relation
	if e := dbutil.QueryAll(db, `
		select *, coalesce((
			select spec from mdl_spec 
			where type='relation' and name=relation
		), '')
		from mdl_rel
		order by relation`,
		func() (err error) {
			pin := rel
			rels = append(rels, &pin)
			return
		}, &rel.Name, &rel.Kind, &rel.Cardinality, &rel.OtherKind, &rel.Spec); e != nil {
		err = e
	} else {
		err = relTemplate.Execute(w, rels)
	}
	return
}

type Relation struct {
	Name, Kind, Cardinality, OtherKind, Spec string
}

func (r *Relation) Text() string {
	// fmt.Relates [many] kinds to [many] kinds.
	var els []string
	els = append(els, "Relates")
	if strings.HasPrefix(r.Cardinality, "any_") {
		els = append(els, "many")
	}
	els = append(els, r.Kind)
	els = append(els, "to")
	if strings.HasSuffix(r.Cardinality, "_any") {
		els = append(els, "many")
	}
	els = append(els, r.OtherKind)
	return strings.Join(els, " ")
}

var relTemplate = template.Must(template.New("rels").Funcs(funcMap).Parse(`
<h1>Relations</h1>
<dl>
	{{- range $i, $_ := . }}
  <dt><a href="/atlas/relations/{{.Name|Safe}}">{{.Name|Title}}</a></dt>
   <dd>{{.Text}}. {{.Spec}}</dd>
	{{- end }}
</dl>
`))
