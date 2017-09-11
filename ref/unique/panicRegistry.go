package unique

// PanicBlocks wraps RegisterBlocks to panic on error.
func PanicBlocks(reg TypeRegistry, blks ...interface{}) {
	if e := RegisterBlocks(reg, blks...); e != nil {
		panic(e)
	}
}

// PanicBlocks wraps RegisterTypes to panic on error.
func PanicTypes(reg TypeRegistry, ptrs ...interface{}) {
	if e := RegisterTypes(reg, ptrs...); e != nil {
		panic(e)
	}
}

// PanicValues wraps RegisterValues to panic on error.
func PanicValues(reg ValueRegistry, vals ...interface{}) {
	if e := RegisterValues(reg, vals...); e != nil {
		panic(e)
	}
}
