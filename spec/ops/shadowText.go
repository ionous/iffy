package ops

import (
	"fmt"
	"io"
)

func (c *ShadowClass) Format(f fmt.State, r rune) {
	io.WriteString(f, fmt.Sprintf("Make%s{", c.Type().Name()))
	for _, n := range c.fields {
		rf := c.rtype.FieldByIndex(n)
		if v := c.getField(&rf, false); v.IsValid() {
			io.WriteString(f, fmt.Sprintf("%s:%#v", rf.Name, v.Interface()))
		}
	}
	io.WriteString(f, "}")
}
