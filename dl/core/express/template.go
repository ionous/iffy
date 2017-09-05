package express

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/spec/ops"
	"strings"
	"unicode"
)

// Templatize
func Templatize(tmpl []Token, cmds *ops.Ops) (ret rt.TextEval, err error) {
	if len(tmpl) == 0 {
		ret = &core.Text{}
	} else {
		var buf rt.ExecuteList
		for _, t := range tmpl {
			if n, e := ReduceToken(t, cmds); e != nil {
				err = e
				break
			} else {
				p := &core.PrintText{n}
				buf = append(buf, p)
			}
		}
		ret = &core.Buffer{buf}
	}
	return
}

func ReduceToken(t Token, cmds *ops.Ops) (ret rt.TextEval, err error) {
	if s := t.Str; t.Plain {
		ret = &core.Text{s}
	} else if parts := strings.Fields(s); len(parts) > 1 {
		ret, err = deduceCmds(cmds, parts)
	} else if dots := strings.FieldsFunc(s, dotFields); len(dots) > 0 {
		ret, err = deduceAccess(dots)
	}
	return
}

func deduceCmds(cmds *ops.Ops, parts []string) (ret rt.TextEval, err error) {
	var root struct{ rt.TextEval }
	b, _ := cmds.NewBuilder(&root)
	if b.Cmd(parts[0]).Begin() {
		for _, val := range parts[1:] {
			b.Val(val)
		}
		b.End()
	}
	if e := b.Build(); e != nil {
		err = e
	} else {
		ret = root.TextEval
	}
	return
}

// given a dotted name separted into parts generate a text eval accessor
// (ex. "example", "example.property")
// note: it doesnt check if the object exists, and it doesnt check if the property exists.
func deduceAccess(dots []string) (ret rt.TextEval, err error) {
	// is upper means a global object
	isUpper := isCapitalized(dots[0])
	// length of 1 means no dots at all:
	if cnt := len(dots); cnt == 1 {
		if !isUpper {
			ret = &core.Get{&core.Object{"@"}, dots[0]}
		} else {
			// but we have no property accessor --
			// and, for reasons, we cant print raw name.
			err = errutil.New("cant say object", dots[0])
		}
	} else {
		var obj rt.ObjectEval
		if isUpper {
			obj = &core.Global{dots[0]}
		} else {
			obj = &core.Object{dots[0]}
		}
		for _, d := range dots[1:] {
			get := &core.Get{obj, d}
			obj, ret = get, get
		}
	}
	return
}

func dotFields(r rune) bool {
	return r == '.'
}

// return true if the passed string starts with an upper case letter
func isCapitalized(n string) (ret bool) {
	for _, r := range n {
		ret = unicode.IsUpper(r)
		break
	}
	return
}
