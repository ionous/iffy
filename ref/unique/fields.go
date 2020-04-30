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
			if q, ok := t.Find(tag); ok && len(q) == 0 {
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

type PropertyFunc func(f *r.StructField, path []int) (done bool)

func WalkProperties(rtype r.Type, fn PropertyFunc) (done bool) {
	return walkProperties(rtype, nil, fn)
}

func walkProperties(rtype r.Type, base []int, fn PropertyFunc) (done bool) {
	if rtype.Kind() != r.Struct {
		panic(rtype.String() + " not a struct")
	} else {
		for i := 0; i < rtype.NumField(); i++ {
			field := rtype.Field(i)
			if IsPublic(field) {
				path := append(base, i)
				if !IsEmbedded(field) {
					if fn(&field, path) {
						break
					}
				} else if walkProperties(field.Type, path, fn) {
					break
				}
			}
		}
	}
	return
}
