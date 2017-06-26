package reflector

import (
	"github.com/stretchr/testify/suite"
	r "reflect"
	"testing"
)

func TestMetadata(t *testing.T) {
	suite.Run(t, new(MetadataSuite))
}

type MetadataSuite struct {
	suite.Suite
}

func (t *MetadataSuite) MakeMetadata(tag string) (out Metadata) {
	s := MakeTag(r.StructTag(tag))
	if len(s) > 0 {
		out = make(Metadata)
		out.AddString(s, "fill")
	}
	return
}

func (t *MetadataSuite) TestKey() {
	m := t.MakeMetadata(`if:"id"`)
	t.Len(m, 1)
	t.Equal("fill", m["id"])
}
func (t *MetadataSuite) TestKeyValue() {
	m := t.MakeMetadata(`if:"plural:tests"`)
	t.Len(m, 1)
	t.Equal("tests", m["plural"])
}
func (t *MetadataSuite) TestMultipleKeys() {
	m := t.MakeMetadata(`if:"id,plural:tests"`)
	t.Len(m, 2)
	t.Equal("fill", m["id"])
	t.Equal("tests", m["plural"])
}
func (t *MetadataSuite) TestMissingCommas() {
	m := t.MakeMetadata(`if:"edge:test:tests"`)
	t.Len(m, 1)
	t.Equal("test:tests", m["edge"])
}

func TestMetaTags(t *testing.T) {
	suite.Run(t, new(MetaTagSuite))
}

type MetaTagSuite struct {
	suite.Suite
}

func (t *MetaTagSuite) TestInvalidTag() {
	tag := MakeTag(r.StructTag(`fi:"id"`))
	t.EqualValues("", tag)
}
func (t *MetaTagSuite) TestKey() {
	tag := MakeTag(r.StructTag(`if:"id"`))
	if id, ok := tag.Find("id"); t.True(ok) {
		t.Empty(id)
	}
}
func (t *MetaTagSuite) TestSimilarTag() {
	tag := MakeTag(r.StructTag(`if:"iddy:5"`))
	_, ok := tag.Find("id")
	t.False(ok)
}
func (t *MetaTagSuite) TestKeyValue() {
	tag := MakeTag(r.StructTag(`if:"plural:tests"`))
	if _, ok := tag.Find("id"); t.False(ok) {
		//
		if res, ok := tag.Find("plural"); t.True(ok) {
			t.Equal("tests", res)
		}
	}
}
func (t *MetaTagSuite) TestMultipleKeys() {
	tag := MakeTag(r.StructTag(`if:"id,other:true,plural:tests"`))
	if id, ok := tag.Find("id"); t.True(ok) {
		t.Empty(id)

		if other, ok := tag.Find("other"); t.True(ok) {
			t.Equal("true", other)
			//
			if res, ok := tag.Find("plural"); t.True(ok) {
				t.Equal("tests", res)
			}
		}
	}
}
func (t *MetaTagSuite) TestMissingCommas() {
	tag := MakeTag(r.StructTag(`if:"edge:test:tests"`))
	if res, ok := tag.Find("edge"); t.True(ok) {
		t.Equal("test:tests", res)
	}
}
