package generic

var (
	True      = BoolOf(true)
	False     = BoolOf(false)
	Zero      = FloatOf(0.0)
	ZeroList  = FloatsOf(nil)
	Empty     = StringOf("")
	EmptyList = StringsOf(nil)
)

const defaultType = "" // empty string
