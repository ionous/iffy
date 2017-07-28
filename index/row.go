package index

type Column int

const (
	Primary Column = iota
	Secondary
	Columns
)

type Row [Columns]string
