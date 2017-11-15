package types

// Operator represents built-in binary functions.
type Operator int

//go:generate stringer -type=Operator
const Precedence = 4

const (
	MUL Operator = (5 << Precedence) | iota
	QUO Operator = (5 << Precedence) | iota
	REM Operator = (5 << Precedence) | iota

	ADD Operator = (4 << Precedence) | iota
	SUB Operator = (4 << Precedence) | iota

	EQL Operator = (3 << Precedence) | iota
	NEQ Operator = (3 << Precedence) | iota
	LSS Operator = (3 << Precedence) | iota
	LEQ Operator = (3 << Precedence) | iota
	GTR Operator = (3 << Precedence) | iota
	GEQ Operator = (3 << Precedence) | iota

	LAND Operator = (2 << Precedence) | iota
	LOR  Operator = (1 << Precedence) | iota
)

// Arity of operator is two.
func (i Operator) Arity() int { return 2 }

// Precedence of the corresponding infix operator.
func (i Operator) Precedence() int { return int(i >> Precedence) }
