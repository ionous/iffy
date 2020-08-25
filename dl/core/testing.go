package core

import (
	"github.com/ionous/iffy/rt"
)

func NewActivity(exe ...rt.Execute) *Activity {
	return &Activity{Exe: exe}
}
