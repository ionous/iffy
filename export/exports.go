package export

import (
	r "reflect"
	"strings"
	"unicode"

	"bitbucket.org/pkg/inflect"
)

type Dict map[string]interface{}

func Tokenize(f *r.StructField) string {
	return "$" + strings.Map(func(c rune) (ret rune) {
		if c == ' ' {
			ret = '_'
		} else {
			ret = unicode.ToUpper(c)
		}
		return
	}, Prettify(f.Name))
}

func Prettify(n string) string {
	return strings.ToLower(inflect.Humanize(n))
}
