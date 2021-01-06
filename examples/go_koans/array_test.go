package go_koans

import "fmt"

func ExampleDeclareAndInitializeIntArray() {
	ia := []int{
		1, 2, 3,
	}
	fmt.Println(ia)
	// Output:
	// [1 2 3]
}
