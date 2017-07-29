package index

import (
	"github.com/ionous/sliceOf"

	"strings"
)

// String only exists in test builds because i dont want the dependency on strings
func (r Row) String() string {
	return strings.Join(sliceOf.String(r.Major, r.Minor), ", ")
}
