package std

import (
	"sort"
)

var objectMap = map[string]interface{}{
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
}

var objectList = func(src map[string]interface{}) (ret []interface{}) {
	for _, v := range src {
		ret = append(ret, v)
	}
	return
}(objectMap)

var nameList = func(src map[string]interface{}) (ret []string) {
	for n, _ := range src {
		ret = append(ret, n)
	}
	sort.Strings(ret)
	return
}(objectMap)
