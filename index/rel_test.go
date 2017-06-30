package index

// func (assert *IndexSuite) TestSomething() {
// 	pairs := []struct {
// 		primary   string
// 		secondary []string
// 	}{
// 		{"claire", sliceOf.String("boomba", "rocky")},
// 		{"grace", sliceOf.String("plume")},
// 		{"hiro", sliceOf.String("loofa")},
// 		{"marja", sliceOf.String("petra")},
// 	}
// 	n := MakeRelation(false, true)
// 	for _, pair := range pairs {
// 		for _, v := range pair.secondary {
// 			data := strings.Join(sliceOf.String(pair.primary, v), "+")
// 			changed := n.Relate(pair.primary, v, data)
// 			assert.True(changed)
// 		}
// 	}

// 	for _, pair := range pairs {
// 		var collect []string
// 		sink := func(other, data string) bool {
// 			collect = append(collect, other)
// 			return false
// 		}
// 		if pets := n.Index[Secondary].Walk(pair.primary, sink); assert.Equal(len(pair.secondary), pets) {
// 			assert.EqualValues(pair.secondary, collect)
// 		}
// 		break
// 	}
// }
