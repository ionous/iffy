package parser_test

import (
	. "github.com/ionous/iffy/parser"
	"github.com/ionous/sliceOf"
	testify "github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

// MyObject provides an example ( for testing ) of mapping an "Object" to a Noun.
type MyObject struct {
	Id         string
	Names      []string
	Classes    []string
	Attributes []string
}

type MyScope []MyObject

func (m MyScope) SearchScope(v NounVisitor) (ret bool) {
	n := MyAdapter{}
	for _, k := range m {
		n.MyObject = &k
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

func (adapt MyAdapter) GetId() string {
	return adapt.Id
}

func (adapt MyAdapter) HasName(name string) bool {
	return MatchAny(name, adapt.Names)
}

func (adapt MyAdapter) HasClass(cls string) bool {
	return MatchAny(cls, adapt.Classes)
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
	scope := MyScope{
		MyObject{Id: "a", Names: sliceOf.String("unique")},
		//
		MyObject{Id: "b", Names: strings.Fields("exact")},
		MyObject{Id: "c", Names: strings.Fields("exact match")},
		//
		MyObject{Id: "d", Names: strings.Fields("inexact match")},
		MyObject{Id: "e", Names: strings.Fields("inexact conflict")},
		//
		MyObject{Id: "f",
			Names:      strings.Fields("filter"),
			Attributes: strings.Fields("attr"),
			Classes:    strings.Fields("class"),
		},
		MyObject{Id: "g",
			Names:      strings.Fields("filter"),
			Attributes: strings.Fields("attr"),
		},
		MyObject{Id: "h",
			Names:   strings.Fields("filter"),
			Classes: strings.Fields("class"),
		},
	}
	if res, e := matching(scope, "unique"); assert.NoError(e) {
		assert.EqualValues(ResolvedObject{
			Id:    "a",
			Words: sliceOf.String("unique"),
		}, res)
	}

	if res, e := matching(scope, "exact match"); assert.NoError(e) {
		assert.EqualValues(ResolvedObject{
			Id:    "c",
			Words: sliceOf.String("exact", "match"),
		}, res)
	}

	if res, e := matchingFilter(scope, "filter", "attr", "class"); assert.NoError(e) {
		assert.EqualValues(ResolvedObject{
			Id:    "f",
			Words: sliceOf.String("filter"),
		}, res)
	}

	if res, e := matching(scope, "inexact"); assert.Error(e) {
		assert.Nil(res)
		assert.EqualValues(AmbiguousObject{
			Ids: sliceOf.String("d", "e"),
			// Words: sliceOf.String("inexact"),
			Depth: 1,
		}, e)
	}

	if res, e := matching(scope, "nothing"); assert.Error(e) {
		assert.Nil(res)
	}
}

func matching(scope Scope, phrase string) (ret Result, err error) {
	match := &Object{}
	words := strings.Fields(phrase)
	return match.Scan(scope, Cursor{Words: words})
}

func matchingFilter(scope Scope, phrase, attr, class string) (ret Result, err error) {
	match := &Object{Filters{&HasAttr{attr}, &HasClass{class}}}
	words := strings.Fields(phrase)
	return match.Scan(scope, Cursor{Words: words})
}
