// Exports golang DSL for use in editing story files.
// Currently, this only generates the imperative commands,
// the modeling parts of the language currently live in the composer javascript
package main

import (
	"encoding/json"
	"fmt"
	r "reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/export"
	"github.com/ionous/iffy/export/tag"
)

// go run spec.go > /Users/ionous/Dev/go/src/github.com/ionous/iffy/cmd/compose/www/src/lang/spec.js
func main() {
	var all []export.Dict
	var slots []r.Type

	groups := make(Groups)

	// slots that commands can fit into
	for _, slot := range export.Slots {
		spec := getSpec(slot.Type)
		if spec.Group != "internal" {
			i := r.TypeOf(slot.Type).Elem()
			//
			slots = append(slots, i)
			if len(spec.Name) == 0 {
				spec.Name = slot.Name
			}
			if len(spec.Desc) == 0 {
				spec.Desc = export.Prettify(slot.Name)
			}
			if len(spec.Group) == 0 {
				spec.Group = slot.Group
			}
			out := export.Dict{
				"name": spec.Name,
				"desc": spec.Desc,
				"uses": "slot",
			}
			addDesc(out, slot.Desc)
			groups.addGroup(out, spec.Group)
			all = append(all, out)
		}
	}

	for _, cmd := range export.Slats {
		spec := cmd.Compose()
		if spec.Group != "internal" {
			rtype := r.TypeOf(cmd).Elem()
			if len(spec.Name) == 0 {
				panic(fmt.Sprintln("missing name for type", rtype.Name()))
			}
			//
			with := make(export.Dict)
			if slotNames := slotsOf(rtype, slots); len(slotNames) > 0 {
				with["slots"] = slotNames
			}
			out := export.Dict{
				"name": spec.Name,
				"uses": "run",
				"with": with,
			}
			// missing spec, missing slots.
			if len(spec.Spec) != 0 {
				out["spec"] = spec.Spec
			} else {
				tokens, params := parse(rtype)
				with["params"] = params
				with["tokens"] = updateTokens(spec.Spec, tokens)
			}
			addDesc(out, spec.Desc)
			groups.addGroup(out, spec.Group)
			all = append(all, out)
		}
	}

	for groupName, _ := range groups {
		out := export.Dict{
			"name": groupName,
			"uses": "group",
		}
		all = append(all, out)
	}

	sort.Slice(all, func(i, j int) (ret bool) {
		uses := strings.Compare(all[i]["uses"].(string), all[j]["uses"].(string))
		switch uses {
		case 0:
			ret = all[i]["name"].(string) < all[j]["name"].(string)
		case -1:
			ret = false
		case 1:
			ret = true
		}
		return

	})

	if b, e := json.MarshalIndent(all, "", "  "); e != nil {
		panic(b)
	} else {
		fmt.Println("/* generated using github.com/ionous/iffy/cmd/spec/spec.go */")
		fmt.Print("const spec = ")
		fmt.Println(string(b))
	}
}

func getSpec(ptrValue interface{}) (ret composer.Spec) {
	if c, ok := ptrValue.(composer.Specification); ok {
		ret = c.Compose()
	}
	return
}

var specType = r.TypeOf((*composer.Spec)(nil)).Elem()

func parse(t r.Type) ([]string, export.Dict) {
	// fix: uppercase $ parameters mixed with text
	// could possibly get from tags on the original command registration.
	// or could use blank text fields and join in-order
	prettyType := export.Prettify(t.Name())

	tokens := []string{prettyType}
	// keyed by token
	params := make(export.Dict)

	export.WalkProperties(t, func(f *r.StructField, path []int) (done bool) {
		tags := tag.ReadTag(f.Tag)
		if _, ok := tags.Find("internal"); !ok {
			prettyField := export.Prettify(f.Name)
			key := export.Tokenize(f)
			typeName, repeats := nameOfType(f.Type)
			tokens = append(tokens, key)
			m := export.Dict{
				"label": prettyField,
				"type":  typeName,
				// optional: tdb
			}
			if repeats {
				m["repeats"] = true
			}
			params[key] = m
		}
		return
	})
	return tokens, params
}

var tokenPlaceholders = regexp.MustCompile(`^\$([0-9]+)$`)

// i dont think this works right because it doesnt change the spec
// i dont think? the composer is either.
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

func (g *Groups) addGroup(out export.Dict, group string) {
	// even no commas results in one group;
	// ideally, id think an empty string would be no groups.... but alas.
	if len(group) > 0 {
		if groups := strings.Split(group, ","); len(groups) > 0 {
			for i, group := range groups {
				(*g)[group] = true
				groups[i] = strings.ToLower(group)
			}
			sort.Strings(groups)
			out["group"] = groups
		}
	}
}

func addDesc(out export.Dict, desc string) {
	if len(desc) > 0 {
		out["desc"] = desc
	}
}

func slotsOf(slat r.Type, slots []r.Type) (ret []string) {
	ptrType := r.PtrTo(slat)
	for _, slot := range slots {
		if ptrType.Implements(slot) {
			slotName := findTypeName(slot)
			ret = append(ret, slotName)
		}
	}
	return
}

func nameOfType(t r.Type) (typeName string, repeats bool) {
	switch kind := t.Kind(); kind {
	case r.Bool:
		typeName = "bool"
	case r.Float32, r.Float64, r.Int, r.Int8, r.Int16, r.Int32, r.Int64:
		// some sort of type  hint? ex. possibly link to custom types
		typeName = "number"
	case r.Slice:
		typeName, _ = nameOfType(t.Elem())
		repeats = true
	case r.Interface: // a reference to another type
		typeName = findTypeName(t)
	case r.String:
		typeName = "text"
	default:
		if kind == r.Ptr && t.Elem().Kind() == r.Struct {
			typeName = findTypeName(t.Elem())
		} else {
			// Array, Map, Ptr, Struct
			// Uint, Uint8, Uint16, Uint32, Uint64
			panic(fmt.Sprintln("unhandled type", t.String()))
		}
	}
	return
}

var reverseLookup map[r.Type]string

func findTypeName(t r.Type) (ret string) {
	if len(reverseLookup) == 0 {
		reverseLookup = make(map[r.Type]string)
		for _, cmd := range export.Slats {
			runType := r.TypeOf(cmd).Elem()
			typeName := cmd.Compose().Name
			reverseLookup[runType] = typeName
		}
		for _, slot := range export.Slots {
			t := r.TypeOf(slot.Type).Elem()
			reverseLookup[t] = slot.Name
		}
	}

	if n, ok := reverseLookup[t]; ok {
		ret = n
	} else {
		panic(fmt.Sprintln("unknown type", t.String()))
	}
	return
}
