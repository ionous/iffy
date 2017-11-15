package chart

// spaces eats whitespace
var spaces SelfStatement = func(self SelfStatement, r rune) (ret State) {
	if isSpace(r) {
		ret = self
	}
	return
}
