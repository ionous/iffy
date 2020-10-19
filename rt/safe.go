package rt

import "github.com/ionous/errutil"

// MissingEval error type for unknown variables while processing loops.
type MissingEval string

// Error returns the name of the unknown variable.
func (e MissingEval) Error() string { return string(e) }

// Run executes the passed statement using the passed runtime;
// does *not* error if the passed exec is nil.
func RunAll(run Runtime, exes []Execute) (err error) {
	for _, exe := range exes {
		if exe != nil {
			if e := exe.Execute(run); e != nil {
				err = e
				break
			}
		}
	}
	return
}

// Run executes the passed statement using the passed runtime;
// does *not* error if the passed exec is nil.
func RunOne(run Runtime, exe Execute) (err error) {
	if exe == nil {
		err = MissingEval("empty execute")
	} else {
		err = exe.Execute(run)
	}
	return
}

// WriteText evaluates t and outputs the results to w.
func WriteText(run Runtime, eval TextEval) (err error) {
	if t, e := GetText(run, eval); e != nil {
		err = e
	} else if w := run.Writer(); w == nil {
		err = errutil.New("missing writer")
	} else {
		_, e := w.WriteString(t)
		err = e
	}
	return
}

// GetBool runs the specified eval, returning an error if the eval is nil.
func GetBool(run Runtime, eval BoolEval) (okay bool, err error) {
	if eval == nil {
		err = MissingEval("empty boolean eval")
	} else {
		okay, err = eval.GetBool(run)
	}
	return
}

// GetNumber runs the specified eval, returning an error if the eval is nil.
func GetNumber(run Runtime, eval NumberEval) (ret float64, err error) {
	if eval == nil {
		err = MissingEval("empty number eval")
	} else {
		ret, err = eval.GetNumber(run)
	}
	return
}

// GetText runs the specified eval, returning an error if the eval is nil.
func GetText(run Runtime, eval TextEval) (ret string, err error) {
	if eval == nil {
		err = MissingEval("empty text eval")
	} else {
		ret, err = eval.GetText(run)
	}
	return
}

// GetOptionalBool runs the optionally specified eval.
func GetOptionalBool(run Runtime, eval BoolEval, fallback bool) (ret bool, err error) {
	if eval == nil {
		ret = fallback
	} else {
		ret, err = eval.GetBool(run)
	}
	return
}

// GetOptionalNumber runs the optionally specified eval.
func GetOptionalNumber(run Runtime, eval NumberEval, fallback float64) (ret float64, err error) {
	if eval == nil {
		ret = fallback
	} else {
		ret, err = eval.GetNumber(run)
	}
	return
}

// GetOptionalText runs the optionally specified eval.
func GetOptionalText(run Runtime, eval TextEval, fallback string) (ret string, err error) {
	if eval == nil {
		ret = fallback
	} else {
		ret, err = eval.GetText(run)
	}
	return
}

// GetOptionalNumber runs the optionally specified eval.
func GetOptionalNumbers(run Runtime, eval NumListEval, fallback []float64) (ret []float64, err error) {
	if eval == nil {
		ret = fallback
	} else {
		ret, err = GetNumList(run, eval)
	}
	return
}

// GetOptionalText runs the optionally specified eval.
func GetOptionalTexts(run Runtime, eval TextListEval, fallback []string) (ret []string, err error) {
	if eval == nil {
		ret = fallback
	} else {
		ret, err = GetTextList(run, eval)
	}
	return
}

// GetNumList returns an new iterator to walk the passed list,
// or an empty iterator if the value is null.
func GetNumList(run Runtime, eval NumListEval) (ret []float64, err error) {
	if eval != nil {
		ret, err = eval.GetNumList(run)
	}
	return
}

// GetTextList returns an new iterator to walk the passed list,
// or an empty iterator if the value is null.
func GetTextList(run Runtime, eval TextListEval) (ret []string, err error) {
	if eval != nil {
		ret, err = eval.GetTextList(run)
	}
	return
}

func CompactNumbers(it Iterator, vals []float64) (ret []float64, err error) {
	for it.HasNext() {
		if n, e := it.GetNext(); e != nil {
			err = e
			break
		} else if v, e := n.GetNumber(); e != nil {
			err = e
			break
		} else {
			vals = append(vals, v)
		}
	}
	if err == nil {
		ret = vals
	}
	return
}

func CompactTexts(it Iterator, vals []string) (ret []string, err error) {
	for it.HasNext() {
		if n, e := it.GetNext(); e != nil {
			err = e
			break
		} else if v, e := n.GetText(); e != nil {
			err = e
			break
		} else {
			vals = append(vals, v)
		}
	}
	if err == nil {
		ret = vals
	}
	return
}
