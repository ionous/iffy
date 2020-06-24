package sds

// sdsPanic
type sdsPanic struct{}

func (on *sdsPanic) Kind() Kind {
	panic("node has unknown kind")
}
func (on *sdsPanic) Note() Note {
	panic("node has unknown annotation")
}
func (on *sdsPanic) Elem() Object {
	panic("node has unknown element data")
}
func (on *sdsPanic) Value() interface{} {
	panic("node has unknown value")
}
func (on *sdsPanic) Next() Object {
	panic("node has no next node")
}
func (on *sdsPanic) Params() []string {
	panic("node not parameterized")
}
func (on *sdsPanic) Param(string) (Note, Object) {
	panic("node has no params")
}
