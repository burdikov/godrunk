package main

import "testing"

func TestCreateDeck(t *testing.T) {
	deck := createDeck()

	if deck == nil {
		t.Error("Deck is nil!")
	}
}
