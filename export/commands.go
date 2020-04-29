package export

import (
	"github.com/ionous/iffy/dl/check"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/rt"
)

type Slot struct {
	Type  interface{} // nil instance, ex. (*core.CompareTo)(nil)
	Desc  string
	Group string // display group(s)
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
	"compare_to": {Type: (*core.CompareTo)(nil),
		Desc: "Comparison Types: Helper used when comparing two numbers, objects, pieces of text, etc."},
}

func Register(reg func(value interface{})) {
	for _, cmd := range Slats {
		reg(cmd)
	}
}

var Slats = []composer.Specification{
	(*core.Determine)(nil), // internal but needed for gob.
	(*check.Test)(nil),
	(*core.AllTrue)(nil),
	(*core.AnyTrue)(nil),

	// FIX: Choose scalar/any?
	(*core.Choose)(nil),
	(*core.ChooseNum)(nil),
	(*core.ChooseText)(nil),
	// FIX: compare scalar?
	(*core.CompareNum)(nil),
	(*core.CompareText)(nil),

	(*core.DoNothing)(nil),
	(*core.ForEachNum)(nil),
	(*core.ForEachText)(nil),

	(*core.GetField)(nil),
	(*core.GetVar)(nil),

	(*core.Is)(nil),
	(*core.IsNot)(nil),

	(*core.Bool)(nil),
	(*core.Number)(nil),
	(*core.Text)(nil),
	(*core.Numbers)(nil),
	(*core.Texts)(nil),

	(*core.SumOf)(nil),
	(*core.DiffOf)(nil),
	(*core.ProductOf)(nil),
	(*core.QuotientOf)(nil),
	(*core.RemainderOf)(nil),

	(*core.Exists)(nil),
	(*core.KindOf)(nil),
	(*core.IsKindOf)(nil),
	(*core.IsExactKindOf)(nil),

	(*core.PrintNum)(nil),
	(*core.PrintNumWord)(nil),

	// FIX: take "List" generically
	(*core.LenOfNumbers)(nil),
	(*core.LenOfTexts)(nil),
	(*core.Range)(nil),

	// FIX: should take a speaker, and we should have a default speaker
	(*core.Say)(nil),
	(*core.Buffer)(nil),
	(*core.Span)(nil),
	(*core.Bracket)(nil),
	(*core.Slash)(nil),
	(*core.Commas)(nil),

	(*core.CycleText)(nil),
	(*core.ShuffleText)(nil),
	(*core.StoppingText)(nil),

	// FIX: set any; with the type picked during assembly
	(*core.SetFieldBool)(nil),
	(*core.SetFieldNum)(nil),
	(*core.SetFieldText)(nil),
	(*core.SetFieldNumList)(nil),
	(*core.SetFieldTextList)(nil),

	(*core.SetVarBool)(nil),
	(*core.SetVarNum)(nil),
	(*core.SetVarText)(nil),
	(*core.SetVarNumList)(nil),
	(*core.SetVarTextList)(nil),

	(*core.IsEmpty)(nil),
	(*core.Includes)(nil),
	(*core.Join)(nil),

	// comparison
	(*core.EqualTo)(nil),
	(*core.NotEqualTo)(nil),
	(*core.GreaterThan)(nil),
	(*core.LessThan)(nil),
	(*core.GreaterOrEqual)(nil),
	(*core.LessOrEqual)(nil),
}
