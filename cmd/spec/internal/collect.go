package internal

import (
	"encoding/json"
	"sort"
	"strings"

	r "reflect"

	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/export"
	"github.com/ionous/iffy/lang"
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
	name := composer.SlotName(slot)
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
		name := composer.SpecName(cmd)
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
		if spec.UsesStr() {
			uses = "str"
		}
		out := export.Dict{
			"name": name,
			"uses": uses,
		}
		// missing spec, missing slots.
		if len(spec.Spec) > 0 {
			out["spec"] = spec.Spec
		} else {
			var tokens []string
			var roles string
			params := export.Dict{}
			if spec.UsesStr() {
				if cs := spec.Strings; len(cs) == 2 && spec.OpenStrings {
					tokens, params = choices(cs, []string{"false", "true"})
				} else if len(cs) > 0 {
					tokens, params = choices(cs, cs)
				}
				if spec.OpenStrings {
					t := export.Tokenize(name)
					tokens = append([]string{t}, tokens...)
					params[t] = lang.Lowerspace(name)
				}
			} else if rtype.Kind() == r.Struct {
				tokens, roles, params = parseSpec(rtype, spec.Fluent)
				if len(roles) > 0 {
					with["roles"] = roles
				}
			} else {
				panic(name)
			}
			with["tokens"] = tokens
			with["params"] = params
		}
		if len(with) > 0 {
			out["with"] = with
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

func choices(vs, ts []string) (retTokens []string, retParams export.Dict) {
	retParams = make(export.Dict)
	for i, el := range ts {
		t := "$" + lang.UpperBreakcase(el)
		retTokens = append(retTokens, t)
		if i+1 < len(vs) {
			retTokens = append(retTokens, " or ")
		}
		retParams[t] = vs[i]
	}
	return
}
