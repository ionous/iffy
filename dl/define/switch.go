package define

import (
	"github.com/ionous/iffy/parser"
	"github.com/ionous/iffy/pat/patspec"
	"github.com/ionous/iffy/rt"
)

type Grammar struct {
	Match parser.Scanner // should be all of / any of.
}

type Pattern struct {
	Pattern patspec.MakePattern
}

type When struct {
	Event  string
	Target string
	Eval   rt.Execute
	// Filters
}

type Plural struct {
	Single, Plural string
}

type The struct {
	// The Spaceport is a room.
	// s.The("story",
	// 	Called("The empty room"),
	// 	HasText("author", T("me")),
	// 	HasText("headline", T("extra extra")))
	// s.The("room",
	// 	Called("somewhere"),
	// 	HasText("description", T("an empty room")),
	// )
}

func (a *The) Define(f *Facts) (err error) {
	return
}
func (a *When) Define(f *Facts) (err error) {
	return
}
func (a *Pattern) Define(f *Facts) (err error) {
	return
}
func (a *Grammar) Define(f *Facts) (nil error) {
	println("xxx grammar?")
	f.Grammar.Match = append(f.Grammar.Match, a.Match)
	return
}
func (a *Plural) Define(f *Facts) (err error) {
	// type Plural PluralRule
	return
}

// s.The("player", Exists(), In("somewhere"))

// s.The("room",
// 		Called("studio"),
// 		When("printing details").Always(
// 			Choose{

// Instead of going somewhere from the spaceport when the player carries something:
//     let N be "[is-are the list of things carried by the player] really suitable gear to take to the moon?" in sentence case;
