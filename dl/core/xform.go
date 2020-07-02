package core

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
)

// MakeSingular
type MakeSingular struct {
	Text rt.TextEval
}

// Compose defines a spec for the composer editor.
func (*MakeSingular) Compose() composer.Spec {
	return composer.Spec{
		Name:  "singularize",
		Group: "format",
		Desc:  "Singularize: Returns the singular form of a plural word. (ex. apple for apples )",
		Spec:  "the singular of {text:text_eval}",
	}
}

func (op *MakeSingular) GetText(run rt.Runtime) (ret string, err error) {
	if t, e := rt.GetText(run, op.Text); e != nil {
		err = e
	} else if len(t) > 0 {
		ret = run.SingularOf(t)
	}
	return
}

// MakePlural
type MakePlural struct {
	Text rt.TextEval
}

// Compose defines a spec for the composer editor.
func (*MakePlural) Compose() composer.Spec {
	return composer.Spec{
		Name:  "pluralize",
		Group: "format",
		Desc:  "Pluralize: Returns the plural form of a singular word. (ex.  apples for apple. )",
		Spec:  "the plural of {text:text_eval}",
	}
}

func (op *MakePlural) GetText(run rt.Runtime) (ret string, err error) {
	if t, e := rt.GetText(run, op.Text); e != nil {
		err = e
	} else if len(t) > 0 {
		ret = run.PluralOf(t)
	}
	return
}
