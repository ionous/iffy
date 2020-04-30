package export

import (
	r "reflect"
	"strings"
	"unicode"

	"bitbucket.org/pkg/inflect"
)

type Dict map[string]interface{}

// Tokenize turns the passed struct field into an composer compatible parameter name.
// ex. FieldName -> $FIELD_NAME
func Tokenize(f *r.StructField) string {
	n := Prettify(f.Name)
	return "$" + strings.Map(func(c rune) (ret rune) {
		if c == ' ' {
			ret = '_'
		} else {
			ret = unicode.ToUpper(c)
		}
		return
	}, n)
}

// Prettify transforms a PascalCased named into lowercase names with spaces.
// ex. "FieldName" into "field name"
func Prettify(n string) string {
	return strings.ToLower(inflect.Humanize(n))
}
