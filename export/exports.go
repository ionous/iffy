package export

import (
	r "reflect"
	"strings"
	"unicode"

	"github.com/ionous/inflect"
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
func Prettify(n string) (ret string) {
	if len(n) != 0 {
		ret = strings.ToLower(inflect.Humanize(n))
	}
	return
}
