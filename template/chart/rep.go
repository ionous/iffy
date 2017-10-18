package chart

import (
	"fmt"
	"strings"
)

// Argument of a function or the subject of a directive.
type Argument interface {
	argNode() // internal marker
}

// Block of text, a directive, or error.
type Block interface {
	blockNode()
}

// ErrorBlock appears in a template when a directive fails to parse.
type ErrorBlock struct{ error }

// TextBlock contains the uninterpreted parts of a template.
type TextBlock struct{ Text string }

// Directive contain instructions which templates turn into code.
type Directive struct {
	Subject    Argument   // The initial value computed by the directive.
	Expression string     // Unparsed text after the main subject; used for modifying the value of the subject via golang expressions.
	Filters    []Function // Chains of function calls. The result of the subject, modified by the expression text, gets passed as the last parameter of the first function. Its result is passed sa the last parameter of the next function, and so on till the end.
}

// Quote in a directive or function parameter.
type Quote struct{ Value string }

// Number in a directive or function parameter.
type Number struct{ Value float64 }

// Reference in a directive or function parameter.
type Reference struct{ Fields []string }

// Function call in a directive or filter chain.
// Functions cannot be used as parameters of other functions unless embedded in sub-directives.
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
