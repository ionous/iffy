package main

import (
	"database/sql"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"strings"
	"text/template"

	"github.com/ionous/iffy/dbutil"
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
	if e := dbutil.QueryAll(db, `
		select kind, path, coalesce(spec, '')
		from mdl_kind
		left join mdl_spec 
			on (type='kind' and name=kind)
		order by path, kind`,
		func() (err error) {
			var prop Prop
			var props []Prop
			var noun string
			var nouns []string
			if e := dbutil.QueryAll(db,
				fmt.Sprintf("select field, value, spec from atlas_fields where kind='%s' order by field", kind.Name),
				func() (err error) {
					props = append(props, prop)
					return
				},
				&prop.Name, &prop.Value, &prop.Spec); e != nil {
				err = e
			} else if e := dbutil.QueryAll(db,
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
		err = kindsTemplate.Execute(w, kinds)
	}
	return
}

var spaces = regexp.MustCompile(`\s+`)
var funcMap = template.FuncMap{
	"Title": strings.Title,
	"Safe": func(s string) string {
		return spaces.ReplaceAllString(s, "-")
	},
	// return true if the struct field in els before idx differs from the one at idx
	"changing": func(idx int, field string, els reflect.Value) (ret bool) {
		if idx == 0 {
			ret = true
		} else {
			curr, prev := els.Index(idx), els.Index(idx-1)
			c := curr.Elem().FieldByName(field).Interface()
			p := prev.Elem().FieldByName(field).Interface()
			ret = c != p
		}
		return
	},
}

var kindsTemplate = template.Must(template.New("kinds").Funcs(funcMap).Parse(`
<h1>Kinds</h1>
{{range $i, $_ := .}}
	{{- if $i}}, {{end -}}
	<a href="#{{ .Name|Safe }}">{{.Name|Title}}</a>{{/**/ -}}
{{- end}}.
{{ range . }}
<h2 id="{{.Name}}">{{.Name|Title}}</h2>
<span>Parent kind: {{/**/ -}}	
{{ if .Parent -}}
	<a href="#{{.Parent|Safe}}">{{.Parent|Title}}</a>.
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
	<dt>{{.Name|Title}}: <span>{{.Value}}.</span></dt>
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
	<a href="/atlas/nouns#{{.|Safe}}">{{.|Title}}</a>
	{{- end -}}.
{{end -}}
{{- end -}}
`))
