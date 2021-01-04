package internal

import (
	"encoding/json"

	"github.com/ionous/errutil"
)

type Patch []PatchCommand

type PatchCommand struct {
	Name      string `json:"patch"`
	Migration `json:"migration"`
}

func (p Patch) Migrate(doc interface{}) (err error) {
	for _, op := range p {
		if e := op.Migrate(doc); e != nil {
			err = e
			break
		}
	}
	return
}

func (c *PatchCommand) UnmarshalJSON(data []byte) (err error) {
	var rep struct {
		Name      string          `json:"patch"`
		Migration json.RawMessage `json:"migration"`
	}
	if e := json.Unmarshal(data, &rep); e != nil {
		err = e
	} else {
		c.Name = rep.Name
		switch n, m := rep.Name, rep.Migration; n {
		case "replace":
			err = c.unmarshal(m, &Replace{})
		case "copy":
			err = c.unmarshal(m, &Copy{})
		default:
			err = errutil.New("unknown migration", n)
		}
	}
	return
}

func (c *PatchCommand) unmarshal(msg json.RawMessage, op Migration) (err error) {
	if e := json.Unmarshal(msg, op); e != nil {
		err = e
	} else {
		c.Migration = op
	}
	return
}
