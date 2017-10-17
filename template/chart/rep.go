package chart

type Spec interface {
	specNode() // internal marker
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
	Subject    Spec
	Expression string
	Filters    []FunctionSpec
}

// TextSpec in a directive or function parameter.
type TextSpec struct{ Value string }

// NumberSpec in a directive or function parameter.
type NumberSpec struct{ Value float64 }

// ReferenceSpec in a directive or function parameter.
type ReferenceSpec struct{ Fields []string }

// FunctionSpec can act as a directive prelude.
// Expressions cannot appear in function args unless embedded in sub-directives.
type FunctionSpec struct {
	Name string
	Args []Spec
}

func (Directive) blockNode()  {}
func (ErrorBlock) blockNode() {}
func (TextBlock) blockNode()  {}

func (*Directive) specNode()     {}
func (*TextSpec) specNode()      {}
func (*NumberSpec) specNode()    {}
func (*FunctionSpec) specNode()  {}
func (*ReferenceSpec) specNode() {}

// { A | B x {y} | C }
// { `"adam!"` | capitalize! | prepend: "Hello " }
