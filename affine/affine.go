package affine

type Affinity string

const (
	Bool       Affinity = "bool"
	Number     Affinity = "number"
	NumList    Affinity = "num_list"
	Text       Affinity = "text"
	TextList   Affinity = "text_list"
	Object     Affinity = "object"
	Record     Affinity = "record"
	RecordList Affinity = "record_list"

	// extended text types
	// Aspect      Affinity = "aspect"
	// Check       Affinity = "check"
	// Kind        Affinity = "singular_kind"
	// Kinds       Affinity = "kind"
	// Noun        Affinity = "noun"
	// Pattern     Affinity = "pattern"
	// PluralKinds Affinity = "plural_kinds"
	// Primitive   Affinity = "primitive"
	// Relation    Affinity = "relation"
	// Scene       Affinity = "scene"
	// Trait       Affinity = "trait"
	// Verb        Affinity = "verb"
)
