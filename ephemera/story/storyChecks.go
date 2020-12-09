package story

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ephemera"
)

// ImportPhrase implements StoryStatement
func (op *TestRule) ImportPhrase(k *Importer) (err error) {
	if n, e := op.TestName.NewName(k); e != nil {
		err = e
	} else if hook, e := op.Hook.ImportProgram(k); e != nil {
		err = e
	} else if prog, e := k.NewGob(hook.SlotType(), hook.CmdPtr()); e != nil {
		err = e
	} else {
		k.NewTestProgram(n, prog)
	}
	return
}

// ImportPhrase implements StoryStatement
func (op *TestScene) ImportPhrase(k *Importer) (err error) {
	if n, e := op.TestName.NewName(k); e != nil {
		err = e
	} else {
		err = k.CollectTest(n, func() error {
			return op.Story.ImportStory(k)
		})
	}
	return
}

func (op *TestStatement) ImportPhrase(k *Importer) (err error) {
	if t := op.Test; t == nil {
		err = ImportError(op, op.At, errutil.Fmt("%w Test", MissingSlot))
	} else if n, e := op.TestName.NewName(k); e != nil {
		err = e
	} else {
		err = t.ImportTest(k, n)
	}
	return
}

type Testing interface {
	ImportTest(k *Importer, testName ephemera.Named) (err error)
}

func (op *TestOutput) ImportTest(k *Importer, testName ephemera.Named) (err error) {
	// note: we use the raw lines here, we don't expect the text output to be a template.
	k.NewTestExpectation(testName, "execute", op.Lines.Str)
	return
}
