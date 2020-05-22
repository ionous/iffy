package chart

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/template/postfix"
	"github.com/ionous/iffy/template/types"
)

// FieldParser reads identifiers separated by dots,
// implements OperandState.
type FieldParser struct {
	err     error
	fields  []string
	pending bool
}

func (p *FieldParser) StateName() string {
	return "fields"
}

// NewRune starts on the dot of a new field.
func (p *FieldParser) NewRune(r rune) (ret State) {
	var id IdentParser
	identChain := MakeChain(&id, Statement("post  field", func(r rune) (ret State) {
		// after we have parsed the identifier, check the incoming rune.
		if n := id.Identifier(); len(n) > 0 {
			p.fields = append(p.fields, n)
			if isDot(r) {
				p.pending = true // a new reference is pending.
				ret = p
			} else {
				p.pending = false
			}
		}
		return
	}))
	// the field was started with a dot that we are asked to parse
	// on subsequent loops we are left just after the dot, and we need that first letter.
	if len(p.fields) != 0 {
		ret = identChain.NewRune(r)
	} else if !isDot(r) {
		p.err = errutil.New("field should start with dot (.)")
	} else {
		ret = identChain
	}
	return
}

// GetFields returns an array of parsed identifiers.
func (p *FieldParser) GetFields() (ret []string, err error) {
	if e := p.err; e != nil {
		err = e
	} else if p.pending {
		err = errutil.New("incomplete fields")
	} else {
		ret = p.fields
	}
	return
}

func (p *FieldParser) GetOperand() (ret postfix.Function, err error) {
	if r, e := p.GetFields(); e != nil {
		err = e
	} else {
		ret = types.Reference(r)
	}
	return
}
