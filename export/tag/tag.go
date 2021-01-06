package tag

import (
	r "reflect"
	"strings"
)

// ReadTag parses tags of the format `if:"one,two,key=value"`.
func ReadTag(f r.StructTag) StructTag {
	return StructTag(f.Get("if"))
}

// StructTag in go are, by convention, key:"value" pairs;
// this extends that for sub-tags within the value part of the pair.
type StructTag string

// String returns the entire value of the golang key:"value" pair.
func (t StructTag) String() string {
	return string(t)
}

// Split returns the comma separate subparts of the tag's value.
func (t StructTag) Split() []string {
	s := string(t)
	return strings.Split(s, ",")
}

// Exists returns true if Find returns true.
func (t StructTag) Exists(key string) (okay bool) {
	_, okay = t.Find(key)
	return
}

// Find finds the named key within the struct tag.
// ex. for the tag, `if:"internal"` Find("internal") returns "internal".
// the tag, `if:"beep=boop"` Find("beep") return "boop".
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
				if s[0] == '=' {
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
