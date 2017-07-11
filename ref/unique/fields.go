package unique

import (
	"github.com/ionous/errutil"
	r "reflect"
)

func Fields(rtype r.Type) *FieldWalker {
	if rtype == nil || rtype.Kind() != r.Struct {
		err := errutil.New("can only walk structs", rtype.Name(), "is", rtype.Kind())
		panic(err)
	}
	fw := &FieldWalker{next: FieldInfo{
		Target: rtype,
	}}
	updateInfo(&fw.next)
	return fw
}

func (w *FieldWalker) HasNext() (okay bool) {
	return w.next.StructField != nil
}

func (w *FieldWalker) GetNext() (ret FieldInfo) {
	ret = w.next
	w.next.Index++
	updateInfo(&w.next)
	return
}

type FieldWalker struct {
	next FieldInfo // the next field
}

type FieldInfo struct {
	*r.StructField        // current field
	Target         r.Type // current target of traversal
	Index          int    // index of field in target
	Path           []int  // path from original target to current target
	parent         r.Type // parent type, it has been found
	pindex         int    // next of parent field in rtype, if parent is non-nil
}

func (f *FieldInfo) IsParent() bool {
	return f.parent != nil && f.Index == f.pindex
}

func (f *FieldInfo) IsPublic() bool {
	return len(f.PkgPath) == 0
}

func (f *FieldInfo) FullPath() []int {
	return append(f.Path, f.Index)
}

// fill the field info for its current index
func updateInfo(n *FieldInfo) (okay bool) {
	for {
		if n.Index < n.Target.NumField() {
			field := n.Target.Field(n.Index)
			if n.parent == nil && field.Anonymous && field.Type.Kind() == r.Struct {
				n.parent = field.Type
				n.pindex = n.Index
			}
			n.StructField = &field
			okay = true
			break
		} else if n.parent == nil {
			n.StructField = nil
			okay = false
			break
		} else {
			// move up
			n.Target, n.parent = n.parent, nil
			n.Path = append(n.Path, n.pindex)
			n.Index = 0
		}
	}
	return
}
