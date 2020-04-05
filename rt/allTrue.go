package rt

// GetAllTrue returns false when one of the passed evals returns false.
// It does not necessarily run all of the tests, it exits as soon as any test return false.
// An empty list returns true.
func GetAllTrue(run Runtime, eval []BoolEval) (okay bool, err error) {
	i, cnt := 0, len(eval)
	for ; i < cnt; i++ {
		if ok, e := GetBool(run, eval[i]); e != nil {
			err = e
			break
		} else if !ok {
			break
		}
	}
	if i == cnt {
		okay = true
	}
	return
}
