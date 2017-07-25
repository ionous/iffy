package rt

type Patterns interface {
	GetBoolMatching(Runtime, Object) (bool, error)
	GetNumMatching(Runtime, Object) (float64, error)
	GetTextMatching(Runtime, Object) (string, error)
	GetObjectMatching(Runtime, Object) (Object, error)
	GetNumStreamMatching(Runtime, Object) (NumberStream, error)
	GetTextStreamMatching(Runtime, Object) (TextStream, error)
	GetObjStreamMatching(Runtime, Object) (ObjectStream, error)
	ExecuteMatching(Runtime, Object) error
}
