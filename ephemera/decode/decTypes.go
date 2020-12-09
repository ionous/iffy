package decode

type SwapType interface {
	Choices() (nameToType map[string]interface{})
}

type StrType interface {
	String() string
	Choices() (closed bool, vals []string)
}

type NumType interface {
	Choices() (closed bool, vals []float64)
}

func IndexOfChoice(op StrType, choice string) (ret int) {
	ret = -1 // if not found
	_, keys := op.Choices()
	for i, k := range keys {
		if choice == k {
			ret = i
			break
		}
	}
	return
}
