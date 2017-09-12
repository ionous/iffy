package pat_test

import (
	"fmt"
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/pat"
	"github.com/ionous/iffy/pat/rule"
	"github.com/ionous/iffy/ref/obj"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rtm"
	"github.com/stretchr/testify/assert"
	"testing"
)

type matchNumber int

func MatchNumber(n int) pat.Filters {
	return pat.Filters{matchNumber(n)}
}

func (m matchNumber) GetBool(run rt.Runtime) (okay bool, err error) {
	var n int
	if obj, ok := run.FindObject("@"); !ok {
		err = fmt.Errorf("context not found")
	} else if e := obj.GetValue("num", &n); e != nil {
		err = e
	} else {
		okay = n == int(m)
	}
	return
}

type SayIt string

func (s SayIt) GetText(run rt.Runtime) (string, error) {
	return string(s), nil
}

// Num specifies a number value.
type Num struct {
	Num float64
}

func Int(i int) *Num {
	return &Num{float64(i)}
}

// GetNumber implements NumberEval providing the dl with a number literal.
func (n *Num) GetNumber(rt.Runtime) (float64, error) {
	return n.Num, nil
}

func ExampleSayMe() {
	classes := make(unique.Types)
	patterns := unique.NewStack(classes)
	rules := rule.MakeRules()
	//
	unique.PanicTypes(classes,
		(*Num)(nil))
	objects := obj.NewObjects()
	run := rtm.New(classes).Objects(objects).Rtm()
	// SayMe converts numbers to text
	// http://learnyouahaskell.com/syntax-in-functions
	type SayMe struct {
		Num float64
	}
	unique.PanicTypes(patterns, (*SayMe)(nil))
	rules.Text.AddRule(ident.IdOf("sayMe"), MatchNumber(1), SayIt("One!"))
	rules.Text.AddRule(ident.IdOf("sayMe"), MatchNumber(2), SayIt("Two!"))
	rules.Text.AddRule(ident.IdOf("sayMe"), MatchNumber(3), SayIt("San!"))
	rules.Text.AddRule(ident.IdOf("sayMe"), MatchNumber(3), SayIt("Three!"))
	rules.Text.AddRule(ident.IdOf("sayMe"), nil, SayIt("Not between 1 and 3"))
	rules.Sort()

	for i := 1; i <= 4; i++ {
		sayMe := run.Emplace(&SayMe{float64(i)})

		if text, e := rules.GetTextMatching(run, sayMe); e != nil {
			fmt.Println("matching:", e)
			break
		} else {
			fmt.Println(fmt.Sprintf("sayMe %d = \"%s\"", i, text))
		}
	}

	// Output:
	// sayMe 1 = "One!"
	// sayMe 2 = "Two!"
	// sayMe 3 = "Three!"
	// sayMe 4 = "Not between 1 and 3"
}

type GetNumber func(rt.Runtime) (float64, error)

func (f GetNumber) GetNumber(run rt.Runtime) (float64, error) {
	return f(run)
}

func TestFactorial(t *testing.T) {
	assert := assert.New(t)
	//
	classes := make(unique.Types)
	patterns := unique.NewStack(classes)
	rules := rule.MakeRules()

	unique.PanicTypes(classes,
		(*Num)(nil))
	objects := obj.NewObjects()
	// Factorial computes an integer multiplied by the factorial of the integer below it.
	type Factorial struct {
		Num float64
	}
	unique.PanicTypes(patterns,
		(*Factorial)(nil))
	//
	rules.Numbers.AddRule(ident.IdOf("factorial"), MatchNumber(0), Int(1))
	//
	rules.Numbers.AddRule(ident.IdOf("factorial"), nil, GetNumber(func(run rt.Runtime) (ret float64, err error) {
		var this int
		if at, ok := run.FindObject("@"); !ok {
			err = fmt.Errorf("context not found")
		} else if e := at.GetValue("num", &this); e != nil {
			err = e
		} else if next, e := rules.GetNumMatching(run, run.Emplace(&Factorial{float64(this - 1)})); e != nil {
			err = e
		} else {
			ret = float64(this) * next
		}
		return
	}))
	// suite?
	run := rtm.New(classes).Objects(objects).Rules(rules).Rtm()
	//
	if n, e := run.GetNumMatching(run, run.Emplace(&Factorial{3})); assert.NoError(e) {
		fac := 3 * (2 * (1 * 1))
		assert.EqualValues(fac, n)
	}
}
