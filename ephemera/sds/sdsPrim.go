package sds

// sdsPrim takes an array of annotation, element
type sdsPrim struct {
	value interface{}
	sdsPanic
}

func (on *sdsPrim) Kind() Kind {
	return Prim
}

func (on *sdsPrim) Value() interface{} {
	return on.value
}
