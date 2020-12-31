package internal

import (
	"encoding/json"
	"sort"
	"strings"

	r "reflect"

	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/export"
)

type Collect struct {
	all    []export.Dict
	stubs  []string
	slots  []r.Type
	groups Groups
}

func (c *Collect) AddGroup(out export.Dict, group string) {
	if c.groups == nil {
		c.groups = make(Groups)
	}
	c.groups.addGroup(out, group)
}

func (c *Collect) AddSlot(slot composer.Slot) {
	i := r.TypeOf(slot.Type).Elem()
	//
	name := typeName(i, slot.Name)
	var desc string
	if len(slot.Desc) > 0 {
		desc = slot.Desc
	} else {
		desc = export.Prettify(name)
	}
	out := export.Dict{
		"name": name,
		"desc": desc,
		"uses": "slot",
	}
	addDesc(out, name, slot.Desc)
	//c.AddGroup(out, spec.Group)
	c.all = append(c.all, out)
	c.slots = append(c.slots, i)
}

func (c *Collect) AddSlat(cmd composer.Composer) {
	if spec := cmd.Compose(); spec.Group != "internal" {
		rtype := r.TypeOf(cmd).Elem()
		name := typeName(rtype, spec.Name)
		var header string
		if spec.Fluent != nil {
			header = rtype.Name()
		} else {
			header = export.Prettify(name)
		}
		//
		with := make(export.Dict)
		if slotNames := slotsOf(rtype, c.slots); len(slotNames) > 0 {
			with["slots"] = slotNames
		}
		uses := "flow"
		if len(spec.Uses) > 0 {
			uses = spec.Uses
		}
		out := export.Dict{
			"name": name,
			"uses": uses,
			"with": with,
		}
		// missing spec, missing slots.
		if len(spec.Spec) != 0 {
			out["spec"] = spec.Spec
		} else {
			tokens, roles, params := parseSpec(rtype, spec.Fluent)
			with["tokens"] = tokens
			with["params"] = params
			if len(roles) > 0 {
				with["roles"] = roles
			}
		}
		if spec.Stub {
			c.stubs = append(c.stubs, name)
		}
		// if Desc doesnt have a colon, should add the name, uppercase if not fluent maybe.
		addDesc(out, header, spec.Desc)
		c.AddGroup(out, spec.Group)
		c.all = append(c.all, out)
	}
}

func (c *Collect) FlushGroups() {
	c.all = c.groups.appendGroups(c.all)
}

func (c *Collect) Sort() {
	sort.Slice(c.all, func(idx, jdx int) (ret bool) {
		i, j := c.all[idx], c.all[jdx]
		uses := strings.Compare(i["uses"].(string), j["uses"].(string))
		switch uses {
		case 0:
			ret = i["name"].(string) < j["name"].(string)
		case -1:
			ret = false
		case 1:
			ret = true
		}
		return

	})
}

func (c *Collect) Marshal() (ret []byte, err error) {
	return json.MarshalIndent(c.all, "", "  ")
}

func (c *Collect) Stubs() (ret []byte) {
	if len(c.stubs) == 0 {
		ret = []byte(`""`)
	} else if b, e := json.MarshalIndent(c.stubs, "", "  "); e != nil {
		ret = []byte("// " + e.Error())
	} else {
		ret = b
	}
	return
}
