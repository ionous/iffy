package sds

// an object containing parameters
type sdsElem struct {
	contents map[string]interface{}
	sdsPanic
}

func (on *sdsElem) Kind() Kind {
	return Elem
}

func (on *sdsElem) Params() (ret []string) {
	if cnt := len(on.contents); cnt > 0 {
		out := make([]string, 0, cnt)
		for k, _ := range on.contents {
			n := Note(k).Field()
			out = append(out, n)
		}
		ret = out
	}
	return
}

func (on *sdsElem) Param(p string) (retn Note, retd Object) {
	found := false
	var val interface{}
	for k, v := range on.contents {
		n := Note(k)
		if p == n.Field() {
			retn, val, found = n, v, true
			break
		}
	}
	if found {
		switch v := val.(type) {
		case map[string]interface{}:
			retd = &sdsElem{contents: v}

		case []interface{}:
			if on.Note().HasAnnotation() {
				retd = &sdsSlice{contents: v}
			} else {
				retd = &sdsPrim{value: v}
			}

		default:
			retd = &sdsPrim{value: v}
		}
	}
	return
}
