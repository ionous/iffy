package reflector

import (
	r "reflect"
	"strings"
)

type Metadata map[string]string

func MergeMetadata(f r.StructField, m *Metadata) {
	if s := f.Tag.Get("if"); len(s) > 0 {
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
func (m Metadata) AddString(tagcontents, fill string) {
	parts := strings.Split(tagcontents, ",")
	for _, s := range parts {
		if len(s) > 0 {
			kv := strings.Split(s, ":")
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
