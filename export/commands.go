package export

import (
	"github.com/ionous/iffy/dl/check"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/std"
	"github.com/ionous/iffy/rt"
)

type Slot struct {
	Type  interface{} // nil instance, ex. (*core.CompareTo)(nil)
	Desc  string
	Group string // display group(s)
}

// a work in progress for sure
type Run struct {
	Type  interface{} // nil instance, ex. (*core.CompareTo)(nil)
	Group string      // display group(s)
	Desc  string
	Slots []string // FIX: really should always be based on implemented interface
	//
	Spec   string   // embedded pre-token string
	Phrase string   // token string
	Locals []string // FIX: names put into scope
}

var Slots = map[string]Slot{
	"execute": {Type: (*rt.Execute)(nil),
		Desc: "Execute: Run a series of statements."},
	"bool_eval": {Type: (*rt.BoolEval)(nil),
		Desc: "Booleans: Statements which return true/false values."},
	"number_eval": {Type: (*rt.NumberEval)(nil),
		Desc: "Numbers: Statements which return a number."},
	"text_eval": {Type: (*rt.TextEval)(nil),
		Desc: "Texts: Statements which return text."},
	"object_eval": {Type: (*rt.ObjectEval)(nil),
		Desc: "Objects: Statements which return an existing object."},
	"num_list_eval": {Type: (*rt.NumListEval)(nil),
		Desc: "Number List: Statements which return a list of numbers."},
	"text_list_eval": {Type: (*rt.TextListEval)(nil),
		Desc: "Text Lists: Statements which return a list of text."},
	"obj_list_eval": {Type: (*rt.ObjListEval)(nil),
		Desc: "Object Lists: Statements which return a list of existing objects."},
	"compare_to": {Type: (*core.CompareTo)(nil),
		Desc: "Comparison Types: Helper used when comparing two numbers, objects, pieces of text, etc."},
}

var Runs = map[string]Run{
	"sum_of": {
		Type:   (*core.SumOf)(nil),
		Group:  "math",
		Desc:   "Add Numbers: Add two numbers.",
		Phrase: "( $1 + $2 )",
	},
	"all_true": {
		Type:  (*core.AllTrue)(nil),
		Group: "logic",
		Desc:  "All True: returns true if all of the evaluations are true.",
	},
	"any_true": {
		Type:  (*core.AnyTrue)(nil),
		Group: "logic",
		Desc:  "Any True: returns true if any of the evaluations are true.",
	},
	"bool_value": {
		Type:  (*core.Bool)(nil),
		Group: "literals",
		Desc:  "Bool Value: specifies an explicit true/false value.",
	},
	// not exported. for internal use only, may go away.
	// {
	// Type: (*core.Buffer)(nil),
	// 	Group: "format",
	// 	Desc:  "collects text said by other statements",
	// },
	"choose": {
		Type: (*core.Choose)(nil),
	},

	// FIX: Choose any?
	"choose_num": {
		Type:  (*core.ChooseNum)(nil),
		Group: "math",
		Desc:  "Choose Number: Pick one of two numbers based on a boolean test.",
	},
	"choose_obj": {
		Type:  (*core.ChooseObj)(nil),
		Group: "objects",
		Desc:  "Choose Object: Pick one of two objects based on a boolean test.",
	},
	"choose_text": {
		Type:  (*core.ChooseText)(nil),
		Group: "format",
		Desc:  "Choose Text: Pick one of two strings based on a boolean test.",
	},
	// not exported. rework as a function?
	// {
	// Type: (*core.ClassName)(nil),
	// 	Group: "objects",
	// 	Desc:  "friendly name of the object's class",
	// },
	// FIX: compare scalar?
	"compare_num": {
		Type:  (*core.CompareNum)(nil),
		Group: "logic",
		Desc:  "Compare Numbers: True if eq,ne,gt,lt,ge,le two numbers.",
	},
	"compare_obj": {
		Type:  (*core.CompareObj)(nil),
		Group: "",
		Desc:  "Compare Objects",
	},
	"compare_text": {
		Type:  (*core.CompareText)(nil),
		Group: "logic",
		Desc:  "Compare Text: True if eq,ne,gt,lt,ge,le two strings ( lexical. )",
	},

	// remove: use templates.
	// {
	// Type: (*core.Comprise)(nil),
	// 	Group: "format",
	// 	Desc:  "writes a prefix and suffix around a body of text so long as the body has content.",
	// },

	"cycle_text": {
		Type:  (*core.CycleText)(nil),
		Group: "cycle",
		Desc:  "Cycle Text: When called multiple times, returns each of its inputs in turn.",
	},

	// use patterns
	// {
	// Type: (*core.DescribeLocation)(nil),
	// 	Group: "locale",
	// 	Desc:  "prints details about the targeted location, including 'paragraphs” for notable objects in the location, and a sentence for otherwise unremarkable objects.",
	// },

	// internal?
	// {
	// Type: (*core.Determine)(nil),
	// 	Group: "patterns",
	// 	Desc:  "runs a pattern",
	// },
	"quotient_of": {
		Type:   (*core.QuotientOf)(nil),
		Group:  "math",
		Phrase: "( $1 / $2 )",
		Desc:   "Divide Numbers: Divide one number by another.",
	},
	"do_nothing": {
		Type:  (*core.DoNothing)(nil),
		Group: "exec",
		Desc:  "Do Nothing: Statement which does nothing.",
	},
	// FIX: add. returns "single value" ( aka Scalar )
	// {
	// 	Name: "Eval",
	// 	Desc: "runs a text template",
	// },
	// FIX: either filter List or add filters for each list type
	"filter": {
		Type:  (*core.Filter)(nil),
		Group: "objects",
		Desc:  "Filter Object List: A list of objects which pass the evaluation.",
	},
	"for_each_num": {
		Type:   (*core.ForEachNum)(nil),
		Group:  "exec",
		Desc:   "For Each Number: Loops over the passed list of numbers, or runs the 'else' statement if empty.",
		Locals: []string{"index", "first", "last", "num"},
	},
	"for_each_obj": {
		Type:   (*core.ForEachObj)(nil),
		Group:  "exec",
		Desc:   "For Each Object: Loops over the passed list of objects, or runs the 'else' statement if empty.",
		Locals: []string{"index", "first", "last", "obj"},
	},
	"for_each_text": {
		Type:   (*core.ForEachText)(nil),
		Group:  "exec",
		Desc:   "For Each Text: Loops over the passed list of text, or runs the 'else' statement if empty.",
		Locals: []string{"index", "first", "last", "text"},
	},
	"get": {
		Type:   (*core.Get)(nil),
		Phrase: "Get $2 of $1",
		Group:  "objects",
		Desc:   "Get Property: Return the value of an object's property.",
	},
	// for templates:
	// {
	// Type: (*core.GetAt)(nil),
	// 	Group: "objects",
	// 	Desc:  "return a named property from the current top object, or look for an object of that name instead",
	// },
	"includes": {
		Type:  (*core.Includes)(nil),
		Group: "strings",
		Desc:  "Includes Text: True if text contains text.",
	},
	// internal
	// {
	// Type: (*core.Is)(nil),
	// 	Group: "logic",
	// 	Desc:  "returns its input",
	// },
	"is_class": {
		Type:   (*core.IsClass)(nil),
		Phrase: "Is $OBJ a kind of $CLASS",
		Group:  "objects",
		Desc:   "Is Kind Of: True if the object is compatible with the named kind.",
	},
	"is_empty": {
		Type:  (*core.IsEmpty)(nil),
		Group: "strings",
		Desc:  "Is Empty: True if the text is empty.",
	},
	"is_exact_class": {
		Type:  (*core.IsExactClass)(nil),
		Group: "objects",
		Desc:  "Is Exact Kind: True if the object is exactly the named kind.",
	},
	"is_not": {
		Type:  (*core.IsNot)(nil),
		Group: "logic",
		Desc:  "Is Not: Returns the opposite value.",
	},
	"join": {
		Type:  (*core.Join)(nil),
		Group: "strings",
		Desc:  "Join Strings: Returns multiple pieces of text as a single new piece of text.",
	},
	// FIX: take "List" generically
	"len": {
		Type:  (*core.Len)(nil),
		Group: "objects",
		Desc:  "Length of Object List: Number of objects.",
	},
	"list_up": {
		Type:  (*core.ListUp)(nil),
		Group: "objects",
		Desc:  "List Up: Generates a list of objects.",
	},
	// patterns:
	// {
	// Type: (*core.LocationOf)(nil),
	// 	Group: "locale",
	// 	Desc:  "returns parent containment ",
	// },

	// functions:
	// {
	// Type: (*core.LowerAn)(nil),
	// 	Group: "format",
	// 	Desc:  "prints the object prefixed by an appropriate indefinite article",
	// },
	// {
	// Type: (*core.LowerThe)(nil),
	// 	Group: "format",
	// 	Desc:  "prints the object prefixed by the appropriate definite article",
	// },
	"remainder_of": {
		Type:   (*core.RemainderOf)(nil),
		Group:  "math",
		Phrase: "( $1 % $2 )",
		Desc:   "Modulus Numbers: Divide one number by another, and return the remainder.",
	},
	"product_of": {
		Type:   (*core.ProductOf)(nil),
		Group:  "math",
		Phrase: "( $1 * $2 )",
		Desc:   "Multiply Numbers: Multiply two numbers.",
	},
	"num_value": {
		Type:  (*core.Number)(nil),
		Group: "literals",
	},
	"numbers": {
		Type:  (*core.Numbers)(nil),
		Group: "literals",
		Desc:  "Number List: Specify a list of multiple numbers.",
	},
	"object_name": {
		Type:  (*core.ObjectName)(nil),
		Group: "objects",
		Desc:  "Named Object: Searches through the scope for a matching name.",
	},
	"object_names": {
		Type:  (*core.ObjectNames)(nil),
		Group: "objects",
		Desc:  "Object List: Searches through the scope for matching names.",
	},
	// hrmmm...
	// {
	// Type: (*core.Player)(nil),
	// 	Group: "locale",
	// 	Desc:  "returns the player's pawn",
	// },
	"pluralize": {
		Type:  (*std.Pluralize)(nil),
		Group: "format",
		Desc:  "Pluralize: Creates plural text from the passed (presumably singular) text.",
	},
	// templates:
	// {
	// Type: (*core.PrintBracket)(nil),
	// 	Group: "format",
	// 	Desc:  "sandwiches text inside parenthesis.",
	// },

	// FIX: optional or?
	"print_list": {
		Type:  (*core.PrintList)(nil),
		Group: "format",
		Desc:  "Say List: Writes words separated with commas, ending with an 'and'.",
	},
	// // FIX: patterns
	// {
	// Type: (*core.PrintNondescriptObjects)(nil),
	// 	Group: "locale",
	// 	Desc:  "PrintNondescriptObjects: prints a bunch of objects using the GroupTogether, PrintGroup, and PrintName patterns.",
	// },
	"print_num": {
		Type:  (*core.PrintNum)(nil),
		Group: "format",
		Desc:  "Say Number: Writes a number using numerals, eg. '1'.",
	},

	// patterns:
	// {
	// Type: (*core.PrintNumWord)(nil),
	// 	Group: "format",
	// 	Desc:  "NumInWords: writes a number using english: eg. 'one'.",
	// },

	// patterns:
	// {
	// 	Type:  (*std.PrintObjects)(nil),
	// 	Group: "locale",
	// 	Desc:  "PrintObjects: prints a bunch of objects with an option header, optional articles, and optionally -- in terse format. when there are no objects, prints an else clause",
	// },

	// templates
	// {
	// Type: (*core.PrintSlash)(nil),
	// 	Group: "format",
	// 	Desc:  "PrintSlash: Writes text separated with a left-leaning slash '/'",
	// },
	"print_span": {
		Type:  (*core.PrintSpan)(nil),
		Group: "format",
		Desc:  "Say Span: Writes text with spaces between words.",
	},
	"range": {
		Type:  (*core.Range)(nil),
		Group: "flow",
		Desc:  "Range of Numbers: Generates a series of numbers.",
	},
	"related_list": {
		Type:  (*core.RelatedList)(nil),
		Group: "objects",
		Desc:  "Related List: Returns a stream of objects related to the requested object.",
	},
	"relation_empty": {
		Type:  (*core.RelationEmpty)(nil),
		Group: "objects",
		Desc:  "Is Relation Empty: Returns true if the requested object has no related objects.",
	},

	// internal for templates
	// {
	// Type: (*core.Render)(nil),
	// 	Name:  "Render",
	// 	Group: "format",
	// 	Desc:  "evalutes a property value ( ex. expands a template )",
	// },

	// FIX: any list
	"reverse": {
		Type:  (*core.Reverse)(nil),
		Group: "objects",
		Desc:  "Reverse Object List: returns the listed objects, last first.",
	},

	// FIX: should take a speaker, and we should have a default speaker
	"say": {
		Type:  (*core.Say)(nil),
		Group: "format",
		Desc:  "Say: writes a piece of text.",
	},

	// FIX: set any; with type picked during assembly
	"set_bool": {
		Type:  (*core.SetBool)(nil),
		Group: "objects",
		Desc:  "Set Bool: Sets the named property to the passed boolean value.",
	},
	"set_num": {
		Type:  (*core.SetNum)(nil),
		Group: "objects",
		Desc:  "Set Number: sets the named property to the passed number.",
	},
	"set_obj": {
		Type:  (*core.SetObj)(nil),
		Group: "objects",
		Desc:  "Set Object: Sets the named property to the passed object ( reference. )",
	},
	"set_state": {
		Type:  (*core.SetState)(nil),
		Group: "objects",
		Desc:  "Set State: Sets the object to the passed state.",
	},
	"set_text": {
		Type:  (*core.SetText)(nil),
		Group: "objects",
		Desc:  "Set Text: Sets the named property to the passed string.",
	},
	"shuffle_text": {
		Type:  (*core.ShuffleText)(nil),
		Group: "format",
		Desc:  "Shuffle Text: When called multiple times returns its inputs at random.",
	},
	"stopping_text": {
		Type:  (*core.StoppingText)(nil),
		Group: "format",
		Desc:  "Stopping Text: When called multiple times returns each of its inputs in turn, sticking to the last one.",
	},
	"diff_of": {
		Type:   (*core.DiffOf)(nil),
		Group:  "math",
		Phrase: "( $1 - $2 )",
		Desc:   "Subtract Numbers: Subtract two numbers.",
	},
	"test": {
		Type:  (*check.Test)(nil),
		Slots: []string{"story_statement"},
	},
	"text_value": {
		Type: (*core.TextValue)(nil),
	},
	"texts": {
		Type:  (*core.Texts)(nil),
		Group: "literals",
		Desc:  "Text List: specifies multiple string values.",
	},

	// {
	// Type: (*core.TopObject)(nil),
	// 	Name:  "TopObject",
	// 	Group: "exec",
	// 	Desc:  "returns the top most object in scope",
	// },

	// a function or pattern
	// {
	// Type: (*core.UpperAn)(nil),
	// 	Group: "format",
	// 	Desc:  "prints the object prefixed by an appropriate indefinite article",
	// },

	// a function of pattern
	// {
	// Type: (*core.UpperThe)(nil),
	// 	Group: "format",
	// 	Desc:  "prints the object prefixed by an appropriate indefinite article",
	// },
}
