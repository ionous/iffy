package assembly

import (
	"fmt"

	"github.com/ionous/iffy/tables"
)

// MissingKinds reports named kinds which don't have a defined ancestry.
// Returns error only on critical errors.
func reportMissingKinds(asm *Assembler) error {
	return reportMissing(asm, "kind",
		`select 1 from mdl_kind k 
			where name= k.kind`)
}

// MissingFields returns named fields which don't have a defined property.
func reportMissingFields(asm *Assembler) error {
	return reportMissing(asm, "field",
		`select 1 from mdl_field p
		where name = p.field`)
}

//
func reportMissingTraits(asm *Assembler) error {
	return reportMissing(asm, "trait",
		`select 1 from mdl_aspect r
		where name = r.trait`)
}

//
func reportMissingAspects(asm *Assembler) error {
	return reportMissing(asm, "aspect",
		`select 1 from mdl_aspect r
		where name = r.aspect`)
}

// fix:
// MissingDefaults returns named kind,field pairs which dont exist paired together
// func MissingDefaults(db *sql.DB, cb func(kind, field string) error) error {
// 	return errutil.New("needs to incorporate fields, traits, and aspects")
// }

func reportMissing(asm *Assembler, cat, exists string) error {
	var name, source, offset string
	q := fmt.Sprintf(`select name, src, offset 
		from eph_named n
		left join eph_source src 
			on (n.idSource = src.rowid)
		where n.category = "%s" 
		and not exists (%s)`, cat, exists)
	return tables.QueryAll(asm.cache.DB(), q,
		func() error {
			asm.reportIssuef(source, offset, "undeclared or missing %s: %q", cat, name)
			return nil
		},
		&name, &source, &offset)
}
