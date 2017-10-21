package chart

type runes struct {
	list []rune
}

func (rs runes) String() string {
	return string(rs.list)
}

func (rs *runes) Accept(r rune, s State) State {
	rs.list = append(rs.list, r)
	return s
}
