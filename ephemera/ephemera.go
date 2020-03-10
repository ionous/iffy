package ephemera

import "github.com/ionous/iffy/tables"

type Recorder struct {
	q   Queue
	src Queued
	// fix
	// we dont really need the "queue" interface anymore --
	// it was good for initial testing/ experimentation --
	// but we can live without
	// itd be nicer to move to table caching instead
	cache *tables.Cache
}

func NewRecorder(srcURI string, q *DbQueue) (ret *Recorder) {
	// fix? should enums ( prim_..., named_... ) be stored as plain strings or as named entities?
	q.Prep("eph_source", []tables.Col{
		{Name: "src", Type: "text"},
	})
	q.Prep("eph_named", []tables.Col{
		{Name: "name", Type: "text"},
		{Name: "category", Type: "text"},
		{Name: "idSource", Type: "int"},
		{Name: "offset", Type: "text"},
	})
	// aliases are used for user input, not story modeling
	q.Prep("eph_alias", []tables.Col{
		{Name: "idNamedAlias", Type: "int"},
		{Name: "idNamedActual", Type: "int"},
	})
	q.Prep("eph_aspect", []tables.Col{
		{Name: "idNamedAspect", Type: "int"},
	})
	q.Prep("eph_default", []tables.Col{
		{Name: "idNamedKind", Type: "int"},
		// field, trait, or aspect
		{Name: "idNamedProp", Type: "int"},
		// future: un/certainty for defaults and values
		// 	{Name: "certainty", Type:  "text",
		// 	Check: "check (certainty in ('usually','always','seldom','never'))"},
		{Name: "value", Type: "blob"},
	})
	q.Prep("eph_kind", []tables.Col{
		{Name: "idNamedKind", Type: "int"},
		{Name: "idNamedParent", Type: "int"},
	})
	q.Prep("eph_noun", []tables.Col{
		{Name: "idNamedNoun", Type: "int"},
		{Name: "idNamedKind", Type: "int"},
	})
	q.Prep("eph_plural", []tables.Col{
		{Name: "idNamedPlural", Type: "int"},
		{Name: "idNamedSingluar", Type: "int"},
	})
	q.Prep("eph_primitive", []tables.Col{
		{Name: "primType", Type: "text"},
		{Name: "idNamedKind", Type: "int"},
		{Name: "idNamedField", Type: "int"},
	})
	q.Prep("eph_relation", []tables.Col{
		{Name: "idNamedRelation", Type: "int"},
		{Name: "idNamedKind", Type: "int"},
		{Name: "idNamedOtherKind", Type: "int"},
		{Name: "cardinality",
			Type:  "text",
			Check: "check (cardinality in ('one_one','one_any','any_one','any_any'))"},
	})
	q.Prep("eph_relative", []tables.Col{
		{Name: "idNamedHead", Type: "int"},
		{Name: "idNamedStem", Type: "int"},
		{Name: "idNamedDependent", Type: "int"},
	})
	q.Prep("eph_trait", []tables.Col{
		{Name: "idNamedTrait", Type: "int"},
		{Name: "idNamedAspect", Type: "int"},
		{Name: "rank", Type: "int"},
	})
	q.Prep("eph_value", []tables.Col{
		{Name: "idNamedNoun", Type: "int"},
		{Name: "idNamedProp", Type: "int"},
		{Name: "value", Type: "blob"},
	})
	q.Prep("eph_verb", []tables.Col{
		{Name: "idNamedStem", Type: "int"},
		{Name: "idNamedRelation", Type: "int"},
		{Name: "verb", Type: "string"},
	})
	//q.Prep("eph_implication"},
	// tables.Col{"idNamedScope"},
	// tables.Col{"idNamedTrait"},
	// tables.Col{"idNamedCertainty"},
	// tables.Col{"idNamedImpliedTrait"}
	//
	if srcId, e := q.Write("eph_source", srcURI); e != nil {
		panic(e) // fix? backwards compat
	} else {
		ret = &Recorder{q, srcId, tables.NewCache(q.db)}
	}
	return
}

// Named records a user-specified string, including its meaning, location and returns a unique identifier for it.
// category is likely one of kind, noun, aspect, attribute, property, relation, alias
// names are not unique, one name can be many types.
// ofs depends on the source, might be item.id$parameter
func (r *Recorder) Named(category, name, ofs string) (ret Named) {
	if namedId, e := r.q.Write("eph_named", name, category, r.src.id, ofs); e != nil {
		panic(e) // fix? backwards compat
	} else {
		ret = Named{namedId.id, name}
	}
	return
}

type Prog struct{ Named }

func (r *Recorder) NewProg(rootType string, blob []byte) (ret Prog) {
	id := r.cache.Must(ephProg, r.src, rootType, blob)
	ret = Prog{Named{id, rootType}}
	return
}

var ephProg = tables.Insert("eph_prog", "idSource", "type", "prog")

var None Named

// NewAlias provides a new name for another name.
func (r *Recorder) NewAlias(alias, actual Named) {
	r.mustWrite("eph_alias", alias, actual)
}

// NewAspect explicitly declares the existence of an aspect.
func (r *Recorder) NewAspect(aspect Named) {
	r.mustWrite("eph_aspect", aspect)
}

// NewCertainty supplies a kind with a default trait.
func (r *Recorder) NewCertainty(certainty string, trait, kind Named) {
	// usually fast horses.
	r.mustWrite("eph_certainty", certainty, trait, kind)
}

// NewDefault supplies a kind with a default value.
func (r *Recorder) NewDefault(kind, field Named, value interface{}) {
	// height horses 5.
	r.mustWrite("eph_default", kind, field, value)
}

// NewKind connects a kind to its parent kind.
func (r *Recorder) NewKind(kind, parent Named) {
	r.mustWrite("eph_kind", kind, parent)
}

// NewNoun connects a noun to its kind.
func (r *Recorder) NewNoun(noun, kind Named) {
	r.mustWrite("eph_noun", noun, kind)
}

// NewPlural maps the plural form of a name to its singular form.
func (r *Recorder) NewPlural(plural, singluar Named) {
	r.mustWrite("eph_plural", plural, singluar)
}

// NewPrimitive property in the named kind.
func (r *Recorder) NewPrimitive(primType string, kind, prop Named) {
	r.mustWrite("eph_primitive", primType, kind, prop)
}

// NewRelation defines a connection between a primary and secondary kind.
func (r *Recorder) NewRelation(relation, primaryKind, secondaryKind Named, cardinality string) {
	r.mustWrite("eph_relation", relation, primaryKind, secondaryKind, cardinality)
}

// NewRelative connects two specific nouns using a verb stem.
func (r *Recorder) NewRelative(primary, stem, secondary Named) {
	r.mustWrite("eph_relative", primary, stem, secondary)
}

// NewTrait records a member of an aspect and its order ( rank. )
func (r *Recorder) NewTrait(trait, aspect Named, rank int) {
	r.mustWrite("eph_trait", trait, aspect, rank)
}

// NewValue assigns the property of a noun a value;
// traits can be assigned by naming the individual trait and setting a true ( or false ) value.
func (r *Recorder) NewValue(noun, prop Named, value interface{}) {
	r.mustWrite("eph_value", noun, prop, value)
}

// NewRelative connects two specific nouns using a verb.
func (r *Recorder) NewVerb(stem, relation Named, verb string) {
	r.mustWrite("eph_verb", stem, relation, verb)
}

func (r *Recorder) NewTest(test Named, prog Prog, expect string) {
	r.cache.Must(ephTest, test, prog, expect)
}

func (r *Recorder) mustWrite(which string, args ...interface{}) {
	if _, e := r.q.Write(which, args...); e != nil {
		panic(e)
	}
}

var ephTest = tables.Insert("eph_test", "idNamedTest", "idProg", "expect")
