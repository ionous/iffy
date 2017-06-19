package spec

import "github.com/ionous/errutil"

func NewArray(spec Spec) Array {
	return NewBlock(spec).Array()
}

func NewBlock(spec Spec) *Context {
	return &Context{spec, true}
}

// Context provides tools for building commands.
type Context struct {
	spec Spec
	Args bool // Args is true when we are ready for positional arguments.
}

// Cmd creates a new command of name with the passed set of non-keywords args. It returns a new context, allowing the user to add additional (key-value, or linear) arguments to the new coommand.
func (ctx *Context) Cmd(name string, args ...interface{}) *Context {
	newSpec := newSpec(ctx.spec, name, args)
	ctx.addArg(newSpec)
	return &Context{newSpec, true}
}

// Value specifies a single literal: whether one primitive value or one array of primitive values. It does not start a new block, because primitive values have no additional parameters.
func (ctx *Context) Value(arg interface{}) {
	ctx.addArg(arg)
}

// Array specifies a new array parameter.
func (ctx *Context) Array() (ret Array) {
	if specs, e := ctx.spec.NewSpecs(); e != nil {
		panic(e)
	} else {
		ctx.addArg(specs)
		ret = Array{ctx, specs, true}
	}
	return ret
}

// Param adds a key-value parameter to the spec.
// The passed name is the key; it returns a Chain for specifying a value.
func (ctx *Context) Param(name string) Chain {
	ctx.Args = false
	return Chain{ctx, name}
}

func newSpec(sf SpecFactory, name string, args []interface{}) (ret Spec) {
	var err error
	if newSpec, e := sf.NewSpec(name); e != nil {
		err = e
	} else {
		for _, arg := range args {
			if e := newSpec.Position(arg); e != nil {
				err = e
				break
			}
		}
		if err == nil {
			ret = newSpec
		}
	}
	if err != nil {
		panic(err)
	}
	return
}

func (ctx *Context) addArg(arg interface{}) {
	var err error
	if !ctx.Args {
		err = errutil.New("cant use positional arguments after key values")
	} else {
		err = ctx.spec.Position(arg)
	}
	if err != nil {
		panic(err)
	}
}

func (ctx *Context) assign(key string, value interface{}) {
	if e := ctx.spec.Assign(key, value); e != nil {
		panic(e)
	}
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
	newSpec := newSpec(c.ctx.spec, name, args)
	c.ctx.assign(c.key, newSpec)
	return &Context{newSpec, true}
}

// Value sets the value of a key specified via Context.Param.
// See also: Context.Value
func (c Chain) Value(arg interface{}) {
	c.ctx.assign(c.key, arg)
}

// Array starts a new list of commands.
func (c Chain) Array() (ret Array) {
	if specs, e := c.ctx.spec.NewSpecs(); e != nil {
		panic(e)
	} else {
		c.ctx.assign(c.key, specs)
		ret = Array{c.ctx, specs, true}
	}
	return
}

// Array adds commands to arrays.
type Array struct {
	ctx   *Context
	specs Specs
	Cmds  bool
}

// Cmd adds a new command to the array.
// See also: Context.Cmd
func (a Array) Cmd(name string, args ...interface{}) *Context {
	newSpec := newSpec(a.ctx.spec, name, args)
	if e := a.specs.AddElement(newSpec); e != nil {
		panic(e)
	}
	return &Context{newSpec, true}
}
