package assembly

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ephemera"
	"github.com/reiver/go-porterstemmer"
)

func NewModelerDB(db *sql.DB) *Modeler {
	dbq := ephemera.NewDBQueue(db)

	// asm_value: view for resolving value ephemera back to strings
	if _, e := db.Exec(
		`create temp view
	asm_value as
		select pv.rowid as idEphValue, nn.name, np.name as prop, pv.value
	from eph_value pv join eph_named nn
		on (pv.idNamedNoun = nn.rowid)
	left join eph_named np
		on (pv.idNamedProp = np.rowid)`); e != nil {
		panic(e)
	}

	// asm_default: view for resolving initial default ephemera back to strings
	if _, e := db.Exec(`create temp view
	asm_default as
		select p.rowid as idEphDefault, nk.name as kind, nf.name as field, p.value as value
	from eph_default p join eph_named nk
		on (p.idNamedKind = nk.rowid)
	left join eph_named nf
 		on (p.idNamedField = nf.rowid)`); e != nil {
		panic(e)
	}

	// asm_default_tree: view for mapping ephemera defaults to mdl_field
	if _, e := db.Exec(
		join(`create temp view 
			asm_default_tree as`,
			searchForFieldsInKind(
				"idEphDefault",       // output column
				"asm_default as src", // asm_default is a view of eph_default with names resolved to strings
				"src.field",          // matched against mdl_field.field
				"src.idEphDefault",   // pull one column up through the hierarchy
				"src.kind",           // the initial kind of the hierarchy ( for every row in src )
			)),
	); e != nil {
		panic(e)
	}

	// asm_value_tree: view for mapping ephemera values to mdl_field
	if _, e := db.Exec(
		join(`create temp view 
			asm_value_tree as`,
			searchForFieldsInKind(
				"idEphValue, noun",        // output columns
				"asm_value as src",        // asm_value is a view of eph_value with names resolved to strings
				"src.prop",                // matched against mdl_field.field
				"src.idEphValue, mn.noun", // pull two columns up through the hierarchy
				// for every named noun or partial noun in src,
				// find the best noun
				// and use that noun's kind as the seed of the hierarchical search.
				`mn.kind join mdl_noun mn /* kind, noun */
					on mn.noun = ( 
					/* find the best noun */
						select m.noun 
						from mdl_name m   /* noun, name, rank */
						where m.name = src.name
						order by m.rank
						limit 1 
					)`),
			`where idModelField != 0 /* we dont need missing fields */`),
	); e != nil {
		panic(e)
	}
	return NewModeler(dbq)
}

func searchForFieldsInKind(correlationOut, srcTable, fieldName, correlatedCols, mkMatch string) string {
	// to get the idModelField as output, we match field ( and kind ) by name
	return fmt.Sprintf(`
	with tree(kind, path, %[5]s, field, idModelField) as
	(select mk.kind, mk.path, %[2]s, %[3]s,
	/* for each implied kind, look for the requested field */
		( select m.rowid from mdl_field m
			where m.kind = mk.kind
			and m.field = %[3]s
		) as idModelField
 	/* for each row in source, look for the implied kind */
	    from %[1]s           /* eph id, name, prop, value */
		join mdl_kind mk
		on mk.kind = %[4]s
	union all
	/* for each parent kind, keep looking for the requested field */
		select super.kind, super.path, %[5]s, tree.field,
			( select m.rowid from mdl_field m
				where m.kind = super.kind
				and m.field = tree.field
			 ) as idModelField
		from tree, mdl_kind super
		where idModelField is null /* keep going til found  */
		and super.kind = substr(tree.path,0,instr(tree.path||",", ","))
	)
	select %[5]s, coalesce(idModelField,0) as idModelField from tree`,
		srcTable, correlatedCols, fieldName, mkMatch, correlationOut)
}

func join(str ...string) string {
	return strings.Join(str, " ")
}

func NewModeler(q ephemera.Queue) *Modeler {
	q.Prep("mdl_kind",
		ephemera.Col{Name: "kind", Type: "text"},
		ephemera.Col{Name: "path", Type: "text"},
		ephemera.Col{Check: "primary key(kind)"},
	)
	q.Prep("mdl_aspect",
		ephemera.Col{Name: "aspect", Type: "text"},
		ephemera.Col{Name: "trait", Type: "text"},
		ephemera.Col{Name: "rank", Type: "int"},
		ephemera.Col{Check: "primary key(aspect, trait)"},
	)
	q.Prep("mdl_rel",
		ephemera.Col{Name: "relation", Type: "text"},
		ephemera.Col{Name: "kind", Type: "text"},
		ephemera.Col{Name: "cardinality", Type: "text"},
		ephemera.Col{Name: "otherKind", Type: "text"},
		ephemera.Col{Check: "primary key(relation)"},
	)
	q.Prep("mdl_field",
		ephemera.Col{Name: "field", Type: "text"},
		ephemera.Col{Name: "kind", Type: "text"},
		ephemera.Col{Name: "type", Type: "text"},
		ephemera.Col{Check: "primary key(kind, field)"},
	)
	q.Prep("mdl_default",
		ephemera.Col{Name: "kind", Type: "text"},
		ephemera.Col{Name: "field", Type: "text"},
		ephemera.Col{Name: "value", Type: "blob"},
	)
	q.Prep("mdl_noun",
		ephemera.Col{Name: "noun", Type: "text"},
		ephemera.Col{Name: "kind", Type: "text"},
		ephemera.Col{Check: "primary key(noun)"},
	)
	q.Prep("mdl_value",
		ephemera.Col{Name: "noun", Type: "text"},
		ephemera.Col{Name: "kind", Type: "text"},
		ephemera.Col{Name: "field", Type: "text"},
		ephemera.Col{Name: "value", Type: "blob"},
	)
	// names are built from noun parts, and possibly from custom aliases.
	// where rank 0 is a better match than rank 1
	q.Prep("mdl_name",
		ephemera.Col{Name: "noun", Type: "text"},
		ephemera.Col{Name: "name", Type: "text"},
		ephemera.Col{Name: "rank", Type: "int"},
	)

	vcols := []ephemera.Col{
		{Name: "relation", Type: "text"},
		{Name: "stem", Type: "text"},
		{Check: "unique(stem)"},
	}
	if q, ok := q.(*ephemera.DbQueue); !ok {
		q.Prep("mdl_verb", vcols...)
	} else {
		// insert without duplication of the relation, stem pair
		q.PrepStatement("mdl_verb",
			`insert into mdl_verb(relation, stem)
				select ?1, ?2
				where not exists (
					select 1 from mdl_verb v
					where v.relation=?1 and v.stem=?2
				)`, vcols)
	}

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

func (m *Modeler) WriteField(field, kind, fieldType string) error {
	_, e := m.q.Write("mdl_field", field, kind, fieldType)
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

func (m *Modeler) WriteRelation(relation, kind, cardinality, other string) error {
	_, e := m.q.Write("mdl_rel", relation, kind, cardinality, other)
	return e
}

func (m *Modeler) WriteTrait(aspect, trait string) error {
	_, e := m.q.Write("mdl_aspect", aspect, trait, 0)
	return e
}

// WriteValue: store the initial value of an instance's field used at start of play.
func (m *Modeler) WriteValue(noun, field string, value interface{}) error {
	_, e := m.q.Write("mdl_value", noun, field, value)
	return e
}

func (m *Modeler) WriteVerb(relation, verb string) error {
	stem := porterstemmer.StemString(verb)
	_, e := m.q.Write("mdl_verb", relation, stem)
	return e
}
