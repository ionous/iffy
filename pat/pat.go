package pat

type Patterns struct {
	Bools
	Numbers
	Text
	Objects
	NumLists
	TextLists
	ObjLists
	Executes
}

func MakePatterns() Patterns {
	return Patterns{
		make(Bools),
		make(Numbers),
		make(Text),
		make(Objects),
		make(NumLists),
		make(TextLists),
		make(ObjLists),
		make(Executes),
	}
}
