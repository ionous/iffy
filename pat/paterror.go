package pat

type NotFound string

func (nf NotFound) Error() string {
	return "pattern not found " + string(nf)
}
