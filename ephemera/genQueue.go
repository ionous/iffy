package ephemera

import "log"

type GenQueue struct {
	keys   map[string][]string
	Tables map[string]interface{}
}

//"insert into foo(id, name) values(?, ?)"
func (jq *GenQueue) Prep(which string, keys ...string) {
	if jq.keys == nil {
		jq.keys = make(map[string][]string)
	}
	if jq.Tables == nil {
		jq.Tables = make(map[string]interface{})
	}
	jq.keys[which] = keys
}

func (jq *GenQueue) Write(which string, args ...interface{}) (ret Queued) {
	keys := jq.keys[which]
	if len(keys) != len(args) {
		log.Fatal("mismatched keys for ", which)
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
