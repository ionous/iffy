package chart

// spaces eats whitespace
var spaces = SelfStatement("spaces", func(self State, r rune) (ret State) {
	if isSpace(r) {
		ret = self
	}
	return
})
