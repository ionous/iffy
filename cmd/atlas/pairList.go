package main

import "text/template"

type Pair struct {
	First, Second string
}

// func listOfNouns(w io.Writer, db *sql.DB) (err error) {
// 	// originally used a channel, but the template iterates over the same elements multiple times
// 	var nouns []Noun
// 	var name, kind, spec string
// 	if e := dbutil.QueryAll(db, `
// 		select noun, kind, coalesce(spec, '')
// 		from mdl_noun
// 		left join mdl_spec
// 			on (type='noun' and name=noun)
// 		order by noun`,
// 		func() (err error) {
// 			var prop, value, rel string
// 			var props []Prop
// 			// var relation string
// 			var relations []string
// 			if e := dbutil.QueryAll(db,
// 				fmt.Sprintf(`
// 					select field, value
// 					from mdl_start
// 					where noun='%s'
// 					order by field`, name),
// 				func() (err error) {
// 					props = append(props, Prop{Name: name, Value: value})
// 					return
// 				},
// 				&prop, &value); e != nil {
// 				err = e
// 			} else if e := dbutil.QueryAll(db,
// 				fmt.Sprintf(`select
// 					distinct relation
// 					from mdl_pair
// 					where noun='%s'
// 					or otherNoun='%s'
// 					order by relation`, name, name),
// 				func() (err error) {
// 					relations = append(relations, rel)
// 					return
// 				}, &rel); e != nil {
// 				err = e
// 			} else {
// 				nouns = append(nouns, Noun{
// 					name, kind, spec,
// 					props,
// 					relations,
// 				})
// 			}
// 			return
// 		}, &name, &kind, &spec); e != nil {
// 		err = e
// 	} else {
// 		err = nounTemplate.Execute(w, nouns)
// 	}
// 	return
// }

type Pairing struct {
	Rel   *Relation
	Pairs []*Pair
}

var pairTemplate = template.Must(template.New("pairs").Funcs(funcMap).Parse(`
<h1>{{.Rel.Name|Title}}</h1>
{{ .Rel.Text }}. {{ .Rel.Spec }}
<table>
	{{- range $i, $el := .Pairs }}
<tr>
  <td>{{ if changing $i "First" $.Pairs }}{{.First|Title}}{{end}}</td>
  <td>{{ if changing $i "Second" $.Pairs }}{{.Second|Title}}{{end}}</td>
</tr>
	{{- end }}
</table>
`))
