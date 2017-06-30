package index

type Column int

const (
	Primary Column = iota
	Secondary
	LineData
	Columns
)

type Line [Columns]string

func MakeLine(a, b, c string) *Line {
	return &Line{a, b, c}
}

func (l *Line) Primary() string {
	return l[Primary]
}
func (l *Line) Secondary() string {
	return l[Secondary]
}
func (l *Line) LineData() string {
	return l[LineData]
}
func (l *Line) match(o *Line, i Column) bool {
	return l[i] == o[i]
}
