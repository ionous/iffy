package chart

// Runes gathers runes for the parsing of string-like data.
type Runes struct {
	list []rune
}

func (rs Runes) Len() int {
	return len(rs.list)
}

func (rs Runes) String() string {
	return string(rs.list)
}

func (rs *Runes) Accept(r rune, s State) State {
	rs.list = append(rs.list, r)
	return s
}
