package assign_test

import (
	"testing"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/assign"
)

func TestAssignments(t *testing.T) {
	const expectsError = errutil.Error("expects error")
	inputs := []interface{}{nil, true, int(12), int64(0), int64(1), float64(4.4), "hello"}

	t.Run("bool", func(t *testing.T) {
		for i, x := range []interface{}{
			false,        // nil
			true,         // bool
			expectsError, // int:12
			false,        // int64:0
			true,         // int64:1
			expectsError, // float64:4.4
			expectsError, // string
		} {
			var res bool
			if e := assign.Value(&res, inputs[i]); e != nil {
				if x != expectsError {
					t.Fatal("unexpected error at", i, e)
				}
			} else if x := x.(bool); x != res {
				t.Fatal("mismatch at", i)
			}
		}
	})
	t.Run("int", func(t *testing.T) {
		for i, x := range []interface{}{
			int(0),       // nil
			expectsError, // bool:true
			int(12),      // int:12
			int(0),       // int64:0
			int(1),       // int64:1
			int(4),       // float64:4.4
			expectsError, // string:"hello"
		} {
			var res int
			if e := assign.Value(&res, inputs[i]); e != nil {
				if x != expectsError {
					t.Fatal("unexpected error at", i, e)
				}
			} else if x := x.(int); x != res {
				t.Fatal("mismatch at", i)
			}
		}
	})
	t.Run("int64", func(t *testing.T) {
		for i, x := range []interface{}{
			int64(0),     // nil
			expectsError, // bool:true
			int64(12),    // int:12
			int64(0),     // int64:0
			int64(1),     // int64:1
			int64(4),     // float64:4.4
			expectsError, // string:"hello"
		} {
			var res int64
			if e := assign.Value(&res, inputs[i]); e != nil {
				if x != expectsError {
					t.Fatal("unexpected error at", i, e)
				}
			} else if x := x.(int64); x != res {
				t.Fatal("mismatch at", i)
			}
		}
	})
	t.Run("float64", func(t *testing.T) {
		for i, x := range []interface{}{
			float64(0),   // nil
			expectsError, // bool:true
			float64(12),  // int:12
			float64(0),   // int64:0
			float64(1),   // int64:1
			float64(4.4), // float64:4.4
			expectsError, // string:"hello"
		} {
			var res float64
			if e := assign.Value(&res, inputs[i]); e != nil {
				if x != expectsError {
					t.Fatal("unexpected error at", i, e)
				}
			} else if x := x.(float64); x != res {
				t.Fatal("mismatch at", i)
			}
		}
	})
	t.Run("string", func(t *testing.T) {
		for i, x := range []interface{}{
			"",           // nil
			expectsError, // bool
			expectsError, // int
			expectsError, // int64:0
			expectsError, // int64:1
			expectsError, // float64
			"hello",      // string
		} {
			var res string
			if e := assign.Value(&res, inputs[i]); e != nil {
				if x != expectsError {
					t.Fatal("unexpected error at", i, e)
				}
			} else if x := x.(string); x != res {
				t.Fatal("mismatch at", i)
			}
		}
	})
}

func TestAssignmentShortcuts(t *testing.T) {
	t.Run("bool", func(t *testing.T) {
		var v bool
		if e := assign.BoolPtr(&v, v); e != nil {
			t.Fatal("expected success", e)
		} else if e := assign.BoolPtr(v, v); e == nil {
			t.Fatal("expected failure")
		}

	})
	t.Run("float", func(t *testing.T) {
		var v float64
		if e := assign.FloatPtr(&v, v); e != nil {
			t.Fatal("expected success", e)
		} else if e := assign.FloatPtr(v, v); e == nil {
			t.Fatal("expected failure")
		}
	})
	t.Run("string", func(t *testing.T) {
		var v string
		if e := assign.StringPtr(&v, v); e != nil {
			t.Fatal("expected success", e)
		} else if e := assign.StringPtr(v, v); e == nil {
			t.Fatal("expected failure")
		}
	})

}
