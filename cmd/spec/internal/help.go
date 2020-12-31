package internal

import (
	"fmt"
	"log"
	r "reflect"
	"regexp"
	"strings"
	"unicode"

	"github.com/ionous/iffy"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/export"
	"github.com/ionous/iffy/export/tag"
	"github.com/ionous/iffy/lang"
)

var tokenPlaceholders = regexp.MustCompile(`^\$([0-9]+)$`)

func firstRuneLower(s string) (ret string) {
	rs := []rune(s)
	rs[0] = unicode.ToLower(rs[0])
	return string(rs)
}

func parseSpec(t r.Type, fluid *composer.Fluency) ([]string, export.Dict) {
	// fix: uppercase $ parameters mixed with text
	// could possibly get from tags on the original command registration.
	// or could use blank text fields and join in-order
	name := firstRuneLower(t.Name())
	if fluid != nil {
		if len(fluid.Name) > 0 {
			name = fluid.Name
		} /*else if syl := strings.IndexFunc(name, func(u rune) bool {
			return unicode.IsUpper(u)
		}); syl > 0 {
			name = name[:syl]
		}*/
		if fluid.RunIn {
			name += ":"
		}
	}

	tokens := []string{name}
	// keyed by token
	params := make(export.Dict)
	commas := " "
	export.WalkProperties(t, func(f *r.StructField, path []int) (done bool) {
		tags := tag.ReadTag(f.Tag)
		if _, ok := tags.Find("internal"); !ok {
			prettyField := firstRuneLower(f.Name)
			key := export.Tokenize(f)
			typeName, repeats := nameOfType(f.Type)
			if fluid == nil || !fluid.RunIn || (len(params) > 0 && f.Type.Kind() != r.Interface) {
				tokens = append(tokens, commas+prettyField+": ", key)
			} else {
				tokens = append(tokens, commas)
				tokens = append(tokens, key)
			}
			m := export.Dict{
				"label": prettyField,
				"type":  typeName,
				// optional: tdb
			}
			if repeats {
				m["repeats"] = true
			}
			params[key] = m
			commas = ", "
		}
		return
	})
	return tokens, params
}

func updateTokens(phrase string, tokens []string) (ret []string) {
	if len(phrase) == 0 {
		ret = tokens
	} else {
		fields := strings.Fields(phrase)
		// replace the fields in a phrase matching $1, etc. with tokens
		// unused, and would need to handle offset of fields in tokens
		// for _, f := range fields {
		// 	if !tokenPlaceholders.MatchString(f) {
		// 		ret = append(ret, f)
		// 	} else if i, e := strconv.Atoi(f[1:]); e != nil {
		// 		panic(e)
		// 	} else {
		// 		t := tokens[i]
		// 		ret = append(ret, t)
		// 	}
		// }
		ret = fields
	}
	return
}

func addDesc(out export.Dict, name, desc string) {
	if len(desc) > 0 {
		if strings.Index(desc, ":") < 0 {
			desc = name + ": " + desc
		}
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

func typeName(t r.Type, name string) (ret string) {
	if len(name) > 0 {
		ret = name
	} else {
		ret = lang.Underscore(t.Name())
	}
	return
}

func findTypeName(t r.Type) (ret string) {
	if len(reverseLookup) == 0 {
		reverseLookup = make(map[r.Type]string)
		for _, slats := range iffy.AllSlats {
			for _, cmd := range slats {
				runType := r.TypeOf(cmd).Elem()
				reverseLookup[runType] = typeName(runType, cmd.Compose().Name)
			}
		}
		for _, slots := range iffy.AllSlots {
			for _, slot := range slots {
				t := r.TypeOf(slot.Type).Elem()
				reverseLookup[t] = typeName(t, slot.Name)
			}
		}
	}

	if n, ok := reverseLookup[t]; ok {
		ret = n
	} else {
		log.Panic("unknown type", t.String())
	}
	return
}
