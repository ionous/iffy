package list_test

// func TestPop(t *testing.T) {
// 	fruit := []string{"Orange", "Lemon", "Mango"}
// 	//
// 	// if got := pop(fruit, true, 1); got != "Orange; Lemon, Mango" {
// 	// 	t.Fatal(got)
// 	// }
// 	if got := pop(fruit, false, 1); got != "Mango; Orange, Lemon" {
// 		t.Fatal(got)
// 	}
// 	// if got := pop(fruit, true, 0); got != "-; Orange, Lemon, Mango" {
// 	// 	t.Fatal(got)
// 	// }
// 	// if got := pop(fruit, false, 0); got != "-; Orange, Lemon, Mango" {
// 	// 	t.Fatal(got)
// 	// }
// }

// // removed; remaining
// func pop(src []string, front bool, amt int) (ret string) {
// 	run := listTime{strings: append([]string{}, src...)}
// 	removed := joinText(&run, &list.Pop{"strings", list.FrontOrBack(front), I(amt)})
// 	remaining := joinStrings(run.strings) // get the variable set by splice
// 	return removed + "; " + remaining
// }
