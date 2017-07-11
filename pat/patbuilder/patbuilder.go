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
	*ref.Classes
	objectClasses *ref.Classes
	patterns      pat.Patterns
}

func NewPatterns(objectClasses *ref.Classes) *Patterns {
	return &Patterns{
		ref.NewClasses(),
		objectClasses,
		pat.Patterns{
			make(pat.BoolMap),
			make(pat.NumberMap),
			make(pat.TextMap),
			make(pat.ObjectMap),
			make(pat.NumListMap),
			make(pat.TextListMap),
			make(pat.ObjListMap),
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
func (b *Patterns) RegisterType(rtype r.Type) (err error) {
	if _, e := b.objectClasses.RegisterClass(rtype); e != nil {
		err = errutil.New("couldnt register pattern to object classes", e)
	} else if _, e := b.Classes.RegisterClass(rtype); e != nil {
		err = errutil.New("couldnt register pattern to pattern classes", e)
	}
	return
}

// FIX? consider using reflection to implement the various pattern builder methods.
// the code duplication is painful.
func (b *Patterns) getPattern(name string) (ret rt.Class, err error) {
	if cls, exists := b.Classes.GetClass(name); !exists {
		err = errutil.New("unknown pattern", name)
	} else {
		ret = cls
	}
	return
}

func (b *Patterns) Bool(name string, filter rt.BoolEval, k rt.BoolEval) (err error) {
	if cls, e := b.getPattern(name); e != nil {
		err = e
	} else {
		id, filters := cls.GetId(), ExpandFilters(filter)
		l := b.patterns.BoolMap[id]
		l = append(l, pat.BoolPattern{filters, k})
		b.patterns.BoolMap[id] = l
	}
	return
}

func (b *Patterns) Number(name string, filter rt.BoolEval, k rt.NumberEval) (err error) {
	if cls, e := b.getPattern(name); e != nil {
		err = e
	} else {
		id, filters := cls.GetId(), ExpandFilters(filter)
		l := b.patterns.NumberMap[id]
		l = append(l, pat.NumberPattern{filters, k})
		b.patterns.NumberMap[id] = l
	}
	return
}

func (b *Patterns) Text(name string, filter rt.BoolEval, k rt.TextEval) (err error) {
	if cls, e := b.getPattern(name); e != nil {
		err = e
	} else {
		id, filters := cls.GetId(), ExpandFilters(filter)
		l := b.patterns.TextMap[id]
		l = append(l, pat.TextPattern{filters, k})
		b.patterns.TextMap[id] = l
	}
	return
}

func (b *Patterns) Object(name string, filter rt.BoolEval, k rt.ObjectEval) (err error) {
	if cls, e := b.getPattern(name); e != nil {
		err = e
	} else {
		id, filters := cls.GetId(), ExpandFilters(filter)
		l := b.patterns.ObjectMap[id]
		l = append(l, pat.ObjectPattern{filters, k})
		b.patterns.ObjectMap[id] = l
	}
	return
}

func (b *Patterns) NumList(name string, filter rt.BoolEval, k rt.NumListEval) (err error) {
	if cls, e := b.getPattern(name); e != nil {
		err = e
	} else {
		id, filters := cls.GetId(), ExpandFilters(filter)
		l := b.patterns.NumListMap[id]
		l = append(l, pat.NumListPattern{filters, k})
		b.patterns.NumListMap[id] = l
	}
	return
}

func (b *Patterns) TextList(name string, filter rt.BoolEval, k rt.TextListEval) (err error) {
	if cls, e := b.getPattern(name); e != nil {
		err = e
	} else {
		id, filters := cls.GetId(), ExpandFilters(filter)
		l := b.patterns.TextListMap[id]
		l = append(l, pat.TextListPattern{filters, k})
		b.patterns.TextListMap[id] = l
	}
	return
}

func (b *Patterns) ObjList(name string, filter rt.BoolEval, k rt.ObjListEval) (err error) {
	if cls, e := b.getPattern(name); e != nil {
		err = e
	} else {
		id, filters := cls.GetId(), ExpandFilters(filter)
		l := b.patterns.ObjListMap[id]
		l = append(l, pat.ObjListPattern{filters, k})
		b.patterns.ObjListMap[id] = l
	}
	return
}

func (b *Patterns) GetPatterns() pat.Patterns {
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
	return b.patterns
}
