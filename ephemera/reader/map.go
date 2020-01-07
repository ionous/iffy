package reader

import "log"

type Map map[string]interface{}

// go can be annoying sometimes
// you cant directly Cast b/t equivalent but inexact types
func Cast(i interface{}) Map {
	return i.(map[string]interface{})
}

func (m Map) StrOf(key string) (ret string) {
	if v, ok := m[key]; ok {
		ret = v.(string)
	}
	return
}

func (m Map) MapOf(key string) (ret Map) {
	if v, ok := m[key]; ok {
		ret = Cast(v)
	}
	return ret
}
func (m Map) SliceOf(key string) []interface{} {
	ret, _ := m[key].([]interface{})
	return ret
}

func (m Map) Expect(key, want string) {
	if have, ok := m[key]; !ok || want != have {
		log.Fatalln("wanted", want, "have", have)
	}
}
