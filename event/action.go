package event

import (
	"github.com/ionous/iffy/rt"
)

type ActionMap map[string]*Action

type Action struct {
	Id             string
	Name           string
	TargetClass    rt.Class
	DataClass      rt.Class
	DefaultActions rt.ExecuteList
}
