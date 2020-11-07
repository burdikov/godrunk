package main

import "testing"

func TestCreateDeck(t *testing.T) {
	deck := createDeck()

	if deck == nil {
		t.Error("Deck is nil!")
	}
}

func BenchmarkCreateDeckPointers(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = createDeck()
	}
}
