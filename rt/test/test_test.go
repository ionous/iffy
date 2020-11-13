package test

import (
	"testing"

	"github.com/kr/pretty"
)

func TestKindsForType(t *testing.T) {
	var ks Kinds
	ks.Add((*GroupPartition)(nil))
	pretty.Println(ks.Kind("GroupPartition"))
}
