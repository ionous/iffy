package trigger

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/event"
	p "github.com/ionous/iffy/parser"
	"github.com/ionous/iffy/rt"
)

type Commands struct {
	*Trigger
}

// Trigger mimics parser.ResolvedTrigger, terminating a matcher sequence, resolving to the named action.
type Trigger struct {
	Raise rt.ObjectEval
}

// Trigger mimics parser.ResolvedTrigger, terminating a matcher sequence, resolving to the named action.
type ResolvedTrigger struct {
	raise  rt.ObjectEval
	events event.EventMap
}

// Scan matches only if the cursor has finished with all words.
func (a *Trigger) Scan(ctx p.Context, scope p.Scope, cs p.Cursor) (ret p.Result, err error) {
	if src, ok := ctx.(Context); !ok {
		err = errutil.Fmt("unknown context %T", ctx)
	} else {
		if w := cs.CurrentWord(); len(w) == 0 {
			ret = ResolvedTrigger{a.Raise, src.events}
		} else {
			err = p.Overflow{p.Depth(cs.Pos)}
		}
	}
	return
}

// WordsMatched always returns 0 -- we expect to be a terminal node in a grammar, the result of parsing through all words and not a word itself.
func (a ResolvedTrigger) WordsMatched() int {
	return 0
}

// Execute generates an event, triggering handlers and default action.
func (a ResolvedTrigger) Execute(run rt.Runtime) (err error) {
	if obj, e := a.raise.GetObject(run); e != nil {
		err = e
	} else {
		err = event.Trigger(run, a.events, obj)
	}
	return
}
