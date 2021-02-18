package go_koans

import (
	"fmt"
	"strings"
)

func Example_loopOverCharacters() {
	for _, ch := range "abc" {
		fmt.Printf("%v-", string(ch))
	}
	// Output:
	// a-b-c-
}

func Example_skipFirstCharacter() {
	fmt.Printf("abc"[1:])
	// Output:
	// bc
}

func Example_joinStrings() {
	fmt.Printf(strings.Join([]string{"ab", "bc", "cd"}, ","))
	// Output:
	// ab,bc,cd
}
func Example_splitString() {
	split := strings.Split("ab,bc,cd", ",")
	for _, str := range split {
		fmt.Println(str)
	}
	// Output:
	// ab
	// bc
	// cd
}
