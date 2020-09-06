package pattern

import (
	"strconv"

	"github.com/ionous/iffy/dl/core"
)

func NewParams(from ...core.Assignment) *Parameters {
	var p Parameters
	for i, from := range from {
		p.Params = append(p.Params, &Parameter{
			Name: "$" + strconv.Itoa(i+1),
			From: from,
		})
	}
	return &p
}

func NewNamedParams(name string, from core.Assignment) *Parameters {
	return &Parameters{[]*Parameter{{
		name, from,
	}}}
}
