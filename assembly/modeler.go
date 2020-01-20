package assembly

import (
	"github.com/ionous/iffy/ephemera"
	"github.com/reiver/go-porterstemmer"
)

func NewModeler(q ephemera.Queue) *Modeler {
	q.Prep("mdl_kind",
		ephemera.Col{Name: "kind", Type: "text"},
		ephemera.Col{Name: "path", Type: "text"},
		ephemera.Col{Check: "primary key(kind)"},
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
	q.Prep("mdl_aspect",
		ephemera.Col{Name: "aspect", Type: "text"},
		ephemera.Col{Name: "trait", Type: "text"},
		ephemera.Col{Name: "rank", Type: "int"},
		ephemera.Col{Check: "primary key(aspect, trait)"},
	)
	q.Prep("mdl_noun",
		ephemera.Col{Name: "kind", Type: "text"},
		// the full text of 'integer' is required for ids
		// also, to get the auto-generated id, we place the declaration in "check".
		ephemera.Col{Check: "id integer primary key"},
	)
	// where rank 0 is a better match than rank 1
	q.Prep("mdl_name",
		ephemera.Col{Name: "idModelNoun", Type: "int"},
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
func (w *Modeler) WriteAncestor(kind, path string) (err error) {
	if _, e := w.q.Write("mdl_kind", kind, path); e != nil {
		err = e
	}
	return
}

func (w *Modeler) WriteField(field, owner, fieldType string) (err error) {
	if _, e := w.q.Write("mdl_field", field, owner, fieldType); e != nil {
		err = e
	}
	return
}

func (w *Modeler) WriteNoun(kind string) (ephemera.Queued, error) {
	return w.q.Write("mdl_noun", kind)
}

// WriteName for noun
func (w *Modeler) WriteName(noun ephemera.Queued, name string, rank int) (ephemera.Queued, error) {
	return w.q.Write("mdl_name", noun, name, rank)
}

func (w *Modeler) WriteRelation(relation, kind, cardinality, other string) (err error) {
	if _, e := w.q.Write("mdl_rel", relation, kind, cardinality, other); e != nil {
		err = e
	}
	return
}

func (w *Modeler) WriteTrait(aspect, trait string) (err error) {
	if _, e := w.q.Write("mdl_aspect", aspect, trait, 0); e != nil {
		err = e
	}
	return
}

func (w *Modeler) WriteVerb(relation, verb string) (err error) {
	stem := porterstemmer.StemString(verb)
	if _, e := w.q.Write("mdl_verb", relation, stem); e != nil {
		err = e
	}
	return
}
