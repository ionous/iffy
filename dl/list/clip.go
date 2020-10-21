package list

// b/c copy() doesnt allocate
func copystrings(src []string) []string {
	out := make([]string, len(src))
	copy(out, src)
	return out
}
func copyfloats(src []float64) []float64 {
	out := make([]float64, len(src))
	copy(out, src)
	return out
}

func clipStart(i, cnt int) (ret int) {
	if i == 0 {
		ret = 0 // unspecified: start at the front of the list
	} else if i > cnt {
		ret = -1 // negative return means an empty list
	} else if i > 0 {
		ret = i - 1 // one based indicies
	} else if ofs := cnt + i; ofs > 0 {
		// offset from the end: slice(-2) extracts the last two elements in the sequence.
		ret = ofs
	} else {
		ret = 0
	}
	return
}

func clipEnd(j, cnt int) (ret int) {
	if j > cnt {
		ret = cnt
	} else if j > 0 {
		ret = j - 1
	} else if ofs := cnt + j; ofs > 0 {
		ret = ofs
	} else {
		ret = -1 // negative return means an empty list
	}
	return
}

// turn a starting index and a number of elements from that index into an ending index
func clipRange(start, rng, cnt int) (ret int) {
	if rng <= 0 {
		ret = start // cut start to start, ie. nothing
	} else if end := start + rng; end < cnt {
		ret = end // up to but not including the end
	} else {
		ret = cnt // every element
	}
	return
}
