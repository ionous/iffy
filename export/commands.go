package export

import (
	"github.com/ionous/iffy/dl/check"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/rt"
)

type Slot struct {
	Name  string
	Type  interface{} // nil instance, ex. (*core.Comparator)(nil)
	Desc  string
	Group string // display group(s)
}

var Slots = []Slot{{
	Name: "execute",
	Type: (*rt.Execute)(nil),
	Desc: "Execute: Run a series of statements.",
}, {
	Name: "bool_eval",
	Type: (*rt.BoolEval)(nil),
	Desc: "Booleans: Statements which return true/false values.",
}, {
	Name: "number_eval",
	Type: (*rt.NumberEval)(nil),
	Desc: "Numbers: Statements which return a number.",
}, {
	Name: "text_eval",
	Type: (*rt.TextEval)(nil),
	Desc: "Texts: Statements which return text.",
}, {
	Name: "num_list_eval",
	Type: (*rt.NumListEval)(nil),
	Desc: "Number List: Statements which return a list of numbers.",
}, {
	Name: "text_list_eval",
	Type: (*rt.TextListEval)(nil),
	Desc: "Text Lists: Statements which return a list of text.",
}, {
	Name: "comparator",
	Type: (*core.Comparator)(nil),
	Desc: "Comparison Types: Helper used when comparing two numbers, objects, pieces of text, etc.",
}, {
	Name: "assignment",
	Type: (*core.Assignment)(nil),
	Desc: "Assignments: Helper used when setting variables.",
}, {
	Name: "testing",
	Type: (*check.Testing)(nil),
	Desc: "Testing: Run a series of tests.",
}}

// Register exported commands used by the data language
func Register(reg func(value interface{})) {
	for _, cmd := range Slats {
		reg(cmd)
	}
}

var Slats = []composer.Specification{
	(*check.TestOutput)(nil),
	//
	(*core.AllTrue)(nil),
	(*core.AnyTrue)(nil),

	// Assign turns an Assignment a normal statement.
	(*core.Assign)(nil),
	(*core.FromBool)(nil),
	(*core.FromNum)(nil),
	(*core.FromText)(nil),
	(*core.FromNumList)(nil),
	(*core.FromTextList)(nil),

	(*core.DetermineAct)(nil),
	(*core.DetermineNum)(nil),
	(*core.DetermineText)(nil),
	(*core.DetermineBool)(nil),
	(*core.DetermineNumList)(nil),
	(*core.DetermineTextList)(nil),
	(*core.Parameters)(nil),
	(*core.Parameter)(nil),

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
