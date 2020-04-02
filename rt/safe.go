package rt

// Run executes the passed statement using the passed runtime;
// does *not* error if the passed exec is nil.
func Run(run Runtime, exec Execute) (err error) {
	if exec != nil {
		err = exec.Execute(run)
	}
	return
}

// GetBool executes the passed eval using the passed runtime;
// errors if the passed eval is nil.
func GetBool(run Runtime, eval BoolEval) (okay bool, err error) {
	if eval == nil {
		err = MissingEval("empty boolean eval")
	} else {
		okay, err = eval.GetBool(run)
	}
	return
}

// GetNumber executes the passed eval using the passed runtime;
// errors if the passed eval is nil.
func GetNumber(run Runtime, eval NumberEval) (ret float64, err error) {
	if eval == nil {
		err = MissingEval("empty number eval")
	} else {
		ret, err = eval.GetNumber(run)
	}
	return
}

// GetText executes the passed eval using the passed runtime;
// errors if the passed eval is nil.
func GetText(run Runtime, eval TextEval) (ret string, err error) {
	if eval == nil {
		err = MissingEval("empty text eval")
	} else {
		ret, err = eval.GetText(run)
	}
	return
}

// MissingEval error type for unknown variables while processing loops.
type MissingEval string

// Error returns the name of the unknown variable.
func (e MissingEval) Error() string { return string(e) }
