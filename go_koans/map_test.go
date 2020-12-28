package go_koans

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeclaringAMap(t *testing.T) {
	var myMap map[int]string
	assert.Nil(t, myMap)
}

func TestInitializeEmptyMap(t *testing.T) {
	myMap := make(map[int]string)
	assert.NotNil(t, myMap)
}

func TestInitializeMapWithEntries(t *testing.T) {
	myMap := map[int]string{
		1: "One",
		2: "Two",
	}
	assert.Equal(t, myMap[1], "One")
	assert.Equal(t, myMap[2], "Two")
}

func TestOverwriteExistingEntry(t *testing.T) {
	myMap := map[int]string{
		1: "One",
		2: "Two",
	}
	myMap[2] = "TWO"
	assert.Equal(t, myMap[1], "One")
	assert.Equal(t, myMap[2], "TWO")
}

func TestAddEntry(t *testing.T) {
	myMap := map[int]string{
		1: "One",
		2: "Two",
	}
	myMap[3] = "Three"
	assert.Equal(t, myMap[1], "One")
	assert.Equal(t, myMap[2], "Two")
	assert.Equal(t, myMap[3], "Three")
}

func TestNumberOfEntries(t *testing.T) {
	myMap := map[int]string{
		1: "One",
		2: "Two",
	}
	assert.Equal(t, 2, len(myMap))
}

func TestDeleteAnEntry(t *testing.T) {
	myMap := map[int]string{
		1: "One",
		2: "Two",
	}
	delete(myMap, 1)
	assert.Equal(t, 1, len(myMap))
}

func TestCheckExistingKey(t *testing.T) {
	myMap := map[int]string{
		1: "One",
		2: "Two",
	}
	valueExists := myMap[1]
	assert.Equal(t, valueExists, "One")
}

func TestNonExistingKeyMeansZeroValue(t *testing.T) {
	myMap := map[int]string{
		1: "One",
		2: "Two",
	}
	valueNotExists := myMap[5]
	assert.Equal(t, valueNotExists, "")
}

func TestKeyMayExistIdiom(t *testing.T) {
	myMap := map[int]string{
		1: "One",
		2: "Two",
	}
	if val, ok := myMap[1]; ok {
		assert.Equal(t, val, "One")
		assert.True(t, ok)
	}
}

func TestLoopOverMap(t *testing.T) {
	myMap := map[int]string{
		1: "One",
		2: "Two",
	}
	for key, value := range myMap {
		assert.NotNil(t, key)
		assert.NotNil(t, value)
	}
}

// 		for key, control := range controllerMap {
