package pattern

import (
	"strconv"

	"github.com/ionous/iffy/dl/core"
)

func NewArgs(from ...core.Assignment) *Arguments {
	var p Arguments
	for i, from := range from {
		p.Args = append(p.Args, &Argument{
			Name: "$" + strconv.Itoa(i+1),
			From: from,
		})
	}
	return &p
}

func NewNamedParams(name string, from core.Assignment) *Arguments {
	return &Arguments{[]*Argument{{
		name, from,
	}}}
}
