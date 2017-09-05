package ops

import (
	r "reflect"
)

type Transform interface {
	// TransformValue should return value if there was no error, but it couldnt convert.
	TransformValue(v interface{}, hint r.Type) (interface{}, error)
}

// DefaultXform acts as no transform.
type DefaultXform struct{}

// TransformValue here returns v, and never error.
func (DefaultXform) TransformValue(v interface{}, hint r.Type) (interface{}, error) {
	return v, nil
}
