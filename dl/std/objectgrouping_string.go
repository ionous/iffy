// Code generated by "stringer -type=ObjectGrouping"; DO NOT EDIT.

package std

import "fmt"

const _ObjectGrouping_name = "GroupWithoutArticlesGroupWithArticlesGroupWithoutObjects"

var _ObjectGrouping_index = [...]uint8{0, 20, 37, 56}

func (i ObjectGrouping) String() string {
	if i < 0 || i >= ObjectGrouping(len(_ObjectGrouping_index)-1) {
		return fmt.Sprintf("ObjectGrouping(%d)", i)
	}
	return _ObjectGrouping_name[_ObjectGrouping_index[i]:_ObjectGrouping_index[i+1]]
}
