// Code generated by "stringer -type=ListenerType"; DO NOT EDIT.

package event

import "fmt"

const _ListenerType_name = "CaptureListenersBubbleListenersListenerTypes"

var _ListenerType_index = [...]uint8{0, 16, 31, 44}

func (i ListenerType) String() string {
	if i < 0 || i >= ListenerType(len(_ListenerType_index)-1) {
		return fmt.Sprintf("ListenerType(%d)", i)
	}
	return _ListenerType_name[_ListenerType_index[i]:_ListenerType_index[i+1]]
}
