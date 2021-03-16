package go_koans

import (
	"fmt"
)

func ExampleDeclaringAMap() {
	var myMap map[int]string
	fmt.Printf("myMap is %v", myMap)
	// Output:
	// myMap is map[]
}

func ExampleInitializeEmptyMap() {
	myMap := make(map[int]string)
	fmt.Printf("myMap is %v", myMap)
	// Output:
	// myMap is map[]
}

func ExampleInitializeMapWithEntries() {
	myMap := map[int]string{
		1: "One",
		2: "Two",
	}
	fmt.Printf("myMap is %v", myMap)
	// Output:
	// myMap is map[1:One 2:Two]
}

func ExampleOverwriteExistingEntry() {
	myMap := map[int]string{
		1: "One",
		2: "Two",
	}
	myMap[2] = "TWO"
	fmt.Printf("myMap is %v", myMap)
	// Output:
	// myMap is map[1:One 2:TWO]
}

func ExampleAddEntry() {
	myMap := map[int]string{
		1: "One",
		2: "Two",
	}
	myMap[3] = "Three"
	fmt.Printf("myMap is %v", myMap)
	// Output:
	// myMap is map[1:One 2:Two 3:Three]
}

func ExampleNumberOfEntries() {
	myMap := map[int]string{
		1: "One",
		2: "Two",
	}
	fmt.Printf("myMap length is %v", len(myMap))
	// Output:
	// myMap length is 2
}

func ExampleDeleteAnEntry() {
	myMap := map[int]string{
		1: "One",
		2: "Two",
	}
	delete(myMap, 1)
	fmt.Printf("myMap is %v", myMap)
	// Output:
	// myMap is map[2:Two]
}

func ExampleLookupExistingKey() {
	myMap := map[int]string{
		1: "One",
		2: "Two",
	}
	fmt.Printf("Value at map key 1 is %v", myMap[1])
	// Output:
	// Value at map key 1 is One
}

func ExampleNonExistentKey() {
	myMap := map[int]bool{
		1: true,
		2: true,
	}
	fmt.Println("When element does not exist it is 'default value' of type:")
	fmt.Println(myMap[5])
	// Output:
	// When element does not exist it is 'default value' of type:
	// false
}

func ExampleMayExistIdiom() {
	myMap := map[int]string{
		1: "One",
		2: "Two",
	}
	if val, ok := myMap[1]; ok {
		fmt.Println("Val is", val)
		fmt.Println("ok is", ok)
	}
	if _, ok := myMap[5]; ok {
		fmt.Println("This block will not run")
	}
	// Output:
	// Val is One
	// ok is true
}

func ExampleLoopOverMap() {
	myMap := map[int]string{
		1: "One",
		2: "Two",
	}
	for key, value := range myMap {
		fmt.Println(key, value)
	}
	// Unordered output:
	// 1 One
	// 2 Two
}
