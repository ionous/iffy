package pat_test

import (
	"fmt"
	"github.com/ionous/iffy/pat"
	"github.com/ionous/iffy/pat/patbuilder"
	"github.com/ionous/iffy/ref"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rtm"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MatchNumber int

func (m MatchNumber) GetBool(run rt.Runtime) (okay bool, err error) {
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
	classes := ref.NewClasses()
	patterns := patbuilder.NewPatterns(classes)
	//
	unique.RegisterTypes(unique.PanicTypes(classes),
		(*Num)(nil))
	objects := ref.NewObjects(classes)
	run := rtm.New(classes).Objects(objects).NewRtm()
	// SayMe converts numbers to text
	// http://learnyouahaskell.com/syntax-in-functions
	type SayMe struct {
		Num float64
	}
	if e := unique.RegisterTypes(patterns, (*SayMe)(nil)); e != nil {
		fmt.Println("new pat:", e)
	} else if e := patterns.AddText("sayMe", MatchNumber(1), SayIt("One!")); e != nil {
		fmt.Println("add one:", e)
	} else if e := patterns.AddText("sayMe", MatchNumber(2), SayIt("Two!")); e != nil {
		fmt.Println("add two:", e)
	} else if e := patterns.AddText("sayMe", MatchNumber(3), SayIt("San!")); e != nil {
		fmt.Println("add san:", e)
	} else if e := patterns.AddText("sayMe", MatchNumber(3), SayIt("Three!")); e != nil {
		fmt.Println("add three:", e)
	} else if e := patterns.AddText("sayMe", nil, SayIt("Not between 1 and 3")); e != nil {
		fmt.Println("add default:", e)
	} else {
		p := patterns.Build()
		for i := 1; i <= 4; i++ {
			if sayMe, e := run.Emplace(&SayMe{float64(i)}); e != nil {
				fmt.Println("emplace:", e)
				break
			} else if text, e := p.GetTextMatching(run, sayMe); e != nil {
				fmt.Println("matching:", e)
				break
			} else {
				fmt.Println(fmt.Sprintf("sayMe %d = \"%s\"", i, text))
			}
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
	classes := ref.NewClasses()
	patterns := patbuilder.NewPatterns(classes)
	unique.RegisterTypes(unique.PanicTypes(classes),
		(*Num)(nil))
	objects := ref.NewObjects(classes)
	// Factorial computes an integer multiplied by the factorial of the integer below it.
	type Factorial struct {
		Num float64
	}
	unique.RegisterTypes(unique.PanicTypes(patterns),
		(*Factorial)(nil))
	//
	var p pat.Patterns
	if e := patterns.AddNumber("factorial", MatchNumber(0), Int(1)); assert.NoError(e) {
		//
		if e := patterns.AddNumber("factorial", nil, GetNumber(func(run rt.Runtime) (ret float64, err error) {
			var this int
			if obj, ok := run.FindObject("@"); !ok {
				err = fmt.Errorf("context not found")
			} else if e := obj.GetValue("num", &this); e != nil {
				err = e
			} else if fact, e := run.Emplace(&Factorial{float64(this - 1)}); e != nil {
				err = e
			} else if next, e := p.GetNumMatching(run, fact); e != nil {
				err = e
			} else {
				ret = float64(this) * next
			}
			return
		})); assert.NoError(e) {
			// suite?
			run := rtm.New(classes).Objects(objects).NewRtm()
			p = patterns.Build()
			//
			if fact, e := run.Emplace(&Factorial{3}); assert.NoError(e) {
				if n, e := p.GetNumMatching(run, fact); assert.NoError(e) {
					fac := 3 * (2 * (1 * 1))
					assert.EqualValues(fac, n)
				}
			}
		}
	}
}
