package ops

import (
	r "reflect"
)

// Transform should return value if there was no error, but it couldnt convert.
type Transform interface {
	TransformValue(v r.Value, hint r.Type) (r.Value, error)
}

// Transform  implements the TransformValue interface for free functions.
type Transformer func(v r.Value, hint r.Type) (r.Value, error)

func (tf Transformer) TransformValue(v r.Value, hint r.Type) (r.Value, error) {
	return tf(v, hint)
}
