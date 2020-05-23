package internal

import (
	"sort"
	"strings"

	"github.com/ionous/iffy/export"
)

type Groups map[string]bool

func (g *Groups) appendGroups(all []export.Dict) []export.Dict {
	for groupName, _ := range *g {
		all = append(all, export.Dict{
			"name": groupName,
			"uses": "group",
		})
	}
	return all
}

func (g *Groups) addGroup(out export.Dict, group string) {
	// even no commas results in one group;
	// ideally, id think an empty string would be no groups.... but alas.
	if len(group) > 0 {
		if groups := strings.Split(group, ","); len(groups) > 0 {
			for i, group := range groups {
				(*g)[group] = true
				groups[i] = strings.ToLower(group)
			}
			sort.Strings(groups)
			out["group"] = groups
		}
	}
}
