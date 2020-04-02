package rt

type EmptyStream struct{}

func (*EmptyStream) Count() int {
	return 0
}

func (*EmptyStream) HasNext() bool {
	return false
}

func (*EmptyStream) GetNumber() (float64, error) {
	panic("Attempted to get a number from an empty stream.")
}

func (*EmptyStream) GetText() (string, error) {
	panic("Attempted to get a number from an empty stream.")
}
