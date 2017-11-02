package postfix

import (
	"testing"
	"unicode"
)

func TestShunt(t *testing.T) {
	test := func(in string, want string) {
		t.Log("in:", in)
		if p, e := build(in); e != nil {
			t.Fatal("error building expression", e)
		} else {
			res := p.String()
			if res != want {
				t.Log("expected:", want)
				t.Fatal("got:", res)
			} else {
				t.Log("ok:", res)
			}
		}
	}
	test("x+y", "xy+")
	test("a+b*c-d", "abc*+d-")
	test("x+y*w", "xyw*+")
	test("(x+y)*w", "xy+w*")
	test("Fa", "aF")
	test("Frx|Gst", "strxFG")
	test("Fa|Gb|Hc", "cbaFGH")
	test("Fr((x+y)*w)|Gst", "strxy+w*FG")
}

type mockop struct {
	char  rune
	arity int
	pred  int // https://golang.org/ref/spec#Operator_precedence
}

func (m mockop) Name() string    { return string(m.char) }
func (m mockop) Arity() int      { return m.arity }
func (m mockop) Precedence() int { return m.pred }
func (m mockop) String() string  { return m.Name() }

var times = mockop{'*', 2, 5}
var plus = mockop{'+', 2, 4}
var minus = mockop{'-', 2, 4}

func build(sym string) (ret Expression, err error) {
	var t Test
	for _, r := range sym {
		if f, e := t.AddRune(r); e != nil {
			err = e
			break
		} else if f != nil {
			if e := t.AddFunction(f); e != nil {
				err = e
				break
			}
		}
	}
	if err == nil {
		if exp, e := t.Flush(); e != nil {
			err = e
		} else {
			ret = exp
		}
	}
	return
}

type Test struct {
	Pipe
}

func (t *Test) AddRune(r rune) (ret Function, err error) {
	switch r {
	case '|':
		err = t.AddPipe()
	case '+':
		ret = plus
	case '-':
		ret = minus
	case '*':
		ret = times
	case '(':
		if e := t.BeginSubExpression(); e != nil {
			err = e
		}
	case ')':
		if e := t.EndSubExpression(); e != nil {
			err = e
		}
	default:
		if !unicode.IsLetter(r) {
			panic("unknown symbol" + string(r))
		}
		if unicode.IsUpper(r) {
			// note: the arity isnt correct here; this is for mockup only.
			ret = mockop{char: r, arity: 3}
		} else {
			ret = mockop{char: r}
		}
	}
	return
}
