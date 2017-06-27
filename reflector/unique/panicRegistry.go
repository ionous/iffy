package unique

import r "reflect"

type _PanicRegistry struct {
	Registry
}

// PanicRegistry wraps a registry inorder to panic on any error it encounters.
func PanicRegistry(r Registry) Registry {
	return _PanicRegistry{r}
}

func (r _PanicRegistry) RegisterType(rtype r.Type) error {
	if e := r.Registry.RegisterType(rtype); e != nil {
		panic(e)
	}
	return nil
}
