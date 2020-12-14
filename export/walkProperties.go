package export

import r "reflect"

type PropertyFunc func(f *r.StructField, path []int) (done bool)

func WalkProperties(rtype r.Type, fn PropertyFunc) (done bool) {
	return walkProperties(rtype, nil, fn)
}

func walkProperties(rtype r.Type, base []int, fn PropertyFunc) (done bool) {
	if walking := rtype.String(); rtype.Kind() != r.Struct {
		panic(walking + " not a struct")
	} else {
		for i := 0; i < rtype.NumField(); i++ {
			field := rtype.Field(i)
			if isPublic(field) {
				path := append(base, i)
				if !isEmbedded(field) {
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

func isEmbedded(f r.StructField) bool {
	return f.Anonymous && f.Type.Kind() == r.Struct
}

func isPublic(f r.StructField) bool {
	return len(f.PkgPath) == 0
}
