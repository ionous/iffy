package postfix

import (
	"testing"
	"unicode"
)

func TestShunt(t *testing.T) {
	succeed(t, "", "")
	succeed(t, "x", "x")
	succeed(t, "x+y", "xy+")
	succeed(t, "a+b*c-d", "abc*+d-")
	succeed(t, "x+y*w", "xyw*+")
	succeed(t, "(x+y)*w", "xy+w*")
	succeed(t, "Fa", "aF")
	succeed(t, "Frx|Gst", "strxFG")
	succeed(t, "Fa|Gb|Hc", "cbaFGH")
	succeed(t, "Fr((x+y)*w)|Gst", "strxy+w*FG")
	succeed(t, "(x+y)", "xy+")
	succeed(t, "(x+y)*(z)", "xy+z*")
	fail(t, "(x+y))") // too many ends
	fail(t, "((x+y)") // unclosed ends
	// we dont have a good way to reliably detect empty statements
	// fail(t, "((x+y))") // empty statement
}

func fail(t *testing.T, in string) {
	t.Log("in:", in)
	if p, e := build(in); e != nil {
		t.Log("okay", e)
	} else {
		t.Fatal("unexpected success", p)
	}
}

func succeed(t *testing.T, in string, want string) {
	t.Log("in:", in)
	if p, e := build(in); e != nil {
		t.Fatal("error building expression", e)
	} else {
		res := p.String()
		if res != want {
			t.Log("expected:", want)
			t.Fatal("got:", len(p), res)
		} else {
			t.Log("ok:", res)
		}
	}
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

			t.AddFunction(f)
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
		t.BeginSubExpression()
	case ')':
		t.EndSubExpression()
	default:
		if !unicode.IsLetter(r) {
			panic("unknown symbol" + string(r))
		} else if unicode.IsUpper(r) {
			// note: the arity isnt correct here; this is for mockup only.
			ret = mockop{char: r, arity: 3}
		} else {
			ret = mockop{char: r}
		}
	}
	return
}
