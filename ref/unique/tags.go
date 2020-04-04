package unique

import (
	r "reflect"
	"strings"
)

type Metadata map[string]string

func Tag(f r.StructTag) StructTag {
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

func MergeMetadata(f *r.StructField, m *Metadata) {
	if s := Tag(f.Tag); len(s) > 0 {
		if len(*m) == 0 {
			*m = make(Metadata)
		}
		(*m).AddString(s, f.Name)
	}
}

// AddString parses the contents of a metadata string
// contents are expected to be comma separated
// each piece of content is separated further into key:value.
// Tags fields are of the format //`if:"key:value"`
// The tagcontents is the part inside the double quotes.
// The fill is used when the value part is empty or missing.
func (m Metadata) AddString(tagcontents StructTag, fill string) {
	parts := tagcontents.Split()
	for _, s := range parts {
		if len(s) > 0 {
			kv := strings.Split(s, ":")
			// no kv?
			if len(kv) == 1 {
				k, v := kv[0], fill
				m[k] = v
			} else {
				k, v := kv[0], strings.Join(kv[1:], ":")
				m[k] = v
			}
		}
	}
	return
}
