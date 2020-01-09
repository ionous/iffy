package assembly

import (
	"github.com/ionous/iffy/ephemera"
)

func NewWriter(q ephemera.Queue) *Writer {
	q.Prep("mdl_ancestry",
		ephemera.Col{Name: "kind", Type: "text"},
		ephemera.Col{Name: "path", Type: "text"},
	)
	//
	q.Prep("mdl_property",
		ephemera.Col{Name: "field", Type: "text"},
		ephemera.Col{Name: "kind", Type: "text"},
		ephemera.Col{Name: "type", Type: "text"},
	)
	q.Prep("mdl_rank",
		ephemera.Col{Name: "aspect", Type: "text"},
		ephemera.Col{Name: "trait", Type: "text"},
		ephemera.Col{Name: "rank", Type: "int"},
	)
	return &Writer{q}
}

type Writer struct {
	q ephemera.Queue
}

// write kind and comma separated ancestors
func (w *Writer) WriteAncestor(kind, path string) {
	w.q.Write("mdl_ancestry", kind, path)
}

func (w *Writer) WriteField(field, owner, fieldType string) {
	w.q.Write("mdl_property", field, owner, fieldType)
}

func (w *Writer) WriteTrait(aspect, trait string) {
	w.q.Write("mdl_rank", aspect, trait, 0)
}
