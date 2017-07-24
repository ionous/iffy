package core

type Commands struct {
	*Add
	*Sub
	*Mul
	*Div
	//
	*AllTrue
	*AnyTrue
	*Bool
	*Buffer
	*SetState
	*Choose
	*ChooseNum
	*ChooseObj
	*ChooseText
	*ClassName
	*CompareNum
	*CompareText
	*CycleText
	*DoNothing
	*EqualTo
	// *Error
	*ForEachNum
	// *ForEachObj
	*ForEachText
	*Get
	// *GoCall
	*GreaterThan
	// *Inc
	*Includes
	*IsEmpty
	*IsSameClass
	*IsSimilarClass
	*IsNot
	// *IsNum
	// *IsObj
	// *IsState
	// *IsText
	// *IsValid
	*Len
	*LesserThan
	*NotEqualTo
	*Num
	*Numbers
	*Object
	// *Objects
	// *ObjListContains
	// *ObjListIsEmpty
	// *Object
	*PrintLine
	*PrintNum
	*PrintNumWord
	*PrintText
	*Range
	*SetBool
	*SetNum
	*SetText
	*SetObj
	*ShuffleText
	*StoppingText
	// *State
	// *StopNow
	*Text
	*Texts
	// *Using
}
