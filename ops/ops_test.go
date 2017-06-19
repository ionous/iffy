package ops

import (
	// "github.com/kr/pretty"
	"github.com/stretchr/testify/assert"
	r "reflect"
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
	},
}

// TODO:
// 1. test unknown commands
// 2. mismatched eleemmnt types
func TestOps(t *testing.T) {
	assert := assert.New(t)
	ops := NewOps()
	ops.RegisterType(r.TypeOf((*Container)(nil)).Elem())
	ops.RegisterType(r.TypeOf((*Contents)(nil)).Elem())
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
			}
		}
		assert.EqualValues(*testData, root)
	}
}
