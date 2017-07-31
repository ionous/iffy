package unique

import (
	r "reflect"
)

func IsPublic(f r.StructField) bool {
	return len(f.PkgPath) == 0
}

func IsEmbedded(f r.StructField) bool {
	return f.Anonymous && f.Type.Kind() == r.Struct
}

// PathOf searches for the passed tag in the passed type, and its embedded fields.
func PathOf(rtype r.Type, tag string) (ret []int, okay bool) {
	return findPathOf(ret, rtype, tag)
}

func findPathOf(base []int, rtype r.Type, tag string) (ret []int, okay bool) {
	for i := 0; i < rtype.NumField(); i++ {
		field := rtype.Field(i)
		if IsPublic(field) {
			t := Tag(field.Tag)
			if _, ok := t.Find(tag); ok {
				ret, okay = append(base, i), true
				break
			} else if IsEmbedded(field) {
				if path, ok := findPathOf(append(base, i), field.Type, tag); ok {
					ret, okay = path, ok
					break
				}
			}
		}
	}
	return
}
