package kindOf

import (
	r "reflect"
)

// BoolLike if rtype is a bool, or an eval which produces a bool.
func BoolLike(rtype r.Type) bool {
	return rtype.Kind() == r.Bool || rtype == TypeBoolEval
}

// StringLike if rtype is a bool, or an eval which produces a bool.
func StringLike(rtype r.Type) bool {
	return rtype.Kind() == r.String || rtype == TypeTextEval
}

// NumberLike if rtype is a bool, or an eval which produces a bool.
func NumberLike(rtype r.Type) bool {
	return Number(rtype) || rtype == TypeNumEval
}
