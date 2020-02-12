package main

import (
	"html/template"
	"reflect"
	"regexp"
	"strings"
)

var spaces = regexp.MustCompile(`\s+`)

var funcMap = template.FuncMap{
	"title": strings.Title,
	"safe": func(s string) string {
		return spaces.ReplaceAllString(s, "-")
	},
	"prefix": strings.HasPrefix,
	"suffix": strings.HasSuffix,
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

var templates *template.Template = template.New("none").Funcs(funcMap)

func registerTemplate(n, t string) {
	templates = template.Must(templates.New(n).Parse(t))
}
