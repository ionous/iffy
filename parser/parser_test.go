package parser_test

// import (
// 	"strings"
// 	// "github.com/ionous/iffy/parser"
// 	"github.com/ionous/sliceOf"
// 	"testing"
// )

// type ParserTest struct {
// 	Output Step
// 	Input  []string
// 	// test will expect the number of clarifications and the number of steps to match
// 	Clarify []string
// }

// type Context struct {
// }

// type Step interface {
// 	// Process the passed input, and return the next step
// 	// returns nil when finished.
// 	NextStep(*Context, string) (*Step, error)
// }

// // Debig
// type Response struct {
// 	Prompt string
// 	Step   *Step
// }

// // probably not a string, but instead:
// type Result struct {
// 	Action string
// 	Nouns  []string
// }

// func (Response) NextStep(context *Context, input string) (ret *Step, err error) {
// 	return
// }

// // for multi-action/ held reponses: take before insert
// // func (Results) NextStep(context *Context, input string) (ret *Step, err error) {
// // 	return
// // }

// func (Result) NextStep(context *Context, input string) (ret *Step, err error) {
// 	return
// }

// func (p *ParserTest) Run(t *testing.T) {
// }

// var shortDirections = sliceOf.String(
// 	"n", "s", "e", "w", "ne", "nw", "se", "sw",
// 	"u", "up", "ceiling", "above", "sky",
// 	"d", "down", "floor", "below", "ground",
// )
// var directions = sliceOf.String(
// 	"north",
// 	"south",
// 	"east",
// 	"west",
// 	"northeast",
// 	"northwest",
// 	"southeast",
// 	"southwest",
// 	"up",   // "up above",
// 	"down", // "ground",
// 	"inside",
// 	"outside")

// var ADirection = append(shortDirections, directions...)

// func TestParser(t *testing.T) {
// 	// first, we want to test a simple set of example actions,
// 	// all of which start the same way, but end with different actions.
// 	// later, we will test disambiguation; errors; multiple objects: etc.
// 	t.Run("looking", func(t *testing.T) {
// 		t.Run("look", func(t *testing.T) {
// 			p := ParserTest{
// 				Verb: sliceOf.String("look", "l"),
// 				Output: Result{
// 					Action: "Look",
// 				},
// 			}
// 			p.Run(t)
// 		})
// 		t.Run("examine", func(t *testing.T) {
// 			p := ParserTest{
// 				Verb:   sliceOf.String("look", "l"),
// 				Phrase: "at something",
// 				Output: Result{
// 					Action: "Examine",
// 					Nouns:  sliceOf.String("something"),
// 				},
// 			}
// 			p.Run(t)
// 		})
// 		t.Run("search", func(t *testing.T) {
// 			p := ParserTest{
// 				Verb:   sliceOf.String("look", "l"),
// 				Phrase: "inside/in/into/through/on something",
// 				Output: Result{
// 					Action: "Search",
// 					Nouns:  sliceOf.String("something"),
// 				},
// 			}
// 			p.Run(t)
// 		})
// 		t.Run("look under", func(t *testing.T) {
// 			p := ParserTest{
// 				Verb:   sliceOf.String("look", "l"),
// 				Phrase: "under something",
// 				Output: Result{
// 					Action: "LookUnder",
// 					Nouns:  sliceOf.String("something"),
// 				},
// 			}
// 			p.Run(t)
// 		})

// 		t.Run("consult", func(t *testing.T) {
// 			p := ParserTest{
// 				Verb:   sliceOf.String("look", "l"),
// 				Phrase: "under something",
// 				Output: Result{
// 					Action: "LookUnder",
// 					Nouns:  sliceOf.String("something"),
// 				},
// 			}
// 			p.Run(t)
// 		})

// 		t.Run("examine dir", func(t *testing.T) {
// 			for _, dir := range ADirection {
// 				p := ParserTest{
// 					Verb:   sliceOf.String("look", "l"),
// 					Phrase: dir,
// 					Output: Result{
// 						Action: "Examine",
// 						Nouns:  sliceOf.String(dir),
// 					},
// 				}
// 			}
// 			p.Run(t)
// 		})
// 		t.Run("examine to dir", func(t *testing.T) {
// 			for _, dir := range ADirection {
// 				p := ParserTest{
// 					Verb:   sliceOf.String("look to", "l to"),
// 					Phrase: dir,
// 					Output: Result{
// 						Action: "Examine",
// 						Nouns:  sliceOf.String(dir),
// 					},
// 				}
// 				p.Run(t)
// 			}
// 		})
// 	})
// }
