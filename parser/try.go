package parser

// import (
// 	"strings"
// )

// // Word matches one word.
// type Word struct {
// 	Word string
// }

// func (try *Word) Try(ctx *Context, cs *Cursor) (okay bool) {
// 	if n, ok := cs.NextWord(); ok {
// 		okay = strings.EqualFold(n, try.Word)
// 	}
// 	return
// }

// // AnyOf matches any one of the passed matchers.
// type AnyOf struct {
// 	AnyOf []Tryer
// }

// func (try *AnyOf) Try(ctx *Context, cs *Cursor) (okay bool) {
// 	for _, s := range try.AnyOf {
// 		cs := *cs // reset the cursor each time
// 		if s.Try(cs) {
// 			okay = true
// 			break
// 		}
// 	}
// 	return
// }

// // AllOf matches the passed matchers in order.
// type AllOf struct {
// 	AllOf []Tryer
// }

// func (try *AllOf) Try(ctx *Context, cs *Cursor) (okay bool) {
// 	for _, s := range try.AllOf {
// 		if s.Try(cs) {
// 			okay = true
// 			break
// 		}
// 	}
// 	return
// }
