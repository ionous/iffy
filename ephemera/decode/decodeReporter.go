package decode

import (
	"github.com/ionous/iffy/ephemera/reader"
)

type IssueReport func(reader.Position, error)

func NewDecoderReporter(source string, report IssueReport) *Decoder {
	dec := &Decoder{source: source, cmds: make(map[string]cmdRec), issueFn: report}
	return dec
}

func (m *Decoder) report(ofs string, err error) {
	m.issueFn(reader.Position{m.source, ofs}, err)
	m.IssueCount++
}
