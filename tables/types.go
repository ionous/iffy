package tables

// primType
const (
	PRIM_TEXT   = "text"   // string
	PRIM_DIGI   = "digi"   // number
	PRIM_EXPR   = "expr"   // text expression
	PRIM_COMP   = "comp"   // number computation
	PRIM_PROG   = "prog"   // program, activity
	PRIM_ASPECT = "aspect" // string
	PRIM_TRAIT  = "trait"  // string
)

// evalType
const (
	EVAL_TEXT = "expr" // string
	EVAL_DIGI = "comp" // number
	EVAL_BOOL = "bool" // boolean
	EVAL_PROG = "prog" // text expression
)

// named category
// FIX -- make these match the names in the file and things are much cleaner
const (
	NAMED_ASPECT    = "aspect"
	NAMED_CERTAINTY = "certainty"
	NAMED_FIELD     = "field"
	NAMED_KIND      = "kind"
	NAMED_PROPERTY  = "prop" // field, trait, or aspect
	NAMED_NOUN      = "noun"
	//NAMED_PATTERN   = "pattern_name"
	//NAMED_VARIABLE   = "variable_name"
	NAMED_RELATION = "relation"
	NAMED_VERB     = "verb"
	NAMED_TEST     = "test"
	NAMED_TRAIT    = "trait"
)

// cardinality
const (
	ONE_TO_ONE   = "one_one"
	ONE_TO_MANY  = "one_any"
	MANY_TO_ONE  = "any_one"
	MANY_TO_MANY = "any_any"
)

// certainty
const (
	USUALLY = "usually"
	ALWAYS  = "always"
	SELDOM  = "seldom"
	NEVER   = "never"
)
