package main

import (
	"database/sql"
	"fmt"
	"io"
	"strings"

	"github.com/ionous/iffy/tables"
)

// some things we could do:
// . an alphabetical index
// . a hierarchical, indented listing
// . headings per kind
//
type Kind struct {
	Name, Path, Spec string
	Props            []Prop // newly introduced fields, or fields with new default values
	Nouns            []string
}

type Prop struct {
	Name, Value string
	Spec        sql.NullString
}

func (k Kind) Parent() string {
	return strings.Split(k.Path, ",")[0]
}

func listOfKinds(w io.Writer, db *sql.DB) (err error) {
	// originally used a channel, but the template iterates over the same elements multiple times
	var kind Kind
	var kinds []Kind
	if e := tables.QueryAll(db, `
		select kind, path, coalesce((
			select spec from mdl_spec 
			where type='kind' and name=kind
			limit 1), '')
		from mdl_kind
		order by path, kind`,
		func() (err error) {
			var prop Prop
			var props []Prop
			var noun string
			var nouns []string
			if e := tables.QueryAll(db,
				fmt.Sprintf("select field, value, spec from atlas_fields where kind='%s' order by field", kind.Name),
				func() (err error) {
					props = append(props, prop)
					return
				},
				&prop.Name, &prop.Value, &prop.Spec); e != nil {
				err = e
			} else if e := tables.QueryAll(db,
				fmt.Sprintf("select noun from mdl_noun where kind='%s' order by noun", kind.Name),
				func() (err error) {
					nouns = append(nouns, noun)
					return
				}, &noun); e != nil {
				err = e
			} else {
				kind.Props = props
				kind.Nouns = nouns
				kinds = append(kinds, kind)
			}
			return
		}, &kind.Name, &kind.Path, &kind.Spec); e != nil {
		err = e
	} else {
		err = templates.ExecuteTemplate(w, "kindList", kinds)
	}
	return
}

func init() {
	registerTemplate("kindList", `
<h1>Kinds</h1>
{{range $i, $_ := .}}
	{{- if $i}}, {{end -}}
	<a href="#{{ .Name|safe }}">{{.Name|title}}</a>{{/**/ -}}
{{- end}}.
{{ range . }}
<h2 id="{{.Name}}">{{.Name|title}}</h2>
<span>Parent kind: {{/**/ -}}  
{{ if .Parent -}}
	<a href="#{{.Parent|safe}}">{{.Parent|title}}</a>.
{{- else -}}
	none.
{{- end -}}
</span> {{/**/ -}}
<span class="spec">{{.Spec}}</span>
{{- /**/ -}}
{{- if len .Props }}

<h3>Properties</h3>{{- /**/}}
<dl>{{- /**/ -}}
{{- range .Props }}
	<dt>{{.Name|title}}: <span>{{.Value}}.</span></dt>
	{{- if .Spec.Valid -}}
	<dd>{{ .Spec.String }}</dd>
	{{- end -}}
{{ end }}
</dl>{{- /**/ -}}
{{- end }}
{{- /*
     */ -}}
{{- if len .Nouns }}

<h3>Nouns</h3>{{- /**/ -}}
{{ range $i, $_ := .Nouns }}
	{{- if $i }},{{ end }}
	<a href="/atlas/nouns#{{.|safe}}">{{.|title}}</a>
	{{- end -}}.
{{end -}}
{{- end -}}
`)
}
