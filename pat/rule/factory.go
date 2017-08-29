package rule

import (
	"github.com/ionous/iffy/pat"
	"sort"
)

type Rules struct {
	pat.Patterns
}

func MakeRules() Rules {
	return Rules{pat.MakePatterns()}
}

// Sort in-place so that lengthier filters are at the font of the list.
func (ps Rules) Sort() {
	for _, l := range ps.Bools {
		sort.SliceStable(l, func(i, j int) bool {
			return len(l[i].Filters) < len(l[j].Filters)
		})
	}
	for _, l := range ps.Numbers {
		sort.SliceStable(l, func(i, j int) bool {
			return len(l[i].Filters) < len(l[j].Filters)
		})
	}
	for _, l := range ps.Text {
		sort.SliceStable(l, func(i, j int) bool {
			return len(l[i].Filters) < len(l[j].Filters)
		})
	}
	for _, l := range ps.Objects {
		sort.SliceStable(l, func(i, j int) bool {
			return len(l[i].Filters) < len(l[j].Filters)
		})
	}
	for _, l := range ps.NumLists {
		sort.SliceStable(l, func(i, j int) bool {
			return len(l[i].Filters) < len(l[j].Filters)
		})
	}
	for _, l := range ps.TextLists {
		sort.SliceStable(l, func(i, j int) bool {
			return len(l[i].Filters) < len(l[j].Filters)
		})
	}
	for _, l := range ps.ObjLists {
		sort.SliceStable(l, func(i, j int) bool {
			return len(l[i].Filters) < len(l[j].Filters)
		})
	}
	for _, l := range ps.Executes {
		sort.SliceStable(l, func(i, j int) bool {
			return len(l[i].Filters) < len(l[j].Filters)
		})
	}
}
