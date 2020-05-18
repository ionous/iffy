package parser_test

import (
	"strings"
	"testing"

	"github.com/ionous/iffy/ident"
	. "github.com/ionous/iffy/parser"
	"github.com/ionous/inflect"
	"github.com/ionous/sliceOf"
)

// MyObject provides an example ( for testing ) of mapping an "Noun" to a NounInstance.
type MyObject struct {
	Id         ident.Id
	Names      []string
	Classes    []string
	Attributes []string
}

func (m *MyObject) String() string {
	return m.Id.Name
}

type MyBounds []*MyObject

func (m MyBounds) Get(r rune) NounInstance {
	return MyAdapter{m[r-'a']}
}

func (m MyBounds) Many(rs ...rune) (ret []NounInstance) {
	for _, r := range rs {
		ret = append(ret, m.Get(r))
	}
	return
}

func (m MyBounds) GetPlayerBounds(string) (Bounds, error) {
	return m, nil
}
func (m MyBounds) GetObjectBounds(ident.Id) (Bounds, error) {
	return m, nil
}
func (m MyBounds) IsPlural(word string) bool {
	return word != inflect.Singularize(word)
}

func (m MyBounds) SearchBounds(v NounVisitor) (ret bool) {
	n := MyAdapter{}
	for _, k := range m {
		n.MyObject = k
		if v(n) {
			ret = true
			break
		}
	}
	return
}

type MyAdapter struct {
	*MyObject
}

func (adapt MyAdapter) Id() ident.Id {
	return adapt.MyObject.Id
}

func (adapt MyAdapter) HasName(name string) bool {
	return MatchAny(name, adapt.Names)
}

func (adapt MyAdapter) HasClass(cls string) bool {
	return MatchAny(cls, adapt.Classes)
}

func (adapt MyAdapter) HasPlural(plural string) bool {
	// we'll use classes as plurals for tests --
	// its possible that might be different for the runtime
	// ex. might check plural / printed names
	return MatchAny(plural, adapt.Classes)
}

func (adapt MyAdapter) HasAttribute(attr string) bool {
	return MatchAny(attr, adapt.Attributes)
}

func MatchAny(n string, l []string) (okay bool) {
	for _, s := range l {
		if strings.EqualFold(n, s) {
			okay = true
			break
		}
	}
	return
}

func TestBounds(t *testing.T) {
	ctx := MyBounds{
		&MyObject{Id: ident.IdOf("a"), Names: sliceOf.String("unique")},
		//
		&MyObject{Id: ident.IdOf("b"), Names: strings.Fields("exact")},
		&MyObject{Id: ident.IdOf("c"), Names: strings.Fields("exact match")},
		//
		&MyObject{Id: ident.IdOf("d"), Names: strings.Fields("inexact match")},
		&MyObject{Id: ident.IdOf("e"), Names: strings.Fields("inexact conflict")},
		//
		&MyObject{Id: ident.IdOf("f"),
			Names:      strings.Fields("filter"),
			Attributes: strings.Fields("attr"),
			Classes:    strings.Fields("class"),
		},
		&MyObject{Id: ident.IdOf("g"),
			Names:      strings.Fields("filter"),
			Attributes: strings.Fields("attr"),
		},
		&MyObject{Id: ident.IdOf("h"),
			Names:   strings.Fields("filter"),
			Classes: strings.Fields("class"),
		},
	}
	if res, e := matching(ctx, "unique"); e != nil {
		t.Fatal("error", e)
	} else if obj, ok := res.(ResolvedObject); !ok {
		t.Fatalf("%T", res)
	} else if obj.NounInstance != ctx.Get('a') {
		t.Fatal("mismatched", obj.NounInstance)
	} else if got, want := strings.Join(obj.Words, ","), "unique"; got != want {
		t.Fatal(got)
	}

	if res, e := matching(ctx, "exact match"); e != nil {
		t.Fatal("error", e)
	} else if obj, ok := res.(ResolvedObject); !ok {
		t.Fatalf("%T", res)
	} else if obj.NounInstance != ctx.Get('c') {
		t.Fatal("mismatched", obj.NounInstance)
	} else if got, want := strings.Join(obj.Words, ","), "exact,match"; got != want {
		t.Fatal(got)
	}

	if res, e := matchingFilter(ctx, "filter", "attr", "class"); e != nil {
		t.Fatal("error", e)
	} else if obj, ok := res.(ResolvedObject); !ok {
		t.Fatalf("%T", res)
	} else if obj.NounInstance != ctx.Get('f') {
		t.Fatal("mismatched", obj.NounInstance)
	} else if got, want := strings.Join(obj.Words, ","), "filter"; got != want {
		t.Fatal(got)
	}

	if res, e := matching(ctx, "inexact"); e == nil || res != nil {
		t.Fatal("expected error", e, res)
	} else if got, want := e.Error(), (AmbiguousObject{
		Nouns: ctx.Many('d', 'e'),
		Depth: 1,
	}).Error(); got != want {
		t.Fatal(got)
	}

	if res, e := matching(ctx, "nothing"); e == nil || res != nil {
		t.Fatal("expected error", e, res)
	}
}

func matching(ctx Context, phrase string) (ret Result, err error) {
	match := &Noun{}
	words := strings.Fields(phrase)
	if bounds, e := ctx.GetPlayerBounds(""); e != nil {
		err = e
	} else {
		ret, err = match.Scan(ctx, bounds, Cursor{Words: words})
	}
	return
}

func matchingFilter(ctx Context, phrase, attr, class string) (ret Result, err error) {
	match := &Noun{Filters{&HasAttr{attr}, &HasClass{class}}}
	words := strings.Fields(phrase)
	if bounds, e := ctx.GetPlayerBounds(""); e != nil {
		err = e
	} else {
		ret, err = match.Scan(ctx, bounds, Cursor{Words: words})
	}
	return
}
