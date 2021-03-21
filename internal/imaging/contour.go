package imaging

type Pos struct {
	X, Y int
}

type BitField struct {
	Field         map[Pos]bool
	Width, Height int
}

func (im BitField) IsSet(x int, y int) bool {
	return im.Field[Pos{x, y}]
}

func FindContour(bitField BitField) BitField {
	// Set bits in the bitField means material
	// Clear bits in the bitField means transparent
	contourPixels := map[Pos]bool{}
	for y := 0; y < bitField.Height; y++ {
		for x := 0; x < bitField.Width; x++ {
			// Only transparent pixels relevant
			if !bitField.IsSet(x, y) {
				// If any if the 8 neighbours is
				// not transparent, fill this one
				if bitField.IsSet(x-1, y-1) ||
					bitField.IsSet(x, y-1) ||
					bitField.IsSet(x+1, y-1) ||
					bitField.IsSet(x-1, y) ||
					bitField.IsSet(x+1, y) ||
					bitField.IsSet(x-1, y+1) ||
					bitField.IsSet(x, y+1) ||
					bitField.IsSet(x+1, y+1) {
					contourPixels[Pos{x, y}] = true
				}
			}
		}
	}
	return BitField{
		Field:  contourPixels,
		Width:  bitField.Width,
		Height: bitField.Height,
	}
}
