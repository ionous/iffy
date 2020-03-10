package main

import (
	"database/sql"
	"fmt"
	"io"

	"github.com/ionous/iffy/tables"
)

func listOfAspects(w io.Writer, db *sql.DB) (err error) {
	// originally used a channel, but the template iterates over the same elements multiple times
	var aspects []Aspect
	var aname, aspec string

	//var name, kind, spec string
	if e := tables.QueryAll(db, `
		select distinct aspect, coalesce((
			select spec from mdl_spec 
			where type='aspect' and name=aspect
			limit 1), '')
		from mdl_aspect
		order by aspect`,
		func() (err error) {
			// 1. collect traits by name
			var trait Trait
			var traits []Trait
			if e := tables.QueryAll(db, fmt.Sprintf(`
				select trait, coalesce((
					select spec from mdl_spec 
					where type='trait' and name=trait
					limit 1), '')
				from mdl_aspect
				where aspect = '%s'`, aname), func() (err error) {
				traits = append(traits, trait)
				return
			}, &trait.Name, &trait.Spec); e != nil {
				err = e
			} else {
				// 2. collect kinds affected by traits
				var kind string
				var kinds []string
				if e := tables.QueryAll(db, fmt.Sprintf(`
					select kind
					from mdl_field
					where type = 'aspect' 
					and field = '%s'`, aname), func() (err error) {
					kinds = append(kinds, kind)
					return
				}, &kind); e != nil {
					err = e
				} else {
					aspects = append(aspects,
						Aspect{
							Name:   aname,
							Spec:   aspec,
							Kinds:  kinds,
							Traits: traits,
						})
				}
			}
			return
		}, &aname, &aspec); e != nil {
		err = e
	} else {
		err = templates.ExecuteTemplate(w, "aspectList", aspects)
	}
	return
}

type Aspect struct {
	Name, Spec string
	Kinds      []string
	Traits     []Trait
}

type Trait struct {
	Name, Spec string
}

func init() {
	registerTemplate("aspectList", `
<h1>Aspects</h1>
	{{- range $i, $_ := . -}}
{{- if $i -}},{{ end }}
  <a href="#{{.Name|safe}}">{{.Name|title}}</a>{{- "" -}}
	{{- end }}.

	{{- range $i, $_ := . }}
<h2 id="{{.Name}}">{{.Name|title}}</h2>
{{- if .Spec }}
{{.Spec}} 
{{- end -}}

		{{- if len .Kinds }}
<h3>Kinds</h3>{{- ""}}
{{- "" -}}
			{{- range $i, $_ := .Kinds -}}
{{- if $i -}},{{ end }}
  <a href="/atlas/kinds#{{.|safe}}">{{.|title}}</a>
			{{- end }}.
		{{- end -}}

		{{- if len .Traits }}
<h3>Traits</h3>{{- "" }}
<dl>
			{{- range $i, $_ := .Traits }}
  <dt>{{.Name|title}}</dt>
  				{{- if .Spec }}
   <dd>{{.Spec}}</dd>
 				{{- end }}
			{{- end }}
</dl>
		{{- end -}}
	{{- end -}}
`)
}
