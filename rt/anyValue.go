package rt

type Value interface {
	BoolEval
	NumberEval
	TextEval
	NumListEval
	TextListEval
}
