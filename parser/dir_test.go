package parser_test

import (
	"github.com/ionous/sliceOf"
)

var Directions = func() (ret []*MyObject) {
	for _, d := range directions {
		obj := &MyObject{Id: d,
			Names:   sliceOf.String(d),
			Classes: sliceOf.String("directions"),
		}
		ret = append(ret, obj)
	}
	return
}()

// // var shortDirections = (
// // 	"n", "s", "e", "w", "ne", "nw", "se", "sw",
// // 	"u", "up", "ceiling", "above", "sky",
// // 	"d", "down", "floor", "below", "ground",
// // )
var directions = sliceOf.String(
	"north",
	"south",
	"east",
	"west",
	"northeast",
	"northwest",
	"southeast",
	"southwest",
	"up",   // "up above",
	"down", // "ground",
	"inside",
	"outside")
