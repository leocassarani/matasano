package main

import (
	"testing"
)

func TestHammingDistance(t *testing.T) {
	a := []byte("this is a test")
	b := []byte("wokka wokka!!!")

	if dist := HammingDistance(a, a); dist != 0 {
		t.Fatalf("expected 0, got %v", dist)
	}

	if dist := HammingDistance(a, b); dist != 37 {
		t.Fatalf("expected 37, got %v", dist)
	}
}
