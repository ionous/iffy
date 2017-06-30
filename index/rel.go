package index

type Relation struct {
	Index [LineData]Index
}

func MakeRelation(uni, que bool) Relation {
	return Relation{
		[LineData]Index{
			Index{Primary, uni, nil},
			Index{Secondary, que, nil},
		},
	}
}

func (r *Relation) Relate(primary, secondary, data string) (changed bool) {
	if len(primary) == 0 {
		changed = r.remove(Secondary, secondary)
	} else if len(secondary) == 0 {
		changed = r.remove(Primary, primary)
	} else {
		l := MakeLine(primary, secondary, data)
		// relating secondary and primary using data.
		if r.Index[Primary].Add(l) {
			r.Index[Secondary].Add(l)
			changed = true
		}
	}
	return
}

func (r *Relation) remove(major Column, key string) (changed bool) {
	minor := (major + 1) & 1
	if cut := r.Index[major].Remove(key); len(cut) > 0 {
		if missing := r.Index[minor].DeleteKeys(cut); len(missing) > 0 {
			// we expect that all minor keys in one are major keys in the other
			panic(missing)
		}
		changed = true
	}
	return
}
