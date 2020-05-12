package assembly

import (
	"database/sql"
	"strings"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/tables"
	"github.com/reiver/go-porterstemmer"
)

func cat(str ...string) string {
	return strings.Join(str, " ")
}

func NewModeler(db *sql.DB) *Modeler {
	return &Modeler{tables.NewCache(db)}
}

type Modeler struct {
	cache *tables.Cache
}

// write kind and comma separated ancestors
func (m *Modeler) WriteAncestor(kind, path string) (err error) {
	_, e := m.cache.Exec(mdl_kind, kind, path)
	return e
}

func (m *Modeler) WriteField(kind, field, fieldType string) error {
	_, e := m.cache.Exec(mdl_field, kind, field, fieldType)
	return e
}

// WriteDefault: if no specific value has been assigned to the an instance of the idModelField's kind,
// the passed default value will be used for that instance's kind.
func (m *Modeler) WriteDefault(kind, field string, value interface{}) error {
	_, e := m.cache.Exec(mdl_default, kind, field, value)
	return e
}

func (m *Modeler) WriteNoun(noun, kind string) error {
	_, e := m.cache.Exec(mdl_noun, noun, kind)
	return e
}

// WriteName for noun
func (m *Modeler) WriteName(noun, name string, rank int) error {
	_, e := m.cache.Exec(mdl_name, noun, name, rank)
	return e
}

// WriteNounWithNames
func (m *Modeler) WriteNounWithNames(noun, kind string) (err error) {
	if e := m.WriteNoun(noun, kind); e != nil {
		err = errutil.Append(err, e)
	} else if e := m.WriteName(noun, noun, 0); e != nil {
		err = errutil.Append(err, e)
	} else {
		split := strings.Fields(noun)
		if cnt := len(split); cnt > 1 {
			for i, k := range split {
				rank := cnt - i
				if e := m.WriteName(noun, k, rank); e != nil {
					err = errutil.Append(err, e)
				}
			}
		}
	}
	return
}

func (m *Modeler) WriteRelation(relation, kind, cardinality, otherKind string) error {
	_, e := m.cache.Exec(mdl_rel, relation, kind, cardinality, otherKind)
	return e
}

func (m *Modeler) WriteTrait(aspect, trait string, rank int) error {
	_, e := m.cache.Exec(mdl_aspect, aspect, trait, rank)
	return e
}

// WriteValue: store the initial value of an instance's field used at start of play.
func (m *Modeler) WriteValue(noun, field string, value interface{}) error {
	_, e := m.cache.Exec(mdl_start, noun, field, value)
	return e
}

func (m *Modeler) WriteVerb(relation, verb string) error {
	const asm_verb = `insert into asm_verb(relation, stem)
				select ?1, ?2
				where not exists (
					select 1 from asm_verb v
					where v.relation=?1 and v.stem=?2
				)`
	stem := porterstemmer.StemString(verb)
	_, e := m.cache.Exec(asm_verb, relation, stem)
	return e
}

func (m *Modeler) WriteRule(name string, prog int64) error {
	_, e := m.cache.Exec(mdl_rule, name, prog)
	return e
}

func (m *Modeler) WriteProg(typeName string, bytes []byte) (int64, error) {
	return m.cache.Exec(mdl_prog, typeName, bytes)
}

var mdl_spec = tables.Insert("mdl_spec", "type", "name", "spec")
var mdl_aspect = tables.Insert("mdl_aspect", "aspect", "trait", "rank")
var mdl_kind = tables.Insert("mdl_kind", "kind", "path")
var mdl_field = tables.Insert("mdl_field", "kind", "field", "type")
var mdl_default = tables.Insert("mdl_default", "kind", "field", "value")
var mdl_rel = tables.Insert("mdl_rel", "relation", "kind", "cardinality", "otherKind")
var mdl_noun = tables.Insert("mdl_noun", "noun", "kind")
var mdl_name = tables.Insert("mdl_name", "noun", "name", "rank")
var mdl_pair = tables.Insert("mdl_pair", "noun", "relation", "otherNoun")
var mdl_start = tables.Insert("mdl_start", "noun", "field", "value")
var mdl_prog = tables.Insert("mdl_prog", "type", "bytes")
var mdl_rule = tables.Insert("mdl_rule", "pattern", "idProg")
