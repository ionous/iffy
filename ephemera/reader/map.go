package reader

// Map string keys to generic values with methods for extracting specific types.
type Map map[string]interface{}

// Box casts the passed value to "Map" because go-lang can be quite annoying.
// when dealing with interface{}, you cant explicitly cast b/t equivalent but inexact types.
func Box(i interface{}) Map {
	// and yet you can implicitly cast b/t those types, b/c reasons.
	return i.(map[string]interface{})
}

// StrOf the value at the passed key as a string.
func (m Map) StrOf(key string) (ret string) {
	if v, ok := m[key]; ok {
		ret = v.(string)
	}
	return
}

// MapOf the value at the passed key as a map.
func (m Map) MapOf(key string) (ret Map) {
	if v, ok := m[key]; ok {
		ret = Box(v)
	}
	return ret
}

// SliceOf the value at the passed key as a slice of interfaces.
func (m Map) SliceOf(key string) []interface{} {
	ret, _ := m[key].([]interface{})
	return ret
}

// Has true if the value at "key" equals "want".
func (m Map) Has(key, want string) (okay bool) {
	if have, ok := m[key]; ok && want == have {
		okay = true
	}
	return
}
