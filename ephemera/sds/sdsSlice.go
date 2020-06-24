package sds

func NewSlice(in []interface{}) Object {
	return &sdsSlice{contents: in}
}

// sdsSlice takes an array of annotation, element
type sdsSlice struct {
	contents []interface{}
	sdsPanic
}

func (on *sdsSlice) Kind() Kind {
	return Slice
}

func (on *sdsSlice) Note() (ret Note) {
	if len(on.contents) > 0 {
		if str, ok := on.contents[0].(string); ok {
			if x := Note(str); x.HasAnnotation() {
				ret = x
			}
		}
	}
	return
}

// elements are often in pairs of notation and data
func (on *sdsSlice) contentIndex() (ret int) {
	if len(on.Note()) > 0 {
		ret = 1
	}
	return
}

func (on *sdsSlice) Elem() (ret Object) {
	if i := on.contentIndex(); i < len(on.contents) {
		el := on.contents[i]
		if m, ok := el.(map[string]interface{}); ok {
			ret = &sdsElem{contents: m}
		}
	}
	return
}

func (on *sdsSlice) Next() (ret Object) {
	next := on.contentIndex() + 1
	if next < len(on.contents) {
		ret = &sdsSlice{contents: on.contents[next:]}
	}
	return
}
