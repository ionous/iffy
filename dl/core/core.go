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
	*SumOf
	*DiffOf
	*ProductOf
	*QuotientOf
	*RemainderOf
	//
	*AllTrue
	*AnyTrue
	*BoolValue
	*Buffer
	*Choose
	*ChooseNum
	*ChooseObj
	*ChooseText
	*ClassName
	*CompareNum
	*CompareObj
	*CompareText
	*Comprise
	*CycleText
	*DoNothing
	// *Error
	*Filter
	*ForEachNum
	*ForEachObj
	*ForEachText
	*Get
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
	*ListUp
	*NumValue
	*Numbers
	*ObjectName
	*ObjectNames
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
	*TextValue
	*Texts
	*TopObject
	// *Using
	// CompareTo:
	*EqualTo
	*GreaterThan
	*LessThan
}
