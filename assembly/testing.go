package assembly

import (
	"strings"

	"github.com/ionous/errutil"
)

// Plurals helps to create a model plural table for testing.
type Plurals struct {
	once map[string]bool
}

// Add a singularly named kind
func (pc *Plurals) AddOne(kind string) {
	if len(kind) > 0 && !pc.once[kind] {
		if strings.HasSuffix(kind, "s") {
			panic(errutil.Fmt("expected singular name in test, got %q", kind))
		}
		if pc.once == nil {
			pc.once = make(map[string]bool)
		}
		pc.once[kind] = true
	}
}

// Add a plural named kind
func (pc *Plurals) AddMany(kinds string) {
	if !strings.HasSuffix(kinds, "s") {
		panic(errutil.Fmt("expected Plurals name in test, got %q", kinds))
	}
	pc.AddOne(kinds[:len(kinds)-1])
}

// WritePlurals to the table.
func (pc *Plurals) WritePlurals(m *Assembler) {
	for one, _ := range pc.once {
		m.WritePlural(one, one+"s")
	}
}

// create some fake model hierarchy using a kind and its comma separated ancestors
// ( both using plural names )
func AddTestHierarchy(m *Assembler, els ...string) (err error) {
	var pc Plurals
	for i, cnt := 0, len(els); i < cnt; i += 2 {
		kind, ancestors := els[i], els[i+1]
		pc.AddMany(kind)
		if e := m.WriteAncestor(kind, ancestors); e != nil {
			err = errutil.Append(err, e)
		}
	}
	pc.WritePlurals(m)
	return
}

// mdl_noun:  kind, ret id noun
// mdl_name: noun, name, rank
func AddTestNouns(m *Assembler, els ...string) (err error) {
	for i, cnt := 0, len(els); i < cnt; i += 2 {
		noun, kind := els[i], els[i+1]
		if e := m.WriteNounWithNames(noun, kind); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

// create some fake model hierarchy of kind(s), field, fieldType
func AddTestFields(m *Assembler, els ...string) (err error) {
	for i, cnt := 0, len(els); i < cnt; i += 3 {
		kind, field, fieldType := els[i], els[i+1], els[i+2]
		if e := m.WriteField(kind, field, fieldType); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

// write aspect, trait pairs
func AddTestTraits(m *Assembler, els ...string) (err error) {
	for i, cnt := 0, len(els); i < cnt; i += 2 {
		kind, field := els[i], els[i+1]
		// rank is not set yet, see AssembleAspects
		if e := m.WriteTrait(kind, field, 0); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

// write some noun, field, value interface{}s
func AddTestStarts(m *Assembler, els ...interface{}) (err error) {
	for i, cnt := 0, len(els); i < cnt; i += 3 {
		noun, field, value := els[i], els[i+1], els[i+2]
		if e := m.WriteStart(noun.(string), field.(string), value); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

// AddTestDefaults writes some kind, field, value interface{}s
func AddTestDefaults(m *Assembler, els ...interface{}) (err error) {
	for i, cnt := 0, len(els); i < cnt; i += 3 {
		kind, field, value := els[i], els[i+1], els[i+2]
		if e := m.WriteDefault(kind.(string), field.(string), value); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

// add modeled data: relation, kind, cardinality, otherKind
func AddTestRelations(m *Assembler, els ...string) (err error) {
	for i, cnt := 0, len(els); i < cnt; i += 4 {
		relation, kind, cardinality, otherKind := els[i+0], els[i+1], els[i+2], els[i+3]
		if e := m.WriteRelation(relation, kind, cardinality, otherKind); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

func AddTestVerbs(m *Assembler, els ...string) (err error) {
	for i, cnt := 0, len(els); i < cnt; i += 2 {
		rel, verb := els[i+0], els[i+1]
		if e := m.WriteVerb(rel, verb); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}
