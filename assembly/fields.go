package assembly

import (
	"database/sql"
	"sort"
	"strings"

	"github.com/ionous/iffy/ephemera"

	"github.com/ionous/errutil"
)

// goal: build table of property, kind, type.
// considerations:
// . property's lowest common ancestor ( lca )
// o ambiguity when properties collapse into root ( and/or an implicit kind )
// . contradiction in type ( alt: runtime can fit into property by type )
// . missing properties ( named but not specified )
// o misspellings, near spellings ( ex. for missing fields )
func DetermineFields(w *Writer, db *sql.DB) (err error) {
	// select primitive aspects which arent named in aspects
	// the primitive field's name is the aspect name
	if missingAspects, e := undeclaredAspects(db); e != nil {
		err = e
	} else {
		var curr, last fieldInfo
		// fix: probably want source line out of this too
		if e := queryAll(db,
			`select nk.name as kind, nf.name as field, p.primType as type, a.path as parents
		from primitive p join named nk
			on (p.idNamedKind = nk.rowid)
		left join named nf
			on (p.idNamedField = nf.rowid)
		left join ancestry a
			on (a.kind = nk.name)
		order by nf.name, nk.name
		`, func() (err error) {
				// we're at a new field, so write the old one.
				if last.Field != curr.Field {
					curr.updateHierarchy()
					if curr.Type == ephemera.PRIM_ASPECT && sort.SearchStrings(missingAspects, curr.Field) >= 0 {
						err = errutil.New("unknown aspect declared as field of kind", curr.Field, curr.Kind)
					} else {
						last.Flush(w)
						last = curr
					}
				} else {
					if len(last.Kind) == 0 {
						curr.updateHierarchy()
						last = curr
					} else {
						// if the types have changed
						// for now we yield a mismatch
						// ( we could allow the same named field in different kinds if we wanted )
						if last.Type != curr.Type {
							e := errutil.New("type mismatch", last.Field, last.Type, "!=", curr.Type)
							err = e
						} else if last.Kind != curr.Kind {
							// warn if we have collapsed into root?
							overlap := findOverlap(last.Hierarchy, curr.updateHierarchy())
							last, last.Hierarchy = curr, overlap
						}
					}
				}
				return
			}, &curr.Kind, &curr.Field, &curr.Type, &curr.Parents); e != nil {
			err = e
		} else {
			last.Flush(w)
		}
	}
	return
}

func undeclaredAspects(db *sql.DB) (ret []string, err error) {
	var str string
	var aspects []string
	if e := queryAll(db,
		`select name from
			( select distinct n.name as name
				from primitive p, named n
				where p.primType = 'aspect'
				and p.idNamedField = n.rowid )
		where name not in 
			( select n.name
				from aspect a, named n
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

// return the chain of lca
func findOverlap(a, b []string) (ret []string) {
	// root is on the right
	adepth, bdepth := len(a), len(b)
	if adepth > bdepth {
		a = a[adepth-bdepth:]
	} else if bdepth > adepth {
		b = b[bdepth-adepth:]
	}
	for i, ael := range a {
		if bel := b[i]; bel == ael {
			ret = b[i:]
			break
		}
	}
	return ret
}

type fieldInfo struct {
	Kind, Field, Type, Parents string   // parents depends on type
	Hierarchy                  []string // hierarchy is derived from type and parents
}

// ancestors holds lca
func (i *fieldInfo) Flush(w *Writer) {
	if len(i.Kind) > 0 {
		lca := i.Hierarchy[0]
		w.WriteField(i.Field, lca, i.Type)
	}
}

func (i *fieldInfo) updateHierarchy() []string {
	i.Hierarchy = append([]string{i.Kind}, strings.Split(i.Parents, ",")...)
	return i.Hierarchy
}
