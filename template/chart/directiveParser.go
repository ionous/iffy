package chart

// // import (
// // 	"github.com/ionous/errutil"
// // 	"github.com/ionous/iffy/template/postfix"
// // )

// // DirectiveParser reads an expression chain.
// type DirectiveParser struct {
// 	pipe postfix.Pipe
// 	err  error
// }

// // // MakeDirectiveParser initializes a new directive parser.
// // func MakeDirectiveParser(seed []rune) (ret DirectiveParser) {
// // 	ret.seed.list = seed
// // 	return
// // }

// // GetExpression returns the parsed chain.
// func (p DirectiveParser) GetExpression() (ret postfix.Expression, err error) {
// 	if e := p.err; e != nil {
// 		err = e
// 	} else {
// 		ret = p.b.Expression()
// 	}
// 	return
// }

// // NewRune starts just after the opening of a directive or its trim.
// // we expect an operand ( with trailing expression ) or a function.
// func (p *DirectiveParser) NewRune(r rune) State {

// 	}

// 	return parseChain(r, &id, Statement(func(r rune) (ret State) {
// 		s
// 		// 		if arg, e := arg.GetArg(); e != nil {
// 		// 			p.err = e // no argument specified
// 		// 		} else if arg == nil {
// 		// 			ret = nil // an unhandled rune
// 		// 		} else if e := p.w.NextArgument(arg); e != nil {
// 		// 			p.err = e
// 		// 		} else {
// 		// 			// after every argument can come operators or close parens or the end
// 		// 			ret = parseChain(r, spaces, Statement(func(r rune) (ret State) {
// 		// 				if
// 		// 				if !isCloseParen(r) {
// 		// 					ret = p.operator(r)
// 		// 				} else if e := p.w.EndSubExpression(); e != nil {
// 		// 					p.err = e
// 		// 				} else {
// 		// 					ret = Statement(p.operator) // eat the closing paren.
// 		// 				}
// 		// 				return
// 		// 			}))
// 		// 		}
// 		return
// 	}))
// }

// // // start on the first rune of an operator:
// // func (p *DirectiveParser) operator(r rune) State {
// // 	var opParser opParser
// // 	return parseChain(r, &opParser, Statement(func(r rune) (ret State) {
// // 		if op, ok := opParser.GetOperator(); !ok {
// // 			if len(opParser.next) > 0 {
// // 				p.err = errutil.New("unhandled rune phrase", string(opParser.next))
// // 			} else {
// // 				ret = nil // an unhandled rune. noting: no operator rune is a subset of any legal phrase that can appear in the same spot as an operator, so there is never an unhandled rune phrase, only a single unhandled rune.
// // 			}
// // 		} else if _, e := p.b.Write(op); e != nil {
// // 			p.err = e
// // 		} else {
// // 			p.nofun = op != PIPE // when op is PIPE, nofun is false: we allow functions.
// // 			ret = parseChain(r, spaces, Statement(func(r rune) (ret State) {
// // 				if !isOpenParen(r) {
// // 					ret = p.argument(r)
// // 				} else if e := p.w.BeginSubExpression(); e != nil {
// // 					p.err = e
// // 				} else {
// // 					ret = Statement(p.argument) // eat the opening rune.
// // 				}
// // 				return
// // 			}))
// // 		}
// // 		return
// // 	}))
// // }
