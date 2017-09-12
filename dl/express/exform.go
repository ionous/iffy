package express

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/spec"
	"github.com/ionous/iffy/spec/ops"
	// "github.com/kr/pretty"
	"go/parser"
	r "reflect"
)

type Xform struct {
	cmds *ops.Ops
	core.Xform
}

func MakeXform(cmds *ops.Ops) Xform {
	return Xform{cmds: cmds}
}

// TransformValue returns src if no error but couldnt convert.
func (xf Xform) TransformValue(val interface{}, hint r.Type) (ret interface{}, err error) {
	if t, ok := tryTokenize(val); ok {
		ret, err = xf.TransformTemplate(t, hint)
	} else {
		ret, err = xf.Xform.TransformValue(val, hint)
	}
	return
}

// TransformTemplate
func (xf Xform) TransformTemplate(tmpl Template, hint r.Type) (ret interface{}, err error) {
	t := ops.NewValue(hint)
	c := xf.cmds.NewFromTarget(t, xf)
	//
	if len(tmpl) == 1 {
		err = xf.convertOne(c, tmpl[0])
	} else {
		err = xf.convertMulti(c, tmpl)
	}
	//
	if err == nil {
		if e := c.Build(); e != nil {
			err = e
		} else {
			ret = t.Field(0).Interface()
		}
	}

	return
}

// one evaluation
func (xf Xform) convertOne(c spec.Block, token Token) (err error) {
	// look for and chomp templates starting with {go}
	if g, ok := token.Go(); !ok {
		err = parseExpr(c, token.Str)
	} else {
		if c.Cmd(g[0]).Begin() {
			err = parseExprs(c, g[1:])
			c.End()
		}
	}
	return
}

// convert multiple evaluations.
// by definition the evaluations separated by plain text.
// ( if there were no text seperators, it would be one evaluation )
func (xf Xform) convertMulti(c spec.Block, tmpl Template) (err error) {
	// because we are mixing text and evals, we expect the whole thing winds up being text. ( otherwise: what would we do with the intervening text. )
	if c.Cmd("buffer").Begin() {
		if c.Cmds().Begin() {
			for _, token := range tmpl {
				if c.Cmd("say").Begin() {
					if token.Plain {
						c.Val(token.Str)
					} else if e := xf.convertOne(c, token); e != nil {
						err = e
						break
					}
					c.End()
				}
			}
			c.End()
		}
		c.End()
	}
	return
}

// add the passed strings as parsed expressions to the passed c
// we could probably support keywords using name:something in the string.
func parseExprs(c spec.Block, strs []string) (err error) {
	for _, s := range strs {
		if e := parseExpr(c, s); e != nil {
			err = e
			break
		}
	}
	return
}

func parseExpr(c spec.Block, s string) (err error) {
	// println("converting expression", s)
	if a, e := parser.ParseExpr(s); e != nil {
		err = e
	} else {
		// println(pretty.Sprint(a))
		err = convertExpr(c, a, true)
	}
	return
}

// tryTokenize attempt to turn the passed val into a string template.
func tryTokenize(val interface{}) (ret Template, okay bool) {
	if s, ok := val.(string); ok {
		ret, okay = Tokenize(s)
	}
	return
}
