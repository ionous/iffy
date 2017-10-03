package core

type Classes struct {
	*NumberCounter
	*TextCounter
	*ObjCounter
	*CycleCounter
	*ShuffleCounter
	*StoppingCounter
}

type Commands struct {
	*Add
	*Sub
	*Mul
	*Div
	*Mod
	//
	*AllTrue
	*AnyTrue
	*Bool
	*Buffer
	*Choose
	*ChooseNum
	*ChooseObj
	*ChooseText
	*ClassName
	*CompareNum
	*CompareText
	*Comprise
	*CycleText
	*DoNothing
	*EqualTo
	// *Error
	*Filter
	*ForEachNum
	*ForEachObj
	*ForEachText
	*Get
	*GreaterThan
	// *Inc
	*Includes
	// *Is/State -> use Get
	*IsEmpty
	*IsExactClass
	*IsClass
	*IsNot
	// *IsNum
	// *IsObj
	// *IsText
	// *IsValid
	*Join
	*Len
	*LesserThan
	*ListUp
	*NotEqualTo
	*Num
	*Numbers
	*Object
	*Objects
	// *ObjListContains
	// *ObjListIsEmpty
	*PrintBracket
	*PrintList
	*PrintNum
	*PrintNumWord
	*PrintSlash
	*PrintSpan
	*Say
	*Range
	*RelatedList
	*RelationEmpty
	*Reverse
	*SetBool
	*SetNum
	*SetText
	*SetObj
	*SetState
	*ShuffleText
	*StoppingText
	// *State
	// *StopNow
	*Text
	*Texts
	*TopObject
	// *Using
}
