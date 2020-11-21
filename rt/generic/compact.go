package generic

func CompactNumbers(it Iterator, vals []float64) (ret Value, err error) {
	for it.HasNext() {
		if n, e := it.GetNext(); e != nil {
			err = e
			break
		} else {
			v := n.Float()
			vals = append(vals, v)
		}
	}
	if err == nil {
		ret = FloatsOf(vals)
	}
	return
}

func CompactTexts(it Iterator, vals []string) (ret Value, err error) {
	for it.HasNext() {
		if n, e := it.GetNext(); e != nil {
			err = e
			break
		} else {
			v := n.String()
			vals = append(vals, v)
		}
	}
	if err == nil {
		ret = StringsOf(vals)
	}
	return
}
