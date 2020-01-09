package assembly

import (
	"github.com/ionous/iffy/ephemera"
)

type Output struct {
}

func (out *Output) Conflict(e error) {
}
func (out *Output) Ambiguity(e error) {
}

func NewWriter(q ephemera.Queue) *Writer {
	q.Prep("ancestry",
		ephemera.Col{Name: "kind", Type: "text"},
		ephemera.Col{Name: "path", Type: "text"},
	)
	//
	q.Prep("property",
		ephemera.Col{Name: "field", Type: "text"},
		ephemera.Col{Name: "kind", Type: "text"},
		ephemera.Col{Name: "type", Type: "text"},
	)
	return &Writer{q}
}

type Writer struct {
	q ephemera.Queue
}

// write kind and comma separated ancestors
func (w *Writer) WriteAncestor(kind, path string) {
	// fix return error
	w.q.Write("ancestry", kind, path)
}

func (w *Writer) WriteField(field, owner, fieldType string) {
	// fix return error
	w.q.Write("property", field, owner, fieldType)
}
