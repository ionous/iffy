package index

import (
	// "github.com/ionous/iffy/index"
	"github.com/ionous/sliceOf"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestTable(t *testing.T) {
	suite.Run(t, new(RelSuite))
}

type RelSuite struct {
	suite.Suite
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

func (assert *RelSuite) MakeRelation() (ret *Table) {
	n := MakeTable(OneToMany)
	for _, pair := range pairs {
		for _, v := range pair.secondary {
			changed, e := n.Relate(pair.primary, v, NoData)
			assert.NoError(e)
			assert.True(changed)
		}
	}
	return &n
}

func (assert *RelSuite) TestConstruct() {
	t, n := assert.T(), assert.MakeRelation()
	t.Log("primary", getKeys(n.Primary, unpackMajor))
	t.Log("secondary", getKeys(n.Secondary, unpackMinor))
	for _, pair := range pairs {
		collect := collect(t, n, pair.primary)
		assert.EqualValues(pair.secondary, collect, pair.primary)
	}
}

func (assert *RelSuite) TestTypes() {
	m := map[Type]struct{ uni, que bool }{
		ManyToMany: {uni: false, que: false},
		ManyToOne:  {uni: true, que: false},
		OneToMany:  {uni: false, que: true},
		OneToOne:   {uni: true, que: true},
	}
	for i := 0; i < 4; i++ {
		i := Type(i)
		n := MakeTable(i)
		assert.Equal(i, n.Type())
		test := m[i]
		assert.Equal(test.uni, n.Primary.Unique)
		assert.Equal(test.que, n.Secondary.Unique)
	}
}

func (assert *RelSuite) TestRemoveMinor() {
	t, n := assert.T(), assert.MakeRelation()
	t.Log(n.Primary.Rows)
	for _, pair := range pairs {
		name, remove := pair.primary, pair.secondary[0]
		before := collect(t, n, name)
		assert.Len(before, len(pair.secondary))
		assert.Contains(before, remove)
		t.Log("removing", remove, "from", name)
		if changed, e := n.Relate("", remove, nil); assert.True(changed) && assert.NoError(e) {
			after := collect(t, n, name)
			t.Log("after", after)
			assert.Len(after, len(pair.secondary)-1)
			assert.NotContains(after, remove)
		}
	}
}

func collect(t *testing.T, n *Table, key string) (ret []string) {
	visits := n.Primary.Walk(key, func(other string) bool {
		ret = append(ret, other)
		return false
	})
	if visits != len(ret) {
		t.Fatal("count should match calls", visits, len(ret))
	}
	return
}
