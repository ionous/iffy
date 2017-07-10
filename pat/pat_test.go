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

// http://learnyouahaskell.com/syntax-in-functions
func ExampleSayMe() {
	b := patbuilder.NewBuilder()
	classes := ref.NewClasses()
	unique.RegisterTypes(unique.PanicTypes(classes),
		(*Num)(nil))
	objects := ref.NewObjects(classes)
	run := rtm.NewRtm(classes, objects, nil)
	//
	if e := b.NewPattern("sayMe"); e != nil {
		fmt.Println("new pat:", e)
	} else if e := b.AddMatch("sayMe", SayIt("One!"), MatchNumber(1)); e != nil {
		fmt.Println("add one:", e)
	} else if e := b.AddMatch("sayMe", SayIt("Two!"), MatchNumber(2)); e != nil {
		fmt.Println("add two:", e)
	} else if e := b.AddMatch("sayMe", SayIt("Three!"), MatchNumber(3)); e != nil {
		fmt.Println("add three:", e)
	} else if e := b.AddMatch("sayMe", SayIt("Not between 1 and 3")); e != nil {
		fmt.Println("add default:", e)
	} else {
		p := b.GetPatterns()
		for i := 1; i <= 4; i++ {
			if obj, e := objects.Emplace(Int(i)); e != nil {
				fmt.Println("emplace:", e)
				break
			} else if t, e := p.GetTextMatching(run, "sayMe", obj); e != nil {
				fmt.Println("matching:", e)
				break
			} else {
				fmt.Println(fmt.Sprintf("sayMe %d = \"%s\"", i, t))
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
	unique.RegisterTypes(unique.PanicTypes(classes),
		(*Num)(nil))
	objects := ref.NewObjects(classes)
	b := patbuilder.NewBuilder()
	//
	var p pat.Patterns
	if e := b.NewPattern("factorial"); assert.NoError(e) {
		//
		if e := b.AddMatch("factorial", Int(1), MatchNumber(0)); assert.NoError(e) {
			//
			if e := b.AddMatch("factorial", GetNumber(func(run rt.Runtime) (ret float64, err error) {
				var this int
				if obj, ok := run.FindObject("@"); !ok {
					err = fmt.Errorf("context not found")
				} else if e := obj.GetValue("num", &this); e != nil {
					err = e
				} else if next, e := objects.Emplace(Int(this - 1)); e != nil {
					err = e
				} else if next, e := p.GetNumMatching(run, "factorial", next); e != nil {
					err = e
				} else {
					ret = float64(this) * next
				}
				return
			})); assert.NoError(e) {
				// suite?
				run := rtm.NewRtm(classes, objects, nil)
				p = b.GetPatterns()
				//
				if obj, e := objects.Emplace(Int(3)); assert.NoError(e) {
					if n, e := p.GetNumMatching(run, "factorial", obj); assert.NoError(e) {
						fac := 3 * (2 * (1 * 1))
						assert.EqualValues(fac, n)
					}
				}
			}
		}
	}
}
