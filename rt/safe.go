package rt

import "io"

// MissingEval error type for unknown variables while processing loops.
type MissingEval string

// Error returns the name of the unknown variable.
func (e MissingEval) Error() string { return string(e) }

// Run executes the passed statement using the passed runtime;
// does *not* error if the passed exec is nil.
func Run(run Runtime, exec Execute) (err error) {
	if exec != nil {
		err = exec.Execute(run)
	}
	return
}

// WriteText evaluates t and outputs the results to w.
func WriteText(run Runtime, eval TextEval) (err error) {
	if t, e := GetText(run, eval); e != nil {
		err = e
	} else {
		_, e := io.WriteString(run, t)
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
func GetOptionalNumber(run Runtime, eval NumberEval, fallback float64) (ret float64, err error) {
	if eval == nil {
		ret = fallback
	} else {
		ret, err = eval.GetNumber(run)
	}
	return
}

func GetTextStream(run Runtime, eval TextListEval) (ret Iterator, err error) {
	if eval == nil {
		ret = &EmptyStream{}
	} else {
		ret, err = eval.GetTextStream(run)
	}
	return
}

func GetNumberStream(run Runtime, eval NumListEval) (ret Iterator, err error) {
	if eval == nil {
		ret = &EmptyStream{}
	} else {
		ret, err = eval.GetNumberStream(run)
	}
	return
}
