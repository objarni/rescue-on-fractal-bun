package go_koans

import "fmt"

func ExampleTraditionalForLoop() {
	for i := 0; i < 100; i += 10 {
		fmt.Print(i)
	}
	// Output:
	// 0102030405060708090
}
