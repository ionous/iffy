package reader

import "github.com/ionous/errutil"

const (
	ItemId    = "id"
	ItemType  = "type"
	ItemValue = "value"
)

type ReadMap func(Map) error
type ReadMaps map[string]ReadMap

//
func Repeats(ms []interface{}, cb ReadMap) (err error) {
	for _, it := range ms {
		if e := cb(Box(it)); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

func At(m Map) string {
	return m.StrOf(ItemId)
}

func BadType(wanted, got, at string) error {
	return errutil.New("unexpected type", got, "wanted", wanted, "at", at)
}

func BadValue(t string, got interface{}, at string) error {
	return errutil.New(t, "has unexpected value", got, "at", at)
}

// helper: check the type of the passed m map
func Type(m Map, expectedType string) (ret string, err error) {
	if t := m.StrOf(ItemType); t != expectedType {
		err = BadType(expectedType, t, At(m))
	} else {
		ret = m.StrOf(ItemId)
	}
	return
}

// expect a map value
func Slat(m Map, expectedType string) (ret Map, err error) {
	if _, e := Type(m, expectedType); e != nil {
		err = e
	} else {
		ret = m.MapOf(ItemValue)
	}
	return
}

func Slot(r Map, expectedType string, slots ReadMaps) (err error) {
	if t := r.StrOf(ItemType); t != expectedType {
		err = BadType(expectedType, t, At(r))
	} else {
		m := r.MapOf(ItemValue)
		t := m.StrOf(ItemType)
		if fn, ok := slots[t]; !ok {
			err = BadType(expectedType, t, At(r))
		} else {
			err = fn(m)
		}
	}
	return
}

// we expect to see one, and only one, of the sub keys in the ItemValue of m.
func Option(r Map, expectedType string, slots ReadMaps) (err error) {
	if t := r.StrOf(ItemType); t != expectedType {
		err = BadType(expectedType, t, At(r))
	} else if m := r.MapOf(ItemValue); len(m) != 1 {
		err = BadValue(t, m, At(r))
	} else {
		// only one in the list.
		for key, value := range m {
			if fn, ok := slots[key]; !ok {
				err = BadValue(t, key, At(r))
			} else if e := fn(Box(value)); e != nil {
				err = e
			}
			break
		}
	}
	return
}

// expect a string variable
func String(m Map, expectedType string) (ret string, err error) {
	if t := m.StrOf(ItemType); t != expectedType {
		err = BadType(expectedType, t, At(m))
	} else if v := m.StrOf(ItemValue); len(v) == 0 {
		err = BadValue(t, v, At(m))
	} else {
		ret = v
	}
	return
}

// expect a string constant
func Const(m Map, expectedType, expectedValue string) (err error) {
	if t := m.StrOf(ItemType); t != expectedType {
		err = BadType(expectedType, t, At(m))
	} else if v := m.StrOf(ItemValue); v != expectedValue {
		err = BadValue(t, v, At(m))
	}
	return
}

//
func Enum(m Map, expectedType string, sub Map) (ret interface{}, err error) {
	if t := m.StrOf(ItemType); t != expectedType {
		err = BadType(expectedType, t, At(m))
	} else {
		n := m.StrOf(ItemValue)
		if i, ok := sub[n]; !ok {
			err = errutil.New("unexpected", expectedType, n)
		} else {
			ret = i
		}
	}
	return
}
