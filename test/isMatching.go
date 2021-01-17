package test

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/pattern"
	"github.com/ionous/iffy/dl/term"
	"github.com/ionous/iffy/rt"
)

// a pattern for matching groups --
// we add rules that if things arent equal we return false
var matchGroups = pattern.ActivityPattern{
	CommonPattern: pattern.CommonPattern{
		Name: "matchGroups",
		Params: []term.Preparer{
			&term.Record{Name: "a", Kind: "GroupSettings"},
			&term.Record{Name: "b", Kind: "GroupSettings"},
		},
		Returns: &term.Bool{Name: "matches"},
	},
	// rules are evaluated in reverse order ( see splitRules )
	Rules: []*pattern.ExecuteRule{{
		Filter:  &core.Always{},
		Execute: matches(true),
	}, {
		Filter: &core.CompareText{
			&core.Unpack{
				Record: &core.Var{Name: "a"},
				Field:  "Label",
			},
			&core.NotEqualTo{},
			&core.Unpack{
				Record: &core.Var{Name: "b"},
				Field:  "Label",
			},
		},
		Execute: matches(false),
	}, {
		Filter: &core.CompareText{
			&core.Unpack{
				Record: &core.Var{Name: "a"},
				Field:  "Innumerable",
			},
			&core.NotEqualTo{},
			&core.Unpack{
				Record: &core.Var{Name: "b"},
				Field:  "Innumerable",
			},
		},
		Execute: matches(false),
	}, {
		Filter: &core.CompareText{
			&core.Unpack{
				Record: &core.Var{Name: "a"},
				Field:  "GroupOptions",
			},
			&core.NotEqualTo{},
			&core.Unpack{
				Record: &core.Var{Name: "b"},
				Field:  "GroupOptions",
			},
		},
		Execute: matches(false),
	}},
}

func matches(b bool) rt.Execute {
	return &core.Assign{core.Variable{Str: "matches"}, &core.FromBool{&core.Bool{b}}}
}
