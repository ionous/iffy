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
	assert.EqualValues([]Ranking{{0, []string{"a"}}},
		matching(scope, "unique"))
	assert.EqualValues([]Ranking{{1, []string{"c"}}},
		matching(scope, "exact match"))
	assert.EqualValues([]Ranking{{0, []string{"d", "e"}}},
		matching(scope, "inexact"))
	assert.EqualValues([]Ranking{{0, []string{"f"}}},
		matchingFilter(scope, "filter", "attr", "class"))
	assert.Empty(
		matching(scope, "nothing"))
}

func matching(scope Scope, phrase string) (ret []Ranking) {
	ctx := Context{
		Scope: scope,
		Match: &Match{Scanner: &Object{}},
		Words: strings.Fields(phrase),
	}
	if ctx.Advance() {
		ret = ctx.Results.Matches
	}
	return
}

func matchingFilter(scope Scope, phrase, attr, class string) (ret []Ranking) {
	ctx := Context{
		Scope: scope,
		Match: &Match{Scanner: &Object{Filters{&HasAttr{attr}, &HasClass{class}}}},
		Words: strings.Fields(phrase),
	}
	if ctx.Advance() {
		ret = ctx.Results.Matches
	}
	return
}
