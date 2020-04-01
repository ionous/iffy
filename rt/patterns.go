package rt

//
type Pattern interface {
	GetBoolMatching(string) (bool, error)
	GetNumMatching(string) (float64, error)
	GetTextMatching(string) (string, error)
	GetNumStreamMatching(string) (NumberStream, error)
	GetTextStreamMatching(string) (TextStream, error)
	ExecuteMatching(string) error
}
