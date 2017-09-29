package pat_test

import (
	"fmt"
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/pat"
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
	if obj, ok := run.TopObject(); !ok {
		err = fmt.Errorf("no top object")
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
	rules := pat.MakeContract(patterns.Types)
	//
	unique.PanicTypes(classes, (*Num)(nil))
	// SayMe converts numbers to text
	// http://learnyouahaskell.com/syntax-in-functions
	type SayMe struct {
		Num float64
	}
	unique.PanicTypes(patterns, (*SayMe)(nil))
	rules.AddTextRule(ident.IdOf("sayMe"), MatchNumber(1), SayIt("One!"))
	rules.AddTextRule(ident.IdOf("sayMe"), MatchNumber(2), SayIt("Two!"))
	rules.AddTextRule(ident.IdOf("sayMe"), MatchNumber(3), SayIt("San!"))
	rules.AddTextRule(ident.IdOf("sayMe"), MatchNumber(3), SayIt("Three!"))
	rules.AddTextRule(ident.IdOf("sayMe"), nil, SayIt("Not between 1 and 3"))

	run, _ := rtm.New(classes).Rules(rules).Rtm()
	for i := 1; i <= 4; i++ {
		sayMe := run.Emplace(&SayMe{float64(i)})

		if text, e := run.GetTextMatching(sayMe); e != nil {
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
	rules := pat.MakeContract(patterns.Types)

	unique.PanicTypes(classes,
		(*Num)(nil))
	// Factorial computes an integer multiplied by the factorial of the integer below it.
	type Factorial struct {
		Num float64
	}
	unique.PanicTypes(patterns,
		(*Factorial)(nil))
	//
	rules.AddNumberRule(ident.IdOf("factorial"), MatchNumber(0), Int(1))
	//
	rules.AddNumberRule(ident.IdOf("factorial"), nil, GetNumber(func(run rt.Runtime) (ret float64, err error) {
		var this int
		if obj, ok := run.TopObject(); !ok {
			err = fmt.Errorf("no top object")
		} else if e := obj.GetValue("num", &this); e != nil {
			err = e
		} else if next, e := run.GetNumMatching(run.Emplace(&Factorial{float64(this - 1)})); e != nil {
			err = e
		} else {
			ret = float64(this) * next
		}
		return
	}))
	if run, e := rtm.New(classes).Rules(rules).Rtm(); assert.NoError(e) {
		if n, e := run.GetNumMatching(run.Emplace(&Factorial{3})); assert.NoError(e) {
			fac := 3 * (2 * (1 * 1))
			assert.EqualValues(fac, n)
		}
	}
}
