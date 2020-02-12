package main

import (
	"database/sql"
	"fmt"
	"io"
	"text/template"

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
		err = pairTemplate.Execute(w, &Pairing{
			Rel:   &pin,
			Pairs: pairs,
		})
	}
	return
}

var pairTemplate = template.Must(template.New("pairs").Funcs(funcMap).Parse(`
<h1>{{.Rel.Name|Title}}</h1>
{{ .Rel.Text }}.{{ if .Rel.Spec }} {{ .Rel.Spec }}{{ end }}
<table>
	{{- range $i, $el := .Pairs }}
<tr>
  <td>{{ if changing $i "First" $.Pairs }}{{.First|Title}}{{end}}</td>
  <td>{{ if changing $i "Second" $.Pairs }}{{.Second|Title}}{{end}}</td>
</tr>
	{{- end }}
</table>
`))
