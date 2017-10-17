package chart

// spaces eats whitespace and then becomes the requested next state
// a bit like using chart history mechanism.
var spaces SelfStatement = func(self SelfStatement, r rune) (ret State) {
	if isSpace(r) {
		ret = self
	}
	return
}
