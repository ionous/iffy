package chart

import (
	"fmt"
	"strings"
)

type Argument interface {
	argNode() // internal marker
}

// Block represents some text or a directive.
type Block interface {
	blockNode()
}

// Template contains text and directives.
type Template struct {
	Blocks []Block
}

// Error acts a block in a template.
type ErrorBlock struct{ error }

// TextBlock in a template.
type TextBlock struct{ Text string }

// Directive can be used as a function argument.
type Directive struct {
	Subject    Argument
	Expression string
	Filters    []Function
}

// Quote in a directive or function parameter.
type Quote struct{ Value string }

// Number in a directive or function parameter.
type Number struct{ Value float64 }

// Reference in a directive or function parameter.
type Reference struct{ Fields []string }

// Function can act as a directive prelude.
// Expressions cannot appear in function args unless embedded in sub-directives.
type Function struct {
	Name string
	Args []Argument
}

func (Directive) blockNode()  {}
func (ErrorBlock) blockNode() {}
func (TextBlock) blockNode()  {}

func (*Directive) argNode() {}
func (*Quote) argNode()     {}
func (*Number) argNode()    {}
func (*Function) argNode()  {}
func (*Reference) argNode() {}

func (d *Directive) String() string {
	var fs []string
	for _, n := range d.Filters {
		fs = append(fs, n.String())
	}
	x := d.Expression
	if len(x) > 0 {
		x = fmt.Sprintf("[%s]", x)
	}
	p := strings.Join(fs, " ")
	if len(p) > 0 {
		p = fmt.Sprintf("(%s)", p)
	}
	return fmt.Sprint("{", d.Subject, x, p, "}")
}
func (q *Quote) String() string {
	return fmt.Sprint("quote:", q.Value)
}
func (n *Number) String() string {
	return fmt.Sprintf("num:'%g'", n.Value)
}
func (f *Function) String() string {
	var fs []string
	for _, n := range f.Args {
		fs = append(fs, fmt.Sprint(n))
	}
	p := strings.Join(fs, ",")
	if len(p) > 0 {
		p = fmt.Sprintf("(%s)", p)
	}
	return fmt.Sprint("call:", f.Name, p)
}
func (r *Reference) String() string {
	return fmt.Sprintf("ref:%d'%v'", len(r.Fields), strings.Join(r.Fields, "."))
}
