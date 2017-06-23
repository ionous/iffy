package cmd

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/spec"
	"github.com/ionous/iffy/spec/builder"
)

// Command represents an arbitrary game statement or operation.
type Command struct {
	Name string
	Args []interface{}
	Keys map[string]interface{}
}

// Commands stores a slice of type Command.
type Commands struct {
	Els []*Command
}

// _CommandFactory implements SpecFactory.
type _CommandFactory struct{}

// _CommandBuilder implements Spec.
type _CommandBuilder struct {
	cmd *Command
}

// _CommandsBuilder implements Specs.
type _CommandsBuilder struct {
	cmds *Commands
}

// RootBuilder provides a package spec implementation for Command/s.
type RootBuilder struct {
	root *Command
	builder.Builder
}

// NewBuilder returns a spec Builder used to assemble commands.
func NewBuilder() (*RootBuilder, bool) {
	root := new(Command)
	spec := &_CommandBuilder{cmd: root}
	return &RootBuilder{
		root:    root,
		Builder: builder.NewBuilder(_CommandFactory{}, spec),
	}, true
}

// Build finializes the builder, returning the root Command.
func (u *RootBuilder) Build() (ret *Command, err error) {
	if u.root == nil {
		err = errutil.New("build can only be called once")
	} else if fini := u.Builder.End(); !fini {
		err = errutil.New("mismatched Block/End(s)")
	} else {
		ret, u.root = u.root, nil
	}
	return
}

// NewSpec implements spec.SpecFactory.
func (_CommandFactory) NewSpec(name string) (spec.Spec, error) {
	return &_CommandBuilder{cmd: &Command{Name: name}}, nil
}

// NewSpecs implements spec.SpecFactory.
func (_CommandFactory) NewSpecs() (spec.Specs, error) {
	cmds := &Commands{}
	return &_CommandsBuilder{cmds: cmds}, nil
}

// Position implements spec.Spec.
func (cb *_CommandBuilder) Position(arg interface{}) (err error) {
	if arg, e := cmdUnpack(arg); e != nil {
		err = e
	} else {
		cb.cmd.Args = append(cb.cmd.Args, arg)
	}
	return
}

// Assign implements spec.Spec.
func (cb *_CommandBuilder) Assign(key string, arg interface{}) (err error) {
	if arg, e := cmdUnpack(arg); e != nil {
		err = e
	} else if cb.cmd.Keys != nil {
		cb.cmd.Keys[key] = arg
	} else {
		cb.cmd.Keys = map[string]interface{}{
			key: arg,
		}
	}
	return
}

// AddElement implements spec.Specs
func (cbs *_CommandsBuilder) AddElement(s spec.Spec) (err error) {
	if cb, ok := s.(*_CommandBuilder); !ok {
		err = errutil.Fmt("unexpected element type %T", s)
	} else {
		// println("appending element")
		cbs.cmds.Els = append(cbs.cmds.Els, cb.cmd)
	}
	return
}

// cmdUnpack helps position and assign push real data into the args
func cmdUnpack(value interface{}) (ret interface{}, err error) {
	switch b := value.(type) {
	case *_CommandBuilder:
		ret = b.cmd
	case *_CommandsBuilder:
		ret = b.cmds
	case float64, string, int, []float64, []string:
		ret = value
	default:
		err = errutil.Fmt("assigning unexpected type %T", value)
	}
	return
}
