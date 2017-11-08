package ops

import (
	r "reflect"
)

// Transform should return value if there was no error, but it couldnt convert.
type Transform interface {
	TransformValue(v r.Value, hint r.Type) (r.Value, error)
}

// TransformFunction implements the Transform interface for free functions.
type TransformFunction struct {
	Transform func(v r.Value, hint r.Type) (r.Value, error)
}

func (tf TransformFunction) TransformValue(v r.Value, hint r.Type) (r.Value, error) {
	return tf.Transform(v, hint)
}
