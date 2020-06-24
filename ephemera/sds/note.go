package sds

// annotation string
type Note string

func (a Note) HasAnnotation() bool {
	str := string(a)
	return findMarker(str) >= 0
}

func (a Note) Field() string {
	f, _, _ := a.Extract()
	return f
}

func (a Note) Extract() (field, typeName, id string) {
	str := string(a)
	if i := findMarker(str); i < 0 {
		field = str
	} else {
		lhs, rhs := str[:i], str[i+2:]
		if j := findAt(rhs); j < 0 {
			typeName = rhs
		} else {
			typeName, id = rhs[:j], rhs[j+1:]
		}
		if len(lhs) > 0 {
			field = lhs
		} else {
			field = typeName
		}
	}
	return
}

// avoid strings package dependency
func findAt(s string) (ret int) {
	ret = -1
	for i, x := range s {
		if x == '@' {
			ret = i
			break
		}
	}
	return
}

// avoid strings package dependency
func findMarker(s string) (ret int) {
	found := false
	ret = -1
	for i, x := range s {
		if x != ':' {
			found = false
		} else if !found {
			found = true
		} else {
			ret = i - 1
			break
		}
	}
	return
}
