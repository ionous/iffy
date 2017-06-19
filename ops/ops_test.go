package ops

import (
	"github.com/ionous/iffy/reflector"
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
	assert := assert.New(t)
	ops := NewOps()
	ops.RegisterType((*Container)(nil))
	ops.RegisterType((*Contents)(nil))
	{
		var root Container
		if c := ops.Build(&root); c.Args {
			c.Param("Value").Value(4)
		}
		assert.EqualValues(4, root.Value)
	}
	{
		var root Container
		if c := ops.Build(&root); c.Args {
			c.Cmd("contents", "all are one")
			c.Cmd("contents").Value("dilute, dilute")
			if c := c.Param("more").Array(); c.Cmds {
				c.Cmd("container").Param("value").Value(5)
				c.Cmd("container").Param("value").Value(7)
			}
		}
		assert.EqualValues(*testData, root)
	}
}

type CommandBlock struct {
	*Container
	*Contents
}

// just make sure we can register a block of commands succesfully.
func TestOpsBlock(t *testing.T) {
	assert := assert.New(t)
	ops := NewOps((*CommandBlock)(nil))
	assert.Contains(ops.names, reflector.MakeId("Container"))
	assert.Contains(ops.names, reflector.MakeId("Contents"))
}
