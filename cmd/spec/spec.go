// Exports golang DSL for use in editing story files.
// Currently, this only generates the imperative commands,
// the modeling parts of the language currently live in the composer javascript
package main

import (
	"fmt"

	"github.com/ionous/iffy"
	"github.com/ionous/iffy/cmd/spec/internal"
)

// go run spec.go > /Users/ionous/Dev/go/src/github.com/ionous/iffy/cmd/compose/www/data/lang/spec.js
func main() {
	var c internal.Collect
	for _, slots := range iffy.AllSlots {
		for _, slot := range slots {
			c.AddSlot(slot)
		}
	}
	for _, slats := range iffy.AllSlats {
		for _, cmd := range slats {
			c.AddSlat(cmd)
		}
	}
	c.FlushGroups()
	c.Sort()
	if b, e := c.Marshal(); e != nil {
		panic(b)
	} else {
		fmt.Println("/* generated using github.com/ionous/iffy/cmd/spec/spec.go */")
		fmt.Print("const spec = ")
		fmt.Println(string(b))
	}
}
