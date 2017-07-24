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
	// someone with a proper name
	"mildred": &Actor{
		Thing{Kind: Kind{Name: "mildred", CommonProper: ProperNamed}},
	},
	"x": &ScrabbleTile{Thing{Kind: Kind{Name: "x"}}},
	"w": &ScrabbleTile{Thing{Kind: Kind{Name: "w"}}},
	"f": &ScrabbleTile{Thing{Kind: Kind{Name: "f"}}},
	"y": &ScrabbleTile{Thing{Kind: Kind{Name: "y"}}},
	"z": &ScrabbleTile{Thing{Kind: Kind{Name: "z"}}},
}

type ScrabbleTile struct {
	Thing
}
