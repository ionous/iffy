package kindOf

import (
	r "reflect"
)

// EvalType determines what kind of eval can produce the passed type.
func EvalType(rtype r.Type) (ret r.Type) {
	switch k := rtype.Kind(); {
	//
	case k == r.Bool:
		ret = TypeBoolEval

	case Number(rtype):
		ret = TypeNumEval

	case k == r.String:
		ret = TypeTextEval

	case k == r.Array || k == r.Slice:
		elem := rtype.Elem()
		switch k := elem.Kind(); {
		case Number(elem):
			ret = TypeNumListEval

		case k == r.String:
			ret = TypeTextListEval
		}
	}
	return
}
