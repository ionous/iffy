package rt

func ExecuteList(run Runtime, x []Execute) (err error) {
	for _, s := range x {
		if e := s.Execute(run); e != nil {
			err = e
			break
		}
	}
	return
}

// Values for SetValues.
type Values map[string]interface{}

// SetValues to the passed object.
// FIX? add an optional map parameter to NewObject?
func SetValues(obj Object, values Values) (err error) {
	for name, v := range values {
		if e := obj.SetValue(name, v); e != nil {
			err = e
			break
		}
	}
	return
}

// Unpack extacts a go value from any eval type.
func Unpack(run Runtime, v interface{}) (ret interface{}, err error) {
	switch eval := v.(type) {
	case BoolEval:
		ret, err = eval.GetBool(run)
	case NumberEval:
		ret, err = eval.GetNumber(run)
	case TextEval:
		ret, err = eval.GetText(run)
	case ObjectEval:
		ret, err = eval.GetObject(run)
	case NumListEval:
		var vals []float64
		if stream, e := eval.GetNumberStream(run); e != nil {
			err = e
		} else {
			for stream.HasNext() {
				if v, e := stream.GetNext(); e != nil {
					err = e
					break
				} else {
					vals = append(vals, v)
				}
			}
			if err == nil {
				ret = vals
			}
		}
	case TextListEval:
		var vals []string
		if stream, e := eval.GetTextStream(run); e != nil {
			err = e
		} else {
			for stream.HasNext() {
				if v, e := stream.GetNext(); e != nil {
					err = e
					break
				} else {
					vals = append(vals, v)
				}
			}
			if err == nil {
				ret = vals
			}
		}
	case ObjListEval:
		var vals []Object
		if stream, e := eval.GetObjectStream(run); e != nil {
			err = e
		} else {
			for stream.HasNext() {
				if v, e := stream.GetNext(); e != nil {
					err = e
					break
				} else {
					vals = append(vals, v)
				}
			}
			if err == nil {
				ret = vals
			}
		}
	default:
		ret = v
	}
	return
}
