package list_test

import (
	"fmt"
	"strings"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/list"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/generic"
)

func ExampleSlices() {
	fruit := []string{"Banana", "Orange", "Lemon", "Apple", "Mango"}
	slice(0, 0, fruit)
	slice(3, 0, fruit)  // start uses a one based index
	slice(3, 5, fruit)  // not including the fifth element
	slice(2, 8, fruit)  // further than the end
	slice(-2, 0, fruit) // starting from, and including, the second to last element
	slice(3, -1, fruit) // the third el, up to, though not including, the last
	slice(30, 0, fruit) // start beyond the end
	slice(0, -9, fruit) // ending before the front
	slice(3, 1, fruit)  // reverse order
	slice(4, -4, fruit) // math reverse order

	// Output:
	//  0, 0: Banana, Orange, Lemon, Apple, Mango
	//  3, 0: Lemon, Apple, Mango
	//  3, 5: Lemon, Apple
	//  2, 8: Orange, Lemon, Apple, Mango
	// -2, 0: Apple, Mango
	//  3,-1: Lemon, Apple
	// 30, 0: -
	//  0,-9: -
	//  3, 1: -
	//  4,-4: -
}

func slice(start, end int, src []string) {
	run := sliceTime{strings: src}
	slice := &list.Slice{"strings", &core.Number{float64(start)}, &core.Number{float64(end)}}
	var s string
	{
		if vs, e := slice.GetTextList(&run); e != nil {
			s = e.Error()
		} else if len(vs) > 0 {
			s = strings.Join(vs, ", ")
		} else {
			s = "-"
		}
	}
	fmt.Println(fmt.Sprintf("%2d,%2d:", start, end), s)
}

type sliceTime struct {
	rt.Panic
	strings []string
}

func (g *sliceTime) GetField(target, field string) (ret rt.Value, err error) {
	if target != object.Variables {
		err = errutil.New("unexpected target", target)
	} else if field == "strings" {
		ret = &generic.StringSlice{Values: g.strings}
	} else {
		err = errutil.New("unexpected field", field)
	}
	return
}
