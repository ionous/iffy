package pat

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	"sort"
)

// Contract builds a set of rules.
// FIX? consider rewriting using reflect so that all patterns are stored together and if the interface doesnt match some expectation it errors.
type Contract struct {
	Types unique.Types // pattern types
	Rulebook
}

func MakeContract(patternTypes unique.Types) Contract {
	return Contract{patternTypes, MakeRulebook()}
}

func (rs Contract) AddBoolRule(id ident.Id, f Filters, k rt.BoolEval) (err error) {
	if rtype, ok := rs.Types[id]; !ok {
		err = errutil.New("no such pattern", id)
	} else {
		rs.Bools[rtype] = append(rs.Bools[rtype], BoolRule{f, k})
	}
	return
}
func (rs Contract) AddNumberRule(id ident.Id, f Filters, k rt.NumberEval) (err error) {
	if rtype, ok := rs.Types[id]; !ok {
		err = errutil.New("no such pattern", id)
	} else {
		rs.Numbers[rtype] = append(rs.Numbers[rtype], NumberRule{f, k})
	}
	return
}
func (rs Contract) AddTextRule(id ident.Id, f Filters, k rt.TextEval) (err error) {
	if rtype, ok := rs.Types[id]; !ok {
		err = errutil.New("no such pattern", id)
	} else {
		rs.Text[rtype] = append(rs.Text[rtype], TextRule{f, k})
	}
	return
}
func (rs Contract) AddObjectRule(id ident.Id, f Filters, k rt.ObjectEval) (err error) {
	if rtype, ok := rs.Types[id]; !ok {
		err = errutil.New("no such pattern", id)
	} else {
		rs.Objects[rtype] = append(rs.Objects[rtype], ObjectRule{f, k})
	}
	return
}
func (rs Contract) AddNumListRule(id ident.Id, f Filters, k rt.NumListEval) (err error) {
	if rtype, ok := rs.Types[id]; !ok {
		err = errutil.New("no such pattern", id)
	} else {
		rs.NumLists[rtype] = append(rs.NumLists[rtype], NumListRule{f, k})
	}
	return
}
func (rs Contract) AddTextListRule(id ident.Id, f Filters, k rt.TextListEval) (err error) {
	if rtype, ok := rs.Types[id]; !ok {
		err = errutil.New("no such pattern", id)
	} else {
		rs.TextLists[rtype] = append(rs.TextLists[rtype], TextListRule{f, k})
	}
	return
}
func (rs Contract) AddObjListRule(id ident.Id, f Filters, k rt.ObjListEval) (err error) {
	if rtype, ok := rs.Types[id]; !ok {
		err = errutil.New("no such pattern", id)
	} else {
		rs.ObjLists[rtype] = append(rs.ObjLists[rtype], ObjListRule{f, k})
	}
	return
}
func (rs Contract) AddExecuteRule(id ident.Id, f Filters, k rt.Execute, flags Flags) (err error) {
	if rtype, ok := rs.Types[id]; !ok {
		err = errutil.New("no such pattern", id)
	} else {
		rs.Executes[rtype] = append(rs.Executes[rtype], ExecuteRule{f, k, flags})
	}
	return
}

// Sort in-place so that lengthier filters are at the front of each list.
func (rs Contract) Sort() {
	for _, l := range rs.Bools {
		sort.SliceStable(l, func(i, j int) bool {
			return len(l[i].Filters) < len(l[j].Filters)
		})
	}
	for _, l := range rs.Numbers {
		sort.SliceStable(l, func(i, j int) bool {
			return len(l[i].Filters) < len(l[j].Filters)
		})
	}
	for _, l := range rs.Text {
		sort.SliceStable(l, func(i, j int) bool {
			return len(l[i].Filters) < len(l[j].Filters)
		})
	}
	for _, l := range rs.Objects {
		sort.SliceStable(l, func(i, j int) bool {
			return len(l[i].Filters) < len(l[j].Filters)
		})
	}
	for _, l := range rs.NumLists {
		sort.SliceStable(l, func(i, j int) bool {
			return len(l[i].Filters) < len(l[j].Filters)
		})
	}
	for _, l := range rs.TextLists {
		sort.SliceStable(l, func(i, j int) bool {
			return len(l[i].Filters) < len(l[j].Filters)
		})
	}
	for _, l := range rs.ObjLists {
		sort.SliceStable(l, func(i, j int) bool {
			return len(l[i].Filters) < len(l[j].Filters)
		})
	}
	for _, l := range rs.Executes {
		sort.SliceStable(l, func(i, j int) bool {
			return len(l[i].Filters) < len(l[j].Filters)
		})
	}
}
