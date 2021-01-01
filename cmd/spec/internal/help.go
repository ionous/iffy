package internal

import (
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

type tokensink struct {
	els   []string
	roles []rune
}

func (ts *tokensink) add(token string, role rune) {
	ts.els = append(ts.els, token)
	ts.roles = append(ts.roles, role)
}

func parseSpec(t r.Type, fluid *composer.Fluid) ([]string, string, export.Dict) {
	const (
		KEY       = 'K' // token $
		ENUM      = 'E'
		CMD       = 'C'
		FUNC      = 'F'
		SEPARATOR = 'Z'
		TERMINAL  = 'T'
		SELECTOR  = 'S'
		AMBIGUOUS = 'Q' // ambiguous
		// unstyled = "_"
	)
	// fix: uppercase $ parameters mixed with text
	// could possibly get from tags on the original command registration.
	// or could use blank text fields and join in-order
	name := firstRuneLower(t.Name())
	if fluid != nil {
		if len(fluid.Name) > 0 {
			name = fluid.Name
		}
	}
	var tokens tokensink
	params := make(export.Dict)
	commas := " "
	export.WalkProperties(t, func(f *r.StructField, path []int) (done bool) {
		tags := tag.ReadTag(f.Tag)
		if _, ok := tags.Find("internal"); !ok {
			var label string
			if l, ok := tags.Find("label"); ok {
				label = l
			} else {
				label = firstRuneLower(f.Name)
			}
			key := export.Tokenize(f)
			typeName, repeats := nameOfType(f.Type)

			// label this field?
			unlabeled := tags.Exists("unlabeled")

			// if havent written the name, we need to write it first.
			if len(tokens.els) == 0 {
				role := AMBIGUOUS //
				if fluid != nil {
					switch fluid.Role {
					case composer.Command:
						role = CMD
					case composer.Function:
						role = FUNC
					case composer.Selector:
						role = SELECTOR
					}
				}
				tokens.add(name, role)
				// if we arent labeling the first parameter...
				// then follow the name directly by a colon.
				// same for implementations of interfaces which are selectors.
				// ex. (put) intoNumList: or (sort) numbers:
				if unlabeled || fluid != nil && fluid.Role == composer.Selector {
					tokens.add(": ", SEPARATOR)
					commas = ""
				}
			}
			if len(commas) > 0 {
				tokens.add(commas, SEPARATOR)
			}
			// note: fluid interfaces might defer their selectors to their children
			if !unlabeled {
				tokens.add(label, SELECTOR)
				tokens.add(": ", SEPARATOR)
			}
			tokens.add(key, KEY)
			m := export.Dict{
				"label": label,
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
	// no fields?
	if len(tokens.els) == 0 {
		tokens.add(name, ENUM)
	} else if tokens.roles[0] == CMD {
		tokens.add(".", TERMINAL)
	}
	return tokens.els, string(tokens.roles), params
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
			if slotName := findTypeName(slot); len(slotName) == 0 {
				log.Panic("unknown slot", slot)
			} else {
				ret = append(ret, slotName)
			}
		}
	}
	return
}

func nameOfType(t r.Type) (typeName string, repeats bool) {
	switch kind := t.Kind(); kind {
	case r.Bool:
		if name := findTypeName(t); len(name) > 0 {
			typeName = name
		} else {
			typeName = "bool"
		}
	case r.Float32, r.Float64, r.Int, r.Int8, r.Int16, r.Int32, r.Int64:
		// some sort of type  hint? ex. possibly link to custom types
		typeName = "number"
	case r.Slice:
		typeName, _ = nameOfType(t.Elem())
		repeats = true
	case r.Interface: // a reference to another type
		if name := findTypeName(t); len(name) == 0 {
			log.Panicln("unknown interface", t)
		} else {
			typeName = name
		}
	case r.String:
		if name := findTypeName(t); len(name) > 0 {
			typeName = name
		} else {
			typeName = "text"
		}
	case r.Struct:
		if name := findTypeName(t); len(name) == 0 {
			typeName = lang.Underscore(t.Name())
		} else {
			typeName = name
		}
	case r.Ptr:
		typeName, repeats = nameOfType(t.Elem())
	default:
		log.Panicln("unhandled type", kind, t.String())
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
	return reverseLookup[t]
}
