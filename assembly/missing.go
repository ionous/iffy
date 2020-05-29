package assembly

import (
	"fmt"

	"github.com/ionous/iffy/tables"
)

// MissingKinds reports named kinds which don't have a defined ancestry.
// Returns error only on critical errors.
func reportMissingKinds(asm *Assembler) error {
	return reportMissing(asm, "kind",
		`select 1 from mdl_kind mk 
			where name= mk.kind`)
}

// fix? what about alternative names ( aliases and partial names ) in the story text?
func reportMissingNouns(asm *Assembler) error {
	return reportMissing(asm, "noun",
		`select 1 from mdl_noun mn 
			where name= mn.noun`)
}

// MissingFields returns named fields which don't have a defined property.
func reportMissingFields(asm *Assembler) error {
	return reportMissing(asm, "field",
		`select 1 from mdl_field mf
		where name = mf.field`)
}

//
func reportMissingTraits(asm *Assembler) error {
	return reportMissing(asm, "trait",
		`select 1 from mdl_aspect ma
		where name = ma.trait`)
}

//
func reportMissingAspects(asm *Assembler) error {
	return reportMissing(asm, "aspect",
		`select 1 from mdl_aspect ma
		where name = ma.aspect`)
}

// fix:
// MissingDefaults returns named kind,field pairs which dont exist paired together
// func MissingDefaults(db *sql.DB, cb func(kind, field string) error) error {
// 	return errutil.New("needs to incorporate fields, traits, and aspects")
// }

func reportMissing(asm *Assembler, cat, exists string) error {
	var name, source, offset string
	q := fmt.Sprintf(`select en.name, es.src, en.offset 
		from eph_named en
		left join eph_source es 
			on (en.idSource = es.rowid)
		where en.category = "%s" 
		and not exists (%s)`, cat, exists)
	return tables.QueryAll(asm.cache.DB(), q,
		func() error {
			asm.reportIssuef(source, offset, "undeclared or missing %s: %q", cat, name)
			return nil
		},
		&name, &source, &offset)
}
