package express

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/lang"
)

type nameCache struct {
	els map[string]interface{}
}

var coreCache nameCache

func (k *nameCache) get(n string) (ret interface{}, okay bool) {
	if len(k.els) == 0 {
		els := make(map[string]interface{})
		for _, v := range core.Slats {
			spec := v.Compose()
			n := lang.Camelize(spec.Name)
			els[n] = v
		}
		k.els = els
	}
	ret, okay = k.els[n]
	return
}
