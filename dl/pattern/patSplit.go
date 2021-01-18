package pattern

import "github.com/ionous/iffy/rt"

// call apply on each list rule ( in reverse order )
// we want the last added rules to win over earlier ones
// and we want to sort "post fix" rules to the end
func splitRules(run rt.Runtime, rules []*Rule) (ret []int, err error) {
	var pre, post int
	cnt := len(rules)
	a := make([]int, cnt)
	for i := cnt - 1; i >= 0; i-- {
		// apply the rule:
		// returns flags if the filters passed, -1 if they did not, error on any error.
		if flags, e := rules[i].GetFlags(run); e != nil {
			err = e
			break
		} else if flags >= 0 {
			if flags == Postfix {
				end := cnt - post - 1
				a[end], post = i, post+1
			} else {
				a[pre], pre = i, pre+1
				// FIX: add Replace as well
				if flags == Terminal {
					break
				}
			}
		}
	}
	if err == nil {
		if pre+post == cnt {
			ret = a
		} else {
			// keep all the prefixed items
			ret = a[:pre]
			// shift the post fixed items into the spot just after the prefixed items
			startOfPost := cnt - post
			for i := startOfPost; i < cnt; i++ {
				ret = append(ret, a[i])
			}
		}
	}
	return
}
