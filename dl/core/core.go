package core

import (
	"github.com/ionous/iffy/dl/composer"
)

var Slots = []composer.Slot{{
	Type: (*Comparator)(nil),
	Desc: "Comparison Types: Helper for comparing values.",
}, {
	Type: (*Assignment)(nil),
	Desc: "Assignments: Helper for setting variables.",
}, {
	Type: (*Fields)(nil),
	Desc: "Helper for setting fields.",
}, {
	Type: (*Brancher)(nil),
	Desc: "Helper for choose action.",
}}

var Slats = []composer.Composer{
	(*Activity)(nil),

	// some boolean tests:
	(*Always)(nil),
	(*AllTrue)(nil),
	(*AnyTrue)(nil),
	(*HasDominion)(nil),
	(*IsNotTrue)(nil), // inverts a bool eval

	// Assign turns an Assignment a normal statement.
	(*Assign)(nil),
	(*Variable)(nil),
	(*FromBool)(nil),
	(*FromNum)(nil),
	(*FromText)(nil),
	(*FromName)(nil),
	(*FromRecord)(nil),
	(*FromObject)(nil),
	(*FromNumbers)(nil),
	(*FromTexts)(nil),
	(*FromRecords)(nil),
	(*CopyFrom)(nil),
	(*Make)(nil),

	// FIX: Choose scalar/any?
	(*ChooseNum)(nil),
	(*ChooseText)(nil),
	// FIX: compare scalar?
	(*CompareNum)(nil),
	(*CompareText)(nil),

	// literals
	(*Bool)(nil),
	(*Number)(nil),
	(*Text)(nil),
	(*Lines)(nil),
	(*Numbers)(nil),
	(*Texts)(nil),

	(*SumOf)(nil),
	(*DiffOf)(nil),
	(*ProductOf)(nil),
	(*QuotientOf)(nil),
	(*RemainderOf)(nil),

	(*SimpleNoun)(nil),
	(*ObjectName)(nil),
	(*ObjectExists)(nil),
	(*NameOf)(nil),
	(*KindOf)(nil),
	(*IsKindOf)(nil),

	(*PrintNum)(nil),
	(*PrintNumWord)(nil),

	// FIX: should take a speaker, and we should have a default speaker
	(*Say)(nil),
	(*Buffer)(nil),
	(*Span)(nil),
	(*Bracket)(nil),
	(*Slash)(nil),
	(*Commas)(nil),
	// text transforms
	(*MakeSingular)(nil),
	(*MakePlural)(nil),
	(*MakeUppercase)(nil),
	(*MakeLowercase)(nil),
	(*MakeTitleCase)(nil),
	(*MakeSentenceCase)(nil),
	(*MakeReversed)(nil),
	//
	(*Matches)(nil),
	// sequences
	(*While)(nil),
	(*CycleText)(nil),
	(*ShuffleText)(nil),
	(*StoppingText)(nil),

	(*Field)(nil),
	(*Unpack)(nil),
	(*Var)(nil),
	(*HasTrait)(nil),

	(*IsEmpty)(nil),
	(*Includes)(nil),
	(*Join)(nil),

	// comparison
	(*EqualTo)(nil),
	(*NotEqualTo)(nil),
	(*GreaterThan)(nil),
	(*LessThan)(nil),
	(*GreaterOrEqual)(nil),
	(*LessOrEqual)(nil),

	(*Arguments)(nil),
	(*Argument)(nil),
	// put at field
	(*PutAtField)(nil),
	(*IntoRec)(nil),
	(*IntoObj)(nil),
	(*IntoObjNamed)(nil),
	// choose action (if)
	(*ChooseAction)(nil),
	(*ChooseMore)(nil),
	(*ChooseNothingElse)(nil),
}
