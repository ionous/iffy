package text

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/template"
)

// DirectiveState tracks the structure of the template keywords: if blocks, else blocks, etc.
type DirectiveState interface {
	next(template.Directive) (DirectiveState, error)
	pop() (DirectiveState, error)
}

// PrevState to return to an earlier DirectiveState.
type PrevState struct {
	prev DirectiveState
}

// Restore a previous state; errors if there was no previous state.
func (t PrevState) pop() (ret DirectiveState, err error) {
	if t.prev == nil {
		err = errutil.New("too many ends!")
	} else {
		ret = t.prev
	}
	return
}
