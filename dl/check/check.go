package check

import "github.com/ionous/iffy/dl/composer"

var Slots = []composer.Slot{{
	Name: "testing",
	Type: (*Testing)(nil),
	Desc: "Testing: Run a series of tests.",
}}

var Slats = []composer.Slat{
	(*TestOutput)(nil),
	//
}
