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

// return the affinity of an affine list, or blank if not a list.
func Element(list Affinity) (ret Affinity) {
	switch a := list; a {
	case TextList:
		ret = Text
	case NumList:
		ret = Number
	case RecordList:
		ret = Record
	}
	return
}

func MatchTypes(a Affinity, at string, b Affinity, bt string) bool {
	return a == b && ((a != Record && a != RecordList) || (at == bt))
}

func IsList(a Affinity) bool {
	elAffinity := Element(a)
	return len(elAffinity) > 0
}
