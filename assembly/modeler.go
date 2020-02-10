package assembly

import (
	"database/sql"
	"strings"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/tables"
	"github.com/reiver/go-porterstemmer"
)

func NewModelerDB(db *sql.DB) *Modeler {
	dbq := ephemera.NewDBQueue(db)
	if e := tables.CreateAssembly(db); e != nil {
		panic(e)
	}
	return NewModeler(dbq)
}

func cat(str ...string) string {
	return strings.Join(str, " ")
}

func NewModeler(q *ephemera.DbQueue) *Modeler {
	q.PrepCols("mdl_kind", []ephemera.Col{
		{Name: "kind", Type: "text"},
		{Name: "path", Type: "text"},
		{Check: "primary key(kind)"},
	})
	q.PrepCols("mdl_aspect", []ephemera.Col{
		{Name: "aspect", Type: "text"},
		{Name: "trait", Type: "text"},
		{Name: "rank", Type: "int"},
		{Check: "primary key(aspect, trait)"},
	})
	q.PrepCols("mdl_rel", []ephemera.Col{
		{Name: "relation", Type: "text"},
		{Name: "kind", Type: "text"},        /* reference to mdl_kind */
		{Name: "cardinality", Type: "text"}, /* one of MANY/ONE */
		{Name: "otherKind", Type: "text"},   /* reference to mdl_kind */
		{Check: "primary key(relation)"},
	})
	q.PrepCols("mdl_field", []ephemera.Col{
		{Name: "kind", Type: "text"}, /* reference to mdl_kind */
		{Name: "field", Type: "text"},
		{Name: "type", Type: "text"}, /* one of PRIM_type */
		{Check: "primary key(kind, field)"},
	})
	q.PrepCols("mdl_default", []ephemera.Col{
		{Name: "kind", Type: "text"},  /* reference to mdl_kind */
		{Name: "field", Type: "text"}, /* partial reference to mdl_field */
		{Name: "value", Type: "blob"},
	})
	q.PrepCols("mdl_noun", []ephemera.Col{
		{Name: "noun", Type: "text"},
		{Name: "kind", Type: "text"}, /* reference to mdl_kind */
		{Check: "primary key(noun)"},
	})
	// names are built from noun parts, and possibly from custom aliases.
	// where rank 0 is a better match than rank 1
	q.PrepCols("mdl_name", []ephemera.Col{
		{Name: "noun", Type: "text"}, /* reference to mdl_noun */
		{Name: "name", Type: "text"},
		{Name: "rank", Type: "int"},
	})
	// insert without duplication of the relation, stem pair
	q.PrepStatement("asm_verb",
		`insert into asm_verb(relation, stem)
				select ?1, ?2
				where not exists (
					select 1 from asm_verb v
					where v.relation=?1 and v.stem=?2
				)`, []ephemera.Col{
			{Name: "relation", Type: "text"}, /* reference to mdl_rel */
			{Name: "stem", Type: "text"},
			{Check: "unique(stem)"},
		})
	q.PrepCols("mdl_start", []ephemera.Col{
		{Name: "noun", Type: "text"},  /* reference to mdl_noun */
		{Name: "field", Type: "text"}, /* partial reference to mdl_field */
		{Name: "value", Type: "blob"},
	})
	q.PrepCols("mdl_pair", []ephemera.Col{
		{Name: "noun", Type: "text"},      /* reference to mdl_noun */
		{Name: "relation", Type: "text"},  /* reference to mdl_rel */
		{Name: "otherNoun", Type: "text"}, /* reference to mdl_noun */
	})
	//
	return &Modeler{q}
}

type Modeler struct {
	q ephemera.Queue
}

// write kind and comma separated ancestors
func (m *Modeler) WriteAncestor(kind, path string) (err error) {
	_, e := m.q.Write("mdl_kind", kind, path)
	return e
}

func (m *Modeler) WriteField(kind, field, fieldType string) error {
	_, e := m.q.Write("mdl_field", kind, field, fieldType)
	return e
}

// WriteDefault: if no specific value has been assigned to the an instance of the idModelField's kind,
// the passed default value will be used for that instance's kind.
func (m *Modeler) WriteDefault(kind, field string, value interface{}) error {
	_, e := m.q.Write("mdl_default", kind, field, value)
	return e
}

func (m *Modeler) WriteNoun(noun, kind string) error {
	_, e := m.q.Write("mdl_noun", noun, kind)
	return e
}

// WriteName for noun
func (m *Modeler) WriteName(noun, name string, rank int) error {
	_, e := m.q.Write("mdl_name", noun, name, rank)
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
	_, e := m.q.Write("mdl_rel", relation, kind, cardinality, otherKind)
	return e
}

func (m *Modeler) WriteTrait(aspect, trait string, rank int) error {
	_, e := m.q.Write("mdl_aspect", aspect, trait, rank)
	return e
}

// WriteValue: store the initial value of an instance's field used at start of play.
func (m *Modeler) WriteValue(noun, field string, value interface{}) error {
	_, e := m.q.Write("mdl_start", noun, field, value)
	return e
}

func (m *Modeler) WriteVerb(relation, verb string) error {
	stem := porterstemmer.StemString(verb)
	_, e := m.q.Write("asm_verb", relation, stem)
	return e
}
