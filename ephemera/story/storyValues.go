package story

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ephemera/decode"
)

func (op *Certainty) ImportString(k *Importer) (ret string, err error) {
	if str := op.Str; decode.IndexOfChoice(op, op.Str) < 0 {
		err = ImportError(op, op.At, errutil.Fmt("%w %q", InvalidValue, str))
	} else {
		ret = str
	}
	return
}

// blocks of text might well be a template.
func (op *Lines) ConvertText() (ret interface{}, err error) {
	return convert_text_or_template(op.Str)
}
