package index

type Column int

const (
	Primary Column = iota
	Secondary
	Columns
)

type KeyData struct {
	Key  [Columns]string
	Data interface{}
}

func MakeKey(a, b string) *KeyData {
	return &KeyData{[Columns]string{a, b}, nil}
}
