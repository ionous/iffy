package iffy

import (
	"encoding/gob"

	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/ephemera/express"
	"github.com/ionous/iffy/pattern"
	"github.com/ionous/iffy/rt"
)

var AllSlots = [][]composer.Slot{rt.Slots, core.Slots}
var AllSlats = [][]composer.Slat{core.Slats, express.Slats}

func RegisterGobs() {
	registerGob()
}

// where should this live?
func init() {
	registerGob()
}

var registeredGob = false

func registerGob() {
	if !registeredGob {
		for _, slats := range AllSlats {
			for _, cmd := range slats {
				gob.Register(cmd)
			}
		}
		for _, rule := range pattern.Rules {
			gob.Register(rule)
		}
		registeredGob = true
	}
}
