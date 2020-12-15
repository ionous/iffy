package story

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/rt"
)

type Activity core.Activity
type Assignment core.Assignment
type BoolEval rt.BoolEval
type Execute rt.Execute
type NumberEval rt.NumberEval
type ObjectEval rt.ObjectEval
type TextEval rt.TextEval

// fix: this doesnt work because story importer doesnt trigger callbacks for str types
func (op *Text) ImportStub(k *Importer) (ret interface{}, err error) {
	var text string
	if t := op.Str; t != "$EMPTY" {
		text = t
	}
	ret = &core.Text{text}
	return
}

// handle the import of text literals, this is a patch for handling "empty" in string values.
func (op *TextValue) ImportStub(k *Importer) (ret interface{}, err error) {
	return op.Text.ImportStub(k)
}
