package story

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/debug"
	"github.com/ionous/iffy/dl/list"
	"github.com/ionous/iffy/dl/rel"
	"github.com/ionous/iffy/ephemera/decode"
	"github.com/ionous/iffy/rt"
)

type Activity core.Activity
type Assignment core.Assignment
type BoolEval rt.BoolEval
type Execute rt.Execute
type NumberEval rt.NumberEval
type TextEval rt.TextEval

type VariableName struct {
	core.Variable
}

type RelationName struct {
	rel.Relation
}

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

// handle the import of boolean flags
func (op *ListEdge) ImportStub(k *Importer) (ret interface{}, err error) {
	ret = list.Edge(op.Str == "$TRUE")
	return
}

// handle the import of boolean flags
func (op *ListOrder) ImportStub(k *Importer) (ret interface{}, err error) {
	ret = list.Order(op.Str == "$TRUE")
	return
}

// handle the import of boolean flags
func (op *ListCase) ImportStub(k *Importer) (ret interface{}, err error) {
	ret = list.Case(op.Str == "$TRUE")
	return
}

// handle the import of int flags
func (op *DebugLevel) ImportStub(k *Importer) (ret interface{}, err error) {
	if !inProg(k) {
		ret = op
	} else if str, found := decode.FindChoice(op, op.Str); !found {
		err = errutil.Fmt("choice %s not found in %T", op.Str, op)
	} else {
		found := -1
		for i, v := range op.Compose().Strings {
			if v == str {
				found = i
				break
			}
		}
		if found < 0 {
			err = errutil.Fmt("index %s not found in %T", op.Str, op)
		} else {
			ret = found
		}
	}
	return
}

// turn comment execution into empty statements
func (op *Comment) ImportStub(k *Importer) (ret interface{}, err error) {
	if !inProg(k) {
		ret = op
	} else {
		ret = &debug.Log{Level: debug.Note, Value: &core.FromText{&core.Text{op.Lines.Str}}}
	}
	return
}

// a hopefully temporary hack
func inProg(k *Importer) (ret bool) {
	for _, k := range k.decoder.Path {
		if k == "story.Activity" {
			ret = true
			break
		}
	}
	return
}
