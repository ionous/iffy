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
	// *Buffer
	// *ChangeState
	*Choose
	*ChooseNum
	*ChooseObj
	*ChooseText
	*CompareNum
	*CompareText
	*CycleText
	// *DoNothing
	*EqualTo
	// *Error
	*ForEachNum
	// *ForEachObj
	*ForEachText
	*Get
	// *GoCall
	*GreaterThan
	// *Inc
	// *IsEmpty
	// *IsFromClass
	// *IsNot
	// *IsNum
	// *IsObj
	// *IsState
	// *IsText
	// *IsValid
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
	*PrintText
	*Range
	*ShuffleText
	*StoppingText
	// *SetNumber
	// *SetObject
	// *SetText
	// *State
	// *StopNow
	*Text
	*Texts
	// *Using
}
