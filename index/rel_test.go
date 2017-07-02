package index

import (
	"github.com/ionous/sliceOf"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

func TestRelation(t *testing.T) {
	suite.Run(t, new(RelSuite))
}

type RelSuite struct {
	suite.Suite
	relation Relation
}

var pairs = []struct {
	primary   string
	secondary []string
}{
	{"claire", sliceOf.String("boomba", "rocky")},
	{"grace", sliceOf.String("plume")},
	{"ipjay", sliceOf.String("loofa")},
	{"marja", sliceOf.String("petra")},
}

func (assert *RelSuite) SetupTest() {
	n := MakeRelation(OneToMany)
	for _, pair := range pairs {
		for _, v := range pair.secondary {
			data := strings.Join(sliceOf.String(pair.primary, v), "+")
			changed := n.Relate(pair.primary, v, data)
			assert.True(changed)
		}
	}
	assert.relation = n
}

func (assert *RelSuite) TestConstruct() {
	t, n := assert.T(), assert.relation
	t.Log("primary", getKeys(n.Index[Primary], Primary))
	t.Log("secondary", getKeys(n.Index[Secondary], Secondary))
	for _, pair := range pairs {
		collect := collect(t, n, pair.primary)
		assert.EqualValues(pair.secondary, collect, pair.primary)
	}
}

func (assert *RelSuite) TestRemoveSecondary() {
	t, n := assert.T(), assert.relation
	t.Log(n.Index[0].Lines)
	for _, pair := range pairs {
		name, remove := pair.primary, pair.secondary[0]
		before := collect(t, n, name)
		assert.Len(before, len(pair.secondary))
		assert.Contains(before, remove)
		t.Log("removing", remove, "from", name)
		if changed := n.Relate("", remove, ""); assert.True(changed) {
			after := collect(t, n, name)
			t.Log("after", after)
			assert.Len(after, len(pair.secondary)-1)
			assert.NotContains(after, remove)
		}
	}
}

func collect(t *testing.T, n Relation, key string) (ret []string) {
	visits := n.Index[Primary].Walk(key, func(other, data string) bool {
		ret = append(ret, other)
		return false
	})
	if visits != len(ret) {
		t.Fatal("count should match calls", visits, len(ret))
	}
	return
}
