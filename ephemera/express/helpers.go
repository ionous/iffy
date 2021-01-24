package express

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/render"
	"github.com/ionous/iffy/rt"
)

func T(s string) rt.TextEval {
	return &core.Text{s}
}
func N(n float64) rt.NumberEval {
	return &core.Number{n}
}
func B(b bool) rt.BoolEval {
	return &core.Bool{b}
}
func O(n string, exact bool) (ret rt.TextEval) {
	if !exact {
		ret = &render.RenderRef{core.Var{n}, render.TryAsBoth}
	} else {
		ret = T(n)
	}
	return ret
}
