package rt

//
type Pattern interface {
	GetBoolMatching(string) (bool, error)
	GetNumMatching(string) (float64, error)
	GetTextMatching(string) (string, error)
	GetNumStreamMatching(string) (Iterator, error)
	GetTextStreamMatching(string) (Iterator, error)
	ExecuteMatching(string) error
}
