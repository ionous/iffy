package text

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/spec/ops"
	"github.com/ionous/iffy/template"
)

// DirectiveState tracks the structure of the template keywords: if blocks, else blocks, etc.
type DirectiveState interface {
	next(template.Directive) (DirectiveState, error)
}

type PrevStates struct {
	list []DirectiveState
}

type Depth int

type Commands struct {
	list []*ops.Command
}

func (t *PrevStates) push(q DirectiveState) {
	t.list = append(t.list, q)
}

// Restore a previous state; errors if there was no previous state.
func (t *PrevStates) pop() (ret DirectiveState, err error) {
	if cnt := len(t.list); cnt == 0 {
		err = errutil.New("too many ends!")
	} else {
		ret, t.list = t.list[cnt-1], t.list[0:cnt-1]
	}
	return
}

func (d Depth) rollup(eng *Engine) (ret DirectiveState, err error) {
	if prev, e := eng.prev.pop(); e != nil {
		err = e
	} else {
		eng.cmds.end() // end span
		for i := 0; i < int(d); i++ {
			eng.cmds.end()
		}
		ret = prev
	}
	return
}

func (cs *Commands) end() {
	cs.list = cs.list[:len(cs.list)-1]
}

func (cs *Commands) begin(cmd *ops.Command) (err error) {
	if e := cs.position(cmd); e != nil {
		err = e
	} else {
		cs.list = append(cs.list, cmd)
	}
	return
}

func (cs *Commands) position(cmd *ops.Command) (err error) {
	if cnt := len(cs.list); cnt == 0 {
		err = errutil.New("list underflow")
	} else {
		top := cs.list[cnt-1]
		if e := top.Position(cmd); e != nil {
			err = e
		}
	}
	return
}
