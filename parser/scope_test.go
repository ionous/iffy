package parser_test

import (
	"bitbucket.org/pkg/inflect"
	"github.com/ionous/iffy/ident"
	. "github.com/ionous/iffy/parser"
	"github.com/ionous/sliceOf"
	testify "github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

// MyObject provides an example ( for testing ) of mapping an "Noun" to a NounVisitor.
type MyObject struct {
	Id         ident.Id
	Names      []string
	Classes    []string
	Attributes []string
}

func (m *MyObject) String() string {
	return m.Id.Name
}

type MyScope []*MyObject

func (m MyScope) Get(r rune) NounVisitor {
	return MyAdapter{m[r-'a']}
}

func (m MyScope) Many(rs ...rune) (ret []NounVisitor) {
	for _, r := range rs {
		ret = append(ret, m.Get(r))
	}
	return
}

func (m MyScope) GetPlayerScope(name ident.Id) (Scope, error) {
	return m, nil
}
func (m MyScope) GetOtherScope(name ident.Id) (Scope, error) {
	return m, nil
}
func (m MyScope) IsPlural(word string) bool {
	return word != inflect.Singularize(word)
}

func (m MyScope) SearchScope(v func(n NounVisitor) bool) (ret bool) {
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

func (adapt MyAdapter) GetId() ident.Id {
	return adapt.Id
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

func TestScope(t *testing.T) {
	assert := testify.New(t)
	ctx := MyScope{
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
	if res, e := matching(ctx, "unique"); assert.NoError(e) {
		assert.EqualValues(ResolvedObject{
			NounVisitor: ctx.Get('a'),
			Words:       sliceOf.String("unique"),
		}, res)
	}

	if res, e := matching(ctx, "exact match"); assert.NoError(e) {
		assert.EqualValues(ResolvedObject{
			NounVisitor: ctx.Get('c'),
			Words:       sliceOf.String("exact", "match"),
		}, res)
	}

	if res, e := matchingFilter(ctx, "filter", "attr", "class"); assert.NoError(e) {
		assert.EqualValues(ResolvedObject{
			NounVisitor: ctx.Get('f'),
			Words:       sliceOf.String("filter"),
		}, res)
	}

	if res, e := matching(ctx, "inexact"); assert.Error(e) {
		assert.Nil(res)
		assert.EqualValues(AmbiguousObject{
			Nouns: ctx.Many('d', 'e'),
			// Words: sliceOf.String("inexact"),
			Depth: 1,
		}, e)
	}

	if res, e := matching(ctx, "nothing"); assert.Error(e) {
		assert.Nil(res)
	}
}

func matching(ctx Context, phrase string) (ret Result, err error) {
	match := &Noun{}
	words := strings.Fields(phrase)
	if scope, e := ctx.GetPlayerScope(ident.None()); e != nil {
		err = e
	} else {
		ret, err = match.Scan(ctx, scope, Cursor{Words: words})
	}
	return
}

func matchingFilter(ctx Context, phrase, attr, class string) (ret Result, err error) {
	match := &Noun{Filters{&HasAttr{attr}, &HasClass{class}}}
	words := strings.Fields(phrase)
	if scope, e := ctx.GetPlayerScope(ident.None()); e != nil {
		err = e
	} else {
		ret, err = match.Scan(ctx, scope, Cursor{Words: words})
	}
	return
}
