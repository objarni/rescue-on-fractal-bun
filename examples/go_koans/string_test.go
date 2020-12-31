package go_koans

import (
	"fmt"
)

func ExampleLoopOverCharacters() {
	for _, ch := range "abc" {
		fmt.Printf("%v-", string(ch))
	}
	// Output:
	// a-b-c-
}

func ExampleSkipFirstCharacter() {
	fmt.Printf("abc"[1:])
	// Output:
	// bc
}
