package ops

import (
	"github.com/ionous/iffy/id"
	"github.com/ionous/iffy/ref/unique"
	"github.com/kr/pretty"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"strconv"
	"testing"
)

type SomeInterface interface {
	DoStuff() string
}

type Container struct {
	One, Two SomeInterface
	Value    int
	More     []SomeInterface
}

type Contents struct {
	Name string
}

func (p *Container) DoStuff() string {
	return strconv.Itoa(p.Value)
}

func (c *Contents) DoStuff() string {
	return c.Name
}

var testData = &Container{
	One: &Contents{"all are one"},
	Two: &Contents{"dilute, dilute"},
	More: []SomeInterface{
		&Container{Value: 5},
		&Container{Value: 7},
	},
}

// TODO:
// 1. test unknown commands
// 2. mismatched element types
func TestOps(t *testing.T) {
	suite.Run(t, new(OpsSuite))
}

type OpsSuite struct {
	suite.Suite
	ops  *Ops
	test *testing.T
}

func (t *OpsSuite) SetupTest() {
	ops := NewOps()
	unique.RegisterTypes(unique.PanicTypes(ops),
		(*Container)(nil), (*Contents)(nil))
	t.ops = ops
	t.test = t.T()
}

func (t *OpsSuite) TestKeyValue() {
	var root Container
	if c, ok := t.ops.NewBuilder(&root); ok {
		c.Param("Value").Val(4)
		//
		if _, e := c.Build(); t.NoError(e) {
			t.EqualValues(4, root.Value)
		}
	}
}

func (t *OpsSuite) TestAllAreOne() {
	var root Container
	if c, ok := t.ops.NewBuilder(&root); ok {
		// the simple way:
		c.Cmd("contents", "all are one")
		// // cause why not:
		if c.Cmd("contents").Begin() {
			c.Val("dilute, dilute").End()
		}
		if c.Param("more").Cmds().Begin() {
			c.Cmd("container", c.Param("value").Val(5))
			c.Cmd("container", c.Param("value").Val(7))
			c.End()
		}
		if _, e := c.Build(); t.NoError(e) {
			t.EqualValues(*testData, root)
			t.test.Log(pretty.Sprint(root))
		}
	}
}

type CommandBlock struct {
	*Container
	*Contents
}

// TestOpsBlock ensures blocks of commands register succesfully.
func TestOpsBlock(t *testing.T) {
	assert := assert.New(t)
	ops := NewOps((*CommandBlock)(nil))
	assert.Contains(ops.Types, id.MakeId("Container"))
	assert.Contains(ops.Types, id.MakeId("Contents"))
}
