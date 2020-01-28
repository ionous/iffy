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
	PRIM_ASPECT = "aspect" // string
	PRIM_TRAIT  = "trait"  // string
)

const (
	NAMED_ASPECT    = "aspect"
	NAMED_CERTAINTY = "certainty"
	NAMED_FIELD     = "field"
	NAMED_KIND      = "kind"
	NAMED_PROPERTY  = "prop" // field, trait, or aspect
	NAMED_NOUN      = "noun"
	NAMED_RELATION  = "relation"
	NAMED_VERB      = "verb"
	NAMED_TRAIT     = "trait"
)

const (
	MANY_TO_ONE  = "any-one"
	ONE_TO_MANY  = "one-any"
	MANY_TO_MANY = "any-any"
	ONE_TO_ONE   = "one-one"
)

const (
	USUALLY = "usually"
	ALWAYS  = "always"
	SELDOM  = "seldom"
	NEVER   = "never"
)

func NewRecorder(srcURI string, q Queue) (ret *Recorder) {
	// fix? should enums ( prim_..., named_... ) be stored as plain strings or as named entities?
	q.Prep("eph_source",
		Col{Name: "src", Type: "text"})
	q.Prep("eph_named",
		Col{Name: "name", Type: "text"},
		Col{Name: "category", Type: "text"},
		Col{Name: "idSource", Type: "int"},
		Col{Name: "offset", Type: "text"})
	// aliases are used for user input, not story modeling
	q.Prep("eph_alias",
		Col{Name: "idNamedAlias", Type: "int"},
		Col{Name: "idNamedActual", Type: "int"})
	q.Prep("eph_aspect",
		Col{Name: "idNamedAspect", Type: "int"})
	q.Prep("eph_default",
		Col{Name: "idNamedKind", Type: "int"},
		// field, trait, or aspect
		Col{Name: "idNamedProp", Type: "int"},
		// future: un/certainty for defaults and values
		// 	Col{Name: "certainty", Type:  "text",
		// 	Check: "check (certainty in ('usually','always','seldom','never'))"},
		Col{Name: "value", Type: "blob"})
	q.Prep("eph_kind",
		Col{Name: "idNamedKind", Type: "int"},
		Col{Name: "idNamedParent", Type: "int"})
	q.Prep("eph_noun",
		Col{Name: "idNamedNoun", Type: "int"},
		Col{Name: "idNamedKind", Type: "int"})
	q.Prep("eph_plural",
		Col{Name: "idNamedPlural", Type: "int"},
		Col{Name: "idNamedSingluar", Type: "int"})
	q.Prep("eph_primitive",
		Col{Name: "primType", Type: "text"},
		Col{Name: "idNamedKind", Type: "int"},
		Col{Name: "idNamedField", Type: "int"})
	q.Prep("eph_relation",
		Col{Name: "idNamedRelation", Type: "int"},
		Col{Name: "idNamedKind", Type: "int"},
		Col{Name: "idNamedOtherKind", Type: "int"},
		Col{Name: "cardinality",
			Type:  "text",
			Check: "check (cardinality in ('one-one','one-any','any-one','any-any'))"})
	q.Prep("eph_relative",
		Col{Name: "idNamedHead", Type: "int"},
		Col{Name: "idNamedVerb", Type: "int"},
		Col{Name: "idNamedDependent", Type: "int"})
	q.Prep("eph_trait",
		Col{Name: "idNamedTrait", Type: "int"},
		Col{Name: "idNamedAspect", Type: "int"},
		Col{Name: "rank", Type: "int"})
	q.Prep("eph_value",
		Col{Name: "idNamedNoun", Type: "int"},
		Col{Name: "idNamedProp", Type: "int"},
		Col{Name: "value", Type: "blob"})
	q.Prep("eph_verb",
		Col{Name: "idNamedVerb", Type: "int"},
		Col{Name: "idNamedRelation", Type: "int"})
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

// NewAlias provides a new name for another name.
func (r *Recorder) NewAlias(alias, actual Named) {
	if _, e := r.q.Write("eph_alias", alias, actual); e != nil {
		panic(e)
	}
}

// NewAspect explicitly declares the existence of an aspect.
func (r *Recorder) NewAspect(aspect Named) {
	if _, e := r.q.Write("eph_aspect", aspect); e != nil {
		panic(e)
	}
}

// NewCertainty supplies a kind with a default trait.
func (r *Recorder) NewCertainty(certainty string, trait, kind Named) {
	// usually fast horses.
	if _, e := r.q.Write("eph_certainty", certainty, trait, kind); e != nil {
		panic(e)
	}
}

// NewDefault supplies a kind with a default value.
func (r *Recorder) NewDefault(kind, field Named, value interface{}) {
	// height horses 5.
	if _, e := r.q.Write("eph_default", kind, field, value); e != nil {
		panic(e)
	}
}

// NewKind connects a kind to its parent kind.
func (r *Recorder) NewKind(kind, parent Named) {
	if _, e := r.q.Write("eph_kind", kind, parent); e != nil {
		panic(e)
	}
}

// NewNoun connects a noun to its kind.
func (r *Recorder) NewNoun(noun, kind Named) {
	if _, e := r.q.Write("eph_noun", noun, kind); e != nil {
		panic(e)
	}
}

// NewPlural maps the plural form of a name to its singular form.
func (r *Recorder) NewPlural(plural, singluar Named) {
	if _, e := r.q.Write("eph_plural", plural, singluar); e != nil {
		panic(e)
	}
}

// NewPrimitive property in the named kind.
func (r *Recorder) NewPrimitive(primType string, kind, prop Named) {
	if _, e := r.q.Write("eph_primitive", primType, kind, prop); e != nil {
		panic(e)
	}
}

// NewRelation defines a connection between a primary and secondary kind.
func (r *Recorder) NewRelation(relation, primaryKind, secondaryKind Named, cardinality string) {
	if _, e := r.q.Write("eph_relation", relation, primaryKind, secondaryKind, cardinality); e != nil {
		panic(e)
	}
}

// NewRelative connects two specific nouns using a verb.
func (r *Recorder) NewRelative(primary, verb, secondary Named) {
	if _, e := r.q.Write("eph_relative", primary, verb, secondary); e != nil {
		panic(e)
	}
}

// NewTrait records a member of an aspect and its order ( rank. )
func (r *Recorder) NewTrait(trait, aspect Named, rank int) {
	if _, e := r.q.Write("eph_trait", trait, aspect, rank); e != nil {
		panic(e)
	}
}

// NewValue assigns the property of a noun a value;
// traits can be assigned by naming the individual trait and setting a true ( or false ) value.
func (r *Recorder) NewValue(noun, prop Named, value interface{}) {
	if _, e := r.q.Write("eph_value", noun, prop, value); e != nil {
		panic(e)
	}
}

// NewRelative connects two specific nouns using a verb.
func (r *Recorder) NewVerb(verb, relation Named) {
	if _, e := r.q.Write("eph_verb", verb, relation); e != nil {
		panic(e)
	}
}
