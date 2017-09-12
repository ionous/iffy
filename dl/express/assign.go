package express

// func ConvertStmt(c spec.Block, l ast.Stmt) (err error) {
// 	switch l := l.(type) {
// 	case *ast.AssignStmt:
// 		if cnt := len(l.Lhs); cnt != len(l.Rhs) {
// 			err = errutil.New("left and right sides dont match")
// 		} else if cnt == 1 {
// 			if e := assign(c, l.lhs[0], l.rhs[0]); e != nil {
// 				err = e
// 			}
// 		} else if cnt > 1 {
// 			if c.Cmds().Begin() {
// 				for i := 0; i < cnt; i++ {
// 					lhs, rhs := l.Lhs[i], l.Rhs[i]
// 					if e := assign(c, lhs, rhs); e != nil {
// 						err = e
// 						break
// 					}
// 				}
// 				c.End()
// 			}
// 		}
// 	}
// 	return
// }

// func assign(c spec.Block, lhs, rhs ast.Expr) (err error) {
// the issue with assignment is that we dont know the type of the right and left sides -- we could assume that the left gives us a property
// but we still dont know about the target ---
// funnily enough the command builder does ---
// that is, when we setField we know the slot

// if n, ok := lhs.(*ast.SelectorExpr); !ok {
// 	// FIX: and more so... we should be an object property
// 	err = errutil.New("error on left, expected object")
// } else if x, ok := n.X.(*ast.Ident); !ok {
// 	err = errutil.Fmt("expected object identifer, got %T", n.X)
// } else if v, e := ConvertExpr(c, rhs); e != nil {
// 	err = errutil.New("error on right", e)
// } else {
// 	obj := makeObject(c, x)
// 	switch v := v.(type) {
// 	case rt.NumberEval:
// 		ret = &core.SetNum{obj, n.Sel.Name, v}
// 	case rt.TextEval:
// 		ret = &core.SetText{obj, n.Sel.Name, v}
// 	default:
// 		err = errutil.New("unknown type %T", v)
// 	}
// }
// return
// }
