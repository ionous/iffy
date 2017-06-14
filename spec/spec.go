package spec

import "github.com/ionous/errutil"

// TODO: detection of primitives arrays, replacing them with literals.

// Command represents some operation to be undertaken by a game.
type Command struct {
	name string
	args []interface{}
	keys map[string]interface{}
}

// Cmd creates a new command using the global context.
// See also: Context.Cmd
func Cmd(name string, params ...interface{}) *Context {
	return global.Cmd(name, params...)
}

var global = NewArray()

// Context provides tools for building commands.
type Context struct {
	cmd  *Command
	Args bool
}

func (c *Context) Command() *Command {
	return c.cmd
}

func (c *Context) addArg(arg interface{}) {
	if len(c.cmd.keys) > 0 {
		err := errutil.New("cant use linear arguments after key-value arguments in", c.cmd.name, c.cmd.keys)
		panic(err)
	}
	c.cmd.args = append(c.cmd.args, arg)
}

// Cmd creates a new command of name with the passed set of non-keywords args. It returns a new context, allowing the user to add additional (key-value, or linear) arguments to the new coommand.
func (c *Context) Cmd(name string, params ...interface{}) *Context {
	cmd := &Command{name, params, nil}
	c.addArg(cmd)
	return &Context{cmd, true}
}

// Value specifies a single literal: whether one primitive value or one array of primitive values. It does not start a new block, because primitive values have no additional parameters.
func (c *Context) Value(arg interface{}) {
	c.addArg(arg)
}

// Array
func (c *Context) Array() Array {
	a := NewArray()
	c.addArg(a.els)
	return a
}

// Param adds a key-value parameter to the cmd command.
// The passed name is the key; it returns a Chain for specifying a value.
func (c *Context) Param(name string) Chain {
	if c.cmd.keys == nil {
		c.cmd.keys = make(map[string]interface{})
	}
	return Chain{c, name}
}

// Chain specifies the value of a Param.
// A user can choose to specify a command or primitive value.
type Chain struct {
	ctx *Context
	key string
}

// Cmd sets the value of a key specified via Context.Param.
// See also: Context.Cmd
func (c Chain) Cmd(name string, args ...interface{}) *Context {
	cmd := &Command{name, args, nil}
	c.ctx.cmd.keys[c.key] = cmd
	return &Context{cmd, true}
}

// Value sets the value of a key specified via Context.Param.
// See also: Context.Value
func (c Chain) Value(arg interface{}) {
	c.ctx.cmd.keys[c.key] = arg
}

// Array starts a new list of commands.
func (c Chain) Array() Array {
	a := NewArray()
	c.ctx.cmd.keys[c.key] = a.els
	return a
}

func NewArray() Array {
	els := &([]*Command{})
	return Array{els, true}
}

// Array adds commands to arrays.
type Array struct {
	els *[]*Command
	Els bool
}

// Cmd adds a new command to the array.
// See also: Context.Cmd
func (a Array) Cmd(name string, args ...interface{}) *Context {
	cmd := &Command{name, args, nil}
	(*a.els) = append((*a.els), cmd)
	return &Context{cmd, true}
}

func (a Array) Commands() []*Command {
	return *a.els
}
