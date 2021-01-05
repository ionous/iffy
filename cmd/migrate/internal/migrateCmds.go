package internal

import (
	"encoding/json"
	"strings"
)

type Target struct {
	Parent string `json:"parent"`
	Field  string `json:"field"`
}

//
type Copy struct {
	From Target `json:"from"`
	To   Target `json:"to"`
}

//
type Replace struct {
	From Target      `json:"from"`
	With interface{} `json:"with"`
}

type JsonData struct {
	Value interface{}
}

func At(parent, field string) Target {
	return Target{parent, field}
}

func Json(s string) (ret interface{}) {
	var data interface{}
	if e := json.Unmarshal([]byte(s), &data); e != nil {
		ret = e
	} else {
		ret = data
	}
	return
}

func (t Target) dequote() string {
	return strings.Replace(t.Parent, "'", `"`, -1)
}

//
type Migration interface {
	Migrate(doc interface{}) (int, error)
}

func (op *Copy) Migrate(doc interface{}) (ret int, err error) {
	return replicate(doc,
		op.From.dequote(), op.From.Field,
		op.To.dequote(), op.To.Field)
}

func (op *Replace) Migrate(doc interface{}) (ret int, err error) {
	if e, ok := op.With.(error); ok {
		err = e
	} else {
		ret, err = op.replace(doc, op.With)
	}
	return
}

func (op *Replace) replace(doc interface{}, val interface{}) (int, error) {
	return replace(doc, op.From.dequote(), op.From.Field, val)
}
