package story

import (
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/ephemera/imp"
	"github.com/ionous/iffy/ephemera/reader"
)

func imp_test_statement(k *imp.Porter, r reader.Map) (err error) {
	if n, e := imp_test_name(k, r.MapOf("name")); e != nil {
		err = e
	} else {
		err = reader.Slot(r, "testing", reader.ReadMaps{
			"test_output": func(m reader.Map) error {
				return imp_test_output(k, n, m)
			},
		})
	}
	return
}

func imp_test_output(k *imp.Porter, test ephemera.Named, r reader.Map) (err error) {
	if m, e := reader.Slat(r, "test_output"); e != nil {
		err = e
	} else if expect, e := imp_lines(k, m.MapOf("$LINES")); e != nil {
		err = e
	} else if prog, e := k.DecodeAny(r); e != nil {
		err = e
	} else if p, e := k.NewProg("test", prog); e != nil {
		err = e
	} else {
		k.NewTest(test, p, expect)
	}
	return
}

func imp_lines(k *imp.Porter, r reader.Map) (ret string, err error) {
	return reader.String(r, "lines")
}
