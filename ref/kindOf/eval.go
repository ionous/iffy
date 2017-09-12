package kindOf

import (
	r "reflect"
)

// EvalType determines what kind of eval can produce the passed type.
func EvalType(rtype r.Type) (ret r.Type) {
	switch k := rtype.Kind(); {
	//
	case k == r.Bool:
		ret = boolEval

	case Number(rtype):
		ret = numEval

	case IdentId(rtype):
		ret = objEval

	case k == r.String:
		ret = textEval

	case k == r.Array || k == r.Slice:
		elem := rtype.Elem()
		switch k := elem.Kind(); {

		case Number(elem):
			ret = numListEval

		case IdentId(elem):
			ret = objListEval

		case k == r.String:
			ret = textListEval

		}
	}
	return
}
