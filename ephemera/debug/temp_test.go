package debug

import (
	"encoding/json"
	"testing"
)

func TestTemp(t *testing.T) {
	if b, e := json.MarshalIndent(FactorialZero, "", "  "); e != nil {
		println(e)
	} else {
		println(string(b))
	}
	if b, e := json.MarshalIndent(FactorialSubtract, "", "  "); e != nil {
		println(e)
	} else {
		println(string(b))
	}
}
