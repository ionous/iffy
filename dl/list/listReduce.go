package list

// import (
// 	"github.com/ionous/iffy/dl/composer"
// 	"github.com/ionous/iffy/dl/core"
// 	"github.com/ionous/iffy/dl/pattern"
// 	"github.com/ionous/iffy/dl/term"
// 	"github.com/ionous/iffy/object"
// 	"github.com/ionous/iffy/rt"
// 	g "github.com/ionous/iffy/rt/generic"
// )

// // A normal reduce would return a value, instead we accumulate into a variable
// type Reduce struct {
// 	FromList, Into, Pattern string // variable names
// }

// func (*Reduce) Compose() composer.Spec {
// 	return composer.Spec{
// 		Name:  "list_map",
// 		Group: "list",
// 		Desc: `Reduce List: Transform the values from one list by combining them into a single value.
// 		The named pattern is called with two parameters: 'in' ( each element of the list ) and 'out' ( ex. a record ).`,
// 	}
// }

// func (op *Reduce) Execute(run rt.Runtime) (err error) {
// 	if e := op.execute(run); e != nil {
// 		err = cmdError(op, e)
// 	}
// 	return
// }

// func (op *Reduce) execute(run rt.Runtime) (err error) {
// 	var one, two core.R
// 	det := makeDet("sort", &core.FromRecord{&one}, &core.FromRecord{&two})

// 	var pat pattern.ActivityPattern
// 	if fromList, e := run.GetField(object.Variables, op.FromList); e != nil {
// 		err = e
// 	} else if outVal, e := run.GetField(object.Variables, op.Into); e != nil {
// 		err = e
// 	} else if e := run.GetEvalByName(op.Pattern, &pat); e != nil {
// 		err = e
// 	} else {
// 		for it := g.ListIt(fromList); it.HasNext(); {
// 			if inVal, e := it.GetNext(); e != nil {
// 				err = e
// 				break
// 			} else if newVal, e := op.reduceOne(run, pat, inVal, outVal); e != nil {
// 				err = e
// 				break
// 			} else {
// 				outVal = newVal
// 			}
// 		}
// 		if err == nil {
// 			err = run.SetField(object.Variables, op.Into, outVal)
// 		}
// 	}
// 	return
// }

// // similar to express buildPattern
// func makeReduce(name string, first, second core.Assignment) rt.BoolEval {
// 	return &pattern.DetermineAct{
// 		Pattern: name,
// 		Arguments: &core.Arguments{
// 			Args: []*core.Argument{{
// 				Name: "$1",
// 				From: first,
// 			}, {
// 				Name: "$2",
// 				From: second,
// 			}},
// 		},
// 	}
// }

// // see also pattern.Stitch
// func (op *Reduce) reduceOne(run rt.Runtime, pat pattern.ActivityPattern, inVal, outVal g.Value) (ret g.Value, err error) {
// 	var parms term.Terms
// 	if e := pat.Prepare(run, &parms); e != nil {
// 		err = e
// 	} else if e := parms.SetField(object.Variables, "in", inVal); e != nil {
// 		err = e
// 	} else if e := parms.SetField(object.Variables, "out", outVal); e != nil {
// 		err = e
// 	} else {
// 		run.PushScope(&parms)
// 		var locals term.Terms
// 		if e := pat.ComputeLocals(run, &locals); e != nil {
// 			err = e
// 		} else {
// 			run.PushScope(&locals)
// 			if e := pat.Execute(run); e != nil {
// 				err = e
// 			} else if v, e := parms.GetField(object.Variables, "out"); e != nil {
// 				err = e
// 			} else {
// 				ret = v
// 			}
// 			run.PopScope()
// 		}
// 		run.PopScope()

// 	}
// 	return
// }
