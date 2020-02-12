package main

import (
	"database/sql"
	"fmt"
	"io"

	"github.com/ionous/iffy/dbutil"
)

type Pairing struct {
	Rel   *Relation
	Pairs []*Pair
}

type Pair struct {
	First, Second string
}

func listOfPairs(w io.Writer, relation string, db *sql.DB) (err error) {
	var rel Relation
	var pair Pair
	var pairs []*Pair

	if e := dbutil.QueryAll(db,
		fmt.Sprintf(`
		select relation, kind, cardinality, otherKind, coalesce((
			select spec from mdl_spec 
			where type='relation' and name=relation
			limit 1), '')
		from mdl_rel
		where relation = '%s'`, relation),
		func() (err error) {
			return
		}, &rel.Name, &rel.Kind, &rel.Cardinality, &rel.OtherKind, &rel.Spec); e != nil {
		err = e
	} else if e := dbutil.QueryAll(db,
		fmt.Sprintf(`
		select noun, otherNoun
		from mdl_pair
		where relation = '%s'`, relation),
		func() (err error) {
			pin := pair
			pairs = append(pairs, &pin)
			return
		}, &pair.First, &pair.Second,
	); e != nil {
		err = e
	} else {
		pin := rel
		err = templates.ExecuteTemplate(w, "pairList", &Pairing{
			Rel:   &pin,
			Pairs: pairs,
		})
	}
	return
}

func init() {
	registerTemplate("relHeader", `Relates
	{{- if HasPrefix .Cardinality "any_" }} many
	{{- end -}}
{{- "" }} <a href="/atlas/kinds#{{.Kind|Safe}}">{{.Kind|Title}}</a> to
	{{- if HasSuffix .Cardinality "_any" }} many
	{{- end -}} 
{{- "" }} <a href="/atlas/kinds#{{.OtherKind|Safe}}">{{.OtherKind|Title}}</a>.
	{{- if .Spec }}
{{ "" }} {{ .Spec }}
	{{- end -}}
`)

	registerTemplate("pairList", `
<h1>{{.Rel.Name|Title}}</h1>
{{ template "relHeader" .Rel }}
<table>
	{{- range $i, $el := .Pairs }}
<tr>
  <td>{{ if changing $i "First" $.Pairs }}<a href="/atlas/nouns#{{.First|Safe}}">{{.First|Title}}</a>{{end}}</td>
  <td>{{ if changing $i "Second" $.Pairs }}<a href="/atlas/nouns#{{.Second|Safe}}">{{.Second|Title}}</a>{{end}}</td>
</tr>
	{{- end }}
</table>
`)
}
