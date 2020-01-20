package assembly

import (
	"database/sql"

	"github.com/ionous/iffy/dbutil"
)

// MissingKinds returns named kinds which don't have a defined ancestry.
func MissingKinds(db *sql.DB, cb func(string) error) error {
	var k string
	return dbutil.QueryAll(db,
		`select distinct name from eph_named n
		where n.category = 'kind'
		and not exists (
			select 1 from mdl_kind a
			where n.name = a.kind
		)`,
		func() error { return cb(k) },
		&k)
}

// MissingFields returns named fields which don't have a defined property.
func MissingFields(db *sql.DB, cb func(string) error) error {
	var k string
	return dbutil.QueryAll(db,
		`select distinct n.name from eph_named n
		where n.category = 'field'
		and not exists (
			select 1 from mdl_field p
			where n.name = p.field
		)`,
		func() error { return cb(k) },
		&k)
}

//
func MissingTraits(db *sql.DB, cb func(string) error) error {
	var k string
	return dbutil.QueryAll(db,
		`select distinct n.name from eph_named n
		where n.category = 'trait'
		and not exists (
			select 1 from mdl_aspect r
			where n.name = r.trait
		)`,
		func() error { return cb(k) },
		&k)
}

//
func MissingAspects(db *sql.DB, cb func(string) error) error {
	var k string
	return dbutil.QueryAll(db,
		`select distinct n.name from eph_named n
		where n.category = 'aspect'
		and not exists (
			select 1 from mdl_aspect r
			where n.name = r.aspect
		)`,
		func() error { return cb(k) },
		&k)
}
