package affine

type Affinity string

const (
	// value types
	Bool     Affinity = "bool"
	Number   Affinity = "number"
	Text     Affinity = "text"
	NumList  Affinity = "num_list"
	TextList Affinity = "text_list"

	// extended types
	Aspect      Affinity = "aspect"
	Check       Affinity = "check"
	Kind        Affinity = "singular_kind"
	Kinds       Affinity = "kind"
	PluralKinds Affinity = "plural_kinds"
	Noun        Affinity = "noun"
	Pattern     Affinity = "pattern"
	Primitive   Affinity = "primitive"
	Relation    Affinity = "relation"
	Verb        Affinity = "verb"
	Scene       Affinity = "scene"
	Trait       Affinity = "trait"
)
