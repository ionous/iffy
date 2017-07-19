package pat

type NotFound string

func (nf NotFound) Error() string {
	return "pattern not found " + string(nf)
}

// Found returns true if error was nil, false but no error if NotFound, error otherwise.
func Found(e error) (okay bool, err error) {
	if e == nil {
		okay = true
	} else if _, ok := e.(NotFound); !ok {
		err = e
	}
	return
}

type Patterns struct {
	BoolMap
	NumberMap
	TextMap
	ObjectMap
	NumListMap
	TextListMap
	ObjListMap
	ExecuteMap
}
