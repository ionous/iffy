package unique

import r "reflect"

type _PanicTypes struct {
	TypeRegistry
}

// PanicTypes wraps a registry inorder to panic on any error it encounters.
func PanicTypes(r TypeRegistry) TypeRegistry {
	return _PanicTypes{r}
}

func (r _PanicTypes) RegisterType(rtype r.Type) error {
	if e := r.TypeRegistry.RegisterType(rtype); e != nil {
		panic(e)
	}
	return nil
}

type _PanicValues struct {
	ValueRegistry
}

// PanicValues wraps a registry inorder to panic on any error it encounters.
func PanicValues(r ValueRegistry) ValueRegistry {
	return _PanicValues{r}
}

func (r _PanicValues) RegisterValue(rtype r.Value) error {
	if e := r.ValueRegistry.RegisterValue(rtype); e != nil {
		panic(e)
	}
	return nil
}
