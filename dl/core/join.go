package core

import (
	"bytes"
	"github.com/ionous/iffy/rt"
)

// Join combines multiple text into a buffer.
// It's similar to Say, except that say prints rather than returns text.
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
