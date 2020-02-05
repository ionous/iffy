package assembly

import (
	"database/sql"
	"strings"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ephemera"
	"github.com/reiver/go-porterstemmer"
)

func NewModelerDB(db *sql.DB) *Modeler {
	dbq := ephemera.NewDBQueue(db)
	if _, e := db.Exec(`
	/* resolve default ephemera to strings.
	 */
	create temp view 
	asm_default as
		select p.rowid as idEphDefault, nk.name as kind, nf.name as prop, p.value as value
	from eph_default p join eph_named nk
		on (p.idNamedKind = nk.rowid)
	left join eph_named nf
 		on (p.idNamedProp = nf.rowid);

	/* resolve value ephemera to strings.
	 */
	create temp view
	asm_value as
		select pv.rowid as idEphValue, nn.name, np.name as prop, pv.value
	from eph_value pv join eph_named nn
		on (pv.idNamedNoun = nn.rowid)
	left join eph_named np
		on (pv.idNamedProp = np.rowid);
	
	/* resolve value ephemera to nouns.
	 */
	create temp view
	asm_noun as 
		select *, ( 
			select n.noun 
			from mdl_name as n
			where asm.name = n.name 
		 	order by rank
			limit 1 
		) as noun
	from asm_value as asm;
	
	/* resolve relative ephemera to strings.
	 */
	create temp view
	asm_relative as
	select rel.rowid as idEphRel, 
		na.name as noun, 
		nv.name as stem,
		nb.name as otherNoun
	from eph_relative rel
	join eph_named na
		on (rel.idNamedHead = na.rowid)
	left join eph_named nv
 		on (rel.idNamedStem = nv.rowid)
 	left join eph_named nb
 		on (rel.idNamedDependent = nb.rowid);
	
	/* resolve relative ephemera to nouns and relations
	use left join(s) to return nulls for missing elements 
	 */
	create temp view
	asm_relation as
	select 	
		idEphRel, 
		stem, relation, cardinality, 
		
		/* first contains the kind of the user specified noun;
			swapped contains the kind of the relation 
		*/
		first.noun as noun, 
		case when instr((
		 			select mk.kind || "," || mk.path || ","
					from mdl_kind mk
					where mk.kind = first.kind
				),  swapped.kind || ",") 
				then first.kind 
		end as kind,

		/* second contains the kind of the other user specified noun;
			swapped contains the other kind of the relation
		 */
		second.noun as otherNoun,
		case when instr((
		 			select mk.kind || "," || mk.path || ","
					from mdl_kind mk
					where mk.kind = second.kind
				),  swapped.otherKind || ",") 
				then second.kind
		end as otherKind
	from (
		select 
			idEphRel,stem,relation,cardinality,
			case swap when 1 then otherNoun else noun end as noun,
			case swap when 1 then noun else otherNoun end as otherNoun,
			case swap when 1 then otherKind else kind end as kind,
			case swap when 1 then kind else otherKind end as otherKind
		from (
			select *, (cardinality = 'one_one') and (noun > otherNoun) as swap
				from asm_relative ar
				left join mdl_verb mv
					using (stem)
				left join mdl_rel mr
					using (relation)
		)
	) as swapped
	left join mdl_noun first
		 on (first.noun = swapped.noun)
	left join mdl_noun second 
		on (second.noun = swapped.otherNoun);

	/* the bits of asm_relation which didnt make it into the start_rel table.
	 */
	create temp view 
	res_mismatch as
	select idEphRel, stem, relation, cardinality, noun, kind, otherNoun, otherKind
	from asm_relation asm
	where max(asm.relation, asm.kind, asm.otherKind) is null
	or case asm.cardinality
		when 'one_one' then
		exists(
			select 1 
			from start_rel rel 
			where (asm.relation = rel.relation) 
			and ((asm.noun = rel.noun) and (asm.otherNoun != rel.otherNoun)
			or (asm.otherNoun = rel.otherNoun) and (asm.noun != rel.noun))
		)
		when 'one_any' then 
		exists(
			/* given otherNoun there is only one valid noun */
			select 1 
			from start_rel rel 
			where (asm.relation = rel.relation)
			and (asm.otherNoun = rel.otherNoun) 
			and (asm.noun != rel.noun)
		)
		when 'any_one' then 
		exists(
			/* given noun there is only one valid otherNoun */
			select 1 
			from start_rel rel 
			where (asm.relation = rel.relation)
			and (asm.noun = rel.noun) 
			and (asm.otherNoun != rel.otherNoun)
		)
	end;
	`); e != nil {
		panic(e)
	}
	//
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
	q.PrepCols("mdl_trait", []ephemera.Col{
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
	q.PrepStatement("mdl_verb",
		`insert into mdl_verb(relation, stem)
				select ?1, ?2
				where not exists (
					select 1 from mdl_verb v
					where v.relation=?1 and v.stem=?2
				)`, []ephemera.Col{
			{Name: "relation", Type: "text"}, /* reference to mdl_rel */
			{Name: "stem", Type: "text"},
			{Check: "unique(stem)"},
		})
	q.PrepCols("start_val", []ephemera.Col{
		{Name: "noun", Type: "text"},  /* reference to mdl_noun */
		{Name: "field", Type: "text"}, /* partial reference to mdl_field */
		{Name: "value", Type: "blob"},
	})
	q.PrepCols("start_rel", []ephemera.Col{
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
	_, e := m.q.Write("mdl_trait", aspect, trait, rank)
	return e
}

// WriteValue: store the initial value of an instance's field used at start of play.
func (m *Modeler) WriteValue(noun, field string, value interface{}) error {
	_, e := m.q.Write("start_val", noun, field, value)
	return e
}

func (m *Modeler) WriteVerb(relation, verb string) error {
	stem := porterstemmer.StemString(verb)
	_, e := m.q.Write("mdl_verb", relation, stem)
	return e
}
