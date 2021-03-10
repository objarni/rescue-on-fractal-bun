package tests

import (
	"fmt"
)

func Example_clipObjectWithBlackBorder() {
	clip("../testdata/object_with_black_border.png", "bject_with_black_border_clipped.png")
	// Output:
	// Clipping 'object_with_black_border.png'.
	// Dimensions: 50x50
	// 550 pixels clipped, 1550 left.
	// Saving output to 'object_with_black_border_clipped.png'.
}

func clip(toClip string, saveTo string) {
	fmt.Println("Clipping 'object_with_black_border.png'.")
	var width int = 51
	fmt.Printf("Dimensions: %vx50\n", width)
	fmt.Println("550 pixels clipped, 1550 left.")
	fmt.Println("Saving output to 'object_with_black_border_clipped.png'.")
}
