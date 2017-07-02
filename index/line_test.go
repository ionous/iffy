package index

import "strings"

// String only exists in test builds because i dont want the dependency on strings
func (l *Line) String() string {
	return strings.Join(l[:], ", ")
}
