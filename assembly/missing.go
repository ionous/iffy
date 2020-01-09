package assembly

import (
	"database/sql"

	"github.com/ionous/iffy/dbutil"
)

// MissingKinds returns named kinds which don't have a defined ancestry.
func MissingKinds(db *sql.DB, cb func(string)) error {
	var k string
	return dbutil.QueryAll(db,
		`select distinct name from eph_named n
		where not exists (
			select 1 from mdl_ancestry a
			where n.name == a.kind
			and n.category = 'kind'
		)`, func() error { cb(k); return nil }, &k)
}

// MissingFields returns named fields which don't have a defined property.
func MissingFields(db *sql.DB, cb func(string)) error {
	var k string
	return dbutil.QueryAll(db,
		`select distinct n.name from eph_named n
		where not exists (
			select 1 from mdl_property p
			where n.name == p.field
			and n.category = 'field'
		)`, func() error { cb(k); return nil }, &k)
}

//
func MissingTraits(db *sql.DB, cb func(string)) error {
	var k string
	return dbutil.QueryAll(db,
		`select distinct n.name from eph_named n
		where not exists (
			select 1 from mdl_rank r
			where n.name == r.trait
			and n.category = 'trait'
		)`, func() error { cb(k); return nil }, &k)
}

//
func MissingAspects(db *sql.DB, cb func(string)) error {
	var k string
	return dbutil.QueryAll(db,
		`select distinct n.name from eph_named n
		where not exists (
			select 1 from mdl_rank r
			where n.name == r.aspect
			and n.category = 'aspect'
		)`, func() error { cb(k); return nil }, &k)
}
