package play

import (
	"github.com/ionous/iffy/event"
	"github.com/ionous/iffy/rt"
)

type ListenTo struct {
	Target  string
	Event   string
	Go      rt.ExecuteList
	Options []EventOptions
}

func (l *ListenTo) Define(f *Facts) (nil error) {
	f.Listeners = append(f.Listeners, *l)
	return
}

type EventOptions interface {
	Options() event.Options
}

type Capture struct{}
type RunAfter struct{}
type TargetOnly struct{}

func (*Capture) Options() event.Options {
	return event.Capture
}
func (*TargetOnly) Options() event.Options {
	return event.TargetOnly
}
func (*RunAfter) Options() event.Options {
	return event.RunAfter
}
