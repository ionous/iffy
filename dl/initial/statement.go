package initial

type Statement interface {
	Assert(*Facts) error
}

type Statements []Statement

func (l Statements) Assess(f *Facts) (err error) {
	for _, s := range l {
		if e := s.Assert(f); e != nil {
			err = e
			break
		}
	}
	return
}
