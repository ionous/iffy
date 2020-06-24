// Package sds reads compacted self-describing json
package sds

type Kind int

const (
	Slice Kind = iota
	Elem
	Prim
)

type Object interface {
	Kind() Kind

	// annotation of element in a slice
	Note() Note

	// object data in a slice
	// ex. slat in a slot
	Elem() Object

	// ex. primitive value of leaf
	Value() interface{}

	// next part of slice
	Next() Object

	// all parameters in an element
	Params() []string
	Param(string) (Note, Object)
}
