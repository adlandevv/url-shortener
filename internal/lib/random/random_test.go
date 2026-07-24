package random

import (
	"testing"
)

func TestNewRandomStringSize(t *testing.T) {
	sizeNeed := 10

	word := NewRandomString(sizeNeed)

	if len(word) != sizeNeed {
		t.Fatalf("word: %s expected size: %d, actual: %d", word, sizeNeed, len(word))
	}
}

func TestNewRandomStringWhenSizeZero(t *testing.T) {
	sizeNeed := 0

	word := NewRandomString(sizeNeed)

	if len(word) != sizeNeed {
		t.Fatalf("word: %s expected size: %d, actual: %d", word, sizeNeed, len(word))
	}
}
