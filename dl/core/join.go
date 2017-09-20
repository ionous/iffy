package core

import (
	"bytes"
	"github.com/ionous/iffy/rt"
)

// Join is really a hack for the fact that Say takes only one string.
// Say only takes one string because the spec generator doesnt transparently handle arrays. FIX: it.
type Join struct {
	Text []rt.TextEval
}

func (p *Join) GetText(run rt.Runtime) (ret string, err error) {
	var buf bytes.Buffer
	for _, t := range p.Text {
		if s, e := t.GetText(run); e != nil {
			err = e
			break
		} else {
			buf.WriteString(s)
		}
	}
	if err == nil {
		ret = buf.String()
	}
	return
}
