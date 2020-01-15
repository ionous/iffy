package ephemera

import "github.com/ionous/errutil"

type GenQueue struct {
	keys   map[string][]string
	Tables map[string]interface{}
}

func (jq *GenQueue) Prep(which string, cols ...Col) {
	if jq.keys == nil {
		jq.keys = make(map[string][]string)
	}
	if jq.Tables == nil {
		jq.Tables = make(map[string]interface{})
	}
	keys := NamesOf(cols)
	jq.keys[which] = keys
}

func (jq *GenQueue) Write(which string, args ...interface{}) (ret Queued, err error) {
	keys := jq.keys[which]
	if len(keys) != len(args) {
		err = errutil.New("mismatched keys for", which)
	} else {
		type Row map[string]interface{}
		type Rows []Row
		row := make(Row)
		for i, k := range keys {
			row[k] = args[i]
		}
		rows, _ := jq.Tables[which].(Rows)
		rows = append(rows, row)
		jq.Tables[which] = rows
		ret = Queued{int64(len(rows))}
	}
	return
}
