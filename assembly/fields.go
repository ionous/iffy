package assembly

import (
	"database/sql"
	"sort"
	"strings"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/tables"
)

// goal: build table of property, kind, type.
// considerations:
// . property's lowest common ancestor ( lca )
// o ambiguity when properties collapse into root ( and/or an implicit kind )
// . contradiction in type ( alt: runtime can fit into property by type )
// . missing properties ( named but not specified )
// o misspellings, near spellings ( ex. for missing fields )
func DetermineFields(m *Modeler, db *sql.DB) (err error) {
	// select primitive aspects which arent named in aspects
	// the primitive field's name is the aspect name
	var out pendingFields
	if missingAspects, e := undeclaredAspects(db); e != nil {
		err = e
	} else if e := out.determineFields(db, missingAspects); e != nil {
		err = e
	} else {
		err = out.write(m)
	}
	return
}

// we cant read and write to the same db simultaneously
type pendingField struct {
	field, target, fieldType string
}
type pendingFields struct {
	list []pendingField
}

func (out *pendingFields) write(m *Modeler) (err error) {
	for _, f := range out.list {
		if e := m.WriteField(f.target, f.field, f.fieldType); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

// fix? is it possible to use upsert to allow us to ratchet up the hierarchy?
func (out *pendingFields) determineFields(db *sql.DB, missingAspects []string) (err error) {
	var curr, last fieldInfo
	// fix: probably want source line out of this too
	if e := tables.QueryAll(db,
		`select nk.name, nf.name, p.primType, a.path
		from eph_primitive p join eph_named nk
			on (p.idNamedKind = nk.rowid)
		left join eph_named nf
			on (p.idNamedField = nf.rowid)
		left join mdl_kind a
			on (a.kind = nk.name)
		order by nf.name, nk.name
		`, func() (err error) {
			// we're at a new field, so write the old one.
			if last.Field != curr.Field {
				curr.updateHierarchy()
				if curr.Type == tables.PRIM_ASPECT && sort.SearchStrings(missingAspects, curr.Field) >= 0 {
					err = errutil.New("unknown aspect declared as field of kind", curr.Field, curr.Kind)
				} else {
					last.Flush(out)
					last = curr
				}
			} else if last.Type != curr.Type {
				// field is the same but the type of the field differs
				// future: allow the same named field in different kinds.
				e := errutil.New("type mismatch", last.Field, last.Type, "!=", curr.Type)
				err = e
			} else if last.Kind != curr.Kind {
				// field and type are the same, kind differs; find a common type for the field.
				// currently, there always will be a valid one: the root.
				// warn if we have collapsed into root?
				_, overlap := findOverlap(last.Hierarchy, curr.updateHierarchy())
				last, last.Hierarchy = curr, overlap
			}
			return
		}, &curr.Kind, &curr.Field, &curr.Type, &curr.Parents); e != nil {
		err = e
	} else {
		last.Flush(out)
	}
	return
}

func undeclaredAspects(db *sql.DB) (ret []string, err error) {
	var str string
	var aspects []string
	if e := tables.QueryAll(db,
		`select name from
			( select distinct n.name as name
				from eph_primitive p, eph_named n
				where p.primType = 'aspect'
				and p.idNamedField = n.rowid )
		where name not in 
			( select n.name
				from eph_aspect a, eph_named n
				where a.idNamedAspect = n.rowid )
		`, func() (err error) {
			aspects = append(aspects, str)
			return
		},
		&str); e != nil {
		err = e
	} else {
		sort.Strings(aspects)
		ret = aspects
	}
	return
}

// given a two hierarchies, return where they overlap.
// if the returned list is the same as b, return 1
// if not, and the returned list is the same as a, return -1
// otherwise, if the return list is shorter than a or b, return 0.
func findOverlap(a, b []string) (retCmp int, retOvr []string) {
	// root is on the right
	if acnt, bcnt := len(a), len(b); acnt != 0 && bcnt != 0 {
		retCmp = 1 // preliminarily, lets assume they are equal
		if acnt > bcnt {
			a = a[acnt-bcnt:]
		} else if bcnt > acnt {
			b = b[bcnt-acnt:]
			retCmp = -1 // a might still be the same
		}
		// now they're the same length;
		// they might be the same hierarchy
		for i, ael := range a {
			if bel := b[i]; bel == ael {
				retOvr = b[i:]
				break
			}
			// if the first element didnt match
			// then neither is a match
			retCmp = 0
		}
	}
	return
}

type fieldInfo struct {
	Kind, Field, Type, Parents string   // parents depends on type
	Hierarchy                  []string // hierarchy is derived from type and parents
}

// ancestors holds lca
func (i *fieldInfo) Flush(out *pendingFields) {
	if len(i.Kind) > 0 {
		lca := i.Hierarchy[0]
		out.list = append(out.list, pendingField{i.Field, lca, i.Type})
	}
}

// split Parents into a slice of strings
func (i *fieldInfo) updateHierarchy() []string {
	i.Hierarchy = append([]string{i.Kind}, strings.Split(i.Parents, ",")...)
	return i.Hierarchy
}
