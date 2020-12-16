package test

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/pattern"
	"github.com/ionous/iffy/dl/term"
)

// a pattern for matching groups --
// we add rules that if things arent equal we return false
var isMatchingGroup = pattern.BoolPattern{
	CommonPattern: pattern.CommonPattern{
		Name: "isMatchingGroup",
		Prologue: []term.Preparer{
			&term.Record{Name: "a", Kind: "GroupSettings"},
			&term.Record{Name: "b", Kind: "GroupSettings"},
		},
	},
	// rules are evaluated in reverse order ( see splitRules )
	Rules: []*pattern.BoolRule{{
		&core.Always{}, &core.Bool{true}}, {
		&core.CompareText{
			&core.Unpack{
				Record: &core.Var{Name: "a"},
				Field:  "Label",
			},
			&core.NotEqualTo{},
			&core.Unpack{
				Record: &core.Var{Name: "b"},
				Field:  "Label",
			},
		}, &core.Bool{Bool: false}}, {
		&core.CompareText{
			&core.Unpack{
				Record: &core.Var{Name: "a"},
				Field:  "Innumerable",
			},
			&core.NotEqualTo{},
			&core.Unpack{
				Record: &core.Var{Name: "b"},
				Field:  "Innumerable",
			},
		}, &core.Bool{Bool: false}}, {
		&core.CompareText{
			&core.Unpack{
				Record: &core.Var{Name: "a"},
				Field:  "GroupOptions",
			},
			&core.NotEqualTo{},
			&core.Unpack{
				Record: &core.Var{Name: "b"},
				Field:  "GroupOptions",
			},
		}, &core.Bool{Bool: false}},
	},
}
