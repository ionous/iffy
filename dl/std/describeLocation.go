package std

import (
	"bytes"
	"github.com/ahmetb/go-linq"
	"github.com/ionous/iffy/rt"
)

// DescribeLocation prints details about the targeted location, including "paragraphs” for notable objects in the location, and a sentence for otherwise unremarkable objects.
// It is a command, for now, and not a pattern due to the complexity of the code required to avoid describing the same object more than once.
type DescribeLocation struct {
	Location rt.Object
}

func (op *DescribeLocation) Execute(run rt.Runtime) (err error) {
	// print the initial description of the location
	if e := rt.Determine(run, &PrintLocation{op.Location}); e != nil {
		err = e
	} else if mentioned, unmentioned, e := op.mentionNotables(run); e != nil {
		err = e
	} else if common, e := run.GetObjStreamMatching(run.Emplace(&UnremarkableObjects{op.Location})); e != nil {
		err = e
	} else {
		n := PrintCommonObjects{}
		linq.From(unmentioned).Concat(linq.From(
			common).Except(linq.From(mentioned))).ToSlice(&n.Objects)
		if e := rt.Determine(run, &n); e != nil {
			err = e
		}
	}
	return
}

func (op *DescribeLocation) mentionNotables(run rt.Runtime) (mentioned, unmentioned []rt.Object, err error) {
	// determine notable objects
	if note, e := run.GetObjStreamMatching(run.Emplace(&NotableObjects{op.Location})); e != nil {
		err = e
	} else {
		for note.HasNext() {
			if obj, e := note.GetObject(); e != nil {
				err = e
				break
			} else {
				// describe the object:
				// objects may not have anything to say, in which case they become "unmentioned"
				var buf bytes.Buffer
				if e := rt.WritersBlock(run, &buf, func() (err error) {
					return rt.Determine(run, &DescribeObject{obj})
				}); e != nil {
					err = e
					break
				} else if buf.Len() == 0 {
					unmentioned = append(unmentioned, obj)
				} else if _, e := run.Writer().Write(buf.Bytes()); e != nil {
					err = e
					break
				} else {
					mentioned = append(mentioned, obj)
				}
			}
		}
	}
	return
}

// PrintLocation -> title and brief
// paragraphs for notable objects
// a sentence for other unremarkable objects

// * collecting “mentioned objects” — ? capture “un/mentioned” objects by trying to “describe object” into a buffer, and seeing if the buffer is empty.
//    * for now, do this in a command i guess -- the only problem is that you cant attach native commands to patterns, you have to invoke some special command.
//    * you’d have to have local variable, or a buffer “if/else” command with an @ for “my own text” ( maybe you could implement "comprise" this way, either put this as a side effect of creating a stream — or, stuff this — need some sort of “add to list” into a variable of describe notable objects.but, we dont really even have a way of passing references other than object — so youd have to pas
// }
