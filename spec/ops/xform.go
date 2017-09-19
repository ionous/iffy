package ops

import (
	r "reflect"
)

type Transform interface {
	// TransformValue should return value if there was no error, but it couldnt convert.
	TransformValue(v r.Value, hint r.Type) (r.Value, error)
}

// DefaultXform acts as no transform.
type DefaultXform struct{}

// TransformValue here returns v, and never error.
func (DefaultXform) TransformValue(v r.Value, hint r.Type) (r.Value, error) {
	return v, nil
}
