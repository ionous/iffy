package ephemera

type Recorder struct {
	q   Queue
	src Queued
	// Source of ephemera for logging, revoking, etc.
}

const (
	PRIM_TEXT   = "text"   // string
	PRIM_DIGI   = "digi"   // number
	PRIM_EXPR   = "expr"   // text expression
	PRIM_COMP   = "comp"   // number computation
	PRIM_PROG   = "prog"   // program
	PRIM_ASPECT = "aspect" // fix? rename "prop"?
	PRIM_TRAIT  = "trait"  // fix? rename "attr"?
)

const (
	NAMED_ASPECT      = "aspect"
	NAMED_CERTAINTY   = "certainty"
	NAMED_FIELD       = "field"
	NAMED_KIND        = "kind"
	NAMED_NOUN        = "noun"
	NAMED_RELATIVIZER = "relativizer"
	NAMED_TRAIT       = "trait"
)

const (
	MANY_TO_ONE  = "any-one"
	ONE_TO_MANY  = "one-any"
	MANY_TO_MANY = "any-any"
	ONE_TO_ONE   = "one-one"
)

func NewRecorder(srcURI string, q Queue) (ret *Recorder) {
	// fix? should enums ( prim_..., named_... ) be stored as plain strings or as named entities?
	q.Prep("eph_source",
		Col{"src", "text"})
	q.Prep("eph_named",
		Col{"name", "text"},
		Col{"category", "text"},
		Col{"idSource", "int"},
		Col{"offset", "text"})
	q.Prep("eph_alias",
		Col{"idNamedAlias", "int"},
		Col{"idNamedActual", "int"})
	q.Prep("eph_aspect",
		Col{"idNamedAspect", "int"})
	q.Prep("eph_certainty",
		Col{"value", "text"},
		Col{"idNamedTrait", "int"},
		Col{"idNamedAspect", "int"})
	q.Prep("eph_kind",
		Col{"idNamedKind", "int"},
		Col{"idNamedParent", "int"})
	q.Prep("eph_noun",
		Col{"idNamedNoun", "int"},
		Col{"idNamedKind", "int"})
	q.Prep("eph_plural",
		Col{"idNamedPlural", "int"},
		Col{"idNamedSingluar", "int"})
	q.Prep("eph_primitive",
		Col{"primType", "text"},
		Col{"idNamedKind", "int"},
		Col{"idNamedField", "int"})
	q.Prep("eph_relation",
		Col{"idNamedRelation", "int"},
		Col{"idNamedPrimary", "int"},
		Col{"idNamedSecondary", "int"},
		Col{"idNamedCardinality", "int"})
	q.Prep("eph_relative",
		Col{"idNamedHead", "int"},
		Col{"idNamedVerb", "int"},
		Col{"idNamedDependent", "int"})
	q.Prep("eph_trait",
		Col{"idNamedTrait", "int"},
		Col{"idNamedAspect", "int"},
		Col{"rank", "int"})
	q.Prep("eph_value",
		Col{"idNamedField", "int"},
		Col{"idNamedNoun", "int"},
		Col{"data", "blob"})
	q.Prep("eph_verb",
		Col{"idNamedVerb", "int"},
		Col{"idNamedRelation", "int"})
	//q.Prep("eph_Implication"},
	// Col{"idNamedScope"},
	// Col{"idNamedTrait"},
	// Col{"idNamedCertainty"},
	// Col{"idNamedImpliedTrait"})
	//
	if srcId, e := q.Write("eph_source", srcURI); e != nil {
		panic(e) // fix? backwards compat
	} else {
		ret = &Recorder{q, srcId}
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

var None Named

// Alias provides a new name for another name.
func (r *Recorder) NewAlias(alias, actual Named) {
	if _, e := r.q.Write("eph_alias", alias, actual); e != nil {
		panic(e)
	}
}

// Aspect explicitly declares the existence of an aspect.
func (r *Recorder) NewAspect(aspect Named) {
	if _, e := r.q.Write("eph_aspect", aspect); e != nil {
		panic(e)
	}
}

// Certainty supplies a kind with a default trait.
func (r *Recorder) NewCertainty(certainty, trait, kind Named) {
	// usually fast horses.
	if _, e := r.q.Write("eph_certainty", certainty, trait, kind); e != nil {
		panic(e)
	}
}

// Kind connects a kind to its parent kind.
func (r *Recorder) NewKind(kind, parent Named) {
	if _, e := r.q.Write("eph_kind", kind, parent); e != nil {
		panic(e)
	}
}

// Noun connects a noun to its kind.
func (r *Recorder) NewNoun(noun, kind Named) {
	if _, e := r.q.Write("eph_noun", noun, kind); e != nil {
		panic(e)
	}
}

// Plural maps the plural form of a name to its singular form.
func (r *Recorder) NewPlural(plural, singluar Named) {
	if _, e := r.q.Write("eph_plural", plural, singluar); e != nil {
		panic(e)
	}
}

// Primitive property in the named kind.
func (r *Recorder) NewPrimitive(primType string, kind, prop Named) {
	if _, e := r.q.Write("eph_primitive", primType, kind, prop); e != nil {
		panic(e)
	}
}

// Relation defines a connection between a primary and secondary kind.
func (r *Recorder) NewRelation(relation, primary, secondary, cardinality Named) {
	if _, e := r.q.Write("eph_relation", relation, primary, secondary, cardinality); e != nil {
		panic(e)
	}
}

// Relative connects two specific nouns using a verb.
func (r *Recorder) NewRelative(primary, verb, secondary Named) {
	if _, e := r.q.Write("eph_relative", primary, verb, secondary); e != nil {
		panic(e)
	}
}

// Trait records a member of an aspect and its order ( rank. )
func (r *Recorder) NewTrait(trait, aspect Named, rank int) {
	if _, e := r.q.Write("eph_trait", trait, aspect, rank); e != nil {
		panic(e)
	}
}

// Value assigns the property of a noun a value;
// traits can be assigned by naming the individual trait and setting a true ( or false ) value.
// ( reverses order of primitive parameters for maximum confusion )
func (r *Recorder) NewValue(prop, noun Named, value interface{}) {
	if _, e := r.q.Write("eph_value", prop, noun, value); e != nil {
		panic(e)
	}
}

// Relative connects two specific nouns using a verb.
func (r *Recorder) NewVerb(verb, relation Named) {
	if _, e := r.q.Write("eph_verb", verb, relation); e != nil {
		panic(e)
	}
}
