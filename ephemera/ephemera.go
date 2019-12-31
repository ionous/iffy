package ephemera

type Recorder struct {
	q   Queue
	src Queued
	// Source of ephemera for logging, revoking, etc.
}

// tdb: should enums be stored as plain strings or as named entities
// ( see also primitives )
// one advantage in naming is we can locate them; one disadvantage is increased db size secondary

func NewRecorder(srcURI string, q Queue) *Recorder {
	q.Prep("Source", "src")
	//
	q.Prep("Named", "name", "category", "idSource", "offset")
	q.Prep("Alias", "idNamedAlias", "idNamedActual")
	q.Prep("Aspect", "idNamedAspect")
	q.Prep("Kind", "idNamedKind", "idNamedParent")
	q.Prep("Noun", "idNamedNoun", "idNamedKind")
	q.Prep("Plural", "idNamedPlural", "idNamedSingluar")
	q.Prep("Primitive", "primType", "idNamedArchetype", "idNamedField")
	q.Prep("Relation", "idNamedRelation", "idNamedPrimary", "idNamedCardinality", "idNamedSecondary")
	q.Prep("Relative", "idNamedRelativizer", "idNamedHead", "idNamedDependent")
	q.Prep("Trait", "idNamedTrait", "idNamedAspect", "rank")
	q.Prep("Value", "idNamedField", "idNamedNoun", "value")
	//q.Prep("Implication", "idNamedScope", "idNamedTrait", "idNamedCertainty", "idNamedImpliedTrait")
	//
	srcId := q.Write("Source", srcURI)
	return &Recorder{q, srcId}
}

// Named records a user-specified string, including its meaning, location and returns a unique identifier for it.
// category is likely one of kind, noun, aspect, attribute, property, relation, alias
// names are not unique, one name can be many types.
// ofs depends on the source, might be item.id$parameter
func (r *Recorder) Named(category, name, ofs string) Named {
	namedId := r.q.Write("Named", name, category, r.src.id, ofs)
	return Named{namedId.id, name}
}

var None Named

// Alias provides a new name for another name.
func (r *Recorder) Alias(alias, actual Named) {
	r.q.Write("Alias", alias, actual)
}

// Aspect explicitly declares the existence of an aspect.
func (r *Recorder) Aspect(aspect Named) {
	r.q.Write("Aspect", aspect)
}

// Certainty supplies a kind with a default trait.
func (r *Recorder) Certainty(certainty, trait, kind Named) {
	// usually fast horses.
	r.q.Write("Certainty", certainty, trait, kind)
}

// Kind connects a kind to its parent kind.
func (r *Recorder) Kind(kind, parent Named) {
	r.q.Write("Kind", kind, parent)
}

// Noun connects a noun to its kind.
func (r *Recorder) Noun(noun, kind Named) {
	r.q.Write("Noun", noun, kind)
}

// Plural maps the plural form of a name to its singular form.
func (r *Recorder) Plural(plural, singluar Named) {
	r.q.Write("Kind", plural, singluar)
}

// Primitive property in the named kind.
func (r *Recorder) Primitive(primType string, kind, prop Named) {
	r.q.Write("Primitive", primType, kind, prop)
}

// Relation defines a connection between a primary and secondary kind.
func (r *Recorder) Relation(relation, primary, cardinality, secondary Named) {
	r.q.Write("Kind", relation, primary, cardinality, secondary)
}

// Relative connects two specific nouns using a relativizer.
func (r *Recorder) Relative(relativizer, primary, secondary Named) {
	r.q.Write("Kind", relativizer, primary, secondary)
}

// Trait records a member of an aspect and its order ( rank. )
func (r *Recorder) Trait(trait, aspect Named, rank int) {
	r.q.Write("Trait", trait, aspect, rank)
}

// Value assigns the property of a noun a value;
// traits can be assigned by naming the individual trait and setting a true ( or false ) value.
// ( reverses order of primitive parameters for maximum confusion )
func (r *Recorder) Value(prop, noun Named, value interface{}) {
	r.q.Write("Value", prop, noun, value)
}
