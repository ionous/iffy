package std

type ObjetctMap map[string]interface{}

func (om ObjetctMap) objects(names []string) (ret []interface{}) {
	for _, n := range names {
		ret = append(ret, om[n])
	}
	return
}

var Thingaverse = ObjetctMap{
	// some unnamed things
	// this relies on the internal means of naming unnamed objects
	"thing#1": &Thing{},
	"thing#2": &Thing{},
	// a named thing
	"pen": &Thing{
		Kind: Kind{Name: "pen"},
	},
	// a thing with a printed name
	"apple": &Thing{
		Kind: Kind{Name: "apple", PrintedName: "empire apple"},
	},
	"box": &Container{
		Thing: Thing{Kind: Kind{Name: "box"}},
		Latch: Latch{Openable: true, Closed: true},
	},
	"cake": &Thing{
		Kind: Kind{Name: "cake"},
	},
	// someone with a proper name
	"mildred": &Actor{
		Thing{Kind: Kind{Name: "mildred", CommonProper: ProperNamed}},
	},
	"x": &ScrabbleTile{Thing{Kind: Kind{Name: "X"}}},
	"w": &ScrabbleTile{Thing{Kind: Kind{Name: "W"}}},
	"f": &ScrabbleTile{Thing{Kind: Kind{Name: "F"}}},
	"y": &ScrabbleTile{Thing{Kind: Kind{Name: "Y"}}},
	"z": &ScrabbleTile{Thing{Kind: Kind{Name: "Z"}}},
}

type ScrabbleTile struct {
	Thing `if:"parent"`
}
