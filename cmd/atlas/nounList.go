package main

import (
	"database/sql"
	"fmt"
	"io"
	"text/template"

	"github.com/ionous/iffy/dbutil"
)

type Noun struct {
	Name, Kind, Spec string
	Props            []Prop
	Relations        []string // list of relations involving the noun
}

func listOfNouns(w io.Writer, db *sql.DB) (err error) {
	// originally used a channel, but the template iterates over the same elements multiple times
	var nouns []Noun
	var name, kind, spec string
	if e := dbutil.QueryAll(db, `
		select noun, kind, coalesce(spec, '')
		from mdl_noun
		left join mdl_spec
			on (type='noun' and name=noun)
		order by noun`,
		func() (err error) {
			var prop, value, rel string
			var props []Prop
			// var relation string
			var relations []string
			if e := dbutil.QueryAll(db,
				fmt.Sprintf(`
					select field, value 
					from mdl_start 
					where noun='%s' 
					order by field`, name),
				func() (err error) {
					props = append(props, Prop{Name: name, Value: value})
					return
				},
				&prop, &value); e != nil {
				err = e
			} else if e := dbutil.QueryAll(db,
				fmt.Sprintf(`select 
					distinct relation 
					from mdl_pair 
					where noun='%s' 
					or otherNoun='%s'
					order by relation`, name, name),
				func() (err error) {
					relations = append(relations, rel)
					return
				}, &rel); e != nil {
				err = e
			} else {
				nouns = append(nouns, Noun{
					name, kind, spec,
					props,
					relations,
				})
			}
			return
		}, &name, &kind, &spec); e != nil {
		err = e
	} else {
		err = nounTemplate.Execute(w, nouns)
	}
	return
}

var nounTemplate = template.Must(template.New("nouns").Funcs(funcMap).Parse(`
<h1>Nouns</h1>
	{{- range $i, $_ := . -}}
{{- if $i -}},{{ end }}
  <a href="#{{.Name|Safe}}">{{.Name|Title}}</a>{{- "" -}}
	{{- end }}.
	{{- range . -}}
{{ "" }}

<h2 id="{{.Name}}">{{.Name|Title}}</h2>
<span>Kind: {{"" -}}	
		{{ if .Kind -}}
<a href="/atlas/kinds#{{.Kind|Safe}}">{{.Kind|Title}}</a>.
		{{- else -}}
none.
		{{- end -}}
</span> {{"" -}}
<span class="spec">{{.Spec}}</span>{{- "" -}}
		{{- if len .Props -}}
{{ "" }}

<h3>Properties</h3>{{- ""}}
<ul>{{- "" -}}
			{{- range $i, $_ := .Props -}}
{{- if $i -}}{{ end }}
  <li>{{.Name|Title}}: <span>{{.Value}}.</span></li>{{- "" -}}
			{{ end }}
</ul>{{- "" -}}
			{{- end }}
		{{- if len .Relations -}}
{{ "" }}

<h3>Relations</h3>{{- "" -}}
			{{- range $i, $_ := .Relations -}}
				{{- if $i }},{{ end }}
  <a href="/atlas/relations/{{.|Safe}}">{{.|Title}}</a>
			{{- end -}}.
		{{- end -}}
	{{- end -}}
`))
