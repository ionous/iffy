package story

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/check"
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/ephemera/imp"
	"github.com/ionous/iffy/ephemera/reader"
)

func imp_test_statement(k *imp.Porter, r reader.Map) (err error) {
	if m, e := reader.Unpack(r, "test_statement"); e != nil {
		err = e
	} else if n, e := imp_test_name(k, m.MapOf("$NAME")); e != nil {
		err = e
	} else {
		err = reader.Slot(m.MapOf("$TEST"), "testing", reader.ReadMaps{
			"test_output": func(m reader.Map) error {
				return imp_test_output(k, n, m)
			},
		})
	}
	return
}

func imp_test_output(k *imp.Porter, test ephemera.Named, r reader.Map) (err error) {
	if m, e := reader.Unpack(r, "test_output"); e != nil {
		err = e
	} else if expect, e := imp_lines(k, m.MapOf("$LINES")); e != nil {
		err = e
	} else if i, e := k.DecodeAny(r); e != nil {
		err = e
	} else if to, ok := i.(*check.TestOutput); !ok {
		err = errutil.Fmt("couldnt decode test output command from %T", to)
	} else {
		// from the gob example:
		// `Pass pointer to interface so Encode sees (and hence sends) a value of
		// interface type. If we passed p directly it would see the concrete type instead.
		// See the blog post, "The Laws of Reflection" for background.`
		var t check.Testing = to
		if p, e := k.NewProg("test_output", &t); e != nil {
			err = e
		} else {
			k.NewTest(test, p, expect)
		}
	}
	return
}

func imp_lines(k *imp.Porter, r reader.Map) (ret string, err error) {
	return reader.String(r, "lines")
}
