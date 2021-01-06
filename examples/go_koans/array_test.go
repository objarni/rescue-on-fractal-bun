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

func ExampleLoopOverArrayIndices() {
	ia := []int{
		1, 2, 3,
	}
	for i := range ia {
		fmt.Print(i)
	}
	// Output:
	// 012
}

func ExampleLoopOverArrayValues() {
	ia := []int{
		1, 2, 3,
	}
	for _, val := range ia {
		fmt.Print(val)
	}
	// Output:
	// 123
}
