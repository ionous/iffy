package qna

import "github.com/ionous/errutil"

// Assign writes v into pv.
func Assign(pv interface{}, v interface{}) (err error) {
	switch out := pv.(type) {
	case *bool:
		switch v := v.(type) {
		case nil:
			*out = false
		case bool:
			*out = v
		case int64: // particularly for sqlite, boolean values can be represented as 1/0
			*out = v != 0
		default:
			err = errutil.Fmt("expected bool, got %T", v)
		}
	case *int:
		switch v := v.(type) {
		case nil:
			*out = 0
		case int:
			*out = int(v)
		case int64:
			*out = int(v)
		case float64:
			*out = int(v)
		default:
			err = errutil.Fmt("expected int, got %T", v)
		}
	case *int64:
		switch v := v.(type) {
		case nil:
			*out = 0.0
		case int:
			*out = int64(v)
		case int64:
			*out = int64(v)
		case float64:
			*out = int64(v)
		default:
			err = errutil.Fmt("expected int64, got %T", v)
		}
	case *float64:
		switch v := v.(type) {
		case nil:
			*out = 0.0
		case int:
			*out = float64(v)
		case int64:
			*out = float64(v)
		case float64:
			*out = float64(v)
		default:
			err = errutil.Fmt("expected float64, got %T", v)
		}
	case *string:
		if v == nil {
			*out = ""
		} else if v, ok := v.(string); !ok {
			err = errutil.Fmt("expected string, got %T", v)
		} else {
			*out = v
		}
	default:
		err = errutil.Fmt("unexpected output type, got %T", pv)
	}
	return
}
