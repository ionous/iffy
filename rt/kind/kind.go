package kind

import (
	r "reflect"
)

// IsInteger returns true if the passed kind is an IntegerType
// https://golang.org/pkg/builtin/#IntegerType
func IsInteger(k r.Kind) (ret bool) {
	switch k {
	case r.Int, r.Int8, r.Int16, r.Int32, r.Int64, r.Uint, r.Uint8, r.Uint16, r.Uint32, r.Uint64:
		ret = true
	}
	return
}

// IsFloat returns true if the passed kind is a FloatType
// https://golang.org/pkg/builtin/#FloatType
func IsFloat(k r.Kind) (ret bool) {
	switch k {
	case r.Float32, r.Float64:
		ret = true
	}
	return
}

// IsNumber returns true if the passed kind is a float or integer
func IsNumber(k r.Kind) bool {
	return IsInteger(k) || IsFloat(k)
}
