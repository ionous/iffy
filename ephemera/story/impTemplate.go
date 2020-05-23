package story

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ephemera/express"
	"github.com/ionous/iffy/ephemera/imp"
	"github.com/ionous/iffy/ephemera/reader"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/template"
)

func imp_render_template(k *imp.Porter, r reader.Map) (ret interface{}, err error) {
	if m, e := reader.Unpack(r, "render_template"); e != nil {
		err = e
	} else if str, e := imp_lines(k, m.MapOf("$EXPRESSION")); e != nil {
		err = e
	} else if xs, e := template.Parse(str); e != nil {
		err = e
	} else if got, e := express.Convert(xs); e != nil {
		err = errutil.New(e, xs)
	} else if eval, ok := got.(rt.TextEval); !ok {
		err = errutil.Fmt("render template has unknown expression %T", got)
	} else {
		ret = &express.Render{eval}
	}
	return
}
