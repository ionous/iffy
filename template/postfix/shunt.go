package postfix

import (
	"github.com/ionous/errutil"
)

// Shunt-ing yard ( dijkstra ) to convert an infix function stream to a postfix function list.
// The key differences to the traditional algorithm are:
//  * Tokens are replaced by Functions; operands are zero-arity functions ( presumably functions which return the value of the operand. )
//  * We use one yard for each sub-expression rather than making parentheses a token or function in the stream symbols.
type Shunt struct {
	out   Expression
	yards Yards
}

// Yards contains a stack of shunting yards: one for each pending sub/expression.
type Yards struct {
	list []Yard
}

// Len returns count of pending sub/expressions.
func (ys Yards) Len() int {
	return len(ys.list)
}

// Top aka innermost sub/expression.
func (ys Yards) Top() (ret *Yard) {
	last := len(ys.list) - 1
	return &ys.list[last]
}

// NewYard starts a new sub/expression.
func (ys *Yards) NewYard() {
	ys.list = append(ys.list, Yard{})
}

// Pop ends a sub/expression.
func (ys *Yards) Pop() (ret Yard) {
	last := len(ys.list) - 1
	ret, ys.list = ys.list[last], ys.list[:last]
	return
}

// Yard containing functions suspended in their transition from infix to postfix.
type Yard struct {
	stack []Function
}

// Flush returns the shunt's postfix ordered output, clearing the shunt.
func (s *Shunt) Flush() (ret Expression, err error) {
	if cnt := s.yards.Len(); cnt > 1 {
		err = errutil.New(cnt-1, "unclosed sub expressions")
	} else if cnt > 0 {
		if e := s.EndSubExpression(); e != nil {
			err = e
		} else {
			ret, s.out = s.out, nil
		}
	}
	return
}

func (s *Shunt) AddExpression(prev []Function) (err error) {
	s.out = append(s.out, prev...)
	return
}

// AddFunction appends the next (infix) operation to the pending postfix expression.
// Zero-arity functions are moved directly to the output, otherwise they are shunted to a yard such that higher precedence functions will pop out of the yard first, leaving
func (s *Shunt) AddFunction(next Function) (err error) {
	if next.Arity() == 0 {
		s.out = append(s.out, next)
	} else {
		if s.yards.Len() == 0 {
			s.yards.NewYard()
		}
		yard := s.yards.Top()
		if cnt := len(yard.stack); cnt == 0 {
			yard.stack = append(yard.stack, next)
		} else {
			top := yard.stack[cnt-1]
			newp, topp := next.Precedence(), top.Precedence()
			// If incoming operator has higher precedence than stack's top,
			// push it on the stack.
			if newp > topp /*|| (newp == topp && next.IsRightAssoc())*/ {
				yard.stack = append(yard.stack, next)
			} else {
				// Pop while incoming operator has lower precedence than stack's top;
				// then push the incoming operator.
				for {
					if cnt := len(yard.stack); cnt == 0 {
						break
					} else {
						top := yard.stack[cnt-1]
						topp := top.Precedence()
						//
						if newp < topp || (newp == topp /*&& !next.IsRightAssoc()*/) {
							s.out = append(s.out, top)
							yard.stack = yard.stack[:cnt-1] // pop
						}
					}
				}
				yard.stack = append(yard.stack, next)
			}
		}
	}
	return
}

// BeginSubExpression delineates an opening parenthesis in a stream of functions.
func (s *Shunt) BeginSubExpression() (err error) {
	s.yards.NewYard()
	return
}

// EndSubExpression indicates a closing parenthesis in a stream of functions.
func (s *Shunt) EndSubExpression() (err error) {
	if s.yards.Len() == 0 {
		err = errutil.New("too many ends")
	} else {
		yard := s.yards.Pop()
		s.out = append(s.out, reverse(yard.stack)...)
	}
	return
}

// in-place reversal of the passed functions.
func reverse(a []Function) []Function {
	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}
	return a
}
