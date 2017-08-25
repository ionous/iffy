package event

import (
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/scope"
)

type Frame struct {
	run    rt.Runtime
	evt    *EventObject
	queue  QueuedActions
	pushed bool
}

func NewFrame(run rt.Runtime, evt *EventObject) (ret *Frame, err error) {
	// create event object
	// FIX: maybe this could be merged with scope?
	if temp, e := run.Emplace(evt); e != nil {
		err = e
	} else {
		run.PushScope(scope.MultiFinder(
			scope.AtFinder(temp),
			scope.ModelFinder(run),
		))
		ret = &Frame{run: run, evt: evt, pushed: true}
	}
	return
}

// Destroy should happen once for every NewFrame
func (ac *Frame) Destroy() {
	if ac.pushed {
		ac.run.PopScope()
		ac.pushed = false
	}
}

// Dispatch the event frame to the passed targets.
func (ac *Frame) DispatchFrame(at, path []Target) (err error) {
	fullPath := [2][]Target{at, path}
	if e := ac.dispatchPhase(CaptureListeners, fullPath, func(i, cnt int) int {
		return cnt - i - 1
	}); e != nil {
		err = e
	} else if !ac.evt.Stopped() && ac.evt.Bubbles {
		if e := ac.dispatchPhase(BubbleListeners, fullPath, func(i, cnt int) int {
			return i
		}); e != nil {
			err = e
		}
	}
	return
}

func (ac *Frame) dispatchPhase(listenerType ListenerType, fullPath [2][]Target, order func(i, cnt int) int) (err error) {
	if evt := ac.evt; !evt.Stopped() {
	OutOfLoop:
		// walk the lists of parent and target handlers
		for i, cnt := 0, len(fullPath); i < cnt && !evt.Stopped(); i++ {
			// allow forward or backward iteration
			which := order(i, cnt)
			path := fullPath[which]
			// the 0th list is always the list of target handlers
			atTarget := which == 0
			if atTarget {
				evt.Phase = AtTarget
			} else if listenerType == CaptureListeners {
				evt.Phase = CapturingPhase
			} else {
				evt.Phase = BubblingPhase
			}
			// walk the list of list of parent or target handlers
			for i, cnt := 0, len(path); i < cnt; i++ {
				// allow forward or backward iteration
				tgt := path[order(i, cnt)]
				// get the handlers just for this phase: bubble/capture.
				hs := tgt.handlers[listenerType]
				for i, cnt := 0, len(hs); i < cnt; i++ {
					// allow forward or backward iteration through the phase
					// FIX: shouldnt the order of handlers always be last registered first?
					h := hs[order(i, cnt)]
					// at target usually this only happens during the bubbles phase.
					if atTarget == h.IsTargetOnly() {
						if h.IsRunAfter() {
							ac.queue = ac.queue.Add(tgt.obj, h.Exec)
						} else {
							// FIX: set hint via hint pointer into scope
							evt.CurrentTarget = tgt.obj
							if e := h.Exec.Execute(ac.run); e != nil {
								err = e
								break OutOfLoop
							} else if evt.StopImmediatePropagation {
								break OutOfLoop
							}
						}
					}
				}
			}
		}
	}
	return
}
