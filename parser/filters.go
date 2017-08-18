package parser

type FilterSpec struct {
	*HasAttr
	*HasClass
}

type Filter interface {
	MatchesNoun(NounVisitor) bool
}

type Filters []Filter

type HasAttr struct {
	Name string
}

type HasClass struct {
	Name string
}

func (f *HasAttr) MatchesNoun(n NounVisitor) bool {
	return n.HasAttribute(f.Name)
}

func (f *HasClass) MatchesNoun(n NounVisitor) bool {
	return n.HasClass(f.Name)
}

func (fs Filters) MatchesNoun(n NounVisitor) bool {
	i, cnt := 0, len(fs)
	for ; i < cnt; i++ {
		if f := fs[i]; !f.MatchesNoun(n) {
			break
		}
	}
	return i == cnt
}
