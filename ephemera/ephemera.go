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

func NewRecorder(srcURI string, q Queue) *Recorder {
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
		Col{"idNamedCardinality", "int"},
		Col{"idNamedSecondary", "int"})
	q.Prep("eph_relative",
		Col{"idNamedRelativizer", "int"},
		Col{"idNamedHead", "int"},
		Col{"idNamedDependent", "int"})
	q.Prep("eph_trait",
		Col{"idNamedTrait", "int"},
		Col{"idNamedAspect", "int"},
		Col{"rank", "int"})
	q.Prep("eph_value",
		Col{"idNamedField", "int"},
		Col{"idNamedNoun", "int"},
		Col{"data", "blob"})
	//q.Prep("eph_Implication"},
	// Col{"idNamedScope"},
	// Col{"idNamedTrait"},
	// Col{"idNamedCertainty"},
	// Col{"idNamedImpliedTrait"})
	//
	srcId := q.Write("eph_source", srcURI)
	return &Recorder{q, srcId}
}

// Named records a user-specified string, including its meaning, location and returns a unique identifier for it.
// category is likely one of kind, noun, aspect, attribute, property, relation, alias
// names are not unique, one name can be many types.
// ofs depends on the source, might be item.id$parameter
func (r *Recorder) Named(category, name, ofs string) Named {
	namedId := r.q.Write("eph_named", name, category, r.src.id, ofs)
	return Named{namedId.id, name}
}

var None Named

// Alias provides a new name for another name.
func (r *Recorder) NewAlias(alias, actual Named) {
	r.q.Write("eph_alias", alias, actual)
}

// Aspect explicitly declares the existence of an aspect.
func (r *Recorder) NewAspect(aspect Named) {
	r.q.Write("eph_aspect", aspect)
}

// Certainty supplies a kind with a default trait.
func (r *Recorder) NewCertainty(certainty, trait, kind Named) {
	// usually fast horses.
	r.q.Write("eph_certainty", certainty, trait, kind)
}

// Kind connects a kind to its parent kind.
func (r *Recorder) NewKind(kind, parent Named) {
	r.q.Write("eph_kind", kind, parent)
}

// Noun connects a noun to its kind.
func (r *Recorder) NewNoun(noun, kind Named) {
	r.q.Write("eph_noun", noun, kind)
}

// Plural maps the plural form of a name to its singular form.
func (r *Recorder) NewPlural(plural, singluar Named) {
	r.q.Write("eph_plural", plural, singluar)
}

// Primitive property in the named kind.
func (r *Recorder) NewPrimitive(primType string, kind, prop Named) {
	r.q.Write("eph_primitive", primType, kind, prop)
}

// Relation defines a connection between a primary and secondary kind.
func (r *Recorder) NewRelation(relation, primary, cardinality, secondary Named) {
	r.q.Write("eph_relation", relation, primary, cardinality, secondary)
}

// Relative connects two specific nouns using a relativizer.
func (r *Recorder) NewRelative(relativizer, primary, secondary Named) {
	r.q.Write("eph_relative", relativizer, primary, secondary)
}

// Trait records a member of an aspect and its order ( rank. )
func (r *Recorder) NewTrait(trait, aspect Named, rank int) {
	r.q.Write("eph_trait", trait, aspect, rank)
}

// Value assigns the property of a noun a value;
// traits can be assigned by naming the individual trait and setting a true ( or false ) value.
// ( reverses order of primitive parameters for maximum confusion )
func (r *Recorder) NewValue(prop, noun Named, value interface{}) {
	r.q.Write("eph_value", prop, noun, value)
}
