package pat

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	r "reflect"
	"sort"
)

// Contract builds a set of rules.
// FIX? consider rewriting using reflect so that all patterns are stored together and if the interface doesnt match some expectation it errors.
type Contract struct {
	Types unique.Types // pattern types
	rules Rulebook
}

func MakeContract(patternTypes unique.Types) Contract {
	return Contract{patternTypes, MakeRulebook()}
}

func (c Contract) AddBoolRule(id ident.Id, f Filters, k rt.BoolEval) (err error) {
	if rtype, ok := c.Types[id]; !ok {
		err = errutil.New("no such pattern", id)
	} else {
		c.rules.Bools[rtype] = append(c.rules.Bools[rtype], BoolRule{f, k})
	}
	return
}
func (c Contract) AddNumberRule(id ident.Id, f Filters, k rt.NumberEval) (err error) {
	if rtype, ok := c.Types[id]; !ok {
		err = errutil.New("no such pattern", id)
	} else {
		c.rules.Numbers[rtype] = append(c.rules.Numbers[rtype], NumberRule{f, k})
	}
	return
}
func (c Contract) AddTextRule(id ident.Id, f Filters, k rt.TextEval) (err error) {
	if rtype, ok := c.Types[id]; !ok {
		err = errutil.New("no such pattern", id)
	} else {
		c.rules.TextPatterns[rtype] = append(c.rules.TextPatterns[rtype], TextRule{f, k})
	}
	return
}
func (c Contract) AddObjectRule(id ident.Id, f Filters, k rt.ObjectEval) (err error) {
	if rtype, ok := c.Types[id]; !ok {
		err = errutil.New("no such pattern", id)
	} else {
		c.rules.Objects[rtype] = append(c.rules.Objects[rtype], ObjectRule{f, k})
	}
	return
}
func (c Contract) AddNumListRule(id ident.Id, f Filters, k rt.NumListEval, flags Flags) (err error) {
	if rtype, ok := c.Types[id]; !ok {
		err = errutil.New("no such pattern", id)
	} else {
		f := ListRule{f, flags}
		c.rules.NumLists[rtype] = append(c.rules.NumLists[rtype], NumListRule{f, k})
	}
	return
}
func (c Contract) AddTextListRule(id ident.Id, f Filters, k rt.TextListEval, flags Flags) (err error) {
	if rtype, ok := c.Types[id]; !ok {
		err = errutil.New("no such pattern", id)
	} else {
		f := ListRule{f, flags}
		c.rules.TextLists[rtype] = append(c.rules.TextLists[rtype], TextListRule{f, k})
	}
	return
}
func (c Contract) AddObjListRule(id ident.Id, f Filters, k rt.ObjListEval, flags Flags) (err error) {
	if rtype, ok := c.Types[id]; !ok {
		err = errutil.New("no such pattern", id)
	} else {
		f := ListRule{f, flags}
		c.rules.ObjLists[rtype] = append(c.rules.ObjLists[rtype], ObjListRule{f, k})
	}
	return
}
func (c Contract) AddExecuteRule(id ident.Id, f Filters, k rt.Execute, flags Flags) (err error) {
	if rtype, ok := c.Types[id]; !ok {
		err = errutil.New("no such pattern", id)
	} else {
		f := ListRule{f, flags}
		c.rules.Executes[rtype] = append(c.rules.Executes[rtype], ExecuteRule{f, k})
	}
	return
}

// Rulebook returns a copy of the rules, sorted based on filter length.
// Fewer filters are earlier in the list.
// NOTE: patterns evaluate in reverse order: prefering long filters, declared later.
func (c Contract) Rulebook() Rulebook {
	b := MakeRulebook()
	for k, l := range c.rules.Bools {
		l := copyRules(l).(BoolRules)
		sort.SliceStable(l, func(i, j int) bool {
			return len(l[i].Filters) < len(l[j].Filters)
		})
		b.Bools[k] = l
	}
	for k, l := range c.rules.Numbers {
		l := copyRules(l).(NumberRules)
		sort.SliceStable(l, func(i, j int) bool {
			return len(l[i].Filters) < len(l[j].Filters)
		})
		b.Numbers[k] = l
	}
	for k, l := range c.rules.TextPatterns {
		l := copyRules(l).(TextRules)
		sort.SliceStable(l, func(i, j int) bool {
			return len(l[i].Filters) < len(l[j].Filters)
		})
		b.TextPatterns[k] = l
	}
	for k, l := range c.rules.Objects {
		l := copyRules(l).(ObjectRules)
		sort.SliceStable(l, func(i, j int) bool {
			return len(l[i].Filters) < len(l[j].Filters)
		})
		b.Objects[k] = l
	}
	for k, l := range c.rules.NumLists {
		l := copyRules(l).(NumListRules)
		sort.SliceStable(l, func(i, j int) bool {
			return len(l[i].Filters) < len(l[j].Filters)
		})
		b.NumLists[k] = l
	}
	for k, l := range c.rules.TextLists {
		l := copyRules(l).(TextListRules)
		sort.SliceStable(l, func(i, j int) bool {
			return len(l[i].Filters) < len(l[j].Filters)
		})
		b.TextLists[k] = l
	}
	for k, l := range c.rules.ObjLists {
		l := copyRules(l).(ObjListRules)
		sort.SliceStable(l, func(i, j int) bool {
			return len(l[i].Filters) < len(l[j].Filters)
		})
		b.ObjLists[k] = l
	}
	for k, l := range c.rules.Executes {
		l := copyRules(l).(ExecuteRules)
		sort.SliceStable(l, func(i, j int) bool {
			return len(l[i].Filters) < len(l[j].Filters)
		})
		b.Executes[k] = l
	}
	return b
}

// where source is BoolRules, etc.rules.
func copyRules(list interface{}) interface{} {
	src := r.ValueOf(list)
	cnt := src.Len()
	dst := r.MakeSlice(src.Type(), cnt, cnt)
	r.Copy(dst, src)
	return dst.Interface()
}
