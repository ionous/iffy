package reader

type Map map[string]interface{}

// go can be annoying.
// when dealing with interface{}, you cant explicitly cast b/t equivalent but inexact types.
func Box(i interface{}) Map {
	// and yet you can implicitly cast b/t those types, b/c reasons.
	return i.(map[string]interface{})
}

func Unbox(m Map) map[string]interface{} {
	return m
}

func (m Map) StrOf(key string) (ret string) {
	if v, ok := m[key]; ok {
		ret = v.(string)
	}
	return
}

func (m Map) MapOf(key string) (ret Map) {
	if v, ok := m[key]; ok {
		ret = Box(v)
	}
	return ret
}
func (m Map) SliceOf(key string) []interface{} {
	ret, _ := m[key].([]interface{})
	return ret
}

func (m Map) Expect(key, want string) (okay bool) {
	if have, ok := m[key]; ok && want == have {
		okay = true
	}
	return
}
