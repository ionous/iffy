package core

// // some sort of assignment like a parameter call.
// type MakeRecord struct {
// 	Kind      string     // name of the record
// 	Arguments *Arguments // arguments to initialize record fields
// }

// func (op *MakeRecord) Compose() composer.Spec {
// 	return composer.Spec{
// 		Name:  "make_record",
// 		Group: "variables",
// 		Desc:  "Make Record: specify a set of fields.",
// 	}
// }

// func (op *MakeRecord) GetEval() interface{} {
// 	return nil // unknown
// }

// func (op *MakeRecord) GetAssignedValue(run rt.Runtime) (ret rt.Value, err error) {
// 	if v, e := op.getAssignedValue(run); e != nil {
// 		err = cmdError(err, e)
// 	} else {
// 		ret = v
// 	}
// 	return
// }

// func (op *MakeRecord) getAssignedValue(run rt.Runtime) (ret rt.Value, err error) {
// 	if rec, e := run.MakeRecord(op.kind); e != nil {
// 		err = e
// 	} else if fixme, ok := rec.(*generic.Record); !ok {
// 		err = errutil.New("unexpected record type %T", rec)
// 	} else {
// 		for _, a := range op.Arguments.Args {
// 			if val, e := a.GetAssignedValue(run); e != nil {
// 				err = errutil.Append(err, e)
// 			} else if e := fixme.SetField(field, val); e != nil {
// 				err = errutil.Append(err, e)
// 			}
// 		}
// 		if err == nil {
// 			ret = rec
// 		}
// 	}
// 	return
// }
