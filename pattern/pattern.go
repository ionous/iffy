package pattern

// Register pattern rules ( which are saved as program in the model database, etc. )
func Register(reg func(value interface{})) {
	for _, cmd := range Rules {
		reg(cmd)
	}
}

// Rule associates a filter with an eval.
// Only when the rule's filter matches, the eval is... evaluated.
// Rules are usually combined into rule sets.
type Rule interface {
	RuleDesc() RuleDesc
}

// RuleDesc, see Rule interface.
type RuleDesc struct {
	Name    string      // under_scored name for the rule
	RuleSet interface{} //
}

// Rules contained by this package.
// fix: would it be better to list rule sets?
// the rule set elements could be used to find the individual rule types.
var Rules = []Rule{
	(*BoolRule)(nil),
	(*NumberRule)(nil),
	(*TextRule)(nil),
	(*NumListRule)(nil),
	(*TextListRule)(nil),
	(*ExecuteRule)(nil),
}
