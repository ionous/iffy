package assembly

import (
	"database/sql"
	"strings"

	"github.com/ionous/errutil"
)

// goal:
// . finalized table of property, kind, type.
// considerations
// . property's lowest common ancestor ( lca )
// o ambiguity when properties collapse into "thing" ( lca implicitly the root )
// . contradiction in type
// . missing properties ( named but not specified )
// o misspellings, near spellings
func DetermineFields(w *Writer, db *sql.DB) (err error) {
	if it, e := db.Query(
		// fix: probably want source line out of this too
		`select nk.name as kind, nf.name as field, p.primType as type, a.path as parents
		from primitive p join named nk
			on (p.idNamedKind = nk.rowid)
		left join named nf
			on (p.idNamedField = nf.rowid)
		left join ancestry a
			on (a.kind = nk.name)
		order by nf.name, nk.name
		`); e != nil {
		err = e
	} else {
		var last fieldInfo // holds our lca
		defer it.Close()
		for it.Next() {
			var curr fieldInfo
			if e := it.Scan(&curr.Kind, &curr.Field, &curr.Type, &curr.Parents); e != nil {
				err = e
				break
			} else {
				// we're at a new field, so write the old one.
				if last.Field != curr.Field {
					last.Flush(w)
					curr.updateHierarchy()
					last = curr
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
							break
						} else if last.Kind != curr.Kind {
							// warn if we have collapsed into root?
							overlap := findOverlap(last.Hierarchy, curr.updateHierarchy())
							last, last.Hierarchy = curr, overlap
						}
					}
				}
			}
		}
		if e := it.Err(); e != nil {
			err = errutil.Append(err, e)
		} else {
			last.Flush(w)
		}
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
