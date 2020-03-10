package tables

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
	NAMED_TEST      = "test"
)

const (
	ONE_TO_ONE   = "one_one"
	ONE_TO_MANY  = "one_any"
	MANY_TO_ONE  = "any_one"
	MANY_TO_MANY = "any_any"
)

const (
	USUALLY = "usually"
	ALWAYS  = "always"
	SELDOM  = "seldom"
	NEVER   = "never"
)
