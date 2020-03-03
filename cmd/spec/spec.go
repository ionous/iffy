// Package main for 'spec".
// Exports golang DSL for use in editing story files.
// Currently, this only generates imperative commands.
package main

import (
	"encoding/json"
	"fmt"
	r "reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"bitbucket.org/pkg/inflect"

	"github.com/ionous/iffy/assembly"
	"github.com/ionous/iffy/ref/unique"
)

// spec types:
// “run”: command structs
// “opt”: would take some custom work to generate. ex. "{proper_noun} or {common_noun}"
// “slot”: eval slots, could either do this manually or generate.
// “str”: handwaves. for enums.
// “num”, “txt”: might? be useful for typedefs or triggered with tags.
//
func main() {

	var all []dict
	var slots []r.Type

	groups := make(Groups)

	// slots that commands can fit into
	for _, slot := range assembly.Slots {
		i := r.TypeOf(slot.Type).Elem()
		slots = append(slots, i)
		name := nameOfType(i)
		desc := prettyName(name)
		uses := "slot"
		out := dict{
			"name": name,
			"desc": desc,
			"uses": uses,
		}
		groups.addGroup(out, slot.Group)
		addDesc(out, slot.Desc)
		all = append(all, out)
	}

	for _, run := range assembly.Runs {
		if run.Type == nil {
			continue
		}
		t := r.TypeOf(run.Type).Elem()
		slotNames := slotsOf(t, slots)
		if len(slotNames) == 0 {
			panic(fmt.Sprintln("missing slot for type", t.Name()))
		}
		typeName := nameOfType(t)
		tokens, params := parse(t)

		with := dict{
			"slots": slotNames,
		}
		out := dict{
			"name": typeName,
			"uses": "run",
			"with": with,
		}
		if !addSpec(out, t) {
			with["params"] = params
			with["tokens"] = updateTokens(run.Phrase, tokens)
		}
		groups.addGroup(out, run.Group)
		addDesc(out, run.Desc)

		all = append(all, out)
	}

	for groupName, _ := range groups {
		out := dict{
			"name": groupName,
			"uses": "group",
		}
		all = append(all, out)
	}

	if b, e := json.MarshalIndent(all, "", "  "); e != nil {
		panic(b)
	} else {
		fmt.Println("/* generated using github.com/ionous/iffy/cmd/spec/spec.go */")
		fmt.Print("const spec = ")
		fmt.Println(string(b))
	}
}

func addSpec(out dict, rtype r.Type) (ret bool) {
	if spec := spec(rtype); len(spec) > 0 {
		out["spec"] = spec
		ret = true
	}
	return
}

func spec(rtype r.Type) (ret string) {
	unique.WalkProperties(rtype,
		func(f *r.StructField, path []int) (done bool) {
			t := unique.Tag(f.Tag)
			ret, done = t.Find("spec")
			return
		})
	return
}

func parse(t r.Type) ([]string, dict) {
	// fix: uppercase $ parameters mixed with text
	// could possibly get from tags on the original command registration.
	// or could use blank text fields and join in-order
	prettyType := prettyName(t.Name())

	tokens := []string{prettyType}
	// keyed by token
	params := make(dict)

	unique.WalkProperties(t, func(f *r.StructField, path []int) (done bool) {
		prettyField := prettyName(f.Name)
		key := nameOfAttr(f)

		var typeString string
		var repeats bool
		switch kind := f.Type.Kind(); kind {
		case r.Bool:
			typeString = "bool"
		case r.Float32, r.Float64, r.Int, r.Int8, r.Int16, r.Int32, r.Int64:
			// some sort of type  hint? ex. possibly link to custom types
			typeString = "number"

		case r.Slice:
			typeString = nameOfType(f.Type.Elem())
			repeats = true
		case r.Interface: // a reference to another type
			typeString = nameOfType(f.Type)
		case r.String:
			typeString = "text"
		default:

			// Array, Map, Ptr, Struct
			// Uint, Uint8, Uint16, Uint32, Uint64
			panic(fmt.Sprintln("unhandled type", t.Name(), f.Name, kind.String()))
		}
		tokens = append(tokens, key)
		m := dict{
			"label": prettyField,
			"type":  typeString,
			// optional: tdb
		}
		if repeats {
			m["repeats"] = true
		}
		params[key] = m
		return
	})
	return tokens, params
}

var tokenPlaceholders = regexp.MustCompile(`^\$([0-9]+)$`)

type dict map[string]interface{}

func updateTokens(phrase string, tokens []string) (ret []string) {
	if len(phrase) == 0 {
		ret = tokens
	} else {
		fields := strings.Fields(phrase)
		for j, f := range fields {
			if tokenPlaceholders.MatchString(f) {
				if i, e := strconv.Atoi(f[1:]); e != nil {
					panic(e)
				} else {
					t := tokens[i]
					fields[j] = t
				}
			}
		}
		ret = fields
	}
	return
}

type Groups map[string]bool

func (g *Groups) addGroup(out dict, group string) {
	// even no commas results in one group;
	// ideally, id think an empty string would be no groups.... but alas.
	if len(group) > 0 {
		if groups := strings.Split(group, ","); len(groups) > 0 {
			for i, group := range groups {
				(*g)[group] = true
				groups[i] = strings.ToLower(group)
			}
			out["group"] = groups
		}
	}
}

func addDesc(out dict, desc string) {
	if len(desc) > 0 {
		out["desc"] = desc
	}
}

func slotsOf(slat r.Type, slots []r.Type) (ret []string) {
	ptrType := r.PtrTo(slat)
	for _, slot := range slots {
		if ptrType.Implements(slot) {
			n := nameOfType(slot)
			ret = append(ret, n)
		}
	}
	return
}

func nameOfType(t r.Type) string {
	return inflect.Underscore(t.Name())
}

func prettyName(n string) string {
	return strings.ToLower(inflect.Humanize(n))
}

func nameOfAttr(f *r.StructField) string {
	return "$" + strings.Map(func(c rune) (ret rune) {
		if c == ' ' {
			ret = '_'
		} else {
			ret = unicode.ToUpper(c)
		}
		return
	}, prettyName(f.Name))
}
