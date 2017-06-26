package text

import (
	"github.com/ionous/iffy/lang"
	"github.com/ionous/iffy/rt"
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

func (the *UpperThe) GetText(run rt.Runtime) (string, error) {
	return articleNamed(run, the.Obj, "The")
}

func (the *LowerThe) GetText(run rt.Runtime) (string, error) {
	return articleNamed(run, the.Obj, "the")
}

func (an *UpperAn) GetText(run rt.Runtime) (ret string, err error) {
	if txt, e := articleNamed(run, an.Obj, ""); e != nil {
		err = e
	} else {
		ret = lang.Capitalize(txt)
	}
	return
}

func (a *LowerAn) GetText(run rt.Runtime) (ret string, err error) {
	if txt, e := articleNamed(run, a.Obj, ""); e != nil {
		err = e
	} else {
		ret = txt
	}
	return
}

// You can only just make out the lamp-post.", or "You can only just make out _ Trevor.", or "You can only just make out the soldiers."
func articleNamed(run rt.Runtime, noun rt.ObjectEval, article string) (ret string, err error) {
	if obj, e := noun.GetObject(run); e != nil {
		err = e
	} else {
		var name string
		if e := obj.GetValue("name", &name); e != nil {
			err = e
		} else {
			var proper bool
			if e := obj.GetValue("proper-named", &proper); e != nil {
				err = e
			} else if proper {
				ret = lang.Titlecase(name)
			} else {
				name = lang.Lowercase(name)
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
	}
	return
}
