package ops

import (
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/ref/unique"
	"github.com/kr/pretty"
	"github.com/stretchr/testify/assert"
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
	cmds := NewOps(nil)
	unique.PanicTypes(cmds,
		(*Container)(nil), (*Contents)(nil))
	t.Run("KeyValue", func(t *testing.T) {
		var root Container
		assert := assert.New(t)
		c := cmds.NewBuilder(&root, DefaultXform{})
		c.Param("Value").Val(4)
		//
		if e := c.Build(); assert.NoError(e) {
			assert.EqualValues(4, root.Value)

		}
	})
	t.Run("AllAreOne", func(t *testing.T) {
		var root Container
		assert := assert.New(t)
		c := cmds.NewBuilder(&root, DefaultXform{})
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
		if e := c.Build(); assert.NoError(e) {
			assert.EqualValues(*testData, root)
			t.Log(pretty.Sprint(root))
		}
	})
}

type CommandBlock struct {
	*Container
	*Contents
}

// TestOpsBlock ensures blocks of commands register successfully.
func TestOpsBlock(t *testing.T) {
	assert := assert.New(t)
	cmds := NewOps(nil)
	unique.PanicBlocks(cmds,
		(*CommandBlock)(nil))
	assert.Contains(cmds.Types, ident.IdOf("Container"))
	assert.Contains(cmds.Types, ident.IdOf("Contents"))
}
