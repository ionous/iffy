package rt

// Finder returns a new runtime that checks the passed finder before checking the runtime.
func Finder(run Runtime, w ObjectFinder) Runtime {
	return _Finder{run, w}
}

// AtFinder returns a new runtime that matches the name "@" against the passed object.
func AtFinder(run Runtime, obj Object) Runtime {
	return Finder(run, _AtFinder{obj})
}

type _Finder struct {
	Runtime
	Finder ObjectFinder
}

func (l _Finder) FindObject(name string) (ret Object, okay bool) {
	if obj, ok := l.Finder.FindObject(name); ok {
		ret, okay = obj, ok
	} else {
		ret, okay = l.Runtime.FindObject(name)
	}
	return
}

type _AtFinder struct {
	obj Object
}

func (l _AtFinder) FindObject(name string) (ret Object, okay bool) {
	if name == "@" {
		ret, okay = l.obj, true
	}
	return
}
