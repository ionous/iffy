package assembly

// goal: build table of mdl_start(noun, field, value) for instances.
// considerations:
// . property's actual kind ( default specified against a derived type )
// . contradiction in specified values
// . contradiction in specified value vs field type ( alt: implicit conversion )
// . missing properties ( kind, field pair doesn't exist in model )
// o certainties: usually, seldom, never, always.
// o misspellings, near spellings ( ex. for missing fields )
func AssembleValues(asm *Assembler) (err error) {
	if e := assembleInitialFields(asm); e != nil {
		err = e
	} else if e := assembleInitialTraits(asm); e != nil {
		err = e
	}
	return
}
