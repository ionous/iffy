package story

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/render"
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/ephemera/express"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/template"
	"github.com/ionous/iffy/template/types"
)

func (op *RenderTemplate) ImportStub(k *Importer) (ret interface{}, err error) {
	if xs, e := template.Parse(op.Template.Str); e != nil {
		err = e
	} else if got, e := express.Convert(xs); e != nil {
		err = errutil.New(e, xs)
	} else if eval, ok := got.(rt.TextEval); !ok {
		err = errutil.Fmt("render template has unknown expression %T", got)
	} else {
		ret = &render.RenderTemplate{eval}
		// pretty.Println(eval)
	}
	return
}

func convert_text_or_template(str string) (ret interface{}, err error) {
	if xs, e := template.Parse(str); e != nil {
		err = e
	} else if str, ok := getSimpleString(xs); ok {
		ret = str // okay; return the string.
	} else {
		if got, e := express.Convert(xs); e != nil {
			err = errutil.New(e, xs)
		} else if eval, ok := got.(rt.TextEval); !ok {
			err = errutil.Fmt("render template has unknown expression %T", got)
		} else if prog, e := ephemera.EncodeGob(&render.RenderTemplate{eval}); e != nil {
			err = e // note: we dont have to encode into render but maybe its nice to have a consistent root type
		} else {
			ret = prog // okay; return bytes.
		}
	}
	return
}

// see if the parsed expression contained anything other than text
// if true, return that text
func getSimpleString(xs template.Expression) (ret string, okay bool) {
	switch len(xs) {
	case 0:
		okay = true
	case 1:
		if quote, ok := xs[0].(types.Quote); ok {
			ret, okay = quote.Value(), true
		}
	}
	return
}
