package test

import (
	"testing"

	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/pattern"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/test/testutil"
	test "github.com/ionous/iffy/test/testutil"
)

func TestMatching(t *testing.T) {
	var kinds testutil.Kinds
	type Things struct{}
	kinds.AddKinds((*Things)(nil), (*GroupSettings)(nil))
	k := kinds.Kind("GroupSettings")

	//
	lt := testTime{Kinds: &kinds,
		PatternMap: testutil.PatternMap{
			"matchGroups": &matchGroups,
		},
	}

	a, b := k.NewRecord(), k.NewRecord()
	runMatching := &pattern.Determine{
		Pattern: "matchGroups", Arguments: core.Args(
			&core.FromValue{g.RecordOf(a)},
			&core.FromValue{g.RecordOf(b)},
		)}
	// default should match
	{
		if ok, e := runMatching.GetBool(&lt); e != nil {
			t.Fatal(e)
		} else if ok.Bool() != true {
			t.Fatal(e)
		}
	}
	// different labels shouldnt match
	{
		if e := test.SetRecord(a, "Label", "beep"); e != nil {
			t.Fatal(e)
		} else if ok, e := runMatching.GetBool(&lt); e != nil {
			t.Fatal(e)
		} else if ok.Bool() != false {
			t.Fatal(e)
		}
	}
	// same labels should match
	{
		if e := test.SetRecord(b, "Label", "beep"); e != nil {
			t.Fatal(e)
		} else if ok, e := runMatching.GetBool(&lt); e != nil {
			t.Fatal(e)
		} else if ok.Bool() != true {
			t.Fatal(e)
		}
	}
	// many fields should match
	{
		if e := test.SetRecord(a, "Innumerable", "IsInnumerable"); e != nil {
			t.Fatal(e)
		} else if e := test.SetRecord(b, "IsInnumerable", true); e != nil {
			t.Fatal(e)
		} else if e := test.SetRecord(a, "GroupOptions", "ObjectsWithArticles"); e != nil {
			t.Fatal(e)
		} else if e := test.SetRecord(b, "GroupOptions", "ObjectsWithArticles"); e != nil {
			t.Fatal(e)
		} else if ok, e := runMatching.GetBool(&lt); e != nil {
			t.Fatal(e)
		} else if ok.Bool() != true {
			t.Fatal(e)
		}
	}
	// names shouldnt be involved
	{
		if e := test.SetRecord(a, "Name", "hola"); e != nil {
			t.Fatal(e)
		} else if ok, e := runMatching.GetBool(&lt); e != nil {
			t.Fatal(e)
		} else if ok.Bool() != true {
			t.Fatal(e)
		}
	}
}
