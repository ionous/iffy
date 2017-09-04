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

type ListenFor struct {
	ListenTo
}

func (l *ListenTo) Define(f *Facts) (nil error) {
	f.ObjectListeners = append(f.ObjectListeners, *l)
	return
}

func (l *ListenFor) Define(f *Facts) (nil error) {
	f.ClassListeners = append(f.ClassListeners, *l)
	return
}

func (l *ListenTo) GetOptions() (ret event.Options) {
	for _, o := range l.Options {
		ret |= o.Options()
	}
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
