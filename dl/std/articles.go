package std

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/lang"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/printer"
	"io"
)

// UpperThe is equivalent to Inform7's [The]
type UpperThe struct {
	Obj rt.ObjectEval
}

// LowerThe is equivalent to Inform7's [the]
type LowerThe struct {
	Obj rt.ObjectEval
}

// UpperAn is equivalent to Inform7's [A/An]
type UpperAn struct {
	Obj rt.ObjectEval
}

// LowerAn is equivalent to Inform7's [a/an]
type LowerAn struct {
	Obj rt.ObjectEval
}

func (the *UpperThe) GetText(run rt.Runtime) (ret string, err error) {
	if obj, e := the.Obj.GetObject(run); e != nil {
		err = e
	} else {
		ret, err = articleName(run, "The", obj)
	}
	return
}

func (the *LowerThe) GetText(run rt.Runtime) (ret string, err error) {
	if obj, e := the.Obj.GetObject(run); e != nil {
		err = e
	} else {
		ret, err = articleName(run, "the", obj)
	}
	return
}

func (an *UpperAn) GetText(run rt.Runtime) (ret string, err error) {
	if obj, e := an.Obj.GetObject(run); e != nil {
		err = e
	} else if txt, e := articleName(run, "", obj); e != nil {
		err = e
	} else {
		ret = lang.Capitalize(txt)
	}
	return
}

func (a *LowerAn) GetText(run rt.Runtime) (ret string, err error) {
	if obj, e := a.Obj.GetObject(run); e != nil {
		err = e
	} else {
		ret, err = articleName(run, "", obj)
	}
	return
}

func (the *UpperThe) Execute(run rt.Runtime) error {
	return core.Print(run, the)
}

func (the *LowerThe) Execute(run rt.Runtime) error {
	return core.Print(run, the)
}

func (an *UpperAn) Execute(run rt.Runtime) error {
	return core.Print(run, an)
}

func (a *LowerAn) Execute(run rt.Runtime) error {
	return core.Print(run, a)
}

func printArticle(run rt.Runtime, article string, obj rt.Object) (err error) {
	if text, e := articleName(run, article, obj); e != nil {
		err = e
	} else if _, e := io.WriteString(run, text); e != nil {
		err = e
	}
	return
}

// You can only just make out the lamp-post.", or "You can only just make out _ Trevor.", or "You can only just make out the soldiers."
func articleName(run rt.Runtime, article string, obj rt.Object) (ret string, err error) {
	// FIX? i think filters would be better -- especically in printWithArticles -- but this matches existing code.
	var buffer printer.Span
	if e := printName(rt.Writer(run, &buffer), obj); e != nil {
		err = e
	} else {
		name := buffer.String()
		var proper bool
		if e := obj.GetValue("proper-named", &proper); e != nil {
			err = e
		} else if proper {
			ret = lang.Titlecase(name)
		} else {
			if len(article) == 0 {
				var indefinite string
				if e := obj.GetValue("indefinite article", &indefinite); e != nil {
					err = e
				} else {
					article = indefinite
					if len(article) == 0 {
						var plural bool
						if e := obj.GetValue("plural-named", &plural); e != nil {
							err = e
						} else {
							if plural {
								article = "some"
							} else if lang.StartsWithVowel(name) {
								article = "an"
							} else {
								article = "a"
							}
						}
					}
				}
			}
			// by now, article should exist; except if err is set.
			if len(article) > 0 {
				ret = article + " " + name
			}
		}
	}
	return
}
