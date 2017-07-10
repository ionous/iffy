package patbuilder

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/id"
	"github.com/ionous/iffy/pat"
	"github.com/ionous/iffy/rt"
	"sort"
)

type Builder struct {
	names    _NameMap
	patterns pat.Patterns
}
type _NameMap map[string]string

func NewBuilder() *Builder {
	return &Builder{
		make(_NameMap),
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

func (b *Builder) NewPattern(name string) (err error) {
	id := id.MakeId(name)
	if old, ok := b.names[id]; !ok {
		b.names[id] = name
	} else if old != name {
		err = errutil.New("conflicting pattern", name, "was", old, id)
	} else {
		err = errutil.New("pattern already exists", name)
	}
	return
}

func (b *Builder) AddMatch(name string, v interface{}, f ...rt.BoolEval) (err error) {
	id := id.MakeId(name)
	if _, exists := b.names[id]; !exists {
		err = errutil.New("unknown pattern", name)
	} else {
		switch k := v.(type) {
		case rt.BoolEval:
			l := b.patterns.BoolMap[id]
			l = append(l, pat.BoolPattern{f, k})
			b.patterns.BoolMap[id] = l
		case rt.NumberEval:
			l := b.patterns.NumberMap[id]
			l = append(l, pat.NumberPattern{f, k})
			b.patterns.NumberMap[id] = l
		case rt.TextEval:
			l := b.patterns.TextMap[id]
			l = append(l, pat.TextPattern{f, k})
			b.patterns.TextMap[id] = l
		case rt.ObjectEval:
			l := b.patterns.ObjectMap[id]
			l = append(l, pat.ObjectPattern{f, k})
			b.patterns.ObjectMap[id] = l
		case rt.NumListEval:
			l := b.patterns.NumListMap[id]
			l = append(l, pat.NumListPattern{f, k})
			b.patterns.NumListMap[id] = l
		case rt.TextListEval:
			l := b.patterns.TextListMap[id]
			l = append(l, pat.TextListPattern{f, k})
			b.patterns.TextListMap[id] = l
		case rt.ObjListEval:
			l := b.patterns.ObjListMap[id]
			l = append(l, pat.ObjListPattern{f, k})
			b.patterns.ObjListMap[id] = l
		default:
			err = errutil.Fmt("unknown pattern eval %T", v)
		}
	}
	return
}

func (b *Builder) GetPatterns() pat.Patterns {
	// sorts in ascending
	// more filters should be left in list
	for _, l := range b.patterns.BoolMap {
		sort.SliceStable(l, func(i, j int) bool {
			return len(l[i].Filters) < len(l[i].Filters)
		})
	}
	for _, l := range b.patterns.NumberMap {
		sort.SliceStable(l, func(i, j int) bool {
			return len(l[i].Filters) < len(l[i].Filters)
		})
	}
	for _, l := range b.patterns.TextMap {
		sort.SliceStable(l, func(i, j int) bool {
			return len(l[i].Filters) < len(l[i].Filters)
		})
	}
	for _, l := range b.patterns.ObjectMap {
		sort.SliceStable(l, func(i, j int) bool {
			return len(l[i].Filters) < len(l[i].Filters)
		})
	}
	for _, l := range b.patterns.NumListMap {
		sort.SliceStable(l, func(i, j int) bool {
			return len(l[i].Filters) < len(l[i].Filters)
		})
	}
	for _, l := range b.patterns.TextListMap {
		sort.SliceStable(l, func(i, j int) bool {
			return len(l[i].Filters) < len(l[i].Filters)
		})
	}
	for _, l := range b.patterns.ObjListMap {
		sort.SliceStable(l, func(i, j int) bool {
			return len(l[i].Filters) < len(l[i].Filters)
		})
	}
	return b.patterns
}
