package spec

import (
	"github.com/ionous/errutil"
)

// Command represents an arbitrary game statement or operation.
type Command struct {
	Name string
	Args []interface{}
	Keys map[string]interface{}
}

type Commands struct {
	Els []*Command
}

// CommandFactory implements SpecFactory.
type CommandFactory struct{}

// CommandBuilder implements Spec.
type CommandBuilder struct {
	CommandFactory
	cmd *Command
}

// CommandsBuilder implements Specs.
type CommandsBuilder struct {
	CommandFactory
	cmds *Commands
}

func NewCommands() *CommandBuilder {
	return &CommandBuilder{cmd: &Command{Name: "root"}}
}

// NewSpec implements SpecFactory.
func (CommandFactory) NewSpec(name string) (Spec, error) {
	return &CommandBuilder{cmd: &Command{Name: name}}, nil
}

// NewSpecs implements SpecFactory.
func (CommandFactory) NewSpecs() (Specs, error) {
	cmds := &Commands{}
	return &CommandsBuilder{cmds: cmds}, nil
}

func (cb *CommandBuilder) Root() *Command {
	return cb.cmd
}

// Position implements Spec.
func (cb *CommandBuilder) Position(arg interface{}) (err error) {
	if arg, e := cmdUnpack(arg); e != nil {
		err = e
	} else {
		cb.cmd.Args = append(cb.cmd.Args, arg)
	}
	return
}

func (cb *CommandBuilder) Assign(key string, arg interface{}) (err error) {
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

func (cbs *CommandsBuilder) Commands() []*Command {
	return cbs.cmds.Els
}

func (cbs *CommandsBuilder) AddElement(s Spec) (err error) {
	if cb, ok := s.(*CommandBuilder); !ok {
		err = errutil.Fmt("unexpected element type %T", s)
	} else {
		cbs.cmds.Els = append(cbs.cmds.Els, cb.cmd)
	}
	return
}

// this is for position and assign, so that we can get the real data into the args
func cmdUnpack(value interface{}) (ret interface{}, err error) {
	switch b := value.(type) {
	case *CommandBuilder:
		ret = b.cmd
	case *CommandsBuilder:
		ret = b.cmds
	case float64, string, int, []float64, []string:
		ret = value
	default:
		err = errutil.Fmt("assigning unexpected type %T", value)
	}
	return
}
