package patbuilder

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/core" // for AllTrue
	"github.com/ionous/iffy/pat"
	"github.com/ionous/iffy/ref"
	"github.com/ionous/iffy/rt"
	r "reflect"
	"sort"
)

type Patterns struct {
	classes  ref.Classes
	patterns pat.Patterns
}

// NewPatterns, adding our patterns to the passed classes
func NewPatterns(classes ref.Classes) *Patterns {
	return &Patterns{
		ref.NewClassStack(classes),
		// FIX? consider using reflection to implement the various patterns.
		// the code duplication is painful.
		pat.Patterns{
			make(pat.BoolMap),
			make(pat.NumberMap),
			make(pat.TextMap),
			make(pat.ObjectMap),
			make(pat.NumListMap),
			make(pat.TextListMap),
			make(pat.ObjListMap),
			make(pat.ExecuteMap),
		},
	}
}

// ExpandFilters turns a single bool eval into an array by looking at its type.
func ExpandFilters(eval rt.BoolEval) (ret []rt.BoolEval) {
	if eval != nil {
		if multi, ok := eval.(*core.AllTrue); ok {
			ret = multi.Test
		} else {
			ret = append(ret, eval)
		}
	}
	return
}

// RegisterType creates a new pattern.
// Compatible with unique.TypeRegistry
func (b *Patterns) RegisterType(rtype r.Type) error {
	_, e := b.classes.RegisterClass(rtype)
	return e
}

func (b *Patterns) getPattern(name string) (ret string, err error) {
	if cls, exists := b.classes.GetClass(name); !exists {
		err = errutil.New("unknown pattern", name)
	} else {
		ret = cls.GetId()
	}
	return
}

func (b *Patterns) AddBool(name string, filter rt.BoolEval, k rt.BoolEval) (err error) {
	if id, e := b.getPattern(name); e != nil {
		err = e
	} else {
		l := b.patterns.BoolMap[id]
		l = append(l, pat.BoolPattern{ExpandFilters(filter), k})
		b.patterns.BoolMap[id] = l
	}
	return
}

func (b *Patterns) AddNumber(name string, filter rt.BoolEval, k rt.NumberEval) (err error) {
	if id, e := b.getPattern(name); e != nil {
		err = e
	} else {
		l := b.patterns.NumberMap[id]
		l = append(l, pat.NumberPattern{ExpandFilters(filter), k})
		b.patterns.NumberMap[id] = l
	}
	return
}

func (b *Patterns) AddText(name string, filter rt.BoolEval, k rt.TextEval) (err error) {
	if id, e := b.getPattern(name); e != nil {
		err = e
	} else {
		l := b.patterns.TextMap[id]
		l = append(l, pat.TextPattern{ExpandFilters(filter), k})
		b.patterns.TextMap[id] = l
	}
	return
}

func (b *Patterns) AddObject(name string, filter rt.BoolEval, k rt.ObjectEval) (err error) {
	if id, e := b.getPattern(name); e != nil {
		err = e
	} else {
		l := b.patterns.ObjectMap[id]
		l = append(l, pat.ObjectPattern{ExpandFilters(filter), k})
		b.patterns.ObjectMap[id] = l
	}
	return
}

func (b *Patterns) AddNumList(name string, filter rt.BoolEval, k rt.NumListEval) (err error) {
	if id, e := b.getPattern(name); e != nil {
		err = e
	} else {
		l := b.patterns.NumListMap[id]
		l = append(l, pat.NumListPattern{ExpandFilters(filter), k})
		b.patterns.NumListMap[id] = l
	}
	return
}

func (b *Patterns) AddTextList(name string, filter rt.BoolEval, k rt.TextListEval) (err error) {
	if id, e := b.getPattern(name); e != nil {
		err = e
	} else {
		l := b.patterns.TextListMap[id]
		l = append(l, pat.TextListPattern{ExpandFilters(filter), k})
		b.patterns.TextListMap[id] = l
	}
	return
}

func (b *Patterns) AddObjList(name string, filter rt.BoolEval, k rt.ObjListEval) (err error) {
	if id, e := b.getPattern(name); e != nil {
		err = e
	} else {
		l := b.patterns.ObjListMap[id]
		l = append(l, pat.ObjListPattern{ExpandFilters(filter), k})
		b.patterns.ObjListMap[id] = l
	}
	return
}

func (b *Patterns) AddExecList(name string, filter rt.BoolEval, k rt.Execute, flags pat.Flags) (err error) {
	if id, e := b.getPattern(name); e != nil {
		err = e
	} else {
		l := b.patterns.ExecuteMap[id]
		l = append(l, pat.ExecutePattern{ExpandFilters(filter), k, flags})
		b.patterns.ExecuteMap[id] = l
	}
	return
}

func (b *Patterns) Build() pat.Patterns {
	// sorts in ascending
	// more filters should be left in list
	for _, l := range b.patterns.BoolMap {
		sort.SliceStable(l, func(i, j int) bool {
			return len(l[i].Filters) < len(l[j].Filters)
		})
	}
	for _, l := range b.patterns.NumberMap {
		sort.SliceStable(l, func(i, j int) bool {
			return len(l[i].Filters) < len(l[j].Filters)
		})
	}
	for _, l := range b.patterns.TextMap {
		sort.SliceStable(l, func(i, j int) bool {
			return len(l[i].Filters) < len(l[j].Filters)
		})
	}
	for _, l := range b.patterns.ObjectMap {
		sort.SliceStable(l, func(i, j int) bool {
			return len(l[i].Filters) < len(l[j].Filters)
		})
	}
	for _, l := range b.patterns.NumListMap {
		sort.SliceStable(l, func(i, j int) bool {
			return len(l[i].Filters) < len(l[j].Filters)
		})
	}
	for _, l := range b.patterns.TextListMap {
		sort.SliceStable(l, func(i, j int) bool {
			return len(l[i].Filters) < len(l[j].Filters)
		})
	}
	for _, l := range b.patterns.ObjListMap {
		sort.SliceStable(l, func(i, j int) bool {
			return len(l[i].Filters) < len(l[j].Filters)
		})
	}
	for _, l := range b.patterns.ExecuteMap {
		sort.SliceStable(l, func(i, j int) bool {
			return len(l[i].Filters) < len(l[j].Filters)
		})
	}
	return b.patterns
}
