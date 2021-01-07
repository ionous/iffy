package core

import (
	"strconv"

	"github.com/ionous/iffy/rt"
)

func NewActivity(exe ...rt.Execute) *Activity {
	return &Activity{Exe: exe}
}

func MakeActivity(exe ...rt.Execute) Activity {
	return Activity{Exe: exe}
}

func Args(from ...Assignment) *Arguments {
	var p Arguments
	for i, from := range from {
		p.Args = append(p.Args, &Argument{
			Name: "$" + strconv.Itoa(i+1),
			From: from,
		})
	}
	return &p
}

func NamedArgs(name string, from Assignment) *Arguments {
	return &Arguments{[]*Argument{{
		name, from,
	}}}
}
