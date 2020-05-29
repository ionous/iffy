package assembly

// goal: build table of mdl_default(kind,field,value) for archetypes.
// considerations:
// . property's actual kind ( default specified against a derived type )
// . contradiction in specified values
// . contradiction in specified value vs field type ( alt: implicit conversion )
// . missing properties ( kind, field pair doesn't exist in model )
// o certainties: usually, seldom, never, always.
// o misspellings, near spellings ( ex. for missing fields )
func AssembleDefaults(asm *Assembler) (err error) {
	if e := assembleDefaultFields(asm); e != nil {
		err = e
	} else if e := assembleDefaultTraits(asm); e != nil {
		err = e
	}
	return
}
