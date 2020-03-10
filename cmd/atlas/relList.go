package main

import (
	"database/sql"
	"io"

	"github.com/ionous/iffy/tables"
)

func listOfRelations(w io.Writer, db *sql.DB) (err error) {
	// 	// originally used a channel, but the template iterates over the same elements multiple times
	var rels []*Relation
	var rel Relation
	if e := tables.QueryAll(db, `
		select relation, kind, cardinality, otherKind, coalesce((
			select spec from mdl_spec 
			where type='relation' and name=relation
			limit 1), '')
		from mdl_rel
		order by relation`,
		func() (err error) {
			pin := rel
			rels = append(rels, &pin)
			return
		}, &rel.Name, &rel.Kind, &rel.Cardinality, &rel.OtherKind, &rel.Spec); e != nil {
		err = e
	} else {
		err = templates.ExecuteTemplate(w, "relList", rels)
	}
	return
}

type Relation struct {
	Name, Kind, Cardinality, OtherKind, Spec string
}

func init() {
	registerTemplate("relList", `
<h1>Relations</h1>
<dl>
	{{- range $i, $_ := . }}
  <dt><a href="/atlas/relations/{{.Name|safe}}">{{.Name|title}}</a></dt>
   <dd>{{ template "relHeader" . }}</dd>
	{{- end }}
</dl>
`)
}
