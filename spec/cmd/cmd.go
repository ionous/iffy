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

// CommandBuilder provides a package spec implementation for Command/s.
type CommandBuilder struct {
	builder.Builder
}

// NewBuilder returns a spec Builder used to assemble commands.
func NewBuilder() (*CommandBuilder, bool) {
	root := new(Command)
	spec := &_Spec{cmd: root}
	return &CommandBuilder{
		builder.NewBuilder(_Factory{}, spec),
	}, true
}

// Build finializes the builder, returning the root Command.
func (u *CommandBuilder) Build() (ret *Command, err error) {
	if res, e := u.Builder.Build(); e != nil {
		err = e
	} else if spec, ok := res.(*_Spec); !ok {
		err = errutil.Fmt("unknown error")
	} else {
		ret = spec.cmd
	}
	return
}

type _Factory struct{}

// NewSpec implements spec.Factory.
func (_Factory) NewSpec(name string) (spec.Spec, error) {
	return &_Spec{cmd: &Command{Name: name}}, nil
}

// NewSpecs implements spec.Factory.
func (_Factory) NewSpecs() (spec.Specs, error) {
	cmds := &Commands{}
	return &_Specs{cmds: cmds}, nil
}

type _Spec struct {
	cmd *Command
}

// Position implements spec.Spec.
func (spec *_Spec) Position(arg interface{}) (err error) {
	if arg, e := cmdUnpack(arg); e != nil {
		err = e
	} else {
		spec.cmd.Args = append(spec.cmd.Args, arg)
	}
	return
}

// Assign implements spec.Spec.
func (spec *_Spec) Assign(key string, arg interface{}) (err error) {
	if arg, e := cmdUnpack(arg); e != nil {
		err = e
	} else if spec.cmd.Keys != nil {
		spec.cmd.Keys[key] = arg
	} else {
		spec.cmd.Keys = map[string]interface{}{
			key: arg,
		}
	}
	return
}

type _Specs struct {
	cmds *Commands
}

// AddElement implements spec.Specs
func (specs *_Specs) AddElement(s spec.Spec) (err error) {
	if spec, ok := s.(*_Spec); !ok {
		err = errutil.Fmt("unexpected element type %T", s)
	} else {
		// println("appending element")
		specs.cmds.Els = append(specs.cmds.Els, spec.cmd)
	}
	return
}

// cmdUnpack helps position and assign push real data into the args
func cmdUnpack(value interface{}) (ret interface{}, err error) {
	switch b := value.(type) {
	case *_Spec:
		ret = b.cmd
	case *_Specs:
		ret = b.cmds
	case float64, string, int, []float64, []string:
		ret = value
	default:
		err = errutil.Fmt("assigning unexpected type %T", value)
	}
	return
}
