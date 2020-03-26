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
		default:
			err = errutil.New("expected bool")
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
			err = errutil.New("expected int")
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
			err = errutil.New("expected int64")
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
			err = errutil.New("expected float64")
		}
	case *string:
		if v == nil {
			*out = ""
		} else if v, ok := v.(string); !ok {
			err = errutil.New("expected string")
		} else {
			*out = v
		}
	default:
		err = errutil.New("unexpected output type")
	}
	return
}
