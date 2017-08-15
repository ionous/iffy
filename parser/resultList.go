package parser

type ResultList struct {
	List   []Result
	Length int
}

func (rs *ResultList) ResultLen() int {
	return rs.Length
}

func (rs *ResultList) AddResult(r Result) {
	if rl, ok := r.(*ResultList); ok {
		rs.List = append(rs.List, rl.List...)
		rs.Length += rl.Length
	} else {
		rs.List = append(rs.List, r)
		rs.Length += r.ResultLen()
	}
	return
}

func (rs *ResultList) Last() (ret Result, okay bool) {
	if cnt := len(rs.List); cnt > 0 {
		ret, okay = rs.List[cnt-1], true
	}
	return
}

func (rs *ResultList) Objects() (objs []string) {
	for _, r := range rs.List {
		switch k := r.(type) {
		case ResolvedObject:
			n := k.Noun
			objs = append(objs, n.GetId())
		case ResolvedMulti:
			for _, n := range k.Nouns {
				objs = append(objs, n.GetId())
			}
		}
	}
	return
}

// func (rs *ResultList) Results() (act string, objs []string, okay bool) {
// 	rl := rs.List
// 	if cnt := len(rl); cnt > 0 {
// 		last := rl[cnt-1]
// 		if a, ok := last.(ResolvedAction); ok {
// 			i, cnt := 0, cnt-1
// 			for ; i < cnt; i++ {
// 				r := rl[i]
// 				switch k := r.(type) {
// 				case ResolvedObject:
// 					objs = append(objs, k.Id)
// 				}
// 			}
// 			if i == cnt {
// 				act = a.Name
// 				okay = true
// 			}
// 		}
// 	}
// 	return
// }
