package express

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/lang"
	"github.com/ionous/iffy/rt"
	"strings"
)

// Dedot takes an object path and turns it into object access.
// (ex. "Example", "example", "example.property")
// If the first identifier in the path is upper-case, the identifier names a global object,
// otherwise it is allowed to identify a property in the local scope.
// Dedot does not guarantee that the object or the named properties exist.
func Dedot(path string) (ret interface{}, okay bool) {
	if dots := strings.FieldsFunc(path, dotFields); len(dots) > 0 {
		var obj rt.ObjectEval
		first, rest := dots[0], dots[1:]
		if isUpper := lang.IsCapitalized(first); isUpper {
			obj = &core.Global{first}
		} else {
			obj = &core.GetAt{first}
		}
		for _, d := range rest {
			// note: for convenience, we store Get{} in an object eval;
			// but, Get{} implements many different kinds of evals.
			obj = &core.Get{obj, d}
		}
		ret, okay = obj, true
	}
	return
}

func dotFields(r rune) bool {
	return r == '.'
}
