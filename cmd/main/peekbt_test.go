package main

import (
	"fmt"
	"testing"
)

func Example_hello() {
	fmt.Println(hello())
	// Output:
	// Welcome to peekbt!
}

func TestHello(t *testing.T) {
	got := hello()
	want := "Welcome to peekbt!"
	if got != want {
		t.Errorf("hello() = %q, want %q", got, want)
	}
}
