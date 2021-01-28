package go_koans

import "fmt"

func ExampleSlicingSentToFuncsModifyOriginalArray() {
	ia := []int{
		1, 2, 3,
	}
	setZero(ia[2:3])
	fmt.Println(ia)
	// Output:
	// slice length: 1
	// [1 2 0]
}

func setZero(slice []int) {
	fmt.Println("slice length:", len(slice))
	for i := range slice {
		slice[i] = 0
	}
}
