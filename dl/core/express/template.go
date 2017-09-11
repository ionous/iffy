package express

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/spec/ops"
	"strings"
)

// Templatize is deprecated. Its better to use expressions --
// tho, there is still a question of command handling.
func Templatize(tmpl Template, cmds *ops.Ops) (ret rt.TextEval, err error) {
	if len(tmpl) == 0 {
		ret = &core.Text{}
	} else {
		var buf rt.ExecuteList
		for _, t := range tmpl {
			if n, e := ReduceToken(t, cmds); e != nil {
				err = e
				break
			} else {
				p := &core.Say{n}
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
	} else if r, ok := Dedot(s); !ok {
		err = errutil.Fmt("couldnt create property path", s)
	} else if r, ok := r.(rt.TextEval); !ok {
		err = errutil.Fmt("couldnt convert %T to text", r)
	} else {
		ret = r
	}
	return
}

func deduceCmds(cmds *ops.Ops, parts []string) (ret rt.TextEval, err error) {
	var root struct{ rt.TextEval }
	c, _ := cmds.NewXBuilder(&root, core.Xform{})
	if c.Cmd(parts[0]).Begin() {
		for _, val := range parts[1:] {
			c.Val(val)
		}
		c.End()
	}
	if e := c.Build(); e != nil {
		err = e
	} else {
		ret = root.TextEval
	}
	return
}
