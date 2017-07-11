package rt

type Patterns interface {
	GetBoolMatching(Object) (bool, error)
	GetNumMatching(Object) (float64, error)
	GetTextMatching(Object) (string, error)
	GetObjectMatching(Object) (Object, error)
	GetNumStreamMatching(Object) (NumberStream, error)
	GetTextStreamMatching(Object) (TextStream, error)
	GetObjStreamMatching(Object) (ObjectStream, error)
	ExecuteMatching(Object) (bool, error)
}
