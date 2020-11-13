package iffy_test

import "testing"

// grouping is a key feature implemented in script.
// here, we want to test its underpinnings.
func TestGroups(t *testing.T) {
	// PrintNondescriptObjects lists a bunch of objects.

	// pattern Group Together maps objects into "partitioning records"
	// the partitions are sorted by some criteria
	// print group  (and print ungrouped ) patterns print the objects

	// note: even without grouping rules, our article printing still gathers unnamed items.
	// t.Run("no grouping", func(t *testing.T) {
	// 	groupTest(t,/*want*/ "Mildred, an empire apple, a pen, and two other things",
	// 		/* returns a list of objects that were declared in the Thingaverse */ sliceOf.String("mildred", "apple", "pen", "thing", "thing"),
	// 		PrintNameRules)
	// })
}

// var Thingaverse = ObjetctMap{
// 	// some unnamed things
// 	// this relies on the internal means of naming unnamed objects
// 	"thing#1": &Thing{},
// 	"thing#2": &Thing{},
// 	// a named thing
// 	"pen": &Thing{
// 		Kind: Kind{Name: "pen"},
// 	},
// 	// a thing with a printed name
// 	"apple": &Thing{
// 		Kind: Kind{Name: "apple", PrintedName: "empire apple"},
// 	},
// 	"box": &Container{
// 		Thing: Thing{Kind: Kind{Name: "box"}},
// 		Latch: Latch{Openable: true, Closed: true},
// 	},
// 	"cake": &Thing{
// 		Kind: Kind{Name: "cake"},
// 	},
// 	// someone with a proper name
// 	"mildred": &Actor{
// 		Thing{Kind: Kind{Name: "mildred", CommonProper: ProperNamed}},
// 	},
// 	"x": &ScrabbleTile{Thing{Kind: Kind{Name: "X"}}},
// 	"w": &ScrabbleTile{Thing{Kind: Kind{Name: "W"}}},
// 	"f": &ScrabbleTile{Thing{Kind: Kind{Name: "F"}}},
// 	"y": &ScrabbleTile{Thing{Kind: Kind{Name: "Y"}}},
// 	"z": &ScrabbleTile{Thing{Kind: Kind{Name: "Z"}}},
// }

// type ScrabbleTile struct {
// 	Thing `if:"parent"`
// }
