package reader

// Position identifies a location in source file.
// fix: should probably be relocated somewhere more appropriate.
type Position struct {
	Source string // ex. story file
	Offset string // ex. node id
}

// IsValid reports whether the position is valid.
func (p *Position) IsValid() bool {
	return len(p.Offset) > 0
}

// String returns a string in one of several forms:
//
//	file:line:column    valid position with file name
//	file:line           valid position with file name but no column (column == 0)
//	line:column         valid position without file name
//	line                valid position without file name and no column (column == 0)
//	file                invalid position with file name
//	-                   invalid position without file name
//
func (p *Position) String() string {
	s := p.Source
	if p.IsValid() {
		if s != "" {
			s += ":" + p.Offset
		}
	}
	if s == "" {
		s = "-"
	}
	return s
}

func (p *Position) LessThan(q *Position) (ret bool) {
	switch {
	case p.Source != q.Source:
		ret = p.Source < q.Source
	default:
		// fix. hmmm.... ids are not necessarily ordered
		// probably? need a rowid sortkey.
		ret = p.Offset < q.Offset
	}
	return
}
