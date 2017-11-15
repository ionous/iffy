package core

import (
	"github.com/ionous/iffy/rt"
	"io"
	//	"github.com/ionous/iffy/rt/printer"
)

// Comprise writes a prefix and suffix around a body of text so long as the body has content.
type Comprise struct {
	Prefix, Body, Suffix rt.TextEval
}

func (p *Comprise) Execute(run rt.Runtime) (err error) {
	if body, e := p.Body.GetText(run); e != nil {
		err = e
	} else if len(body) > 0 {
		if prefix, e := p.Prefix.GetText(run); e != nil {
			err = e
		} else if suffix, e := p.Suffix.GetText(run); e != nil {
			err = e
		} else {
			w := run.Writer()
			for _, v := range []string{prefix, body, suffix} {
				if _, e := io.WriteString(w, v); e != nil {
					err = e
					break
				}
			}
		}
	}
	return
}
