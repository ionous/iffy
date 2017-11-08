package express

// ParseDirective handles multi-part expressions:
// namely "go <return the value of a function>", and
// and, "directive <return the value of a pattern>".
// func ParseDirective(c spec.Block, parts []string, hint r.Type) (err error) {
// 	// not wild about this being here --
// 	// but they are technically expressions and not templates.
// 	if len(parts) == 1 {
// 		err = ParseExpression(c, parts[0], hint)
// 	} else {
// 		switch op, rest := parts[0], parts[1:]; op {
// 		case "go":
// 			op, rest := rest[0], rest[1:]
// 			if c.Cmd(op).Begin() {
// 				if len(rest) > 0 {
// 					err = ParseDirective(c, rest, hint)
// 				}
// 				c.End()
// 			}
// 		case "determine":
// 			if c.Cmd(op).Begin() {
// 				pat, rest := rest[0], rest[1:]
// 				if c.Cmd(pat).Begin() {
// 					if len(rest) > 0 {
// 						err = ParseDirective(c, rest, hint)
// 					}
// 					c.End()
// 				}
// 				c.End()
// 			}
// 		default:
// 			err = errutil.New("unknown multi-part expression", parts)
// 		}
// 	}
// 	return
// }
