package export

import (
	"github.com/ionous/iffy/dl/check"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/next"
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
	"num_list_eval": {Type: (*rt.NumListEval)(nil),
		Desc: "Number List: Statements which return a list of numbers."},
	"text_list_eval": {Type: (*rt.TextListEval)(nil),
		Desc: "Text Lists: Statements which return a list of text."},
	"compare_to": {Type: (*next.CompareTo)(nil),
		Desc: "Comparison Types: Helper used when comparing two numbers, objects, pieces of text, etc."},
}

var Runs = []composer.Specification{
	(*check.Test)(nil),
	(*next.Determine)(nil), // internal but needed for gob.

	(*next.AllTrue)(nil),
	(*next.AnyTrue)(nil),

	// FIX: Choose scalar/any?
	(*next.Choose)(nil),
	(*next.ChooseNum)(nil),
	(*next.ChooseText)(nil),
	// FIX: compare scalar?
	(*next.CompareNum)(nil),
	(*next.CompareText)(nil),

	(*next.DoNothing)(nil),
	(*next.ForEachNum)(nil),
	(*next.ForEachText)(nil),

	(*next.GetField)(nil),

	(*next.GetVar)(nil),

	(*next.Is)(nil),
	(*next.IsNot)(nil),

	(*next.Bool)(nil),
	(*next.Number)(nil),
	(*next.Text)(nil),
	(*next.Numbers)(nil),
	(*next.Texts)(nil),

	(*next.SumOf)(nil),
	(*next.DiffOf)(nil),
	(*next.ProductOf)(nil),
	(*next.QuotientOf)(nil),
	(*next.RemainderOf)(nil),

	(*next.Exists)(nil),
	(*next.KindOf)(nil),
	(*next.IsKindOf)(nil),
	(*next.IsExactKindOf)(nil),

	(*next.PrintNum)(nil),
	(*next.PrintNumWord)(nil),

	// FIX: take "List" generically
	(*next.LenOfNumbers)(nil),
	(*next.LenOfTexts)(nil),
	(*next.Range)(nil),

	// FIX: should take a speaker, and we should have a default speaker
	(*next.Say)(nil),
	(*next.Buffer)(nil),
	(*next.Span)(nil),
	(*next.Bracket)(nil),
	(*next.Slash)(nil),
	(*next.Commas)(nil),

	(*next.CycleText)(nil),
	(*next.ShuffleText)(nil),
	(*next.StoppingText)(nil),

	// FIX: set any; with type picked during assembly
	(*next.SetFieldBool)(nil),
	(*next.SetFieldNum)(nil),
	(*next.SetFieldText)(nil),
	(*next.SetFieldNumList)(nil),
	(*next.SetFieldTextList)(nil),

	(*next.SetVarBool)(nil),
	(*next.SetVarNum)(nil),
	(*next.SetVarText)(nil)
	(*next.SetVarNumList)(nil),
	(*next.SetVarTextList)(nil),

	(*next.IsEmpty)(nil),
	(*next.Includes)(nil),
	(*next.Join)(nil),
}
