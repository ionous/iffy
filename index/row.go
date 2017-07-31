package index

type Row struct {
	Major, Minor string
}

func (r Row) String() string {
	return r.Major + "," + r.Minor
}
