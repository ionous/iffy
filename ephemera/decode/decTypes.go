package decode

type SwapType interface {
	Choices() (nameToType map[string]interface{})
}

type StrType interface {
	String() string
	Choices() (closed bool, choices map[string]string)
}

type NumType interface {
	Choices() (closed bool, vals []float64)
}

func FindChoice(op StrType, choice string) (ret string, okay bool) {
	closed, keys := op.Choices()
	str, ok := keys[choice]
	return str, ok || !closed
}
