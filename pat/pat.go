package pat

type NotFound string

func (nf NotFound) Error() string {
	return "pattern not found " + string(nf)
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
