package core

import (
	"strings"

	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/lang"
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
		Spec:  "the singular {text:text_eval}",
	}
}

func (op *MakeSingular) GetText(run rt.Runtime) (ret string, err error) {
	if t, e := rt.GetText(run, op.Text); e != nil {
		err = cmdError(op, e)
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
		err = cmdError(op, e)
	} else if len(t) > 0 {
		ret = run.PluralOf(t)
	}
	return
}

type MakeUppercase struct {
	Text rt.TextEval
}
type MakeLowercase struct {
	Text rt.TextEval
}

// Start each word with a capital letter.
type MakeTitleCase struct {
	Text rt.TextEval
}

// Start each sentence with a capital letter.
type MakeSentenceCase struct {
	Text rt.TextEval
}

type MakeReversed struct {
	Text rt.TextEval
}

// Compose defines a spec for the composer editor.
func (*MakeUppercase) Compose() composer.Spec {
	return composer.Spec{
		Name:  "make_uppercase",
		Group: "format",
		Desc: `Uppercase: returns new text, with every letter turned into uppercase. 
		For example, "APPLE" from "apple".`,
		Spec: "{text:text_eval} in uppercase",
	}
}

// Compose defines a spec for the composer editor.
func (*MakeLowercase) Compose() composer.Spec {
	return composer.Spec{
		Name:  "make_lowercase",
		Group: "format",
		Desc: `Lowercase: returns new text, with every letter turned into lowercase. 
		For example, "shout" from "SHOUT".`,
		Spec: "{text:text_eval} in lowercase",
	}
}

// Compose defines a spec for the composer editor.
func (*MakeTitleCase) Compose() composer.Spec {
	return composer.Spec{
		Name:  "make_title_case",
		Group: "format",
		Desc: `Title case: returns new text, starting each word with a capital letter. 
		For example, "Empire Apple" from "empire apple".`,
		Spec: "{text:text_eval} in title-case",
	}
}

// Compose defines a spec for the composer editor.
func (*MakeSentenceCase) Compose() composer.Spec {
	return composer.Spec{
		Name:  "make_sentence_case",
		Group: "format",
		Desc: `Sentence case: returns new text, start each sentence with a capital letter. 
		For example, "Empire Apple." from "Empire apple.".`,
		Spec: "{text:text_eval} in sentence-case",
	}
}

// Compose defines a spec for the composer editor.
func (*MakeReversed) Compose() composer.Spec {
	return composer.Spec{
		Name:  "make_reversed",
		Group: "format",
		Desc: `Reverse text: returns new text flipped back to front. 
		For example, "elppA" from "Apple", or "noon" from "noon".`,
		Spec: "{text:text_eval} in reverse",
	}
}

func (op *MakeLowercase) GetText(run rt.Runtime) (ret string, err error) {
	if t, e := rt.GetText(run, op.Text); e != nil {
		err = cmdError(op, e)
	} else {
		ret = strings.ToLower(t)
	}
	return
}

func (op *MakeUppercase) GetText(run rt.Runtime) (ret string, err error) {
	if t, e := rt.GetText(run, op.Text); e != nil {
		err = cmdError(op, e)
	} else {
		ret = strings.ToUpper(t)
	}
	return
}

func (op *MakeTitleCase) GetText(run rt.Runtime) (ret string, err error) {
	if t, e := rt.GetText(run, op.Text); e != nil {
		err = cmdError(op, e)
	} else if len(t) > 0 {
		ret = lang.Titlecase(t)
	}
	return
}

func (op *MakeSentenceCase) GetText(run rt.Runtime) (ret string, err error) {
	if t, e := rt.GetText(run, op.Text); e != nil {
		err = cmdError(op, e)
	} else if len(t) > 0 {
		ret = lang.SentenceCase(t)
	}
	return
}

func (op *MakeReversed) GetText(run rt.Runtime) (ret string, err error) {
	if t, e := rt.GetText(run, op.Text); e != nil {
		err = cmdError(op, e)
	} else {
		a := []rune(t)
		for i := len(a)/2 - 1; i >= 0; i-- {
			opp := len(a) - 1 - i
			a[i], a[opp] = a[opp], a[i]
		}
		ret = string(a)
	}
	return
}
