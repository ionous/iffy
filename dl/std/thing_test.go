package std

import (
	"sort"
)

var thingMap = map[string]*Thing{
	// some unnamed things
	// this relies on the internal means of naming unnamed objects
	"thing#1": &Thing{},
	"thing#2": &Thing{},
	// a named thing
	"pen": &Thing{
		Kind: Kind{Name: "pen"},
	},
	// a thing with a printed name
	"sword": &Thing{
		Kind: Kind{Name: "sword", PrintedName: "plastic sword"},
	},
}

var thingList = func(src map[string]*Thing) (ret []interface{}) {
	for _, v := range src {
		ret = append(ret, v)
	}
	return
}(thingMap)

var nameList = func(src map[string]*Thing) (ret []string) {
	for n, _ := range src {
		ret = append(ret, n)
	}
	sort.Strings(ret)
	return
}(thingMap)
