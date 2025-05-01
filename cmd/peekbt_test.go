package main

import "testing"

func Example_() {
	goMain([]string{"peekbt!"})
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