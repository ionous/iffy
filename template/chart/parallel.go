package chart

// MakeParallel region; run all of the passed states until they all return nil.
func MakeParallel(name string, rs ...State) State {
	return SelfStatement(name, func(self State, r rune) (ret State) {
		var cnt int
		for _, s := range rs {
			if next := s.NewRune(r); next != nil {
				rs[cnt] = next
				cnt++
			}
		}
		if cnt > 0 {
			rs = rs[:cnt]
			ret = self
		}
		return
	})
}
