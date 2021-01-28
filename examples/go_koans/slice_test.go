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

func ExampleAppend() {
	ia := []int{
		1, 2, 3,
	}
	newSlice := append(ia, 4)
	fmt.Println(newSlice)
	fmt.Println("len: ", len(newSlice))
	fmt.Println("cap: ", cap(newSlice))
	// Output:
	// 1 2 3 4
	// len: 4
	// cap: 6
}
