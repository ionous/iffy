package chart

import (
	"bytes"
	"fmt"
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/template"
	"github.com/ionous/iffy/template/postfix"
)

// Directive containing the parsed content of a template.
// Both or either of the key and the expression can be empty.
type Directive struct {
	Key string
	postfix.Expression
}

func UnknownDirective(v Directive) error {
	return errutil.Fmt("unknown key '%s'", v.Key)
}

func UnexpectedExpression(v Directive) (err error) {
	if len(v.Expression) > 0 {
		err = errutil.Fmt("unexpected expression following key '%s'", v.Key)
	}
	return
}
func ExpectedExpression(v Directive) (err error) {
	if len(v.Expression) == 0 {
		err = errutil.Fmt("expected expression following key '%s'", v.Key)
	}
	return
}

// String of a directive in the format:
// {key:expression} or {expression}
func (d Directive) String() (ret string) {
	if len(d.Key) > 0 {
		ret = fmt.Sprintf("{%s:%s}", d.Key, d.Expression)
	} else if q, ok := d.isQuote(); ok {
		ret = string(q)
	} else {
		ret = fmt.Sprintf("{%s}", d.Expression)
	}
	return
}

func (d Directive) isQuote() (ret template.Quote, okay bool) {
	if cnt := len(d.Expression); cnt == 1 {
		ret, okay = d.Expression[0].(template.Quote)
	}
	return
}

// Format a string from slice of directives.
func Format(ds []Directive) string {
	var buf bytes.Buffer
	for _, d := range ds {
		buf.WriteString(fmt.Sprint(d))
	}
	return buf.String()
}
