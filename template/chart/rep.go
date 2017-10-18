package chart

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
	Filters    []FunctionArg
}

// QuotedArg in a directive or function parameter.
type QuotedArg struct{ Value string }

// NumberArg in a directive or function parameter.
type NumberArg struct{ Value float64 }

// ReferenceArg in a directive or function parameter.
type ReferenceArg struct{ Fields []string }

// FunctionArg can act as a directive prelude.
// Expressions cannot appear in function args unless embedded in sub-directives.
type FunctionArg struct {
	Name string
	Args []Argument
}

func (Directive) blockNode()  {}
func (ErrorBlock) blockNode() {}
func (TextBlock) blockNode()  {}

func (*Directive) argNode()    {}
func (*QuotedArg) argNode()    {}
func (*NumberArg) argNode()    {}
func (*FunctionArg) argNode()  {}
func (*ReferenceArg) argNode() {}

// { A | B x {y} | C }
// { `"adam!"` | capitalize! | prepend: "Hello " }
