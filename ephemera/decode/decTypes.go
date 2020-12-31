package decode

type SwapType interface {
	Choices() (nameToType map[string]interface{})
}

type StrType interface {
	Choices() (closed bool, choices map[string]string)
}

type NumType interface {
	Choices() (closed bool, vals []float64)
}

func FindChoice(op StrType, choice string) (ret string, found bool) {
	closed, keys := op.Choices()
	if str, ok := keys[choice]; ok {
		ret, found = str, ok
	} else if !ok && !closed {
		ret = choice
	}
	return
}
