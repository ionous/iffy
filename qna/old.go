package qna

// patternParamAt: ps.Prep(db,
// 	`select param from mdl_pat where pattern=? and idx=?`),
//
// func makeKeyWithIndex(obj string, idx int) keyType {
// 	return keyType{obj, "$" + strconv.Itoa(idx)}
// }
//
// returns the name of a field based on an index
// ex. especially for resolving positional pattern parameters into names.
// func (n *Fields) GetFieldByIndex(obj string, idx int) (ret string, err error) {
// 	if idx <= 0 {
// 		err = errutil.New("GetFieldByIndex out of range", idx)
// 	} else {
// 		// first, lookup the parameter name
// 		key := makeKeyWithIndex(obj, idx)
// 		// we use the cache to keep $(idx) -> param name.
// 		val, ok := n.pairs[key]
// 		if !ok {
// 			val, err = n.getCachingQuery(key, n.patternParamAt, obj, idx)
// 		}
// 		if field, ok := val.(string); !ok {
// 			err = rt.UnknownField{key.target, key.field}
// 		} else {
// 			ret = field
// 		}
// 	}
// 	return
// }
