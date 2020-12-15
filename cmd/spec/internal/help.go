package internal

import (
	"fmt"
	r "reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/ionous/iffy"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/export"
	"github.com/ionous/iffy/export/tag"
)

var tokenPlaceholders = regexp.MustCompile(`^\$([0-9]+)$`)

func getSpec(ptrValue interface{}) (ret composer.Spec) {
	if c, ok := ptrValue.(composer.Composer); ok {
		ret = c.Compose()
	}
	return
}

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
		for _, slats := range iffy.AllSlats {
			for _, cmd := range slats {
				runType := r.TypeOf(cmd).Elem()
				typeName := cmd.Compose().Name
				reverseLookup[runType] = typeName
			}
		}
		for _, slots := range iffy.AllSlots {
			for _, slot := range slots {
				t := r.TypeOf(slot.Type).Elem()
				reverseLookup[t] = slot.Name
			}
		}
	}

	if n, ok := reverseLookup[t]; ok {
		ret = n
	} else {
		panic(fmt.Sprintln("unknown type", t.String()))
	}
	return
}
