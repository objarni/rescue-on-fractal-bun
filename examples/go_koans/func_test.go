package go_koans

import (
	"fmt"
)

type LittleThing struct {
	i int
	s string
}

func modifyThing(thing LittleThing) {
	thing.i = thing.i + 1
	thing.s += "def"
}

func ExampleStructsArePassByValue() {
	aLittleThing := LittleThing{i: 1, s: "abc"}
	fmt.Println(aLittleThing)
	// Output:
	// {1 abc}
}
