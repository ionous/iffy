package trigger

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/parser"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rtm"
	"strings"
)

// Parse, for testing's sake, expects no ambiguity, no missing objects.
func Parse(run *rtm.Rtm, input string) (err error) {
	in := strings.Fields(input)
	ctx := Context{run, run.Events}
	if scope, e := ctx.GetPlayerScope(""); e != nil {
		err = e
	} else if res, e := run.Scan(ctx, scope, parser.Cursor{Words: in}); e != nil {
		err = errutil.New(e, "for", input)
	} else if list, ok := res.(*parser.ResultList); !ok {
		err = errutil.New("expected result list %T", res)
	} else if last, ok := list.Last(); !ok {
		err = errutil.New("result list was empty")
	} else if act, ok := last.(ResolvedTrigger); !ok {
		err = errutil.New("expected resolved action %T", last)
	} else if objs, e := objectify(run, list.Objects()); e != nil {
		err = e
	} else {
		pose := Pose{objs}
		run := rt.AtFinder(run, pose)
		err = act.Execute(run)
	}
	return
}

// turn ids into objects.
func objectify(run rt.Runtime, ids []ident.Id) (ret []rt.Object, err error) {
	for _, id := range ids {
		if obj, ok := run.GetObject(id.Name); !ok {
			err = errutil.New("couldnt find object", id)
			break
		} else {
			ret = append(ret, obj)
		}
	}
	return
}
