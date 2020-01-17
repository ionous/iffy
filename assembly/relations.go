package assembly

import (
	"database/sql"
	"strings"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dbutil"
	_ "github.com/mattn/go-sqlite3"
)

type relInfo struct {
	relation, cardinality string
	kind, otherKind       hierarchy
}

func (p *relInfo) flush(w *Modeler) (err error) {
	if len(p.relation) > 0 {
		if !p.kind.valid || !p.otherKind.valid {
			err = errutil.New("couldnt determine valid lowest common ancestor")
		} else {
			err = w.WriteRelation(p.relation, p.kind.name, p.cardinality, p.otherKind.name)
		}
	}
	return
}

// by default, lca contains kind --
// but it might not once new kinds are merged into it.
type hierarchy struct {
	name    string
	parents string   // mdl hierarchy of kind
	lca     []string // root is on the right.
	valid   bool     // valid if lca is a named part
}

// normalize name, parents into an array of kinds.
func (h *hierarchy) getAncestry() []string {
	return append([]string{h.name}, strings.Split(h.parents, ",")...)
}

func (h *hierarchy) set(lca []string) {
	h.lca, h.valid = lca, len(lca) > 1
}
func (h *hierarchy) update(other *hierarchy) {
	if h.name != other.name {
		cmp, lca := findOverlap(h.lca, other.getAncestry())
		h.name = other.name
		h.valid = cmp != 0
		h.lca = lca
	}
}

// in, eph_relation: R, K, cardinality, Q
// out, mdl_rel: R, K(lca), Q(lca), cardinality
// fix? right now the coalesce allows missing kinds through,
// the behavior otherwise is Scan error on column index 5, and not particularly helpful
func DetermineRelations(w *Modeler, db *sql.DB) (err error) {
	var curr, last relInfo
	// we select by R, sorted by R, C, K, Q
	// when C differs, we error.
	// when K differs, we LCA K.
	// when Q differs, we LCA Q.
	if e := dbutil.QueryAll(db,
		`select
			nr.name,
			r.cardinality,
			nk.name, coalesce(ak.path, ""),
			nq.name, coalesce(aq.path, "")
		from eph_relation r join eph_named nr
			on (r.idNamedRelation = nr.rowid)
		left join eph_named nk
			on (r.idNamedKind = nk.rowid)
		left join eph_named nq
			on (r.idNamedOtherkind = nq.rowid)
		left join mdl_ancestry ak
			on (ak.kind = nk.name)
		left join mdl_ancestry aq
			on (aq.kind = nq.name)
		order by nr.name, r.cardinality, nk.name, nq.name
		`, func() (err error) {
			// when R differs, write to the output.
			if last.relation != curr.relation {
				if e := last.flush(w); e != nil {
					err = e
				} else {
					// move curr into last for the next queried row.
					curr.kind.set(curr.kind.getAncestry())
					curr.otherKind.set(curr.otherKind.getAncestry())
					last = curr
				}
			} else if last.cardinality != curr.cardinality {
				// same relation can't have different cardinality(s)
				err = errutil.New("cardinality mismatch", curr.relation, last.cardinality, curr.cardinality)
			} else {
				last.kind.update(&curr.kind)
				last.otherKind.update(&curr.otherKind)
			}
			return
		},
		&curr.relation,
		&curr.cardinality,
		&curr.kind.name, &curr.kind.parents,
		&curr.otherKind.name, &curr.otherKind.parents,
	); e != nil {
		err = e
	} else {
		err = last.flush(w)
	}
	return
}
