package rt

//
type Pattern interface {
	GetBoolMatching(string) (bool, error)
	GetNumMatching(string) (float64, error)
	GetTextMatching(string) (string, error)
	GetObjectMatching(string) (string, error)
	GetNumStreamMatching(string) (NumberStream, error)
	GetTextStreamMatching(string) (TextStream, error)
	GetObjStreamMatching(string) (ObjectStream, error)
	ExecuteMatching(string) error
}
