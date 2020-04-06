package assign

// BoolPtr, a shortcut when the expected output type matches the input type.
func BoolPtr(pv interface{}, b bool) (err error) {
	if outptr, ok := pv.(*bool); !ok {
		err = Mismatch("BoolPtr", outptr, pv)
	} else {
		(*outptr) = b
	}
	return
}

// FloatPtr, a shortcut when the expected output type matches the input type.
func FloatPtr(pv interface{}, v float64) (err error) {
	if outptr, ok := pv.(*float64); !ok {
		err = Mismatch("FloatPtr", outptr, pv)
	} else {
		(*outptr) = v
	}
	return
}

// StringPtr, a shortcut when the expected output type matches the input type.
func StringPtr(pv interface{}, str string) (err error) {
	if outptr, ok := pv.(*string); !ok {
		err = Mismatch("StringPtr", outptr, pv)
	} else {
		(*outptr) = str
	}
	return
}

// Value writes any i into any pv, or at least tries to.
func Value(pv interface{}, i interface{}) (err error) {
	switch out := pv.(type) {
	case *bool:
		if v, e := ToBool(i); e != nil {
			err = e
		} else {
			*out = v
		}
	case *int:
		if v, e := ToInt(i); e != nil {
			err = e
		} else {
			*out = v
		}
	case *int64:
		if v, e := ToInt64(i); e != nil {
			err = e
		} else {
			*out = v
		}
	case *float64:
		if v, e := ToFloat(i); e != nil {
			err = e
		} else {
			*out = v
		}
	case *string:
		if v, e := ToString(i); e != nil {
			err = e
		} else {
			*out = v
		}
	default:
		err = Mismatch("Value", pv, i)
	}
	return
}

// ToBool converts nil, int64, and bool values to bool.
func ToBool(i interface{}) (ret bool, err error) {
	switch i := i.(type) {
	case nil:
		ret = false
	case bool:
		ret = i
	case int64: // particularly for sqlite, boolean values can be represented as 1/0
		ret = i != 0
	default:
		err = Mismatch("ToBool", ret, i)
	}
	return
}

// ToInt converts nil and numeric values to int.
func ToInt(i interface{}) (ret int, err error) {
	switch i := i.(type) {
	case nil:
		ret = 0
	case int:
		ret = int(i)
	case int64:
		ret = int(i)
	case float64:
		ret = int(i)
	default:
		err = Mismatch("ToInt", ret, i)
	}
	return
}

// ToInt64 converts nil and numeric values to int.
func ToInt64(i interface{}) (ret int64, err error) {
	switch i := i.(type) {
	case nil:
		ret = 0.0
	case int:
		ret = int64(i)
	case int64:
		ret = int64(i)
	case float64:
		ret = int64(i)
	default:
		err = Mismatch("ToInt64", ret, i)
	}
	return
}

// ToFloat converts nil and numeric values to float.
func ToFloat(i interface{}) (ret float64, err error) {
	switch i := i.(type) {
	case nil:
		ret = 0.0
	case int:
		ret = float64(i)
	case int64:
		ret = float64(i)
	case float64:
		ret = float64(i)
	default:
		err = Mismatch("ToFloat", ret, i)
	}
	return
}

// ToString converts nil and string values to string.
func ToString(i interface{}) (ret string, err error) {
	if i == nil {
		ret = ""
	} else if i, ok := i.(string); !ok {
		err = Mismatch("ToString", ret, i)
	} else {
		ret = i
	}
	return
}
