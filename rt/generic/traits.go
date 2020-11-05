package generic

import (
	"sort"

	"github.com/ionous/iffy/affine"
)

// pairs of aspect and traits so we can move from one to the other.
type trait struct {
	Trait  string
	Aspect string
}

func makeTraits(aspect *Kind, traits []trait) []trait {
	for i, cnt := 0, aspect.NumField(); i < cnt; i++ {
		// fix: check if Type() is set to trait?
		if ft := aspect.Field(i); ft.Affinity != affine.Bool || ft.Type != "trait" {
			panic(aspect.name + "aspect has non trait fields")
		} else {
			traits = append(traits, trait{Trait: ft.Name, Aspect: aspect.name})
		}
	}
	return traits
}

func sortTraits(ts []trait) {
	sort.Slice(ts, func(i, j int) bool {
		it, jt := ts[i], ts[j]
		return it.Trait < jt.Trait
	})
}

// find aspect from trait name in a sorted list of traits
func findAspect(trait string, ts []trait) (ret string) {
	cnt := len(ts)
	if i := sort.Search(cnt, func(i int) bool {
		return ts[i].Trait >= trait
	}); i < cnt {
		// note: search returns the insertion point, not the found element.
		if found := ts[i]; found.Trait == trait {
			ret = found.Aspect
		}
	}
	return
}
