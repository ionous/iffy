package tag

import (
	r "reflect"
	"strings"
)

func ReadTag(f r.StructTag) StructTag {
	return StructTag(f.Get("if"))
}

type StructTag string

func (t StructTag) String() string {
	return string(t)
}

func (t StructTag) Split() []string {
	s := string(t)
	return strings.Split(s, ",")
}

func (t StructTag) Find(key string) (ret string, okay bool) {
	s := string(t)
	for {
		// find key in the string;
		// do this in a loop to handle errant substring matches: "iffy:" vs "if:"
		if i := strings.Index(s, key); i < 0 {
			break
		} else {
			// evaluate everything after
			s = s[i+len(key):]
			// at the end of string, or a separator then we have a key without value
			if len(s) == 0 || s[0] == ',' {
				okay = true
				break
			} else {
				// after the key is a value separator? extract the value
				if s[0] == ':' {
					if end := strings.Index(s, ","); end < 0 {
						ret = s[1:] // no new separator before the end, use everything after the :
					} else {
						ret = s[1:end] // found a separator before the end, split out our value
					}
					okay = true // either way, we are done.
					break
				}
			}
		}
	}
	return
}
