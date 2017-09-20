package template

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/spec"
	"strings"
)

type conditional struct {
	tcontext
	tprev
	tdepth
}

func (n conditional) advance(c spec.Block, t Token) (ret tstate, err error) {
	// println("condition", n.depth, t.Str)
	if plain(c, t) {
		ret = n
	} else {
		if t.Str == "end" {
			if prev, e := n.pop(); e != nil {
				err = e
			} else {
				n.rollup(c)
				ret = prev
			}
		} else {
			parts := strings.Fields(t.Str)
			switch op, rest := parts[0], parts[1:]; op {
			case "else", "otherwise":
				endJoin(c)
				if res, e := n.continuation(c, rest, true); e != nil {
					err = errutil.New(op, e, t.Pos)
				} else {
					startJoin(c)
					ret = res
				}
			case "unless":
				endJoin(c)
				if res, e := n.continuation(c, rest, false); e != nil {
					err = errutil.New(op, e, t.Pos)
				} else {
					startJoin(c)
					ret = res
				}
			default:
				ret, err = n.defaultAdvance(n, c, t)
			}
		}
	}
	return
}

func (n conditional) continuation(c spec.Block, rest []string, useIf bool) (ret tstate, err error) {
	if len(rest) == 0 {
		// println("terminal", n.tdepth)
		ret = terminal{n.tcontext, tprev{n}, n.tdepth}
	} else if e := n.condition(c, rest, useIf); e != nil {
		err = e
	} else {
		// println("------", n.tdepth)
		ret = conditional{n.tcontext, tprev{n}, n.tdepth + 1}
	}
	return
}
